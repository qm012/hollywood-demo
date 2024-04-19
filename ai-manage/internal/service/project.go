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

type ProjectService interface {
	CreateAdmin(req *request.CreateAdminProjectReq) *dun.StatusCode
	UpdateAdmin(req *request.UpdateAdminProjectReq) *dun.StatusCode
	SearchAdmin(req *request.SearchAdminProjectReq) (*dun.PageInfo, *dun.StatusCode)
	DeleteAdmin(req *request.DeleteAdminProjectReq) *dun.StatusCode

	searchByIDs(ids []primitive.ObjectID) (map[primitive.ObjectID]*model.Project, error)
	inc(ID primitive.ObjectID, quantity int) error
	findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*model.Project, error)
	find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (projects model.Projects, err error)
}

type projectService struct {
	*mongo.Collection
}

var (
	projectServiceOnce sync.Once
	ps                 *projectService
)

func NewProjectService() ProjectService {
	projectServiceOnce.Do(func() {
		ps = &projectService{
			Collection: global.Mongo.Database(constant.DatabaseAI).Collection("project"),
		}
	})
	return ps
}

func (p *projectService) CreateAdmin(req *request.CreateAdminProjectReq) *dun.StatusCode {
	// 处理项目名称重复
	project, err := p.findOne(context.Background(), bson.D{{Key: "name", Value: req.Name}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if project != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同项目名称，请更换一个项目名称"))
	}

	now := time.Now().UnixMilli()
	document := &model.Project{
		ID:             primitive.NewObjectID(),
		Name:           req.Name,
		PromptQuantity: 0,
		ModifiedAt:     now,
		CreatedAt:      now,
	}

	insertOneResult, err := p.Collection.InsertOne(context.Background(), document)
	if err != nil {
		global.Logger.Error("保存项目失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	global.Logger.Info("保存项目成功", zap.Any("insertOneResult", insertOneResult))
	return nil
}

func (p *projectService) UpdateAdmin(req *request.UpdateAdminProjectReq) *dun.StatusCode {
	// 验证数据是否存在
	id := req.IdParamOmit.GetId()
	oldProject, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldProject == nil {
		return dun.StatusCodeDataNotFound
	}

	// 处理项目名称重复
	project, err := p.findOne(context.Background(), bson.D{{Key: "name", Value: req.Name}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if project != nil && project.ID != oldProject.ID {
		return dun.NewStatusCode(http.StatusInternalServerError, fmt.Sprintf("已经存在相同项目名称，请更换一个项目名称"))
	}

	var update bson.D
	{
		fields := make(bson.D, 0, 5)
		fields = append(fields, bson.E{Key: "name", Value: req.Name})
		fields = append(fields, bson.E{Key: "modifiedAt", Value: time.Now().UnixMilli()})
		update = bson.D{{"$set", fields}}
	}

	updateResult, err := p.Collection.UpdateByID(context.Background(), id, update)
	if err != nil {
		global.Logger.Error("更新项目失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	global.Logger.Info("更新项目成功", zap.Any("updateResult", updateResult))
	return nil
}

func (p *projectService) SearchAdmin(req *request.SearchAdminProjectReq) (*dun.PageInfo, *dun.StatusCode) {
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

	projects, err := p.find(context.Background(), filter, opt)
	if err != nil {
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}

	resps := make([]*response.SearchAdminProjectResp, 0, req.PageSize)
	for _, project := range projects {
		resps = append(resps, &response.SearchAdminProjectResp{
			ID:             project.ID.Hex(),
			Name:           project.Name,
			PromptQuantity: project.PromptQuantity,
			ModifiedAt:     project.ModifiedAt,
			CreatedAt:      project.CreatedAt,
		})
	}

	return dun.NewPageInfo(count, resps).SetPageSize(req.PageNum, req.PageSize), nil
}

func (p *projectService) DeleteAdmin(req *request.DeleteAdminProjectReq) *dun.StatusCode {
	// 验证数据是否存在
	id := req.IdParamOmit.GetId()
	oldProject, err := p.findOne(context.Background(), bson.D{{Key: "_id", Value: id}})
	if err != nil {
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	if oldProject == nil {
		return dun.StatusCodeDataNotFound
	}

	if oldProject.PromptQuantity > 0 {
		return dun.NewStatusCode(http.StatusInternalServerError, "该项目下还有prompt，无法删除")
	}

	filter := make(bson.D, 0, 2)
	filter = append(filter, bson.E{Key: "_id", Value: id})
	_, err = p.Collection.DeleteOne(context.Background(), filter)
	if err != nil {
		global.Logger.Error("删除项目失败", zap.Error(err))
		return dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	return nil
}

func (p *projectService) searchByIDs(ids []primitive.ObjectID) (map[primitive.ObjectID]*model.Project, error) {
	if len(ids) == 0 {
		return map[primitive.ObjectID]*model.Project{}, nil
	}
	filter := bson.D{{Key: "_id", Value: bson.M{"$in": ids}}}
	projects, err := p.find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var resultMap = make(map[primitive.ObjectID]*model.Project, len(ids))
	for _, project := range projects {
		resultMap[project.ID] = project
	}
	return resultMap, nil
}

func (p *projectService) inc(ID primitive.ObjectID, quantity int) error {
	update := bson.M{
		"$inc": bson.M{
			"promptQuantity": quantity,
		},
	}
	// 更新数量
	updateResult, err := p.Collection.UpdateByID(context.Background(), ID, update)
	if err != nil {
		global.Logger.Error(" (p *projectService) inc 更新使用的prompt数量失败", zap.Error(err), zap.Any("ID", ID), zap.Any("quantity", quantity))
		return err
	}
	global.Logger.Info("(p *projectService) inc(ID primitive.ObjectID, quantity int) 成功", zap.Any("ID", ID), zap.Any("quantity", quantity), zap.Any("updateResult", updateResult))
	return nil
}

func (p *projectService) findOne(ctx context.Context, filter interface{}, opts ...*options.FindOneOptions) (*model.Project, error) {
	// 反序列化数据
	project := new(model.Project)
	if err := p.Collection.FindOne(ctx, filter, opts...).Decode(project); err != nil {

		if err == mongo.ErrNoDocuments {
			return nil, nil
		}

		global.Logger.Error("从mongo查询项目数据失败", zap.Error(err), zap.Any("filter", filter))
		return nil, err
	}

	return project, nil
}

func (p *projectService) find(ctx context.Context, filter interface{}, opts ...*options.FindOptions) (projects model.Projects, err error) {
	defer func() {
		if err != nil {
			global.Logger.Error("从mongo查询项目列表数据失败", zap.Any("filter", filter), zap.Error(err))
			return
		}
		global.Logger.Info("从mongo查询项目列表数据成功", zap.Any("filter", filter))
	}()
	cursor, err := p.Collection.Find(ctx, filter, opts...)
	if err != nil {
		return nil, err
	}
	if err = cursor.Err(); err != nil {
		return nil, err
	}

	projects = make(model.Projects, 0, 30)
	if err = cursor.All(context.Background(), &projects); err != nil {
		return nil, err
	}

	if err = cursor.Close(context.Background()); err != nil {
		return nil, err
	}

	return projects, nil
}
