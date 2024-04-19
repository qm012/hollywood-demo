package constant

import (
	"math/rand"
	"time"
	"vland.live/app/global"
)

type playerOptionRoundEvents []*PlayerOptionRoundEvent

func (p playerOptionRoundEvents) GetByID(ID string) *PlayerOptionRoundEvent {
	if ID == "" {
		return nil
	}
	for _, i2 := range p {
		if i2.ID == ID {
			return i2
		}
	}
	return nil
}

type PlayerOptionRoundEventType string

const (
	PlayerOptionRoundEventTypeCapture   PlayerOptionRoundEventType = "capture"   // 拍摄事件
	PlayerOptionRoundEventTypePromotion PlayerOptionRoundEventType = "promotion" // 推广事件
)

type OptionOutcome struct {
	Text      string `json:"text" bson:"text"`
	ShortText string `json:"short_text" bson:"short_text"`
}

type PlayerOptionRoundEventPrompt struct {
	ChoicesID string `json:"choices_id" bson:"choices_id"`
	OutcomeID string `json:"outcome_id" bson:"outcome_id"`
}

// PlayerOptionRoundEvent 玩家每一回合的随机事件
type PlayerOptionRoundEvent struct {
	ID          string                       `json:"id" bson:"ID"`                   // 主键
	EventType   PlayerOptionRoundEventType   `json:"event_type" bson:"event_type"`   // 事件类型
	Description string                       `json:"description" bson:"description"` // 描述
	Prompt      PlayerOptionRoundEventPrompt `json:"prompt" bson:"prompt"`           // 当前事件对应ID
}

// Probability 根据事件类型计算概率
// 公式：6. 事件本身难度值(10-15) <=（玩家属性  + 选项结果的随机数）= 成功
//  1. 事件本身难度值(10-15) => 显示到Weather下面
//  2. 出结果的随机值（1-20）=>不展示
//  3. 用户的属性值（默认1-10）
//     公式：2+3 和 1 比大小
//     过期：1. 大于等于则成功 => 5% 为大成功
//     过期：2. 小于则失败 => %5 为大失败
//
// 只出一次 1：大失败 20：大成功，后续再比较数值
//
//	 params:
//		@difficultyValue 事件类型本身难度值
//		@attributeValue 用户属性值
func (p *PlayerOptionRoundEvent) Probability(difficultyValue, attributeValue int, specialNpcName string) *OptionOutcome {
	var (
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
		// 出结果随机一个值
		randomValue = r.Intn(20) + 1

		greatFailureValue = 1
		greatSuccessValue = 20

		text      string
		shortText string
	)
	var (
		bigSuccess = "Skill check: big success"
		bigFailure = "Skill check: big failure"
		success    = "Skill check: success"
		failure    = "Skill check: failure"
	)
	//Skill check: big success
	//Skill check: big failure
	//skill check：success
	//skill check：failure
	switch p.ID {
	case PlayerCaptureOptionRoundEventTwoActorsInteraction.ID,
		PlayerCaptureOptionRoundEventPlayerActorInteraction.ID,
		PlayerPromotionOptionRoundEventMeetAndGreet.ID:
		if randomValue == greatFailureValue { // 大失败
			text = "The film quality decreases significantly. but there is a certain rise in popularity."
			shortText = bigFailure
		} else if randomValue == greatSuccessValue { // 大成功
			text = "The film quality improves significantly."
			shortText = bigSuccess
		} else {
			if randomValue+attributeValue >= difficultyValue {
				// 成功
				text = "The film quality improves."
				shortText = success
			} else {
				// 失败
				text = "The film quality decreases."
				shortText = failure
			}
		}
	case PlayerCaptureOptionRoundEventSpecialActorEvent.ID,
		PlayerCaptureOptionRoundEventSpecialNpcAndNpcAndPlayerEvent.ID:
		if randomValue == greatFailureValue { // 大失败
			text = specialNpcName + "'s affection towards you has decreased significantly."
			shortText = bigFailure
		} else if randomValue == greatSuccessValue { // 大成功
			text = specialNpcName + "'s affection towards you has increased significantly."
			shortText = bigSuccess
		} else {
			if randomValue+attributeValue >= difficultyValue {
				// 成功
				text = specialNpcName + "'s affection towards you has increased."
				shortText = success
			} else {
				// 失败
				text = specialNpcName + "'s affection towards you has decreased."
				shortText = failure
			}
		}
	}

	return &OptionOutcome{
		Text:      text,
		ShortText: shortText,
	}
}

