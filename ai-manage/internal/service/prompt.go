package service

import (
	"context"
	"github.com/qm012/dun"
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"
	"vland.live/app/global"
	"vland.live/app/internal/constant"
	"vland.live/app/internal/model"
	"vland.live/app/internal/model/request"
	"vland.live/app/internal/model/response"
)

type PromptService interface {
	CreateAdmin(req *request.CreateAdminPromptReq) *dun.StatusCode
	UpdateAdmin(req *request.UpdateAdminPromptReq) *dun.StatusCode
	UpdateAdminLocked(req *request.UpdateAdminPromptLockedReq) *dun.StatusCode
	SearchAdmin(req *request.SearchAdminPromptReq) (*dun.PageInfo, *dun.StatusCode)
	SearchAdminDetail(req *request.SearchAdminPromptDetailReq) (*response.SearchAdminPromptDetailResp, *dun.StatusCode)
	DeleteAdmin(req *request.DeleteAdminPromptReq) *dun.StatusCode

	SaveAdminVersion(req *request.SaveAdminVersionReq) (*response.SaveAdminPromptVersionResp, *dun.StatusCode)
	CreateAdminPromptVersion(req *request.CreateAdminPromptVersionReq) (*response.CreateAdminPromptVersionResp, *dun.StatusCode)
	UpdateAdminPromptVersionIsProduction(req *request.UpdateAdminPromptVersionIsProductionReq) *dun.StatusCode
	UpdateAdminPromptVersionName(req *request.UpdateAdminPromptVersionNameReq) *dun.StatusCode
	DeleteAdminPromptVersion(req *request.DeleteAdminPromptVersionReq) *dun.StatusCode

	Chat(req *request.ChatPromptReq) (*response.ChatPromptResp, *dun.StatusCode)

	findOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) (*model.Prompt, error)
	find(ctx context.Context, filter any, opts ...*options.FindOptions) (prompts model.Prompts, err error)
}

type promptService struct {
	*mongo.Collection
	ProjectService
}

var (
	promptServiceOnce sync.Once
	promptservice     *promptService
)

func NewPromptService() PromptService {
	promptServiceOnce.Do(func() {
		promptservice = &promptService{
			Collection:     global.Mongo.Database(constant.DatabaseAI).Collection("prompt"),
			ProjectService: NewProjectService(),
		}
	})
	return promptservice
}

