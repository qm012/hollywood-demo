package constant

import (
	"strings"
)

type DynamicFieldType string

const (
	DynamicFieldTypeString DynamicFieldType = "string"
	DynamicFieldTypeArray  DynamicFieldType = "array"
)

type DynamicField struct {
	Type  string
	Value string
}

type GptOptionsV1 struct {
	Dialogue  []string `json:"dialogue" bson:"dialogue"`
	Monologue string   `json:"monologue" bson:"monologue"`
	Event     string   `json:"event_v1" bson:"event_v1"`
	A         string   `json:"a" bson:"a"`
	B         string   `json:"b" bson:"b"`
	C         string   `json:"c" bson:"c"`
	D         string   `json:"d" bson:"d"`
}

func (g *GptOptionsV1) FormatDialogue(names ...string) string {

	var dialogue string
	var name string
	if len(names) > 0 {
		name = names[0]
	}

	if length := len(g.Dialogue); length == 1 {
		dialogue = name + ": " + g.Dialogue[0]
	} else if length > 1 {
		slice := make([]string, 0, len(g.Dialogue))
		for _, i2 := range g.Dialogue {
			if i2 == "" {
				continue
			}
			slice = append(slice, i2)
		}

		dialogue = strings.Join(slice, "\n")
	}

	if g.Monologue != "" {
		dialogue = name + ": " + g.Monologue
	}

	return dialogue
}

type GptOutcomeV1 struct {
	Reaction             string `json:"reaction" bson:"reaction"`
	EventProgress        string `json:"event_progress" bson:"event_progress"`
	Outcome              string `json:"outcome" bson:"outcome"`
	Headline             string `json:"headline" bson:"headline"`
	ExtendedConversation string `json:"extended_conversation" bson:"extended_conversation"`
}

type GptNewsV1 struct {
	Post     string   `json:"post" bson:"post"`
	Comments []string `json:"comments" bson:"comments"`
}
