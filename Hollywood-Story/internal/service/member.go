package service

import (
	"context"
	"github.com/goccy/go-json"
	"github.com/qm012/dun"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"sync"
	"time"
	"vland.live/app/global"
	"vland.live/app/internal/constant"
	"vland.live/app/internal/model"
	"vland.live/app/internal/model/request"
	"vland.live/app/internal/model/response"
	"vland.live/app/internal/util"
)

type MemberService interface {
	InputInfo(req *request.InputMemberInfoReq) *dun.StatusCode
	GetMemberInfo(req *request.GetMemberInfoReq) (*response.GetMemberInfoResp, *dun.StatusCode)
	RefreshAttributes(req *request.RefreshAttributesReq) (*response.RefreshMemberAttributesResp, *dun.StatusCode)
	StartOrNextRound(req *request.StartOrNextRoundReq) (*response.StartOrNextRoundResp, *dun.StatusCode)
	ClickButtonOutcome(req *request.ClickButtonOutcomeReq) (*response.ClickButtonOutcomeResp, *dun.StatusCode)
	ClickButtonNewsByOutcome(req *request.ClickButtonNewsByOutcomeReq) (*response.ClickButtonNewsByOutcomeResp, *dun.StatusCode)

	UpdateScreenplayReq(req *request.UpdateMemberScreenplayReq) *dun.StatusCode

	findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*model.Member, error)
}

type memberService struct {
	*mongo.Collection
	ScreenplayService
}

var (
	memberServiceOnce sync.Once
	ms                *memberService
)

func NewMemberService() MemberService {
	memberServiceOnce.Do(func() {
		ms = &memberService{
			Collection:        global.Mongo.Database(constant.DatabaseHollywood).Collection("member"),
			ScreenplayService: NewScreenplayService(),
		}
	})
	return ms
}

func (m *memberService) InputInfo(req *request.InputMemberInfoReq) *dun.StatusCode {

	// 删除旧信息
	_, err := m.Collection.DeleteOne(context.Background(), bson.D{{Key: "device_id", Value: req.DeviceID}})
	if err != nil {
		global.Logger.Error("删除旧信息失败", zap.Error(err))
		return dun.StatusCodeSystemError
	}

	specialNPCName := req.SpecialNPCName
	if specialNPCName == "" {
		specialNPCName = "Noah"
	}

	now := time.Now().UnixMilli()
	document := &model.Member{
		ID:             primitive.NewObjectID(),
		DeviceID:       req.DeviceID,
		Nickname:       req.Nickname,
		SpecialNPCName: specialNPCName,
		Age:            req.Age,
		Gender:         req.Gender,
		Personality:    req.Personality,
		Occupation:     req.Occupation,
		Attribute:      constant.TableRandomAttributes(),
		CurrentRound:   nil,
		Film:           nil,
		ModifiedAt:     now,
		CreatedAt:      now,
	}

	insertOneResult, err := m.Collection.InsertOne(context.Background(), document)
	if err != nil {
		global.Logger.Error("设置用户信息失败", zap.Error(err))
		return dun.StatusCodeSystemError
	}

	global.Logger.Info("设置用户信息成功", zap.Any("insertOneResult", insertOneResult))
	return nil
}

