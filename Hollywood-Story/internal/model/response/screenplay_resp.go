package response

import "vland.live/app/internal/model"

type SearchAdminScreenplayResp struct {
	ID             string   `json:"id"`              // 主键
	Key            string   `json:"key"`             // 键
	Name           string   `json:"name" `           // 名称
	Labels         []string `json:"labels"`          // 标签
	Synopsis       string   `json:"synopsis"`        // 概要
	WelcomeMessage string   `json:"welcome_message"` // 欢迎语
	ModifiedAt     int64    `json:"modified_at"`     // 更新时间
	CreatedAt      int64    `json:"created_at"`      // 创建时间
}

type SearchScreenplayResp struct {
	Key      string       `json:"key"`      // 键
	Name     string       `json:"name" `    // 名称
	Actors   model.Actors `json:"actors"`   // 演员列表
	Labels   []string     `json:"labels"`   // 标签
	Synopsis string       `json:"synopsis"` // 概要
}
