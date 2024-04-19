package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qm012/dun"
	"net/http"
	"sync"
	"vland.live/app/internal/model/request"
	"vland.live/app/internal/service"
)

type MemberCollector interface {
	Enter(ctx *gin.Context)
	RefreshAttributes(ctx *gin.Context)
	StartOrNextRound(ctx *gin.Context)
	ClickButtonOutcome(ctx *gin.Context)
	ClickButtonNewsByOutcome(ctx *gin.Context)
	UpdateScreenplayReq(ctx *gin.Context)
	Search(ctx *gin.Context)
}

type memberCollector struct {
	service.MemberService
}

var (
	memberCollectorOnce sync.Once
	mc                  MemberCollector
)

func NewMemberCollector() MemberCollector {
	memberCollectorOnce.Do(func() {
		mc = &memberCollector{
			MemberService: service.NewMemberService(),
		}
	})
	return mc
}

func (m *memberCollector) Enter(ctx *gin.Context) {
	var req = &request.InputMemberInfoReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := m.MemberService.InputInfo(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (m *memberCollector) UpdateScreenplayReq(ctx *gin.Context) {
	var req = &request.UpdateMemberScreenplayReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if statusCode := m.MemberService.UpdateScreenplayReq(req); statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx)
}

func (m *memberCollector) RefreshAttributes(ctx *gin.Context) {
	var req = &request.RefreshAttributesReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	resp, statusCode := m.MemberService.RefreshAttributes(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, resp)
}

func (m *memberCollector) StartOrNextRound(ctx *gin.Context) {
	var req = &request.StartOrNextRoundReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	resp, statusCode := m.MemberService.StartOrNextRound(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, resp)
}
func (m *memberCollector) ClickButtonOutcome(ctx *gin.Context) {
	var req = &request.ClickButtonOutcomeReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	resp, statusCode := m.MemberService.ClickButtonOutcome(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, resp)
}

func (m *memberCollector) ClickButtonNewsByOutcome(ctx *gin.Context) {
	var req = &request.ClickButtonNewsByOutcomeReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	resp, statusCode := m.MemberService.ClickButtonNewsByOutcome(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, resp)
}

func (m *memberCollector) Search(ctx *gin.Context) {
	var req = &request.GetMemberInfoReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	info, statusCode := m.MemberService.GetMemberInfo(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}

	dun.Success(ctx, info)
}