func (m *memberService) GetMemberInfo(req *request.GetMemberInfoReq) (*response.GetMemberInfoResp, *dun.StatusCode) {
	member, err := m.searchByDeviceID(req.DeviceID)
	if err != nil {
		global.Logger.Error("根据设备ID获取用户信息失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	if member == nil {
		return nil, dun.StatusCodeDataNotFound
	}

	return &response.GetMemberInfoResp{
		ID:             member.ID.Hex(),
		DeviceID:       member.DeviceID,
		Nickname:       member.Nickname,
		SpecialNPCName: member.SpecialNPCName,
		Age:            member.Age,
		Gender:         member.Gender,
		Personality:    member.Personality,
		Occupation:     member.Occupation,
		Attribute:      member.Attribute,
		CurrentRound:   member.CurrentRound,
		Film:           member.Film,
	}, nil
}

func (m *memberService) RefreshAttributes(req *request.RefreshAttributesReq) (*response.RefreshMemberAttributesResp, *dun.StatusCode) {
	member, err := m.searchByDeviceID(req.DeviceID)
	if err != nil {
		global.Logger.Error("根据设备ID获取用户信息失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	if member == nil {
		return nil, dun.StatusCodeDataNotFound
	}

	if member.CurrentRound != nil {
		return nil, constant.ErrRefreshAttribute
	}

	attribute := constant.TableRandomAttributes()
	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "attribute", Value: attribute})
		fields = append(fields, bson.E{Key: "modified_at", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := m.Collection.UpdateOne(context.Background(), bson.D{{Key: "device_id", Value: req.DeviceID}}, update)
	if err != nil {
		global.Logger.Error("刷新属性值失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	global.Logger.Info("刷新属性值成功", zap.Any("updateResult", updateResult))
	return &response.RefreshMemberAttributesResp{
		Text:   attribute.Text,
		Values: attribute.Values,
	}, nil
}

func (m *memberService) StartOrNextRound(req *request.StartOrNextRoundReq) (*response.StartOrNextRoundResp, *dun.StatusCode) {
	member, err := m.searchByDeviceID(req.DeviceID)
	if err != nil {
		global.Logger.Error("根据设备ID获取用户信息失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	if member == nil {
		return nil, dun.StatusCodeDataNotFound
	}

	if member.Film == nil {
		return nil, dun.StatusCodeDataNotFound
	}

	if req.SpecialNPCName != "" {
		member.SpecialNPCName = req.SpecialNPCName
	}

	roundEvent := req.GetEvent()
	if roundEvent == nil {
		return nil, constant.ErrEventID
	}
	// 记录回合数
	var currentRoundNum = 1
	if member.CurrentRound != nil {
		currentRoundNum = member.CurrentRound.Number
		currentRoundNum++
	} else {
		member.CurrentRound = &model.MemberCurrentRound{}
	}

	// 当前轮的数字
	member.CurrentRound.Number = currentRoundNum
	// 当前回合事件
	member.CurrentRound.Event = roundEvent

	// 根据回合数处理主题
	var (
		eventTheme     = constant.EventTopic.Random()
		welcomeMessage = ""
	)

	//First day of the movie shooting（第一回合必出）
	//Last day of the movie shooting（最后一回合必出）
	switch currentRoundNum {
	case 1:
		eventTheme = "First day of the movie shooting"
		welcomeMessage = member.WelcomeMessage()
	case constant.GlobalCaptureRoundNum:
		eventTheme = "Last day of the movie shooting"
	}

	var (
		location = constant.Location.Random()

		//actionTag  = constant.ActionTags.Random()
		weather = constant.Weathers.Random()

		promptChoicesID = roundEvent.Prompt.ChoicesID
		variables       = member.VariablesChoices(location, weather, eventTheme)

		difficultyValue = member.CurrentRound.SetDifficultyValue()
	)

	// 请求gpt
	apiResp, err := util.APIInternalPrompt(promptChoicesID, variables)
	if err != nil {
		global.Logger.Error("Choices 请求内部调用失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	replyContent := apiResp.Content
	global.Logger.Info("Choices 请求内部调用", zap.Any("promptChoicesID", promptChoicesID), zap.Any("variables", variables), zap.Any("replyContent", replyContent), zap.Any("事件类型", member.CurrentRound.Event.Description))
	var gptOptions = &constant.GptOptions{}
	if err = json.Unmarshal([]byte(replyContent), gptOptions); err != nil {
		global.Logger.Error("解析replyContent gptOptions失败", zap.Any("replyContent", replyContent), zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	// 将当轮的选项结果存起来
	gptOptionsV1 := gptOptions.V1()
	member.CurrentRound.GptOptions = gptOptionsV1
	// 删除相遇的对方
	member.Film.DeleteMeetingActors(member.CurrentRound.Actors)

	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "current_round", Value: member.CurrentRound})
		fields = append(fields, bson.E{Key: "film", Value: member.Film})
		fields = append(fields, bson.E{Key: "modified_at", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := m.Collection.UpdateOne(context.Background(), bson.D{{Key: "device_id", Value: req.DeviceID}}, update)
	if err != nil {
		global.Logger.Error("更新回合数据失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	global.Logger.Info("更新回合数据成功", zap.Any("updateResult", updateResult))
	return &response.StartOrNextRoundResp{
		EventTheme:      eventTheme,
		Actors:          member.CurrentRound.Actors,
		DifficultyValue: difficultyValue,
		RoundNumber:     currentRoundNum,
		WelcomeMessage:  welcomeMessage,
		Location:        location,
		Weather:         weather,
		AttributeValues: member.CurrentRound.AttributeValues,
		GptOptions:      gptOptions.V1(),
	}, nil
}

func (m *memberService) ClickButtonOutcome(req *request.ClickButtonOutcomeReq) (*response.ClickButtonOutcomeResp, *dun.StatusCode) {
	member, err := m.searchByDeviceID(req.DeviceID)
	if err != nil {
		global.Logger.Error("根据设备ID获取用户信息失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	if member == nil {
		return nil, dun.StatusCodeDataNotFound
	}

	playerOption := constant.PlayerCaptureOptions.GetByID(req.OptionID)
	if playerOption == nil {
		global.Logger.Error("ClickButtonOutcome 没有找到选项", zap.Any("req", req))
		return nil, constant.ErrOptionID
	}

	if currentRound := member.CurrentRound; currentRound == nil ||
		currentRound.Event == nil {
		return nil, constant.ErrNeedStartGame
	}

	var (
		promptOutcomeID = member.CurrentRound.Event.Prompt.OutcomeID
		variables       = member.VariablesOutcome(req.OptionID)
	)
	// 请求gpt
	apiResp, err := util.APIInternalPrompt(promptOutcomeID, variables)
	if err != nil {
		global.Logger.Error("ClickButtonOutcome 请求内部调用失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	replyContent := apiResp.Content
	// 解析json
	var gptOutcome = &constant.GptOutcome{}
	if err = json.Unmarshal([]byte(replyContent), gptOutcome); err != nil {
		global.Logger.Error("解析replyContent GptOutcome失败", zap.Any("replyContent", replyContent), zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	// 将gpt返回的结果存起来出新闻用
	var gptOutcomeV1 = gptOutcome.V1()
	member.CurrentRound.GptOutcome = gptOutcomeV1

	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "current_round", Value: member.CurrentRound})
		fields = append(fields, bson.E{Key: "modified_at", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := m.Collection.UpdateOne(context.Background(), bson.D{{Key: "device_id", Value: req.DeviceID}}, update)
	if err != nil {
		global.Logger.Error("更新回合数据失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	global.Logger.Info("更新回合数据成功", zap.Any("updateResult", updateResult))

	return &response.ClickButtonOutcomeResp{
		OptionOutcome: member.CurrentRound.OptionOutcome,
		GptOutcome:    gptOutcomeV1,
	}, nil
}

func (m *memberService) ClickButtonNewsByOutcome(req *request.ClickButtonNewsByOutcomeReq) (*response.ClickButtonNewsByOutcomeResp, *dun.StatusCode) {
	member, err := m.searchByDeviceID(req.DeviceID)
	if err != nil {
		global.Logger.Error("根据设备ID获取用户信息失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	if member == nil {
		return nil, dun.StatusCodeDataNotFound
	}
	if currentRound := member.CurrentRound; currentRound == nil ||
		currentRound.GptOutcome == nil {
		return nil, constant.ErrNeedOption
	}

	variables := member.VariablesOutcomeNews(member.CurrentRound.GptOutcome)

	// 请求gpt
	apiResp, err := util.APIInternalPrompt(constant.OutcomeNewsPromptID, variables)
	if err != nil {
		global.Logger.Error("ClickButtonNewsByOutcome News 请求内部调用失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	replyContent := apiResp.Content
	// 解析json
	var gptNews = &constant.GptNews{}
	if err = json.Unmarshal([]byte(replyContent), gptNews); err != nil {
		global.Logger.Error("解析replyContent GptNews 失败", zap.Any("replyContent", replyContent), zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	gptNewsV1 := gptNews.V1()
	// 将gptNew存起来
	member.CurrentRound.GptNews = gptNewsV1

	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "current_round", Value: member.CurrentRound})
		fields = append(fields, bson.E{Key: "modified_at", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := m.Collection.UpdateOne(context.Background(), bson.D{{Key: "device_id", Value: req.DeviceID}}, update)
	if err != nil {
		global.Logger.Error("更新回合数据失败", zap.Error(err))
		return nil, dun.StatusCodeSystemError
	}
	global.Logger.Info("更新回合数据成功", zap.Any("updateResult", updateResult))
	return &response.ClickButtonNewsByOutcomeResp{
		GptNews: gptNewsV1,
	}, nil
}

func (m *memberService) UpdateScreenplayReq(req *request.UpdateMemberScreenplayReq) *dun.StatusCode {
	member, err := m.searchByDeviceID(req.DeviceID)
	if err != nil {
		global.Logger.Error("根据设备ID获取用户信息失败", zap.Error(err))
		return dun.StatusCodeSystemError
	}
	if member == nil {
		return dun.StatusCodeDataNotFound
	}
	screenplay, err := m.ScreenplayService.searchByKey(req.Key)
	if err != nil {
		global.Logger.Error("获取剧本信息失败", zap.Error(err))
		return dun.StatusCodeSystemError
	}
	if screenplay == nil {
		return dun.StatusCodeDataNotFound
	}

	film := &model.MemberFilm{
		Screenplay: screenplay,
		Actors:     req.Actors,
		ActorMap:   nil,
	}
	film.Init()
	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "film", Value: film})
		fields = append(fields, bson.E{Key: "modified_at", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := m.Collection.UpdateOne(context.Background(), bson.D{{Key: "device_id", Value: req.DeviceID}}, update)
	if err != nil {
		global.Logger.Error("更新用户剧本数据失败", zap.Error(err))
		return dun.StatusCodeSystemError
	}
	global.Logger.Info("更新用户剧本数据成功", zap.Any("updateResult", updateResult))
	return dun.StatusCodeSuccess
}

func (m *memberService) searchByDeviceID(deviceID string) (*model.Member, error) {
	if deviceID == "" {
		return nil, nil
	}
	filter := bson.D{{Key: "device_id", Value: deviceID}}
	member, err := m.findOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return member, nil
}

func (m *memberService) findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*model.Member, error) {
	// 反序列化数据
	member := new(model.Member)
	if err := m.Collection.FindOne(ctx, filter, opts...).Decode(member); err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		global.Logger.Error("从mongo查询用户数据失败", zap.Error(err), zap.Any("filter", filter))
		return nil, err
	}

	return member, nil
}
