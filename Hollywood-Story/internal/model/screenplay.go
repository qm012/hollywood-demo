package model

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"strings"
)

type Screenplays []*Screenplay

type Screenplay struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`                          // 主键
	Key            string             `json:"key" bson:"key"`                         // 键
	Name           string             `json:"name" bson:"name"`                       // 名称
	Labels         []string           `json:"labels" bson:"labels"`                   // 标签
	Synopsis       string             `json:"synopsis" bson:"synopsis"`               // 概要
	WelcomeMessage string             `json:"welcome_message" bson:"welcome_message"` // 欢迎语
	ModifiedAt     int64              `json:"modified_at" bson:"modified_at"`         // 更新时间
	CreatedAt      int64              `json:"created_at" bson:"created_at"`           // 创建时间
}

// Format 格式化剧本的信息
func (s *Screenplay) Format() string {
	format := fmt.Sprintf(`Title：%s
	Label: %s
	Synopsis:%s`, s.Name, strings.Join(s.Labels, ","), s.Synopsis)
	return format
}

// FormatSimple 格式化剧本的信息:简易版
func (s *Screenplay) FormatSimple() string {
	format := fmt.Sprintf(`Title：%s
	Label: %s`, s.Name, strings.Join(s.Labels, ","))
	return format
}
