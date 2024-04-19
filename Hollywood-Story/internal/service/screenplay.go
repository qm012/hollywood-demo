package service

import (
	"context"
	"fmt"
	"github.com/qm012/dun"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"net/http"
	"sync"
	"time"
	"vland.live/app/global"
	"vland.live/app/internal/constant"
	"vland.live/app/internal/model"
	"vland.live/app/internal/model/request"
	"vland.live/app/internal/model/response"
)

type ScreenplayService interface {
	CreateAdmin(req *request.CreateAdminScreenplayReq) *dun.StatusCode
	UpdateAdmin(req *request.UpdateAdminScreenplayReq) *dun.StatusCode
	SearchAdmin(req *request.SearchAdminScreenplayReq) (*dun.PageInfo, *dun.StatusCode)
	DeleteAdmin(req *request.DeleteAdminScreenplayReq) *dun.StatusCode

	Search(req *request.SearchScreenplayReq) (*dun.PageInfo, *dun.StatusCode)
	searchByID(id primitive.ObjectID) (*model.Screenplay, error)
	searchByKey(key string) (*model.Screenplay, error)
	findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*model.Screenplay, error)
	find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (Screenplays model.Screenplays, err error)
}

type screenplayService struct {
	*mongo.Collection
}

var (
	screenplayServiceOnce sync.Once
	sps                   *screenplayService
)

func NewScreenplayService() ScreenplayService {
	screenplayServiceOnce.Do(func() {
		sps = &screenplayService{
			Collection: global.Mongo.Database(constant.DatabaseHollywood).Collection("screenplay"),
		}
	})
	return sps
}

