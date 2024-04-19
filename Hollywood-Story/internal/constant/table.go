package constant

import (
	"fmt"
	"math/rand"
	"time"
)

type Attribute struct {
	Text   string         `json:"text" bson:"text"`     // 格式化的属性文本
	Values map[string]int `json:"values" bson:"values"` // 属性对应的随机数
}

type AttributeValues []*AttributeValue

func (a AttributeValues) Keys() Slice {
	var keys = make(Slice, 0, len(a))
	for _, value := range a {
		keys = append(keys, value.Key)
	}
	return keys
}

type AttributeValue struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

func (a *AttributeValue) Format() string {
	return fmt.Sprintf("[%s - %v]", a.Key, a.Value)
}

func (a *AttributeValue) FormatSimple() string {
	return fmt.Sprintf("%s", a.Key)
}

func (a *AttributeValue) Copy() *AttributeValue {
	return &AttributeValue{
		Key:   a.Key,
		Value: a.Value,
	}
}

// 属性
const (
	attributeLabelPlaywriting   = "play_writing"  // 剧作
	attributeLabelActing        = "acting"        // 演技
	attributeLabelCollaboration = "collaboration" // 协同
	attributeLabelLikability    = "likability"    // 名气/亲和力
	attributeLabelInterest      = "interest"      // 兴趣
)

// 子属性
const (
	attributeLabelChildRomanceWriting = "romance_writing" // 爱情
	attributeLabelChildSciFiWriting   = "sciFi_writing"   // 科幻
	attributeLabelChildActionWriting  = "action_writing"  // 动作
	attributeLabelChildComedyWriting  = "comedy_writing"  // 喜剧
	attributeLabelChildMysteryWriting = "mystery_writing" // 悬疑

	attributeLabelChildFacialExpression = "facial_expression" // 表情
	attributeLabelChildDialogue         = "dialogue"          // 台词
	attributeLabelChildSinging          = "singing"           // 声乐
	attributeLabelChildDancing          = "dancing"           // 舞蹈
	attributeLabelChildBodyExpression   = "body_expression"   // 形体

	attributeLabelChildCollaboration = "collaboration" // 协同

	attributeLabelChildLikability = "likability" // 名气/亲和力

	attributeLabelChildPhilosophy     = "philosophy"      // 哲学
	attributeLabelChildFortuneTelling = "fortune_telling" // 算命
	attributeLabelChildAesthetic      = "aesthetic"       // 美学
)

var (
	elementMap = map[string]map[string]struct{}{
		attributeLabelPlaywriting: {
			attributeLabelChildRomanceWriting: {},
			attributeLabelChildSciFiWriting:   {},
			attributeLabelChildActionWriting:  {},
			attributeLabelChildComedyWriting:  {},
			attributeLabelChildMysteryWriting: {},
		},
		attributeLabelActing: {
			attributeLabelChildFacialExpression: {},
			attributeLabelChildDialogue:         {},
			attributeLabelChildSinging:          {},
			attributeLabelChildDancing:          {},
			attributeLabelChildBodyExpression:   {},
		},
		attributeLabelCollaboration: {
			attributeLabelChildCollaboration: {},
		},
		attributeLabelLikability: {
			attributeLabelChildLikability: {},
		},
		attributeLabelInterest: {
			attributeLabelChildPhilosophy:     {},
			attributeLabelChildFortuneTelling: {},
			attributeLabelChildAesthetic:      {},
		},
	}
)

func getValuesByElementMap() map[string]int {
	var (
		r        = rand.New(rand.NewSource(time.Now().UnixNano()))
		maxValue = 10
		values   = make(map[string]int, 20)
	)

	for _, v := range elementMap {
		for k1 := range v {
			values[k1] = r.Intn(maxValue) + 1
		}
	}
	return values
}

func randomByKey(m map[string]struct{}) string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	index := rand.Intn(len(keys)) // 生成0到长度之间的随机数作为索引
	return keys[index]
}

