package constant

import (
	"github.com/bwmarrin/discordgo"
)

type playerCaptureOptions []*PlayerCaptureOption

// PlayerCaptureOption 玩家拍摄选项
type PlayerCaptureOption struct {
	ID    string
	Name  string
	Emoji string
}

var (
	PlayerCaptureOptions = playerCaptureOptions{PlayerCaptureOptionGoodness, PlayerCaptureOptionEvil, PlayerCaptureOptionNonIntervene, PlayerCaptureOptionPrank}
)

var (
	PlayerCaptureOptionGoodness = &PlayerCaptureOption{
		ID:   "player_capture_select_option_goodness",
		Name: "玩家行善",
		//Emoji: "❤️",
	}
	PlayerCaptureOptionEvil = &PlayerCaptureOption{
		ID:   "player_capture_select_option_evil",
		Name: "玩家作恶",
		//Emoji: "🦹",
	}
	PlayerCaptureOptionNonIntervene = &PlayerCaptureOption{
		ID:   "player_capture_select_option_nonintervene",
		Name: "玩家旁观（不参与）",
		//Emoji: "\U0001FAE3",
	}
	PlayerCaptureOptionPrank = &PlayerCaptureOption{
		ID:   "player_capture_select_option_prank",
		Name: "玩家搞怪（恶作剧）",
		//Emoji: "😈",
	}
)

func (p playerCaptureOptions) GetByID(ID string) *PlayerCaptureOption {
	if ID == "" {
		return nil
	}
	for _, playerCaptureOption := range p {
		if playerCaptureOption.ID == ID {
			return playerCaptureOption
		}
	}
	return nil
}

func (p playerCaptureOptions) SelectMessageComponent(attributeValues AttributeValues) (*[]discordgo.MessageComponent, map[string]*AttributeValue) {

	//var operationOptions = []discordgo.SelectMenuOption{
	//	{
	//		Label:       PlayerCaptureOptionGoodness.String(),
	//		Value:       string(PlayerCaptureOptionGoodness),
	//		Description: "",
	//		Emoji: discordgo.ComponentEmoji{
	//			Name: PlayerCaptureOptionGoodness.Emoji(),
	//		},
	//		Default: false,
	//	},
	//	{
	//		Label:       PlayerCaptureOptionEvil.String(),
	//		Value:       string(PlayerCaptureOptionEvil),
	//		Description: "",
	//		Emoji: discordgo.ComponentEmoji{
	//			Name: PlayerCaptureOptionEvil.Emoji(),
	//		},
	//		Default: false,
	//	},
	//	{
	//		Label:       PlayerCaptureOptionNonIntervene.String(),
	//		Value:       string(PlayerCaptureOptionNonIntervene),
	//		Description: "",
	//		Emoji: discordgo.ComponentEmoji{
	//			Name: PlayerCaptureOptionNonIntervene.Emoji(),
	//		},
	//		Default: false,
	//	},
	//	{
	//		Label:       PlayerCaptureOptionPrank.String(),
	//		Value:       string(PlayerCaptureOptionPrank),
	//		Description: "",
	//		Emoji: discordgo.ComponentEmoji{
	//			Name: PlayerCaptureOptionPrank.Emoji(),
	//		},
	//		Default: false,
	//	},
	//}

	actionsRowComponents := make([]discordgo.MessageComponent, 0, len(p))
	buttonAttributeMap := make(map[string]*AttributeValue, 4)
	for _, playerCaptureOption := range p {

		var label string

		switch playerCaptureOption {
		case PlayerCaptureOptionGoodness:

			attributeValue := attributeValues[0]
			buttonAttributeMap[PlayerCaptureOptionGoodness.ID] = attributeValue.Copy()

			label = "A: " + attributeValue.Format()
		case PlayerCaptureOptionEvil:

			attributeValue := attributeValues[1]
			buttonAttributeMap[PlayerCaptureOptionEvil.ID] = attributeValue.Copy()

			label = "B: " + attributeValue.Format()
		case PlayerCaptureOptionNonIntervene:

			attributeValue := attributeValues[2]
			buttonAttributeMap[PlayerCaptureOptionNonIntervene.ID] = attributeValue.Copy()

			label = "C: " + attributeValue.Format()
		case PlayerCaptureOptionPrank:

			attributeValue := attributeValues[3]
			buttonAttributeMap[PlayerCaptureOptionPrank.ID] = attributeValue.Copy()

			label = "D: " + attributeValue.Format()
		}

		actionsRowComponents = append(actionsRowComponents, discordgo.Button{
			Label:    label,
			Style:    discordgo.PrimaryButton,
			Disabled: false,
			//Emoji: &discordgo.ComponentEmoji{
			//	Name: playerCaptureOption.Emoji,
			//},
			CustomID: playerCaptureOption.ID,
		})
	}

	return &[]discordgo.MessageComponent{
		//discordgo.ActionsRow{
		//	Components: []discordgo.MessageComponent{
		//		discordgo.SelectMenu{
		//			// Select menu, as other components, must have a customID, so we set it to this value.
		//			CustomID:    CustomIDPlayerCaptureSelectOption.ID,
		//			Placeholder: CustomIDPlayerCaptureSelectOption.Placeholder,
		//			Options:     operationOptions,
		//			Disabled:    false,
		//		},
		//	},
		//},
		discordgo.ActionsRow{
			Components: actionsRowComponents,
		},
	}, buttonAttributeMap
}

