package constant

import (
	"fmt"
	"math/rand"
	"time"
)

type Slice []string

func (s Slice) Random() string {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return s[r.Intn(len(s))]
}

func (s Slice) RandomNoRepeat(num int) Slice {
	if num <= 1 {
		num = 1
	}

	if num >= len(s) {
		num = len(s)
	}

	copiedValue := make(map[string]struct{}, len(s))
	for _, value := range s {
		copiedValue[value] = struct{}{}
	}

	results := make(Slice, 0, num)

	for num > 0 {
		keys := make([]string, 0, len(copiedValue))
		for key := range copiedValue {
			keys = append(keys, key)
		}
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		randomKey := keys[r.Intn(len(copiedValue))]

		results = append(results, randomKey)
		delete(copiedValue, randomKey)
		num--
	}

	return results
}

type location struct {
	values map[string]Slice
}

func (l *location) Random() string {
	valueLength := len(l.values)
	keys := make([]string, 0, valueLength)
	for key := range l.values {
		keys = append(keys, key)
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randomKey := keys[r.Intn(valueLength)]
	return fmt.Sprintf("%s-%s", randomKey, l.values[randomKey].Random())
}

var (
	// Location 地点数据
	Location = &location{
		values: map[string]Slice{
			"Film studio": {"public dressing room", "Leading dressing room", "rehearsal room", "screenwriter office", "director office", "producer office", "commissary", "sound stage", "studio backlot", "warehouse"},
			//"Library":       {"library"},
			//"café":          {"Lounge", "private booth", "kitchen"},
			//"gym":           {"lobby", "private room", "dressing room"},
			//"Act studio":    {"classroom"},
			//"Script school": {"classroom"},
			//"Bar":           {"lounge seating", "private booth"},
			//"cinema":        {"entrance", "lobby", "auditorium", "projection room"},
			//"club":          {"dance floor", "private booth", "terrace"},
			//"theater":       {"stage", "seating", "backstage"},
		},
	}

	// EventTopic 事件主题
	EventTopic = Slice{"", "Supernatural", "Love drama", "Crime committed", "Sex scandal", "Blackmail", "Flirty", "Underground deal", "unexpected revealing", "nonsense", "the mob connections", "illegal connection", "Comedic Scene", "Fraud", "Social media bully", "social media harassment", "Harassment", "romantic", "sex scene", "Action Scene", "Horror Scene", "unexpected accident", "Unwanted guests", "Money scam"}

	// ActionTags 行事标签
	ActionTags = Slice{"do evil", "do good", "do nonsense", "satirize", "ridicule", "confront", "manipulate", "take advantage", "attack", "defend", "give in", "calm", "agree", "support", "accept", "commit"}

	// ActionPromotionTags 推广事件的形式标签
	ActionPromotionTags = Slice{"be friendly", "be aggressive", "remain silent", "talk nonsense"}

	// Weathers 天气
	Weathers = Slice{"Sunny", "Cloudy", "Light rain", "Heavy rain", "Stormy rain", "Thunderstorm", "Foggy", "Bright sunshine", "Extremely hot", "Dust storm", "Extreme weather"}
)