func (s *screenplayService) CreateAdmin(req *request.CreateAdminScreenplayReq) *dun.StatusCode {
	// 处理名称重复
	screenplay, err := s.findOne(context.Background(), bson.D{{Key: "name", Value: req.Name}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if screenplay != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同剧本名称，请更换一个剧本名称"))
	}

	// 处理key重复
	screenplay, err = s.findOne(context.Background(), bson.D{{Key: "key", Value: req.Key}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if screenplay != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同剧本key，请更换一个剧本key"))
	}

	now := time.Now().UnixMilli()
	document := &model.Screenplay{
		ID:             primitive.NewObjectID(),
		Key:            req.Key,
		Name:           req.Name,
		Labels:         req.Labels,
		Synopsis:       req.Synopsis,
		WelcomeMessage: req.WelcomeMessage,
		ModifiedAt:     now,
		CreatedAt:      now,
	}

	insertOneResult, err := s.Collection.InsertOne(context.Background(), document)
	if err != nil {
		global.Logger.Error("保存剧本失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	global.Logger.Info("保存剧本成功", zap.Any("insertOneResult", insertOneResult))
	return nil
}

func (s *screenplayService) UpdateAdmin(req *request.UpdateAdminScreenplayReq) *dun.StatusCode {
	// 验证数据是否存在
	id := req.IdParamOmit.GetId()
	oldScreenplay, err := s.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldScreenplay == nil {
		return dun.StatusCodeDataNotFound
	}

	// 处理剧本名称重复
	screenplay, err := s.findOne(context.Background(), bson.D{{Key: "name", Value: req.Name}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if screenplay != nil && screenplay.ID != oldScreenplay.ID {
		return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同剧本名称，请更换一个剧本名称"))
	}

	// 处理剧本key重复
	screenplay, err = s.findOne(context.Background(), bson.D{{Key: "key", Value: req.Key}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if screenplay != nil && screenplay.ID != oldScreenplay.ID {
		return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同剧本key，请更换一个剧本key"))
	}

	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "key", Value: req.Key})
		fields = append(fields, bson.E{Key: "name", Value: req.Name})
		fields = append(fields, bson.E{Key: "labels", Value: req.Labels})
		fields = append(fields, bson.E{Key: "synopsis", Value: req.Synopsis})
		fields = append(fields, bson.E{Key: "welcome_message", Value: req.WelcomeMessage})
		fields = append(fields, bson.E{Key: "modified_at", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := s.Collection.UpdateByID(context.Background(), id, update)
	if err != nil {
		global.Logger.Error("更新剧本失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("更新剧本成功", zap.Any("updateResult", updateResult))
	return nil
}

func (s *screenplayService) SearchAdmin(req *request.SearchAdminScreenplayReq) (*dun.PageInfo, *dun.StatusCode) {
	filter := req.Filter()
	// 查询数量
	count, err := s.Collection.CountDocuments(context.Background(), filter)
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

	screenplays, err := s.find(context.Background(), filter, opt)
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	resps := make([]*response.SearchAdminScreenplayResp, 0, req.PageSize)
	for _, screenplay := range screenplays {
		resps = append(resps, &response.SearchAdminScreenplayResp{
			ID:             screenplay.ID.Hex(),
			Key:            screenplay.Key,
			Name:           screenplay.Name,
			Labels:         screenplay.Labels,
			Synopsis:       screenplay.Synopsis,
			WelcomeMessage: screenplay.WelcomeMessage,
			ModifiedAt:     screenplay.ModifiedAt,
			CreatedAt:      screenplay.CreatedAt,
		})
	}

	return dun.NewPageInfo(count, resps).SetPageSize(req.PageNum, req.PageSize), nil
}

func (s *screenplayService) DeleteAdmin(req *request.DeleteAdminScreenplayReq) *dun.StatusCode {
	// 验证数据是否存在
	id := req.IdParamOmit.GetId()
	oldScreenplay, err := s.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldScreenplay == nil {
		return dun.StatusCodeDataNotFound
	}

	filter := make(bson.D, 0, 2)
	filter = append(filter, bson.E{Key: "_id", Value: id})
	_, err = s.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		global.Logger.Error("删除剧本失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (s *screenplayService) Search(req *request.SearchScreenplayReq) (*dun.PageInfo, *dun.StatusCode) {
	filter := req.Filter()
	// 查询数量
	count, err := s.Collection.CountDocuments(context.Background(), filter)
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

	screenplays, err := s.find(context.Background(), filter, opt)
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	resps := make([]*response.SearchScreenplayResp, 0, req.PageSize)
	for _, screenplay := range screenplays {
		resps = append(resps, &response.SearchScreenplayResp{
			Key:      screenplay.Key,
			Name:     screenplay.Name,
			Actors:   model.GetActorsByNum(req.ActorNum),
			Labels:   screenplay.Labels,
			Synopsis: screenplay.Synopsis,
		})
	}

	return dun.NewPageInfo(count, resps).SetPageSize(req.PageNum, req.PageSize), nil
}

func (s *screenplayService) searchByID(id primitive.ObjectID) (*model.Screenplay, error) {
	if id.IsZero() {
		return nil, nil
	}
	filter := bson.D{{Key: "_id", Value: id}}
	screenplay, err := s.findOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return screenplay, nil
}

func (s *screenplayService) searchByKey(key string) (*model.Screenplay, error) {
	if key == "" {
		return nil, nil
	}
	filter := bson.D{{Key: "key", Value: key}}
	screenplay, err := s.findOne(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	return screenplay, nil
}

func (s *screenplayService) findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*model.Screenplay, error) {
	// 反序列化数据
	screenplay := new(model.Screenplay)
	if err := s.Collection.FindOne(ctx, filter, opts...).Decode(screenplay); err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		global.Logger.Error("从mongo查询剧本数据失败", zap.Error(err), zap.Any("filter", filter))
		return nil, err
	}

	return screenplay, nil
}

func (s *screenplayService) find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (screenplays model.Screenplays, err error) {
	defer func() {
		if err != nil {
			global.Logger.Error("从mongo查询剧本列表数据失败", zap.Any("filter", filter), zap.Error(err))
			return
		}
		global.Logger.Info("从mongo查询剧本列表数据成功", zap.Any("filter", filter))
	}()
	cursor, err := s.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	screenplays = make(model.Screenplays, 0, 30)
	if err = cursor.All(context.Background(), &screenplays); err != nil {
		return nil, err
	}

	if err = cursor.Close(context.Background()); err != nil {
		return nil, err
	}

	return screenplays, nil
}
