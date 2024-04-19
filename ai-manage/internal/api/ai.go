package api

import (
	"github.com/gin-gonic/gin"
	"github.com/qm012/dun"
	"io"
	"net/http"
	"sync"
	"vland.live/app/internal/constant"
	"vland.live/app/internal/model/request"
	"vland.live/app/internal/service"
)

type AICollector interface {
	V1ChatCompletions(ctx *gin.Context)
}

type aiCollector struct {
	service.AIService
}

var (
	aiCollectorOnce sync.Once
	aic             AICollector
)

func NewAICollector() AICollector {
	aiCollectorOnce.Do(func() {
		aic = &aiCollector{
			AIService: service.NewAIService(),
		}
	})
	return aic
}

func (a *aiCollector) V1ChatCompletions(ctx *gin.Context) {
	var req = &request.ChatCompletionReq{}
	if err := ctx.ShouldBind(req); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if err := req.Verify(); err != nil {
		dun.Failed400(ctx, dun.NewStatusCode(http.StatusBadRequest, err.Error()))
		return
	}

	if req.Stream { // 流式返回
		// SSE header
		ctx.Header("Content-Type", "text/event-stream")
		ctx.Header("Cache-Control", "no-cache")
		ctx.Header("Connection", "keep-alive")
		ctx.Header("Transfer-Encoding", "chunked")

		req.Recv = make(chan any, 1024)
		go a.AIService.ChatStream(req)

		//var err error
		ctx.Stream(func(w io.Writer) bool {
			for data := range req.Recv {
				switch content := data.(type) {
				case string:
					ctx.SSEvent(constant.SSEventChatName, content)
					return true
				case error:
					ctx.SSEvent(constant.SSEventError, content)
					return false
				}
			}
			return false
		})
		//if err != nil {
		//	dun.Failed200(ctx, dun.NewStatusCode(http.StatusInternalServerError, err.Error()))
		//	return
		//}
		return
	}

	// 普通返回
	resp, statusCode := a.AIService.Chat(req)
	if statusCode != nil {
		dun.Failed200(ctx, statusCode)
		return
	}
	dun.Success(ctx, resp)
}
