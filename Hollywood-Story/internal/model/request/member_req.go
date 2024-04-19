package request

import (
	"vland.live/app/internal/constant"
	"vland.live/app/internal/model"
)

type InputMemberInfoReq struct {
	DeviceID       string   `json:"device_id" binding:"required"`                 // 暂时以设备ID用作用户唯一ID
	Nickname       string   `json:"nickname" binding:"required,max=100"`          // 用户输入的昵称
	SpecialNPCName string   `json:"special_npc_name" binding:"omitempty,max=100"` // 特殊 NPC name
	Age            int      `json:"age" binding:"required,max=200"`               // 年龄
	Gender         string   `json:"gender" binding:"required,max=100"`            // 性别
	Personality    []string `json:"personality" binding:"omitempty,max=100"`      // 性格
	Occupation     string   `json:"occupation" binding:"omitempty,max=100"`       // 职业
}

type UpdateMemberScreenplayReq struct {
	DeviceID string       `json:"device_id" binding:"required"`      // 暂时以设备ID用作用户唯一ID
	Key      string       `json:"key" binding:"required,max=100"`    // 剧本键
	Actors   model.Actors `json:"actors" binding:"required,max=100"` // 演员列表
}

type GetMemberInfoReq struct {
	DeviceID string `form:"device_id" binding:"required,max=1000"` // 暂时以设备ID用作用户唯一ID
}

type RefreshAttributesReq struct {
	DeviceID string `json:"device_id" binding:"required,max=1000"` // 暂时以设备ID用作用户唯一ID
}

type StartOrNextRoundReq struct {
	DeviceID       string                              `json:"device_id" binding:"required,max=1000"`                  // 暂时以设备ID用作用户唯一ID
	EventType      constant.PlayerOptionRoundEventType `json:"event_type" binding:"omitempty,oneof=capture promotion"` // 事件类型
	EventID        string                              `json:"event_id" binding:"required,max=100"`                    // 事件ID
	SpecialNPCName string                              `json:"special_npc_name" binding:"omitempty,max=100"`           // 特殊 NPC name
}

func (s *StartOrNextRoundReq) GetEvent() *constant.PlayerOptionRoundEvent {

	if s.EventType == "" {
		// 默认赋值
		s.EventType = constant.PlayerOptionRoundEventTypeCapture
	}

	var roundEvent *constant.PlayerOptionRoundEvent
	switch s.EventType {
	case constant.PlayerOptionRoundEventTypeCapture:
		roundEvent = constant.PlayerCaptureOptionRoundEvents.GetByID(s.EventID)
	case constant.PlayerOptionRoundEventTypePromotion:
		roundEvent = constant.PlayerPromotionOptionRoundEvents.GetByID(s.EventID)
	}

	return roundEvent
}

type ClickButtonOutcomeReq struct {
	DeviceID string `json:"device_id" binding:"required,max=1000"` // 暂时以设备ID用作用户唯一ID
	OptionID string `json:"option_id" binding:"required,max=100"`  // 选项ID
}

type ClickButtonNewsByOutcomeReq struct {
	DeviceID string `json:"device_id" binding:"required,max=1000"` // 暂时以设备ID用作用户唯一ID
}
