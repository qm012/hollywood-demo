package model

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
	"vland.live/app/global"
	"vland.live/app/internal/constant"
)

type MemberCurrentRound struct {
	Number             int                                 `json:"number" bson:"number"`                             // 第几回合
	Event              *constant.PlayerOptionRoundEvent    `json:"event_v1" bson:"event"`                            // 当前回合的事件类型，TODO 后续处理，要基于阶段来确认事件类型
	Actors             Actors                              `json:"actors" bson:"actors"`                             // 当前回合的演员
	GptOptions         *constant.GptOptionsV1              `json:"gpt_options" bson:"gpt_options"`                   // 当前回合的gpt返回选项
	GptOutcome         *constant.GptOutcomeV1              `json:"gpt_outcome" bson:"gpt_outcome"`                   // 当前回合的gpt返回得结果
	GptNews            *constant.GptNewsV1                 `json:"gpt_news" bson:"gpt_news"`                         // 当前回合的gpt返回结果的新闻评论
	DifficultyValue    int                                 `json:"difficulty_value" bson:"difficulty_value"`         // 当前回合的事件类型难度值
	OptionOutcome      *constant.OptionOutcome             `json:"option_outcome" bson:"option_outcome"`             // 当前回合的事件类型的选项结果
	AttributeValues    constant.AttributeValues            `json:"attribute_values" bson:"attribute_values"`         // 当前回合的4个随机属性
	ButtonAttributeMap map[string]*constant.AttributeValue `json:"button_attribute_map" bson:"button_attribute_map"` // 当前回合button对应的属性
}

// MemberFilm 影片信息
type MemberFilm struct {
	Screenplay *Screenplay       `json:"screenplay" bson:"screenplay"` // 剧本信息
	Actors     Actors            `json:"actors" bson:"actors"`         // 演员列表
	ActorMap   map[string]Actors `json:"actor_map" bson:"actor_map"`   // 每次用户玩的时候的演员对应的还能聊的人 key：演员 value：这个演员的可对话演员列表（根据需要，每轮对话的人都会从对方的列表中删除）
}

// Init 初始化每个演员的对话者列表
func (m *MemberFilm) Init() {
	var actorMap = make(map[string]Actors, 30)
	for _, actor := range m.Actors {
		actorMap[actor.ID] = m.Actors.Copy().DeleteByID(actor.ID)
	}
	m.ActorMap = actorMap
}

// DeleteMeetingActors 删除相遇的对方
//
//	基于每个回合的演员列表，需要删除对方，避免在以后的回合中之前回合已经对话的演员还可以相互遇见
func (m *MemberFilm) DeleteMeetingActors(roundActors Actors) {
	for _, roundActorA := range roundActors {
		speakers, ok := m.ActorMap[roundActorA.ID]
		if !ok {
			continue
		}
		for _, roundActorB := range roundActors {
			if roundActorA.ID == roundActorB.ID {
				continue
			}
			m.ActorMap[roundActorA.ID] = speakers.DeleteByID(roundActorB.ID)
		}
	}
}

// GetNonMeet2Actors 获取两个不相遇的演员
func (m *MemberFilm) GetNonMeet2Actors() Actors {
	var (
		nonMeetActors = make(Actors, 0, 2)
		randActors    = m.Actors.Copy() // 使用copy后的演员列表来随机
	)

	// 打乱顺序
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(randActors), func(i, j int) {
		randActors[i], randActors[j] = randActors[j], randActors[i]
	})

	// 随机一个
	for _, randActor := range randActors {
		// 使用随机的randActor到map中找到可以对话的人
		speakers, ok := m.ActorMap[randActor.ID]
		if !ok {
			continue
		}
		if len(speakers) == 0 {
			continue
		}
		speaker := speakers.GetSingleByRandom()
		if speaker == nil {
			continue
		}
		nonMeetActors = append(nonMeetActors, randActor.Copy(), speaker.Copy())
		// 找到后则退出
		break
	}

	if len(nonMeetActors) < 2 {
		if global.Logger != nil {
			global.Logger.Error("获取两个不相遇的演员数量不足",
				zap.Any("len(randActors)", len(randActors)),
				zap.Any("nonMeetActors", nonMeetActors),
				zap.Any("u.Film.ActorMap", m.ActorMap),
				zap.Any("可能原因", "某些演员的对话者列表已经为空(从u.Film.ActorMap查找定位数据)"))
		}
	}

	return nonMeetActors
}

