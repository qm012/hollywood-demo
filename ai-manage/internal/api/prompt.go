package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qm012/dun"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"sync"
	"vland.live/app/internal/model/request"
	"vland.live/app/internal/service"
)

type PromptCollector interface {
	CreateAdmin(ctx *gin.Context)
	UpdateAdmin(ctx *gin.Context)
	UpdateAdminLocked(ctx *gin.Context)
	DeleteAdmin(ctx *gin.Context)
	SearchAdmin(ctx *gin.Context)
	SearchAdminDetail(ctx *gin.Context)
	SaveAdminVersion(ctx *gin.Context)
	CreateAdminPromptVersion(ctx *gin.Context)
	UpdateAdminPromptVersionIsProduction(ctx *gin.Context)
	UpdateAdminPromptVersionName(ctx *gin.Context)
	DeleteAdminPromptVersion(ctx *gin.Context)
	Chat(ctx *gin.Context)
	ChatExternal(ctx *gin.Context) // 外部调用
}

type promptCollector struct {
	service.PromptService
}

var (
	promptCollectorOnce sync.Once
	promptcollector     PromptCollector
)

func NewPromptCollector() PromptCollector {
	promptCollectorOnce.Do(func() {
		promptcollector = &promptCollector{
			PromptService: service.NewPromptService(),
		}
	})
	return promptcollector
}

func (p *promptCollector) CreateAdmin(ctx *gin.Context) {
	var req = &request.CreateAdminPromptReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := p.PromptService.CreateAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *promptCollector) UpdateAdmin(ctx *gin.Context) {
	var req = &request.UpdateAdminPromptReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	if statusCode := p.PromptService.UpdateAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *promptCollector) UpdateAdminLocked(ctx *gin.Context) {
	var req = &request.UpdateAdminPromptLockedReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	if statusCode := p.PromptService.UpdateAdminLocked(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *promptCollector) DeleteAdmin(ctx *gin.Context) {
	var req = &request.DeleteAdminPromptReq{}
	if err := ctx.ShouldBindUri(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := p.PromptService.DeleteAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *promptCollector) SearchAdmin(ctx *gin.Context) {
	var req = &request.SearchAdminPromptReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	info, statusCode := p.PromptService.SearchAdmin(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, info)
}

func (p *promptCollector) SearchAdminDetail(ctx *gin.Context) {
	var req = &request.SearchAdminPromptDetailReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	resp, statusCode := p.PromptService.SearchAdminDetail(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, resp)
}

func (p *promptCollector) SaveAdminVersion(ctx *gin.Context) {
	var req = &request.SaveAdminVersionReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if err := req.Verify(); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	resp, statusCode := p.PromptService.SaveAdminVersion(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, resp)
}

func (p *promptCollector) CreateAdminPromptVersion(ctx *gin.Context) {
	var req = &request.CreateAdminPromptVersionReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	resp, statusCode := p.PromptService.CreateAdminPromptVersion(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, resp)
}

func (p *promptCollector) UpdateAdminPromptVersionIsProduction(ctx *gin.Context) {
	var req = &request.UpdateAdminPromptVersionIsProductionReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}
	versionID := ctx.Param("version_id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	req.VersionID = versionID
	if statusCode := p.PromptService.UpdateAdminPromptVersionIsProduction(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *promptCollector) UpdateAdminPromptVersionName(ctx *gin.Context) {
	var req = &request.UpdateAdminPromptVersionNameReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}
	versionID := ctx.Param("version_id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	req.VersionID = versionID
	if statusCode := p.PromptService.UpdateAdminPromptVersionName(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *promptCollector) DeleteAdminPromptVersion(ctx *gin.Context) {
	var req = &request.DeleteAdminPromptVersionReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}
	versionID := ctx.Param("version_id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	req.VersionID = versionID
	if statusCode := p.PromptService.DeleteAdminPromptVersion(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *promptCollector) Chat(ctx *gin.Context) {
	var req = &request.ChatPromptReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	ID := ctx.Param("id")
	if _, err := primitive.ObjectIDFromHex(ID); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	req.IdParamOmit.Id = ID
	resp, statusCode := p.PromptService.Chat(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}
	dun.Success(ctx, resp)
}

func (p *promptCollector) ChatExternal(ctx *gin.Context) {
	apiKey := ctx.GetHeader("APIKEY")
	if apiKey != "1NVL791XZCEWGxY7co6zgP5EkB1jf2dv9m" {
		ctx.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	p.Chat(ctx)
}