func (p *promptService) CreateAdmin(req *request.CreateAdminPromptReq) *dun.StatusCode {
	//// 处理prompt名称重复
	//project, err := p.findOne(context.Background(), bson.D{{Key: "name", Value: req.Name}})
	//if err != nil {
	//	return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	//}
	//if project != nil {
	//	return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同prompt名称，请更换一个prompt名称"))
	//}

	var (
		now       = time.Now().UnixMilli()
		projectID = req.GetProjectID()
	)

	document := &model.Prompt{
		ID:         primitive.NewObjectID(),
		ProjectID:  projectID,
		Name:       req.Name,
		Versions:   model.PromptVersions{},
		Locked:     false,
		ModifiedAt: now,
		CreatedAt:  now,
	}

	insertOneResult, err := p.Collection.InsertOne(context.Background(), document)
	if err != nil {
		global.Logger.Error("保存prompt失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if err = p.ProjectService.inc(projectID, constant.PromptIncByDefault); err != nil {
		global.Logger.Error("增加项目的prompt数量失败", zap.Any("req", req), zap.Error(err))
	}
	global.Logger.Info("保存prompt成功", zap.Any("insertOneResult", insertOneResult))
	return nil
}

func (p *promptService) UpdateAdmin(req *request.UpdateAdminPromptReq) *dun.StatusCode {
	var (
		id        = req.IdParamOmit.GetId()
		projectID = req.GetProjectID()
	)
	// 验证数据是否存在
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return dun.StatusCodeDataNotFound
	}

	//// 处理prompt名称重复
	//prompt, err := p.findOne(context.Background(), bson.D{{Key: "name", Value: req.Name}})
	//if err != nil {
	//	return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	//}
	//if prompt != nil && oldPrompt.ID != prompt.ID {
	//	return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同prompt名称，请更换一个prompt名称"))
	//}

	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "projectID", Value: projectID})
		fields = append(fields, bson.E{Key: "name", Value: req.Name})
		fields = append(fields, bson.E{Key: "modifiedAt", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := p.Collection.UpdateByID(context.Background(), id, update)
	if err != nil {
		global.Logger.Error("更新prompt失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("更新prompt成功", zap.Any("updateResult", updateResult))
	if oldPrompt.ProjectID != projectID {
		// 将被替换的project减少1个使用数量
		if err = p.ProjectService.inc(oldPrompt.ProjectID, constant.PromptDecrByDefault); err != nil {
			global.Logger.Error("减少项目的prompt数量失败", zap.Any("req", req), zap.Any("oldPrompt.ProjectID", oldPrompt.ProjectID.Hex()), zap.Error(err))
		}
		// 将新的project增加一个使用数量
		if err = p.ProjectService.inc(projectID, constant.PromptIncByDefault); err != nil {
			global.Logger.Error("增加项目的prompt数量失败", zap.Any("req", req), zap.Any("projectID", projectID), zap.Error(err))
		}
	}
	return nil
}

func (p *promptService) UpdateAdminLocked(req *request.UpdateAdminPromptLockedReq) *dun.StatusCode {
	var (
		id = req.IdParamOmit.GetId()
	)
	// 验证数据是否存在
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return dun.StatusCodeDataNotFound
	}
	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "locked", Value: !oldPrompt.Locked})
		fields = append(fields, bson.E{Key: "modifiedAt", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := p.Collection.UpdateByID(context.Background(), id, update)
	if err != nil {
		global.Logger.Error("更新prompt locked状态失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("更新prompt locked状态成功", zap.Any("updateResult", updateResult))
	return nil
}

func (p *promptService) SearchAdmin(req *request.SearchAdminPromptReq) (*dun.PageInfo, *dun.StatusCode) {
	filter := req.Filter()
	// 查询数量
	count, err := p.Collection.CountDocuments(context.Background(), filter)
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if count == 0 {
		return dun.NewPageInfo(count, []struct{}{}).SetPageSize(req.PageNum, req.PageSize), nil
	}

	var mongoSort bson.D
	{
		sort := req.Sort.Mongo()
		if req.SortField != "" {
			mongoSort = append(mongoSort, bson.E{Key: req.SortField, Value: sort})
		}
		mongoSort = append(mongoSort, bson.E{Key: "_id", Value: sort})
	}
	opt := options.Find().
		SetLimit(int64(req.PageSize)).
		SetSkip(int64(req.Offset(count))).
		SetSort(mongoSort)

	prompts, err := p.find(context.Background(), filter, opt)
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	var (
		resps      = make([]*response.SearchAdminPromptResp, 0, req.PageSize)
		projectIDs = prompts.GetProjectIDs()
	)

	// 查找项目数据
	projectMap, err := p.ProjectService.searchByIDs(projectIDs)
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	for _, prompt := range prompts {

		var (
			projectName string
			versionName string
		)
		project, ok := projectMap[prompt.ProjectID]
		if ok {
			projectName = project.Name
		}

		for _, promptVersion := range prompt.Versions {
			if promptVersion.IsProduction {
				versionName = promptVersion.Name
				break
			}
		}

		resps = append(resps, &response.SearchAdminPromptResp{
			ID:          prompt.ID.Hex(),
			Name:        prompt.Name,
			ProjectName: projectName,
			VersionName: versionName,
			Locked:      prompt.Locked,
			ModifiedAt:  prompt.ModifiedAt,
			CreatedAt:   prompt.CreatedAt,
		})
	}

	return dun.NewPageInfo(count, resps).SetPageSize(req.PageNum, req.PageSize), nil
}

func (p *promptService) SearchAdminDetail(req *request.SearchAdminPromptDetailReq) (*response.SearchAdminPromptDetailResp, *dun.StatusCode) {
	// 验证数据是否存在
	id := req.IdParamOmit.GetId()
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return nil, dun.StatusCodeDataNotFound
	}
	// 将版本数据倒序排列
	sort.SliceStable(oldPrompt.Versions, func(i, j int) bool {
		return oldPrompt.Versions[i].ID.Timestamp().After(oldPrompt.Versions[j].ID.Timestamp())
	})

	// 处理版本数据
	var versions = make([]*response.SearchAdminPromptDetailVersion, 0, len(oldPrompt.Versions))
	for _, promptVersion := range oldPrompt.Versions {

		// 处理请求openai的数据和变量数据
		var (
			chatReq   *response.SearchAdminPromptDetailVersionChatReq
			variables = make([]*response.SearchAdminPromptDetailVersionVariable, 0, len(promptVersion.Variables))
		)

		if promptVersionChatReq := promptVersion.ChatReq; promptVersionChatReq != nil {
			var responseFormat openai.ChatCompletionResponseFormatType
			if promptVersionChatReq.ResponseFormat != nil {
				responseFormat = promptVersionChatReq.ResponseFormat.Type
			}

			var messages = make([]response.ChatCompletionMessage, 0, len(promptVersionChatReq.Messages))

			for _, message := range promptVersionChatReq.Messages {
				messages = append(messages, response.ChatCompletionMessage{
					Role:    message.Role,
					Content: message.Content,
					Name:    message.Name,
				})
			}

			chatReq = &response.SearchAdminPromptDetailVersionChatReq{
				Model:            promptVersionChatReq.Model,
				Messages:         messages,
				MaxTokens:        promptVersionChatReq.MaxTokens,
				Temperature:      promptVersionChatReq.Temperature,
				TopP:             promptVersionChatReq.TopP,
				FrequencyPenalty: promptVersionChatReq.FrequencyPenalty,
				ResponseFormat:   responseFormat,
			}
		}

		for _, promptVersionVariable := range promptVersion.Variables {
			variables = append(variables, &response.SearchAdminPromptDetailVersionVariable{
				Key:   promptVersionVariable.Key,
				Value: promptVersionVariable.Value,
			})
		}

		versions = append(versions, &response.SearchAdminPromptDetailVersion{
			ID:           promptVersion.ID.Hex(),
			Name:         promptVersion.Name,
			IsProduction: promptVersion.IsProduction,
			ChatReq:      chatReq,
			Variables:    variables,
			Modifier:     promptVersion.Modifier,
			Creator:      promptVersion.Creator,
			ModifiedAt:   promptVersion.ModifiedAt,
			CreatedAt:    promptVersion.CreatedAt,
		})
	}

	// 查找项目数据
	projectMap, err := p.ProjectService.searchByIDs([]primitive.ObjectID{oldPrompt.ProjectID})
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	var (
		projectName string
	)
	project, ok := projectMap[oldPrompt.ProjectID]
	if ok {
		projectName = project.Name
	}
	resp := &response.SearchAdminPromptDetailResp{
		ID:          oldPrompt.ID.Hex(),
		ProjectName: projectName,
		Name:        oldPrompt.Name,
		Versions:    versions,
		Locked:      oldPrompt.Locked,
		ModifiedAt:  oldPrompt.ModifiedAt,
		CreatedAt:   oldPrompt.CreatedAt,
	}

	return resp, nil
}

func (p *promptService) DeleteAdmin(req *request.DeleteAdminPromptReq) *dun.StatusCode {

	// 验证数据是否存在
	id := req.IdParamOmit.GetId()
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return dun.StatusCodeDataNotFound
	}

	if oldPrompt.Locked {
		return dun.NewStatusCode(http.StatusInternalServerError, "被锁住的prompt无法删除")
	}

	filter := make(bson.D, 0, 2)
	filter = append(filter, bson.E{Key: "_id", Value: id})
	_, err = p.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		global.Logger.Error("删除prompt失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if err = p.ProjectService.inc(oldPrompt.ProjectID, constant.PromptDecrByDefault); err != nil {
		global.Logger.Error("减少项目的prompt数量失败", zap.Any("req", req), zap.Error(err))
	}
	return nil
}

func (p *promptService) SaveAdminVersion(req *request.SaveAdminVersionReq) (*response.SaveAdminPromptVersionResp, *dun.StatusCode) {

	// 验证数据是否存在
	var (
		id        = req.IdParamOmit.GetId()
		now       = time.Now().UnixMilli()
		variables = make([]*model.PromptVersionVariable, 0, len(req.Variables))
		resp      = &response.SaveAdminPromptVersionResp{
			VersionID: req.VersionID,
		}
	)
	// 处理变量数据
	for _, variable := range req.Variables {
		variables = append(variables, &model.PromptVersionVariable{
			Key:   variable.Key,
			Value: variable.Value,
		})
	}

	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return nil, dun.StatusCodeDataNotFound
	}

	var responseFormat *openai.ChatCompletionResponseFormat
	if req.ResponseFormat != "" {
		responseFormat = &openai.ChatCompletionResponseFormat{
			Type: req.ResponseFormat,
		}
	}

	switch req.OperationType {
	case request.SaveVersionOperationTypeOverride:

		var (
			versionID = req.GetVersionID()
		)

		if !oldPrompt.Versions.Contains(versionID) {
			return nil, dun.StatusCodeDataNotFound
		}

		messages := make([]openai.ChatCompletionMessage, 0, len(req.Messages))
		for _, message := range req.Messages {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:         message.Role,
				Content:      message.Content,
				MultiContent: nil,
				Name:         message.Name,
				FunctionCall: nil,
				ToolCalls:    nil,
				ToolCallID:   "",
			})
		}
		chatReq := openai.ChatCompletionRequest{
			Model:            req.Model,
			Messages:         messages,
			MaxTokens:        req.MaxTokens,
			Temperature:      req.Temperature,
			TopP:             req.TopP,
			N:                0,
			Stream:           false,
			Stop:             nil,
			PresencePenalty:  0,
			ResponseFormat:   responseFormat,
			Seed:             nil,
			FrequencyPenalty: req.FrequencyPenalty,
			LogitBias:        nil,
			User:             "",
			Functions:        nil,
			FunctionCall:     nil,
			Tools:            nil,
			ToolChoice:       nil,
		}
		// 更新数据
		filter := bson.D{{Key: "_id", Value: id},
			{Key: "versions._id", Value: versionID},
		}
		update := bson.D{{"$set", bson.D{
			//{Key: "modifier", Value: req.Operator},
			{Key: "modifiedAt", Value: now},
			{Key: "versions.$.name", Value: req.Name},
			{Key: "versions.$.chatReq", Value: chatReq},
			{Key: "versions.$.variables", Value: variables},
			{Key: "versions.$.modifier", Value: req.Operator},
			{Key: "versions.$.modifiedAt", Value: now},
		}}}
		updateResult, err := p.Collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			global.Logger.Error("覆盖prompt version 失败", zap.Any("filter", filter), zap.Any("update", update), zap.Error(err))
			return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
		}
		global.Logger.Info("覆盖prompt version 成功", zap.Any("updateResult", updateResult))
	case request.SaveVersionOperationTypeAddition:

		var isProduction bool
		if len(oldPrompt.Versions) == 0 {
			isProduction = true
		}

		messages := make([]openai.ChatCompletionMessage, 0, len(req.Messages))
		for _, message := range req.Messages {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:         message.Role,
				Content:      message.Content,
				MultiContent: nil,
				Name:         message.Name,
				FunctionCall: nil,
				ToolCalls:    nil,
				ToolCallID:   "",
			})
		}
		versionID := primitive.NewObjectID()
		resp = &response.SaveAdminPromptVersionResp{
			VersionID: versionID.Hex(),
		}
		promptVersion := &model.PromptVersion{
			ID:           versionID,
			Name:         req.Name,
			IsProduction: isProduction,
			ChatReq: &openai.ChatCompletionRequest{
				Model:            req.Model,
				Messages:         messages,
				MaxTokens:        req.MaxTokens,
				Temperature:      req.Temperature,
				TopP:             req.TopP,
				N:                0,
				Stream:           false,
				Stop:             nil,
				PresencePenalty:  0,
				ResponseFormat:   responseFormat,
				Seed:             nil,
				FrequencyPenalty: req.FrequencyPenalty,
				LogitBias:        nil,
				User:             "",
				Functions:        nil,
				FunctionCall:     nil,
				Tools:            nil,
				ToolChoice:       nil,
			},
			Variables:  variables,
			Modifier:   req.Operator,
			Creator:    req.Operator,
			ModifiedAt: now,
			CreatedAt:  now,
		}

		var (
			filter = bson.D{{Key: "_id", Value: id}}
			fields = bson.D{bson.E{Key: "modifiedAt", Value: now}}
			update = bson.M{
				"$set":      fields,
				"$addToSet": bson.M{"versions": bson.M{"$each": model.PromptVersions{promptVersion}}},
			}
		)

		updateResult, err := p.Collection.UpdateOne(context.Background(), filter, update)
		if err != nil {
			global.Logger.Error("新增prompt version 失败", zap.Any("filter", filter), zap.Any("update", update), zap.Error(err))
			return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
		}
		global.Logger.Info("新增prompt version 成功", zap.Any("updateResult", updateResult))
	}

	return resp, nil
}