type Member struct {
	ID             primitive.ObjectID  `json:"id" bson:"_id"`                            // 主键
	DeviceID       string              `json:"device_id" bson:"device_id"`               // 暂时以设备ID用作用户唯一ID
	Nickname       string              `json:"nickname" bson:"nickname"`                 // 用户输入的昵称
	SpecialNPCName string              `json:"special_npc_name" bson:"special_npc_name"` // 特殊NPC name
	Age            int                 `json:"age" bson:"age"`                           // 年龄
	Gender         string              `json:"gender" bson:"gender"`                     // 性别
	Personality    []string            `json:"personality" bson:"personality"`           // 性格
	Occupation     string              `json:"occupation" bson:"occupation"`             // 职业
	Attribute      *constant.Attribute `json:"attribute" bson:"attribute"`               // 用户属性值
	CurrentRound   *MemberCurrentRound `json:"current_round" bson:"current_round"`       // 当前回合的数据
	Film           *MemberFilm         `json:"film" bson:"film"`                         // 影片信息
	ModifiedAt     int64               `json:"modified_at" bson:"modified_at"`           // 更新时间
	CreatedAt      int64               `json:"created_at" bson:"created_at"`             // 创建时间
}

func (m *Member) Format() string {
	var builder strings.Builder
	//builder.WriteString(fmt.Sprintf("name:%s, age:%d, gender:%s, occupation:%s, personality:%s",
	//	m.Nickname, m.Age, m.Gender, m.Occupation, strings.Join(m.Personality, ",")))
	builder.WriteString(fmt.Sprintf("name:%s, age:%d, gender:%s",
		m.Nickname, m.Age, m.Gender))
	return builder.String()
}

func (m *MemberCurrentRound) SetDifficultyValue() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成10到15之间的随机整数
	value := r.Intn(6) + 10
	m.DifficultyValue = value
	return value
}

func (m *Member) WelcomeMessage() string {
	return fmt.Sprintf(m.Film.Screenplay.WelcomeMessage, m.SpecialNPCName, m.Nickname)
}

func (m *Member) VariablesChoices(location, weather, eventTheme string) map[string]string {

	var specialStr = ""
	var roundActors Actors
	switch m.CurrentRound.Event.ID {
	case constant.PlayerCaptureOptionRoundEventTwoActorsInteraction.ID:
		roundActors = m.Film.GetNonMeet2Actors()
	case constant.PlayerCaptureOptionRoundEventPlayerActorInteraction.ID:
		roundActors = Actors{m.Film.Actors.GetSingleByRandom()}
	case constant.PlayerCaptureOptionRoundEventSpecialActorEvent.ID:
		specialNpc1 := SpecialNpcs[0]
		roundActors = Actors{
			{
				ID: "specialNpc1",
				//Name:       u.NPCName,
				Name:       m.SpecialNPCName,
				Occupation: specialNpc1.Occupation,
				Identity:   "",
				//Event:       "",
				Personality: strings.Join(specialNpc1.Personality, ","),
				Hobbies:     specialNpc1.Hobbies,
			},
		}
		specialStr = `
- event type:
    the NPC drops by the film shooting scene to visit the player. Note that the NPC is not in this film crew.
- event style:
    romantic`
	case constant.PlayerCaptureOptionRoundEventSpecialNpcAndNpcAndPlayerEvent.ID:
		specialNpc1 := SpecialNpcs[0]
		roundActors = Actors{
			{
				ID:         "specialNpc1",
				Name:       m.SpecialNPCName,
				Occupation: specialNpc1.Occupation,
				Identity:   "",
				//Event:       "",
				Personality: strings.Join(specialNpc1.Personality, ","),
				Hobbies:     specialNpc1.Hobbies,
			},
			m.Film.Actors.GetSingleByRandom(),
		}
	case constant.PlayerPromotionOptionRoundEventMeetAndGreet.ID:
		roundActors = m.Film.GetNonMeet2Actors()
	}

	if len(roundActors) == 0 {
		global.Logger.Error("没有找到对应的演员，数据错误")
		return nil
	}

	m.CurrentRound.Actors = roundActors
	// 用户每一轮的随机属性
	m.CurrentRound.AttributeValues = m.Attribute.AttributeValueByRandom(4)

	var (
		randomActionTags     = constant.ActionTags.RandomNoRepeat(4)
		_, attributeValueMap = constant.PlayerCaptureOptions.SelectMessageComponent(m.CurrentRound.AttributeValues)
	)

	if m.CurrentRound.Event.ID == constant.PlayerPromotionOptionRoundEventMeetAndGreet.ID {
		randomActionTags = constant.ActionPromotionTags.RandomNoRepeat(4)
	}

	// 给当前回合的button属性值存起来，用户点击时出结果用
	m.CurrentRound.ButtonAttributeMap = attributeValueMap

	a, _ := attributeValueMap[constant.PlayerCaptureOptionGoodness.ID]
	b, _ := attributeValueMap[constant.PlayerCaptureOptionEvil.ID]
	c, _ := attributeValueMap[constant.PlayerCaptureOptionNonIntervene.ID]
	d, _ := attributeValueMap[constant.PlayerCaptureOptionPrank.ID]

	variables := map[string]string{

		"{action_tag_a}": randomActionTags[0],
		"{attribute_a}":  a.FormatSimple(),

		"{action_tag_b}": randomActionTags[1],
		"{attribute_b}":  b.FormatSimple(),

		"{action_tag_c}": randomActionTags[2],
		"{attribute_c}":  c.FormatSimple(),

		"{action_tag_d}": randomActionTags[3],
		"{attribute_d}":  d.FormatSimple(),

		"{npc_infos}":    roundActors.Format(),
		"{player_infos}": m.Format(),
		"{special_info}": specialStr,
		"{location}":     location,
		"{weather}":      weather,
		"{event_theme}":  eventTheme,
		"{movie_script}": m.Film.Screenplay.Format(),
	}

	switch m.CurrentRound.Event.ID {
	case constant.PlayerCaptureOptionRoundEventSpecialActorEvent.ID:
		// 特殊npc不需要eventTheme
		delete(variables, "{event_theme}")
	}
	return variables
}

