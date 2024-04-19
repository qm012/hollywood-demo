package request

import (
	"errors"
	"github.com/qm012/dun"
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vland.live/app/initialize/params"
)

type CreateAdminPromptReq struct {
	ProjectID string `json:"project_id" binding:"required,objectId"` // 项目ID
	Name      string `json:"name" binding:"required,max=50"`         // prompt 名称
}

func (c *CreateAdminPromptReq) GetProjectID() primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(c.ProjectID)
	return ID
}

type UpdateAdminPromptLockedReq struct {
	params.IdParamOmit
}

type UpdateAdminPromptVersionIsProductionReq struct {
	params.IdParamOmit        // 主键
	VersionID          string // 版本ID
	Operator           string `json:"operator" binding:"required,max=50"`
}

func (u *UpdateAdminPromptVersionIsProductionReq) GetVersionID() primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(u.VersionID)
	return ID
}

type UpdateAdminPromptVersionNameReq struct {
	params.IdParamOmit        // 主键
	VersionID          string // 版本ID
	Name               string `json:"name" binding:"required,max=50"`
	Operator           string `json:"operator" binding:"required,max=50"`
}

func (u *UpdateAdminPromptVersionNameReq) GetVersionID() primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(u.VersionID)
	return ID
}

type DeleteAdminPromptVersionReq struct {
	params.IdParamOmit        // 主键
	VersionID          string // 版本ID
	Operator           string `json:"operator" binding:"required,max=50"`
}

func (u *DeleteAdminPromptVersionReq) GetVersionID() primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(u.VersionID)
	return ID
}

type UpdateAdminPromptReq struct {
	params.IdParamOmit

	ProjectID string `json:"project_id" binding:"required,objectId"` // 项目ID
	Name      string `json:"name" binding:"required,max=50"`         // prompt 名称
}

func (u *UpdateAdminPromptReq) GetProjectID() primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(u.ProjectID)
	return ID
}

type DeleteAdminPromptReq struct {
	params.IdParamOmit
}

type SearchAdminPromptDetailReq struct {
	params.IdParamOmit
}

type SearchAdminPromptReq struct {
	ProjectID string `form:"project_id" binding:"omitempty,objectId"` // 项目ID
	Name      string `form:"name" binding:"omitempty,max=50"`         // prompt 名称
	dun.PageSearch
}

func (s *SearchAdminPromptReq) GetProjectID() primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(s.ProjectID)
	return ID
}

func (s *SearchAdminPromptReq) Filter() bson.D {
	filter := make(bson.D, 0, 2)
	if len(s.Name) != 0 {

		tempFilter := make(bson.A, 0, 2)

		tempFilter = append(tempFilter, bson.M{
			"name": bson.M{"$regex": primitive.Regex{
				Pattern: s.Name,
				Options: "i", // 不区分大小写 i
			}},
		})
		filter = append(filter, bson.E{Key: "$or", Value: tempFilter})
	}

	if s.ProjectID != "" {
		filter = append(filter, bson.E{Key: "projectID", Value: s.GetProjectID()})
	}

	return filter
}

type SaveVersionOperationType string

const (
	SaveVersionOperationTypeOverride SaveVersionOperationType = "override" // 覆盖
	SaveVersionOperationTypeAddition SaveVersionOperationType = "addition" // 新增
)

type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	// This property isn't in the official documentation, but it's in
	// the documentation for the official library for python:
	// - https://github.com/openai/openai-python/blob/main/chatml.md
	// - https://github.com/openai/openai-cookbook/blob/main/examples/How_to_count_tokens_with_tiktoken.ipynb
	Name string `json:"name"`
}

type SaveAdminVersionVariable struct {
	Key   string `json:"key" binding:"required,max=100"`
	Value string `json:"value" binding:"omitempty,max=100000"`
}

type SaveAdminVersionReq struct {
	params.IdParamOmit                                         // prompt主键
	Variables          []*SaveAdminVersionVariable             `json:"variables" binding:"omitempty,max=100"`   // 变量列表
	VersionID          string                                  `json:"version_id" binding:"omitempty,objectId"` // 版本ID，当操作类型=SaveVersionOperationTypeOverride，则需要填写
	Name               string                                  `json:"name" binding:"omitempty,max=20"`         // 版本名称，当操作类型=SaveVersionOperationTypeAddition，则需要填写
	Model              string                                  `json:"model" binding:"required"`
	Messages           []ChatCompletionMessage                 `json:"messages" binding:"required,min=1,max=1000,dive"`
	MaxTokens          int                                     `json:"max_tokens" binding:"omitempty,max=4096"`
	Temperature        float32                                 `json:"temperature" binding:"omitempty"`
	TopP               float32                                 `json:"top_p" binding:"omitempty"`
	FrequencyPenalty   float32                                 `json:"frequency_penalty" binding:"omitempty"`
	ResponseFormat     openai.ChatCompletionResponseFormatType `json:"response_format" binding:"omitempty,oneof=json_object text"`
	OperationType      SaveVersionOperationType                `json:"operation_type" binding:"required,oneof=override addition"` // 操作类型
	Operator           string                                  `json:"operator" binding:"required,max=50"`
}

func (s *SaveAdminVersionReq) GetVersionID() primitive.ObjectID {
	ID, _ := primitive.ObjectIDFromHex(s.VersionID)
	return ID
}

var completionModels = map[string]bool{
	openai.GPT3Dot5Turbo:        true,
	openai.GPT3Dot5Turbo0301:    true,
	openai.GPT3Dot5Turbo0613:    true,
	openai.GPT3Dot5Turbo1106:    true,
	openai.GPT3Dot5Turbo16K:     true,
	openai.GPT3Dot5Turbo16K0613: true,
	openai.GPT4:                 true,
	openai.GPT4TurboPreview:     true,
	openai.GPT4VisionPreview:    true,
	openai.GPT40314:             true,
	openai.GPT40613:             true,
	openai.GPT432K:              true,
	openai.GPT432K0314:          true,
	openai.GPT432K0613:          true,
}

func (s *SaveAdminVersionReq) Verify() error {

	if !completionModels[s.Model] {
		return errors.New("不支持的模型")
	}
	switch s.OperationType {
	case SaveVersionOperationTypeOverride:
		if s.VersionID == "" {
			return errors.New("版本ID不能为空")
		}
	case SaveVersionOperationTypeAddition:
		if s.Name == "" {
			return errors.New("版本名称不能为空")
		}
	}
	return nil
}

type CreateAdminPromptVersionReq struct {
	params.IdParamOmit        // prompt主键
	Name               string `json:"name" binding:"required,max=20"`     // 版本名称
	Operator           string `json:"operator" binding:"required,max=50"` // 操作人
}

type ChatPromptReq struct {
	params.IdParamOmit
	Variables      map[string]string       `json:"variables" binding:"required,max=10000"`
	AppendMessages []ChatCompletionMessage `json:"append_messages" binding:"omitempty,max=1000,dive"`
}
