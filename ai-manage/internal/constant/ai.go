package constant

type AIPlatform string

const (
	AIPlatformOpenai AIPlatform = "openai"
	AIPlatformAzure  AIPlatform = "azure"
)

const (
	SSEventChatName = "message"
	SSEventError    = "error"
	SSEventDone     = "[DONE]"
)