func (a *Attribute) AttributeValueByRandom(num int) AttributeValues {

	if num <= 1 {
		num = 1
	}

	if num >= len(a.Values) {
		num = len(a.Values)
	}
	// 以下是随机4个属性
	//copiedValue := make(map[string]int, len(a.Values))
	//for key, value := range a.Values {
	//	copiedValue[key] = value
	//}
	//
	//results := make([]*AttributeValue, 0, num)
	//
	//for num > 0 {
	//	keys := make([]string, 0, len(copiedValue))
	//	for key := range copiedValue {
	//		keys = append(keys, key)
	//	}
	//	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//	randomKey := keys[r.Intn(len(copiedValue))]
	//
	//	results = append(results, &AttributeValue{
	//		Key:   randomKey,
	//		Value: copiedValue[randomKey],
	//	})
	//	delete(copiedValue, randomKey)
	//	num--
	//}
	//
	//return results

	// 下面则根据规则处理
	// - 随机规则：每四个选项（即一回合事件的选项）中，必有且仅有两个选项的【属性标签】是【剧作】和【演技】的【子属性】。另外两个选项的【属性标签】则从除【剧作】和【演技】外的【属性】中抽取。

	results := make([]*AttributeValue, 0, num)
	// 剧作
	{
		key := randomByKey(elementMap[attributeLabelPlaywriting])
		results = append(results, &AttributeValue{
			Key:   key,
			Value: a.Values[key],
		})
	}
	// 演技
	{
		key := randomByKey(elementMap[attributeLabelActing])
		results = append(results, &AttributeValue{
			Key:   key,
			Value: a.Values[key],
		})
	}

	var (
		tempMap = map[string]struct{}{
			attributeLabelCollaboration: {},
			attributeLabelLikability:    {},
			attributeLabelInterest:      {},
		}
		r = rand.New(rand.NewSource(time.Now().UnixNano()))
	)

	for i := 0; i < 2; i++ {
		keys := make([]string, 0, len(tempMap))
		for key := range tempMap {
			keys = append(keys, key)
		}
		randomKey := keys[r.Intn(len(tempMap))]
		key := randomByKey(elementMap[randomKey])
		results = append(results, &AttributeValue{
			Key:   key,
			Value: a.Values[key],
		})
		delete(tempMap, randomKey)
	}

	return results
}

var (
	attributes = `┌───────────────┬─────────────────────┬───────┐
│              Attributes             │ Value │
├───────────────┼─────────────────────┼───────┤
│               │   Romance writing   │  %3d  │
├               ┼─────────────────────┼───────┤
│               │    Sci-Fi writing   │  %3d  │
├               ┼─────────────────────┼───────┤
│ Play Writing  │    Action writing   │  %3d  │
├               ┼─────────────────────┼───────┤
│               │    Comedy writing   │  %3d  │
├               ┼─────────────────────┼───────┤
│               │   Mystery writing   │  %3d  │
├───────────────┼─────────────────────┼───────┤
│               │  Facial Expression  │  %3d  │
├               ┼─────────────────────┼───────┤
│               │       Dialogue      │  %3d  │
├               ┼─────────────────────┼───────┤
│    Acting     ┼       Singing       │  %3d  │
├               ┼─────────────────────┼───────┤
│               ┼ 		Dancing       │  %3d  │
├               ┼─────────────────────┼───────┤
│               ┼   Body Expression   │  %3d  │
├───────────────┼─────────────────────┼───────┤
│             Collaboration           │  %3d  │
├───────────────┼─────────────────────┼───────┤
│              Likability             │  %3d  │
├───────────────┼─────────────────────┼───────┤
│               │      philosophy     │  %3d  │
├               ┼─────────────────────┼───────┤
│   Interest    ┼   Fortune Telling   │  %3d  │
├               ┼─────────────────────┼───────┤
│               ┼ 	    aesthetic     │  %3d  │
└───────────────┴─────────────────────┴───────┘`
)

func TableRandomAttributes() *Attribute {

	values := getValuesByElementMap()

	text := fmt.Sprintf(attributes,
		values[attributeLabelChildRomanceWriting],
		values[attributeLabelChildSciFiWriting],
		values[attributeLabelChildActionWriting],
		values[attributeLabelChildComedyWriting],
		values[attributeLabelChildMysteryWriting],
		values[attributeLabelChildFacialExpression],
		values[attributeLabelChildDialogue],
		values[attributeLabelChildSinging],
		values[attributeLabelChildDancing],
		values[attributeLabelChildBodyExpression],
		values[attributeLabelChildCollaboration],
		values[attributeLabelChildLikability],
		values[attributeLabelChildPhilosophy],
		values[attributeLabelChildFortuneTelling],
		values[attributeLabelChildAesthetic],
	)

	text = fmt.Sprintf("```%s```", text)

	return &Attribute{
		Text:   text,
		Values: values,
	}
}
