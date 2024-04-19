package response

import "github.com/sashabaranov/go-openai"

type SearchAdminPromptResp struct {
	ID          string `json:"id"`           // 主键
	Name        string `json:"name"`         // prompt 名称
	ProjectName string `json:"project_name"` // 项目名称
	VersionName string `json:"version_name"` // 版本名称
	Locked      bool   `json:"locked"`       // 是否锁住，true:已锁，false:未锁 上锁后的prompt，删除不可点击
	ModifiedAt  int64  `json:"modified_at"`  // 更新时间
	CreatedAt   int64  `json:"created_at"`   // 创建时间
}

type SaveAdminPromptVersionResp struct {
	VersionID string `json:"version_id"`
}

type CreateAdminPromptVersionResp struct {
	VersionID string `json:"version_id"`
}

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name"`
}

type SearchAdminPromptDetailVersionChatReq struct {
	Model            string                                  `json:"model" `
	Messages         []ChatCompletionMessage                 `json:"messages"`
	MaxTokens        int                                     `json:"max_tokens"`
	Temperature      float32                                 `json:"temperature"`
	TopP             float32                                 `json:"top_p"`
	FrequencyPenalty float32                                 `json:"frequency_penalty"`
	ResponseFormat   openai.ChatCompletionResponseFormatType `json:"response_format"`
}

type SearchAdminPromptDetailVersionVariable struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type SearchAdminPromptDetailVersion struct {
	ID           string                                    `json:"id"`            // 主键
	Name         string                                    `json:"name"`          // 版本名称
	IsProduction bool                                      `json:"is_production"` // 如果是生产版本，则为true；否则为false
	ChatReq      *SearchAdminPromptDetailVersionChatReq    `json:"chat_req"`      // 聊天请求参数
	Variables    []*SearchAdminPromptDetailVersionVariable `json:"variables"`     // 变量数据
	Modifier     string                                    `json:"modifier"`      // 更新人
	Creator      string                                    `json:"creator"`       // 创建人
	ModifiedAt   int64                                     `json:"modified_at"`   // 更新时间
	CreatedAt    int64                                     `json:"created_at"`    // 创建时间
}

type SearchAdminPromptDetailResp struct {
	ID          string                            `json:"id"`           // 主键
	ProjectName string                            `json:"project_name"` // 项目名称
	Name        string                            `json:"name"`         // prompt 名称
	Versions    []*SearchAdminPromptDetailVersion `json:"versions"`     // 版本数据
	Locked      bool                              `json:"locked"`       // 是否锁住，true:已锁，false:未锁 上锁后的prompt，删除不可点击
	ModifiedAt  int64                             `json:"modified_at"`  // 更新时间
	CreatedAt   int64                             `json:"created_at"`   // 创建时间
}

type ChatPromptResp struct {
	Content string       `json:"content"`
	Usage   openai.Usage `json:"usage"`
}
