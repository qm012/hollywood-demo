package constant

import (
	"fmt"
	"strings"
)

type GptOptions struct {
	Dialogue  any    `json:"Dialogue" bson:"dialogue"`
	Monologue string `json:"Monologue" bson:"monologue"`
	Event     string `json:"Event" bson:"event"`
	A         string `json:"A" bson:"a"`
	B         string `json:"B" bson:"b"`
	C         string `json:"C" bson:"c"`
	D         string `json:"D" bson:"d"`
}

func (g *GptOptions) V1() *GptOptionsV1 {

	var dialogue = make([]string, 0, 10)

	switch t := g.Dialogue.(type) {
	case string:
		dialogue = append(dialogue, t)
	case []interface{}:

		for _, i2 := range t {
			if i2 == nil {
				continue
			}
			v := fmt.Sprint(i2)
			if v == "" {
				continue
			}
			dialogue = append(dialogue, v)
		}
	}

	return &GptOptionsV1{
		Dialogue:  dialogue,
		Monologue: g.Monologue,
		Event:     g.Event,
		A:         g.A,
		B:         g.B,
		C:         g.C,
		D:         g.D,
	}
}

func (g *GptOptions) FormatDialogue(names ...string) string {

	var dialogue string
	var name string
	if len(names) > 0 {
		name = names[0]
	}
	switch t := g.Dialogue.(type) {
	case string:
		dialogue = name + ": " + t
	case []interface{}:

		slice := make([]string, 0, len(t))
		for _, i2 := range t {
			if i2 == nil {
				continue
			}
			v := fmt.Sprint(i2)
			if v == "" {
				continue
			}
			slice = append(slice, v)
		}

		dialogue = strings.Join(slice, "\n")
	}

	if g.Monologue != "" {
		dialogue = name + ": " + g.Monologue
	}

	return dialogue
}

func (g *GptOptions) Format(attributeValues AttributeValues, names ...string) string {

	format := `
	**Dialogue:**
	%s
	
	**Event:** %s
	
	**A:%s** %s
	**B:%s** %s
	**C:%s** %s
	**D:%s** %s
	`
	var dialogue = g.FormatDialogue(names...)

	return fmt.Sprintf(format, dialogue, g.Event,
		attributeValues[0].Format(), g.A,
		attributeValues[1].Format(), g.B,
		attributeValues[2].Format(), g.C,
		attributeValues[3].Format(), g.D,
	)
}

type GptOutcome struct {
	Reaction             string `json:"Reaction" bson:"reaction"`
	EventProgress        string `json:"Event Progress" bson:"event_progress"`
	Outcome              string `json:"Outcome" bson:"outcome"`
	Headline             string `json:"Headline" bson:"headline"`
	ExtendedConversation string `json:"Extended Conversation" bson:"extended_conversation"`
}

func (g *GptOutcome) V1() *GptOutcomeV1 {
	return &GptOutcomeV1{
		Reaction:             g.Reaction,
		EventProgress:        g.EventProgress,
		Outcome:              g.Outcome,
		Headline:             g.Headline,
		ExtendedConversation: g.ExtendedConversation,
	}
}

func (g *GptOutcome) Format() string {
	format := `%s
`
	return fmt.Sprintf(format, g.EventProgress)
}

type GptNews struct {
	Post     string   `json:"Post"`
	Comments []string `json:"Comments"`
}

func (g *GptNews) V1() *GptNewsV1 {
	return &GptNewsV1{
		Post:     g.Post,
		Comments: g.Comments,
	}
}

func (g *GptNews) Format() string {
	format := `**Post：**%s
**Comments:**
%s
`
	var comments strings.Builder
	for i, comment := range g.Comments {
		comments.WriteString(fmt.Sprintf(" **%v：**%v\n", i+1, comment))
	}
	return fmt.Sprintf(format, g.Post, comments.String())
}