func (p *promptService) CreateAdminPromptVersion(req *request.CreateAdminPromptVersionReq) (*response.CreateAdminPromptVersionResp, *dun.StatusCode) {
	// 验证数据是否存在
	var (
		id        = req.IdParamOmit.GetId()
		versionID = primitive.NewObjectID()
		now       = time.Now().UnixMilli()
	)
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return nil, dun.StatusCodeDataNotFound
	}

	var isProduction bool
	if len(oldPrompt.Versions) == 0 {
		isProduction = true
	}

	promptVersion := &model.PromptVersion{
		ID:           versionID,
		Name:         req.Name,
		IsProduction: isProduction,
		ChatReq:      nil,
		Modifier:     req.Operator,
		Creator:      req.Operator,
		ModifiedAt:   now,
		CreatedAt:    now,
	}

	var (
		filter = bson.D{{Key: "_id", Value: id}}
		fields = bson.D{bson.E{Key: "modifiedAt", Value: now}}
		update = bson.M{
			"$set":      fields,
			"$addToSet": bson.M{"versions": bson.M{"$each": model.PromptVersions{promptVersion}}},
		}
	)

	updateResult, err := p.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		global.Logger.Error("创建prompt version 失败", zap.Any("filter", filter), zap.Any("update", update), zap.Error(err))
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("创建prompt version 成功", zap.Any("updateResult", updateResult))
	return &response.CreateAdminPromptVersionResp{
		VersionID: versionID.Hex(),
	}, nil
}

