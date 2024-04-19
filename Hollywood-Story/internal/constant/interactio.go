package constant

import (
	"github.com/bwmarrin/discordgo"
)

func InteractionRoundEnd() []discordgo.MessageComponent {
	var operationOptions = []discordgo.SelectMenuOption{
		{
			Label:       "Changes in Film Attributes Generate Only One News Article",
			Value:       "Changes in Film Attributes Generate Only One News Article",
			Description: "",
			Default:     false,
		},
		{
			Label:       "Set injuries or deaths turn into news from physical conflicts",
			Value:       "Set injuries or deaths turn into news from physical conflicts",
			Description: "",
			Default:     false,
		},
		{
			Label:       "Player Choices: 'Leak to Paparazzi' or 'Upload to Social Media' Generate News",
			Value:       "Player Choices: 'Leak to Paparazzi' or 'Upload to Social Media' Generate News",
			Description: "",
			Default:     false,
		},
		{
			Label:       "Event Options: 'Leak to Paparazzi' or 'Upload to Social Media' Generate News",
			Value:       "Event Options: 'Leak to Paparazzi' or 'Upload to Social Media' Generate News",
			Description: "",
			Default:     false,
		},
		{
			Label:       "Completion of the Film Itself Becomes a Celebrity News Item",
			Value:       "Completion of the Film Itself Becomes a Celebrity News Item",
			Description: "",
			Default:     false,
		},
	}
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.SelectMenu{
					// Select menu, as other components, must have a customID, so we set it to this value.
					CustomID:    "view_celebrity_news",
					Placeholder: "View celebrity news.",
					Options:     operationOptions,
					Disabled:    false,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    CustomIDButtonRoundEndShareTwitter.Label,
					Style:    discordgo.LinkButton,
					Disabled: false,
					Emoji: &discordgo.ComponentEmoji{
						Name: CustomIDButtonRoundEndShareTwitter.Emoji,
					},
					URL: "https://twitter.com/intent/tweet?text=Hello%20world",
				},
				discordgo.Button{
					Label:    "🎟️Restart game",
					Style:    discordgo.SuccessButton,
					Disabled: false,
					CustomID: CustomIDButtonInputInfo.ID,
				},
				discordgo.Button{
					Label:    CustomIDButtonBotInfo.Label,
					Style:    discordgo.SecondaryButton,
					Disabled: false,
					CustomID: CustomIDButtonBotInfo.ID,
				},
			},
		},
	}
}

func InteractionPlayerCaptureNextRound() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    CustomIDPlayerCaptureNextRound.Label,
					Style:    discordgo.SuccessButton,
					Disabled: false,
					Emoji: &discordgo.ComponentEmoji{
						Name: CustomIDPlayerCaptureNextRound.Emoji,
					},
					CustomID: CustomIDPlayerCaptureNextRound.ID,
				},
				discordgo.Button{
					Label:    "End (Ignore this button, nothing will happen).",
					Style:    discordgo.DangerButton,
					Disabled: false,
					Emoji: &discordgo.ComponentEmoji{
						Name: "⏹️",
					},
					CustomID: "1",
				},
			},
		},
	}
}

func InteractionInputUserInfo() *discordgo.InteractionResponse {
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseModal,
		Data: &discordgo.InteractionResponseData{
			CustomID: CustomIDButtonSubmitInputInfo.ID,
			Title:    CustomIDButtonSubmitInputInfo.Placeholder,
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    CustomIDButtonSubmitInputInfoUsername.ID,
							Label:       CustomIDButtonSubmitInputInfoUsername.Label,
							Style:       discordgo.TextInputShort,
							Placeholder: CustomIDButtonSubmitInputInfoUsername.Placeholder,
							Required:    true,
							MaxLength:   20,
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.TextInput{
							CustomID:    CustomIDButtonSubmitInputInfoAge.ID,
							Label:       CustomIDButtonSubmitInputInfoAge.Label,
							Style:       discordgo.TextInputShort,
							Placeholder: CustomIDButtonSubmitInputInfoAge.Placeholder,
							Value:       "",
							Required:    true,
							MinLength:   1,
							MaxLength:   3,
						},
					},
				},
				//discordgo.ActionsRow{
				//	Components: []discordgo.MessageComponent{
				//		discordgo.TextInput{
				//			CustomID:    CustomIDButtonSubmitInputInfoNPCName.ID,
				//			Label:       CustomIDButtonSubmitInputInfoNPCName.Label,
				//			Style:       discordgo.TextInputShort,
				//			Placeholder: CustomIDButtonSubmitInputInfoNPCName.Placeholder,
				//			Value:       "",
				//			Required:    true,
				//			MinLength:   1,
				//			MaxLength:   10,
				//		},
				//	},
				//},
			},
		},
	}

	return response
}
