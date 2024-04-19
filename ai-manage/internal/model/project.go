package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Projects []*Project

type Project struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`                         // 主键
	Name           string             `json:"name" bson:"name"`                      // 项目名称
	PromptQuantity int                `json:"prompt_quantity" bson:"promptQuantity"` // prompt 数量
	ModifiedAt     int64              `json:"modified_at" bson:"modifiedAt"`         // 更新时间
	CreatedAt      int64              `json:"created_at" bson:"createdAt"`           // 创建时间
}
