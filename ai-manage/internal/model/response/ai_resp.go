package response

import "github.com/sashabaranov/go-openai"

type ChatResp struct {
	Content string       `json:"content"`
	Usage   openai.Usage `json:"usage"`
}

type TranscriptionsResp struct {
	Text string `json:"text"`
}
