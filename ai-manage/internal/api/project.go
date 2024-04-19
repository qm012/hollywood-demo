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

type ProjectCollector interface {
	CreateAdmin(ctx *gin.Context)
	UpdateAdmin(ctx *gin.Context)
	DeleteAdmin(ctx *gin.Context)
	SearchAdmin(ctx *gin.Context)
}

type projectCollector struct {
	service.ProjectService
}

var (
	projectCollectorOnce sync.Once
	pc                   ProjectCollector
)

func NewProjectCollector() ProjectCollector {
	projectCollectorOnce.Do(func() {
		pc = &projectCollector{
			ProjectService: service.NewProjectService(),
		}
	})
	return pc
}

func (p *projectCollector) CreateAdmin(ctx *gin.Context) {
	var req = &request.CreateAdminProjectReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := p.ProjectService.CreateAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *projectCollector) UpdateAdmin(ctx *gin.Context) {
	var req = &request.UpdateAdminProjectReq{}
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
	if statusCode := p.ProjectService.UpdateAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *projectCollector) DeleteAdmin(ctx *gin.Context) {
	var req = &request.DeleteAdminProjectReq{}
	if err := ctx.ShouldBindUri(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := p.ProjectService.DeleteAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (p *projectCollector) SearchAdmin(ctx *gin.Context) {
	var req = &request.SearchAdminProjectReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	info, statusCode := p.ProjectService.SearchAdmin(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, info)
}