var (
	PlayerCaptureOptionRoundEvents = playerOptionRoundEvents{
		PlayerCaptureOptionRoundEventTwoActorsInteraction,
		PlayerCaptureOptionRoundEventSpecialActorEvent,
		PlayerCaptureOptionRoundEventPlayerActorInteraction,
		PlayerCaptureOptionRoundEventSpecialNpcAndNpcAndPlayerEvent,
	}
	PlayerPromotionOptionRoundEvents = playerOptionRoundEvents{
		PlayerPromotionOptionRoundEventMeetAndGreet,
	}
)

// 默认使用生产的
var (
	PlayerCaptureOptionRoundEventTwoActorsInteraction = &PlayerOptionRoundEvent{
		ID:          "TwoActorsInteraction",
		EventType:   PlayerOptionRoundEventTypeCapture,
		Description: "Two actors interact, and the player chooses to intervene in the interaction.",
		Prompt: PlayerOptionRoundEventPrompt{
			ChoicesID: "6583a9c332314355db5f7b03",
			OutcomeID: "6583b09932314355db5f7b09",
		},
	}
	PlayerCaptureOptionRoundEventSpecialActorEvent = &PlayerOptionRoundEvent{
		ID:          "SpecialActorEvent",
		EventType:   PlayerOptionRoundEventTypeCapture,
		Description: "Special actors visit the set as an event, hostile actors attempt to sabotage the film shoot, and friendly actors help with the film shoot.",
		Prompt: PlayerOptionRoundEventPrompt{
			ChoicesID: "6583c9ba32314355db5f7b18",
			OutcomeID: "6583ca3d32314355db5f7b1b",
		},
	}
	PlayerCaptureOptionRoundEventPlayerActorInteraction = &PlayerOptionRoundEvent{
		ID:          "PlayerActorInteraction",
		EventType:   PlayerOptionRoundEventTypeCapture,
		Description: "Player engages in a solo interaction with a specific actor, and the player chooses their response.",
		Prompt: PlayerOptionRoundEventPrompt{
			ChoicesID: "6583b52232314355db5f7b0c",
			OutcomeID: "6583b5cb32314355db5f7b0f",
		},
	}
	PlayerCaptureOptionRoundEventSpecialNpcAndNpcAndPlayerEvent = &PlayerOptionRoundEvent{
		ID:          "SpecialNpcAndNpcAndPlayerInteraction",
		EventType:   PlayerOptionRoundEventTypeCapture,
		Description: "test",
		Prompt: PlayerOptionRoundEventPrompt{
			ChoicesID: "6583b6bb32314355db5f7b12",
			OutcomeID: "6583b6c432314355db5f7b14",
		},
	}

	PlayerPromotionOptionRoundEventMeetAndGreet = &PlayerOptionRoundEvent{
		ID:          "MeetAndGreet",
		EventType:   PlayerOptionRoundEventTypePromotion,
		Description: "test",
		Prompt: PlayerOptionRoundEventPrompt{
			ChoicesID: "65b34219f87b25725998f423",
			OutcomeID: "65b34233f87b25725998f425",
		},
	}
)

var (
	OutcomeNewsPromptID = "6583af4c32314355db5f7b06"
)

func (p playerOptionRoundEvents) Random() *PlayerOptionRoundEvent {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return p[r.Intn(len(p))]
}

func InitCaptureEvent() {

	if global.Config.App.Develop() {
		PlayerCaptureOptionRoundEventTwoActorsInteraction.Prompt = PlayerOptionRoundEventPrompt{
			ChoicesID: "6583a9c332314355db5f7b03",
			OutcomeID: "6583b09932314355db5f7b09",
		}
		PlayerCaptureOptionRoundEventSpecialActorEvent.Prompt = PlayerOptionRoundEventPrompt{
			ChoicesID: "6583c9ba32314355db5f7b18",
			OutcomeID: "6583ca3d32314355db5f7b1b",
		}
		PlayerCaptureOptionRoundEventPlayerActorInteraction.Prompt = PlayerOptionRoundEventPrompt{
			ChoicesID: "6583b52232314355db5f7b0c",
			OutcomeID: "6583b5cb32314355db5f7b0f",
		}
		PlayerCaptureOptionRoundEventSpecialNpcAndNpcAndPlayerEvent.Prompt = PlayerOptionRoundEventPrompt{
			ChoicesID: "6583b6bb32314355db5f7b12",
			OutcomeID: "6583b6c432314355db5f7b14",
		}
		PlayerPromotionOptionRoundEventMeetAndGreet.Prompt = PlayerOptionRoundEventPrompt{
			ChoicesID: "65b34219f87b25725998f423",
			OutcomeID: "65b34233f87b25725998f425",
		}

		OutcomeNewsPromptID = "6583af4c32314355db5f7b06"
	}

}