func (m *Member) VariablesOutcome(optionID string) map[string]string {

	m.CurrentRound.OptionOutcome = m.CurrentRound.Event.Probability(m.CurrentRound.DifficultyValue,
		m.CurrentRound.ButtonAttributeMap[optionID].Value,
		m.SpecialNPCName)

	variables := map[string]string{
		"{npc_infos}":        m.CurrentRound.Actors.Format(),
		"{player_infos}":     m.Format(),
		"{event_desc}":       m.CurrentRound.GptOptions.Event,
		"{conversation}":     m.CurrentRound.GptOptions.FormatDialogue(strings.Join(m.CurrentRound.Actors.Names(), "\n")),
		"{players_reaction}": m.CurrentRound.GetGptOptionsReaction(optionID),
		"{film_attr}":        m.CurrentRound.OptionOutcome.Text,
	}

	return variables
}

func (m *Member) VariablesOutcomeNews(gptOutcome *constant.GptOutcomeV1) map[string]string {
	variables := map[string]string{
		"{npc_infos}":    m.CurrentRound.Actors.Format(),
		"{player_infos}": m.Format(),
		"{movie_script}": m.Film.Screenplay.FormatSimple(),
		//"{event_desc}":       u.CurrentRoundGptOptions.Event,// 使用事件版
		"{event_desc}":       m.CurrentRound.GptOptions.FormatDialogue(m.CurrentRound.Actors.Names()...), // 使用对话版
		"{players_reaction}": gptOutcome.Reaction,
		"{event_outcome}":    gptOutcome.Outcome,
		"{post_headline}":    gptOutcome.Headline,
	}
	return variables
}

func (m *MemberCurrentRound) GetGptOptionsReaction(ID string) string {
	if ID == "" {
		return ""
	}
	if m.GptOptions == nil {
		return ""
	}
	switch ID {
	case constant.PlayerCaptureOptionGoodness.ID:
		return m.GptOptions.A
	case constant.PlayerCaptureOptionEvil.ID:
		return m.GptOptions.B
	case constant.PlayerCaptureOptionNonIntervene.ID:
		return m.GptOptions.C
	case constant.PlayerCaptureOptionPrank.ID:
		return m.GptOptions.D
	}
	return ""
}