func (p *promptService) UpdateAdminPromptVersionIsProduction(req *request.UpdateAdminPromptVersionIsProductionReq) *dun.StatusCode {
	// 验证数据是否存在
	var (
		id        = req.IdParamOmit.GetId()
		versionID = req.GetVersionID()
		now       = time.Now().UnixMilli()
	)
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return dun.StatusCodeDataNotFound
	}

	oldPromptVersion := oldPrompt.Versions.GetByID(versionID)
	if oldPromptVersion == nil {
		return dun.StatusCodeDataNotFound
	}

	if oldPromptVersion.IsProduction {
		return dun.NewStatusCode(http.StatusInternalServerError, "已经是生产状态，无需设置")
	}

	// 更新数据
	filter := bson.D{{Key: "_id", Value: id}, {Key: "versions._id", Value: versionID}}
	update := bson.D{{"$set", bson.D{
		{Key: "modifiedAt", Value: now},
		{Key: "versions.$.isProduction", Value: true},
		{Key: "versions.$.modifier", Value: req.Operator},
		{Key: "versions.$.modifiedAt", Value: now},
	}}}
	_, err = p.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		global.Logger.Error("变更prompt 版本的生产状态失败", zap.Any("filter", filter), zap.Any("update", update), zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	// 将其他的数据变更为非生产版本
	updateFilter := bson.D{{Key: "_id", Value: id},
		{Key: "versions", Value: bson.M{"$elemMatch": bson.M{"_id": bson.M{"$ne": versionID}}}},
	}
	updateMany := bson.D{{"$set", bson.D{
		{Key: "modifiedAt", Value: now},
		{Key: "versions.$[item].isProduction", Value: false},
		{Key: "versions.$[item].modifier", Value: req.Operator},
		{Key: "versions.$[item].modifiedAt", Value: now},
	}}}
	opts := options.Update().SetArrayFilters(options.ArrayFilters{
		Filters: []interface{}{
			bson.M{"item._id": bson.M{"$ne": versionID}},
		},
	})
	updateResult, err := p.Collection.UpdateMany(context.Background(), updateFilter, updateMany, opts)
	if err != nil {
		global.Logger.Error("将其他的数据变更为非生产版本失败", zap.Any("filter", updateFilter), zap.Any("update", updateMany), zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("变更prompt 版本的生产状态成功", zap.Any("updateResult", updateResult))
	return nil
}

func (p *promptService) UpdateAdminPromptVersionName(req *request.UpdateAdminPromptVersionNameReq) *dun.StatusCode {
	// 验证数据是否存在
	var (
		id        = req.IdParamOmit.GetId()
		versionID = req.GetVersionID()
		now       = time.Now().UnixMilli()
	)
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return dun.StatusCodeDataNotFound
	}

	// 更新数据
	filter := bson.D{{Key: "_id", Value: id}, {Key: "versions._id", Value: versionID}}
	update := bson.D{{"$set", bson.D{
		{Key: "modifiedAt", Value: now},
		{Key: "versions.$.name", Value: req.Name},
		{Key: "versions.$.modifier", Value: req.Operator},
		{Key: "versions.$.modifiedAt", Value: now},
	}}}
	updateResult, err := p.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		global.Logger.Error("变更prompt 版本的名称失败", zap.Any("filter", filter), zap.Any("update", update), zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("变更prompt 版本的名称成功", zap.Any("updateResult", updateResult))
	return nil
}

func (p *promptService) DeleteAdminPromptVersion(req *request.DeleteAdminPromptVersionReq) *dun.StatusCode {

	// 验证数据是否存在
	var (
		id        = req.IdParamOmit.GetId()
		versionID = req.GetVersionID()
	)
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return dun.StatusCodeDataNotFound
	}

	// 处理数据
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.M{
		"$pull": bson.M{"versions": bson.M{"_id": versionID}},
	}

	updateResult, err := p.Collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		global.Logger.Error("删除prompt版本失败", zap.Any("filter", filter), zap.Any("update", update), zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("删除prompt版本成功", zap.Any("updateResult", updateResult))
	return nil
}

func (p *promptService) Chat(req *request.ChatPromptReq) (*response.ChatPromptResp, *dun.StatusCode) {
	var id = req.IdParamOmit.GetId()
	// 验证数据是否存在
	oldPrompt, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldPrompt == nil {
		return nil, dun.StatusCodeDataNotFound
	}
	var releaseVersion *model.PromptVersion
	for _, version := range oldPrompt.Versions {
		if version.IsProduction {
			releaseVersion = version
		}
	}
	if releaseVersion == nil {
		global.Logger.Error("releaseVersion == nil", zap.Any("req", req))
		return nil, dun.StatusCodeSystemError
	}
	var chatReq = releaseVersion.ChatReq
	// 替换变量
	chatMessages := make([]openai.ChatCompletionMessage, 0, len(chatReq.Messages))
	for _, message := range chatReq.Messages {
		for key, variable := range req.Variables {
			message.Content = strings.ReplaceAll(message.Content, key, variable)
		}
		chatMessages = append(chatMessages, message)
	}
	chatReq.Messages = chatMessages

	// 自定义历史
	{
		messages := make([]openai.ChatCompletionMessage, 0, len(req.AppendMessages))
		for _, message := range req.AppendMessages {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:         message.Role,
				Content:      message.Content,
				MultiContent: nil,
				Name:         message.Name,
				FunctionCall: nil,
				ToolCalls:    nil,
				ToolCallID:   "",
			})
		}

		chatReq.Messages = append(chatReq.Messages, messages...)
	}

	global.Logger.Info("(p *promptService) Chat(req *request.ChatPromptReq) 得到releaseVersion", zap.Any("releaseVersion", releaseVersion), zap.Any("chatReq", chatReq))

	resp, err := global.AIAzureClient.CreateChatCompletion(context.Background(), *chatReq)
	if err != nil {
		global.Logger.Error("(p *promptService) Chat(req *request.ChatPromptReq)  失败", zap.Any("req", req), zap.Error(err))
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	replyContent := resp.Choices[0].Message.Content
	chatResp := &response.ChatPromptResp{
		Content: replyContent,
		Usage:   resp.Usage,
	}
	return chatResp, nil
}

func (p *promptService) findOne(ctx context.Context, filter any, opts ...*options.FindOneOptions) (*model.Prompt, error) {
	// 反序列化数据
	prompt := new(model.Prompt)
	if err := p.Collection.FindOne(ctx, filter, opts...).Decode(prompt); err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		global.Logger.Error("从mongo查询prompt数据失败", zap.Error(err), zap.Any("filter", filter))
		return nil, err
	}

	return prompt, nil
}

func (p *promptService) find(ctx context.Context, filter any, opts ...*options.FindOptions) (prompts model.Prompts, err error) {
	defer func() {
		if err != nil {
			global.Logger.Error("从mongo查询prompt列表数据失败", zap.Any("filter", filter), zap.Error(err))
			return
		}
		global.Logger.Info("从mongo查询prompt列表数据成功", zap.Any("filter", filter))
	}()
	cursor, err := p.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	prompts = make(model.Prompts, 0, 30)
	if err = cursor.All(context.Background(), &prompts); err != nil {
		return nil, err
	}

	if err = cursor.Close(context.Background()); err != nil {
		return nil, err
	}

	return prompts, nil
}
