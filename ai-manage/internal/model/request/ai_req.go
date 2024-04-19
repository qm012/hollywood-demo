package request

import (
	"errors"
	"github.com/sashabaranov/go-openai"
	"vland.live/app/internal/constant"
)

type ChatCompletionReq struct {
	AIPlatform       constant.AIPlatform                     `json:"ai_platform" binding:"required,oneof=openai azure"`
	Stream           bool                                    `json:"stream"`
	Model            string                                  `json:"model" binding:"required"`
	Messages         []ChatCompletionMessage                 `json:"messages" binding:"required,min=1,max=1000,dive"`
	MaxTokens        int                                     `json:"max_tokens" binding:"omitempty,max=4096"`
	Temperature      float32                                 `json:"temperature" binding:"omitempty"`
	TopP             float32                                 `json:"top_p" binding:"omitempty"`
	FrequencyPenalty float32                                 `json:"frequency_penalty" binding:"omitempty"`
	ResponseFormat   openai.ChatCompletionResponseFormatType `json:"response_format" binding:"required,oneof=json_object text"`
	Operator         string                                  `json:"operator" binding:"required,max=50"`

	Recv chan any `json:"-"`
}

func (c *ChatCompletionReq) Verify() error {
	if !completionModels[c.Model] {
		return errors.New("不支持的模型")
	}
	return nil
}

func (c *ChatCompletionReq) ChatReq() openai.ChatCompletionRequest {
	var (
		messages = make([]openai.ChatCompletionMessage, 0, len(c.Messages))
	)

	for _, message := range c.Messages {
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
		Model:           c.Model,
		Messages:        messages,
		MaxTokens:       c.MaxTokens,
		Temperature:     c.Temperature,
		TopP:            c.TopP,
		N:               0,
		Stream:          c.Stream,
		Stop:            nil,
		PresencePenalty: 0,
		ResponseFormat: &openai.ChatCompletionResponseFormat{
			Type: c.ResponseFormat,
		},
		Seed:             nil,
		FrequencyPenalty: c.FrequencyPenalty,
		LogitBias:        nil,
		User:             "",
		Functions:        nil,
		FunctionCall:     nil,
		Tools:            nil,
		ToolChoice:       nil,
	}
	return chatReq
}
