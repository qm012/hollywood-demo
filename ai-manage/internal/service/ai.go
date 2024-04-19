package service

import (
	"context"
	"errors"
	"github.com/qm012/dun"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"io"
	"net/http"
	"sync"
	"vland.live/app/global"
	"vland.live/app/internal/model/request"
	"vland.live/app/internal/model/response"
)

type AIService interface {
	ChatStream(req *request.ChatCompletionReq)
	Chat(req *request.ChatCompletionReq) (*response.ChatResp, *dun.StatusCode)
}

type aiService struct {
}

var (
	aiServiceOnce sync.Once
	as            *aiService
)

func NewAIService() AIService {
	aiServiceOnce.Do(func() {
		as = &aiService{}
	})
	return as
}

func (a *aiService) ChatStream(req *request.ChatCompletionReq) {
	var (
		stream   *openai.ChatCompletionStream
		aiClient *openai.Client = global.AIAzureClient
	)
	//switch req.AIPlatform {
	//case constant.AIPlatformOpenai:
	//	aiClient = global.AIOpenaiClient
	//case constant.AIPlatformAzure:
	//	aiClient = global.AIAzureClient
	//}

	stream, err := aiClient.CreateChatCompletionStream(context.Background(), req.ChatReq())
	if err != nil {
		global.Logger.Error("(a *aiService) ChatStream 失败")
		req.Recv <- err
		return
	}

	defer func() {
		stream.Close()
		close(req.Recv)
	}()

	for {
		var resp openai.ChatCompletionStreamResponse
		resp, err = stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				return
			}

			global.Logger.Error("(a *aiService) ChatStream 失败", zap.Any("req", req), zap.Error(err))
			req.Recv <- err
			return
		}
		content := resp.Choices[0].Delta.Content
		req.Recv <- content
	}
}

func (a *aiService) Chat(req *request.ChatCompletionReq) (*response.ChatResp, *dun.StatusCode) {
	var (
		aiClient *openai.Client = global.AIAzureClient
	)
	//switch req.AIPlatform {
	//case constant.AIPlatformOpenai:
	//	aiClient = global.AIOpenaiClient
	//case constant.AIPlatformAzure:
	//	aiClient = global.AIAzureClient
	//}
	resp, err := aiClient.CreateChatCompletion(context.Background(), req.ChatReq())
	if err != nil {
		global.Logger.Error("(a *aiService) Chat 失败", zap.Any("req", req), zap.Error(err))
		return nil, dun.NewStatusCode(http.StatusInternalServerError, err.Error())
	}
	replyContent := resp.Choices[0].Message.Content
	chatResp := &response.ChatResp{
		Content: replyContent,
		Usage:   resp.Usage,
	}
	return chatResp, nil
}
