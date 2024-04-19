package model

import (
	"github.com/sashabaranov/go-openai"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Prompts []*Prompt
type PromptVersions []*PromptVersion

func (p Prompts) GetProjectIDs() []primitive.ObjectID {
	var (
		length     = len(p)
		projectIDs = make([]primitive.ObjectID, 0, length)
	)

	for _, prompt := range p {
		projectIDs = append(projectIDs, prompt.ProjectID)
	}

	return projectIDs
}

func (p PromptVersions) Contains(id primitive.ObjectID) bool {

	if id.IsZero() {
		return false
	}

	for _, promptVersion := range p {
		if promptVersion.ID == id {
			return true
		}
	}

	return false
}

func (p PromptVersions) GetByID(id primitive.ObjectID) *PromptVersion {

	if id.IsZero() {
		return nil
	}

	for _, promptVersion := range p {
		if promptVersion.ID == id {
			return promptVersion
		}
	}

	return nil
}

type PromptVersionVariable struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

type PromptVersion struct {
	ID           primitive.ObjectID            `json:"id" bson:"_id"`                     // 主键
	Name         string                        `json:"name" bson:"name"`                  // 版本名称
	IsProduction bool                          `json:"is_production" bson:"isProduction"` // 如果是生产版本，则为true；否则为false
	ChatReq      *openai.ChatCompletionRequest `json:"chat_req" bson:"chatReq"`           // 聊天请求参数
	Variables    []*PromptVersionVariable      `json:"variables" bson:"variables"`        // 变量数据
	Modifier     string                        `json:"modifier" bson:"modifier"`          // 更新人
	Creator      string                        `json:"creator" bson:"creator"`            // 创建人
	ModifiedAt   int64                         `json:"modified_at" bson:"modifiedAt"`     // 更新时间
	CreatedAt    int64                         `json:"created_at" bson:"createdAt"`       // 创建时间
}

type Prompt struct {
	ID         primitive.ObjectID `json:"id" bson:"_id"`                 // 主键
	ProjectID  primitive.ObjectID `json:"project_id" bson:"projectID"`   // 项目ID
	Name       string             `json:"name" bson:"name"`              // prompt 名称
	Versions   PromptVersions     `json:"versions" bson:"versions"`      // 版本数据
	Locked     bool               `json:"locked" bson:"locked"`          // 是否锁住，true:已锁，false:未锁 上锁后的prompt，删除不可点击
	ModifiedAt int64              `json:"modified_at" bson:"modifiedAt"` // 更新时间
	CreatedAt  int64              `json:"created_at" bson:"createdAt"`   // 创建时间
}
