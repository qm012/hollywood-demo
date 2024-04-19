package response

import (
	"vland.live/app/internal/constant"
	"vland.live/app/internal/model"
)

type GetMemberInfoResp struct {
	ID             string                    `json:"id"`               // 主键
	DeviceID       string                    `json:"device_id"`        // 暂时以设备ID用作用户唯一ID
	Nickname       string                    `json:"nickname"`         // 用户输入的昵称
	SpecialNPCName string                    `json:"special_npc_name"` // 特殊 NPC name
	Age            int                       `json:"age"`              // 年龄
	Gender         string                    `json:"gender"`           // 性别
	Personality    []string                  `json:"personality"`      // 性格
	Occupation     string                    `json:"occupation"`       // 职业
	Attribute      *constant.Attribute       `json:"attribute"`        // 用户属性值
	CurrentRound   *model.MemberCurrentRound `json:"current_round"`    // 当前回合的数据
	Film           *model.MemberFilm         `json:"film"`             // 影片信息
}

type RefreshMemberAttributesResp struct {
	Text   string         `json:"text"`   // 格式化的属性文本
	Values map[string]int `json:"values"` // 属性对应的随机数
}

type StartOrNextRoundResp struct {
	EventTheme      string                   `json:"event_theme"`      // 事件主题
	Actors          model.Actors             `json:"actors"`           // 当前回合的演员
	DifficultyValue int                      `json:"difficulty_value"` // 当前回合的事件类型难度值
	RoundNumber     int                      `json:"round_number"`     // 当前的回合数
	WelcomeMessage  string                   `json:"welcome_message"`  // 欢迎语，只有第一轮才有
	Location        string                   `json:"location"`         // 地点
	Weather         string                   `json:"weather"`          // 天气
	AttributeValues constant.AttributeValues `json:"attribute_values"` // 4个属性列表，需要按顺序处理
	GptOptions      *constant.GptOptionsV1   `json:"gpt_options"`      // gpt 返回的选项
}

type ClickButtonOutcomeResp struct {
	OptionOutcome *constant.OptionOutcome `json:"option_outcome"` // 选项结果
	GptOutcome    *constant.GptOutcomeV1  `json:"gpt_outcome"`    // Gpt返回的结果
}

type ClickButtonNewsByOutcomeResp struct {
	GptNews *constant.GptNewsV1 `json:"gpt_news"` // gpt返回的新闻评论
}
