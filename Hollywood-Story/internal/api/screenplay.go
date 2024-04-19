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

type ScreenplayCollector interface {
	CreateAdmin(ctx *gin.Context)
	UpdateAdmin(ctx *gin.Context)
	DeleteAdmin(ctx *gin.Context)
	SearchAdmin(ctx *gin.Context)

	Search(ctx *gin.Context)
}

type screenplayCollector struct {
	service.ScreenplayService
}

var (
	screenplayCollectorOnce sync.Once
	spc                     ScreenplayCollector
)

func NewScreenplayCollector() ScreenplayCollector {
	screenplayCollectorOnce.Do(func() {
		spc = &screenplayCollector{
			ScreenplayService: service.NewScreenplayService(),
		}
	})
	return spc
}

func (s *screenplayCollector) CreateAdmin(ctx *gin.Context) {
	var req = &request.CreateAdminScreenplayReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := s.ScreenplayService.CreateAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (s *screenplayCollector) UpdateAdmin(ctx *gin.Context) {
	var req = &request.UpdateAdminScreenplayReq{}
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
	if statusCode := s.ScreenplayService.UpdateAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (s *screenplayCollector) DeleteAdmin(ctx *gin.Context) {
	var req = &request.DeleteAdminScreenplayReq{}
	if err := ctx.ShouldBindUri(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := s.ScreenplayService.DeleteAdmin(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (s *screenplayCollector) SearchAdmin(ctx *gin.Context) {
	var req = &request.SearchAdminScreenplayReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	info, statusCode := s.ScreenplayService.SearchAdmin(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, info)
}

func (s *screenplayCollector) Search(ctx *gin.Context) {
	var req = &request.SearchScreenplayReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	info, statusCode := s.ScreenplayService.Search(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, info)
}