////Probability 选项对应的概率
//func (p *PlayerCaptureOption) Probability() PlayOptionOutcome {
//	time.Sleep(time.Microsecond)
//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
//	probabilityValue := r.Float64()
//	switch p {
//	case PlayerCaptureOptionGoodness:
//		if probabilityValue < 0.75 {
//			return PlayerCaptureOptionOutcomeSuccess
//		} else if probabilityValue < 0.90 {
//			return PlayerCaptureOptionOutcomeFailure
//		} else if probabilityValue < 0.95 {
//			return PlayerCaptureOptionOutcomeCriticalSuccess
//		} else {
//			return PlayerCaptureOptionOutcomeCriticalFailure
//		}
//	case PlayerCaptureOptionEvil:
//		if probabilityValue < 0.75 {
//			return PlayerCaptureOptionOutcomeFailure
//		} else if probabilityValue < 0.90 {
//			return PlayerCaptureOptionOutcomeSuccess
//		} else if probabilityValue < 0.95 {
//			return PlayerCaptureOptionOutcomeCriticalSuccess
//		} else {
//			return PlayerCaptureOptionOutcomeCriticalFailure
//		}
//	case PlayerCaptureOptionNonIntervene:
//		if probabilityValue < 0.50 {
//			return PlayerCaptureOptionOutcomeSuccess
//		} else {
//			return PlayerCaptureOptionOutcomeFailure
//		}
//	case PlayerCaptureOptionPrank:
//		if probabilityValue < 0.25 {
//			return PlayerCaptureOptionOutcomeSuccess
//		} else if probabilityValue < 0.50 {
//			return PlayerCaptureOptionOutcomeFailure
//		} else if probabilityValue < 0.75 {
//			return PlayerCaptureOptionOutcomeCriticalSuccess
//		} else {
//			return PlayerCaptureOptionOutcomeCriticalFailure
//		}
//	}
//	return nil
//}
//
//// PlayerCaptureOptionOutcome 玩家拍摄选项结果（基于每个选项有不同的概率）
//type PlayerCaptureOptionOutcome string
//
//const (
//	PlayerCaptureOptionOutcomeSuccess         PlayerCaptureOptionOutcome = "success"
//	PlayerCaptureOptionOutcomeFailure         PlayerCaptureOptionOutcome = "failure"
//	PlayerCaptureOptionOutcomeCriticalSuccess PlayerCaptureOptionOutcome = "criticalSuccess"
//	PlayerCaptureOptionOutcomeCriticalFailure PlayerCaptureOptionOutcome = "criticalFailure"
//)
//
//func (p PlayerCaptureOptionOutcome) String() string {
//	switch p {
//	case PlayerCaptureOptionOutcomeSuccess:
//		return "Movie quality goes UP!"
//	case PlayerCaptureOptionOutcomeFailure:
//		return "Movie quality goes DOWN!"
//	case PlayerCaptureOptionOutcomeCriticalSuccess:
//		return "Movie quality drastically goes UP!"
//	case PlayerCaptureOptionOutcomeCriticalFailure:
//		return "Movie quality drastically goes down! Movie gets more coverage on SNS."
//	}
//	return ""
//}
