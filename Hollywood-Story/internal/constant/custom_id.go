package constant

const (
	CommandName        = "hollywood"
	CommandNameOption1 = "start"
)

type CustomID struct {
	ID          string
	Label       string
	Emoji       string
	Placeholder string
}

var (
	CustomIDButtonInputInfo = CustomID{
		Label: "🎟️Starts to enter information",
		ID:    "input_user_info",
	}
	CustomIDButtonSubmitInputInfo = CustomID{
		Label:       "Submit entry information",
		ID:          "submit_input_user_info",
		Placeholder: "Please set your personal information",
	}
	CustomIDButtonSubmitInputInfoUsername = CustomID{
		Label:       "nickname",
		ID:          "input_user_info_nickname",
		Placeholder: "Please enter your nickname",
	}
	CustomIDButtonSubmitInputInfoAge = CustomID{
		Label:       "age",
		ID:          "input_user_info_age",
		Placeholder: "Please enter your age（Must be a number.）",
	}
	CustomIDButtonSubmitInputInfoNPCName = CustomID{
		Label:       "npc_name",
		ID:          "input_user_info_npc_name",
		Placeholder: "Please enter your npc name",
	}
	CustomIDButtonStart = CustomID{
		Label: "🎬Start the game",
		ID:    "start",
	}
	CustomIDButtonBotInfo = CustomID{
		Label: "❔Bot introduction",
		ID:    "bot_info",
	}
	CustomIDSelectPlayerInfoGender = CustomID{
		Label:       "",
		ID:          "player_info_gender",
		Placeholder: "Please select a gender",
	}
	CustomIDSelectPlayerInfoPersonality = CustomID{
		Label:       "",
		ID:          "player_info_personality",
		Placeholder: "Please select a personality",
	}
	CustomIDSelectPlayerInfoScreenplay = CustomID{
		Label:       "",
		ID:          "player_info_screenplay",
		Placeholder: "Please select a screenplay",
	}
	CustomIDSelectPlayerInfoOccupation = CustomID{
		Label:       "",
		ID:          "player_info_occupation",
		Placeholder: "Please select a occupation",
	}
	CustomIDPlayerCaptureSelectOption = CustomID{
		Label:       "",
		ID:          "player_capture_select_option",
		Placeholder: "Please select your operation 👇",
	}
	CustomIDRefreshPlayerAttributes = CustomID{
		ID:    "refresh_player_attributes",
		Label: "Attributes value",
		Emoji: "🔄",
	}
	CustomIDPlayerCaptureNextRound = CustomID{
		ID:          "player_capture_next_round",
		Label:       "Next Round",
		Emoji:       "⏭️",
		Placeholder: "",
	}
	CustomIDButtonRoundEndShareTwitter = CustomID{
		Label:       "Share on Twitter",
		Emoji:       "💯",
		Placeholder: "",
	}
)
