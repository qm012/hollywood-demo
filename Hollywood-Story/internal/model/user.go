package model

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"math/rand"
	"strings"
	"time"
	"vland.live/app/global"
	"vland.live/app/internal/constant"
)

// UserFilm 影片信息
type UserFilm struct {
	Screenplay *constant.Screenplay // 剧本信息
	Actors     Actors               // 演员列表
	ActorMap   map[string]Actors    // 每次用户玩的时候的演员对应的还能聊的人 key：演员 value：这个演员的可对话演员列表（根据需要，每轮对话的人都会从对方的列表中删除）
}

// Init 初始化每个演员的对话者列表
func (u *UserFilm) Init() {
	var actorMap = make(map[string]Actors, 30)
	for _, actor := range u.Actors {
		actorMap[actor.ID] = u.Actors.Copy().DeleteByID(actor.ID)
	}
	u.ActorMap = actorMap
}

// DeleteMeetingActors 删除相遇的对方
//
//	基于每个回合的演员列表，需要删除对方，避免在以后的回合中之前回合已经对话的演员还可以相互遇见
func (u *UserFilm) DeleteMeetingActors(roundActors Actors) {
	for _, roundActorA := range roundActors {
		speakers, ok := u.ActorMap[roundActorA.ID]
		if !ok {
			continue
		}
		for _, roundActorB := range roundActors {
			if roundActorA.ID == roundActorB.ID {
				continue
			}
			u.ActorMap[roundActorA.ID] = speakers.DeleteByID(roundActorB.ID)
		}
	}
}

// GetNonMeet2Actors 获取两个不相遇的演员
func (u *UserFilm) GetNonMeet2Actors() Actors {
	var (
		nonMeetActors = make(Actors, 0, 2)
		randActors    = u.Actors.Copy() // 使用copy后的演员列表来随机
	)

	// 打乱顺序
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(randActors), func(i, j int) {
		randActors[i], randActors[j] = randActors[j], randActors[i]
	})

	// 随机一个
	for _, randActor := range randActors {
		// 使用随机的randActor到map中找到可以对话的人
		speakers, ok := u.ActorMap[randActor.ID]
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
				zap.Any("u.Film.ActorMap", u.ActorMap),
				zap.Any("可能原因", "某些演员的对话者列表已经为空(从u.Film.ActorMap查找定位数据)"))
		}
	}

	return nonMeetActors
}

// AppointActors 获取一个随机演员的对话者列表
func (u *UserFilm) AppointActors(num int) Actors {
	// 随机一个
	randomActor := u.Actors.GetSingleByRandom()
	if randomActor == nil {
		global.Logger.Error("获取不到随机的演员")
		return nil
	}
	actorID := randomActor.ID
	return u.AppointActorsByID(actorID, num)
}

// AppointActorsByID 随机获取指定演员的 num 个对话演员列表
//
//	如果 len(actorIDs) == 0 则随机一个演员
//	如果 len(actorIDs) > 0  则使用 actorIDs[0]
func (u *UserFilm) AppointActorsByID(actorID string, num int) Actors {

	speakers, ok := u.ActorMap[actorID]
	if !ok {
		return Actors{}
	}

	if len(speakers) == 0 {
		return Actors{}
	}

	return speakers.GetMultiByRandom(num)
}

type User struct {
	ID                             string
	Nickname                       string                              // 用户输入的昵称
	NPCName                        string                              // NPC name
	Age                            int                                 // 年龄
	Gender                         string                              // 性别
	Personality                    []string                            // 性格
	Occupation                     string                              // 职业
	Attribute                      *constant.Attribute                 // 用户属性值
	Username                       string                              // discord的名称
	CurrentRound                   int                                 // 第几回合
	CurrentRoundEvent              *constant.PlayerOptionRoundEvent    // 当前回合的事件类型，TODO 后续处理，要基于阶段来确认事件类型
	CurrentRoundActors             Actors                              // 当前回合的演员
	CurrentRoundGptOptions         *constant.GptOptions                // 当前回合的gpt返回选项
	CurrentRoundDifficultyValue    int                                 // 当前回合的事件类型难度值
	CurrentRoundOptionOutcome      *constant.OptionOutcome             // 当前回合的事件类型的选项结果
	CurrentRoundAttributeValues    constant.AttributeValues            // 当前回合的4个随机属性
	CurrentRoundButtonAttributeMap map[string]*constant.AttributeValue // 当前回合button对应的属性
	Film                           *UserFilm                           // 影片信息
	CreatedAt                      int64                               // 创建时间，单位：秒
}

// NewUser 创建一个用户
func NewUser(ID, username string) *User {
	user := &User{
		ID:        ID,
		Username:  username,
		CreatedAt: time.Now().Unix(),
	}
	return user
}

func (u *User) SetDifficultyValue() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	// 生成10到15之间的随机整数
	value := r.Intn(6) + 10
	u.CurrentRoundDifficultyValue = value
	return value
}

func (u *User) RoundText() string {

	return fmt.Sprintf("The movie 《%s》is filming on the %dth day.", u.Film.Screenplay.Name, u.CurrentRound)
	//return "Round " + strconv.Itoa(u.CurrentRound) + "\n" + u.GetPhaseName()
	//return "Round " + constant.ChineseNumMap[u.CurrentRound]
}

func (u *User) SetCurrentRoundGptOptions(gptOptions *constant.GptOptions) {
	u.CurrentRoundGptOptions = gptOptions
}

func (u *User) GetCurrentRoundGptOptions(ID string) string {
	if ID == "" {
		return ""
	}
	if u.CurrentRoundGptOptions == nil {
		return ""
	}
	switch ID {
	case constant.PlayerCaptureOptionGoodness.ID:
		return u.CurrentRoundGptOptions.A
	case constant.PlayerCaptureOptionEvil.ID:
		return u.CurrentRoundGptOptions.B
	case constant.PlayerCaptureOptionNonIntervene.ID:
		return u.CurrentRoundGptOptions.C
	case constant.PlayerCaptureOptionPrank.ID:
		return u.CurrentRoundGptOptions.D
	}
	return ""
}

// GetPhaseName 获取阶段选项名称
func (u *User) GetPhaseName() string {
	var name = "Capture phase"
	if u.CurrentRound > constant.GlobalCaptureRoundNum {
		name = "Promotion phase"
	}
	return name
}

func (u *User) Format() string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("name:%s, age:%d, gender:%s, occupation:%s, personality:%s",
		u.Nickname, u.Age, u.Gender, u.Occupation, strings.Join(u.Personality, ",")))
	return builder.String()
}

func (u *User) WelcomeMessage() string {
	return fmt.Sprintf(u.Film.Screenplay.WelcomeMessage, u.NPCName, u.Nickname)
}

func (u *User) VariablesChoices(location, weather, eventTheme string) (*[]discordgo.MessageComponent, map[string]string) {

	var specialStr = ""
	var roundActors Actors
	switch u.CurrentRoundEvent.ID {
	case constant.PlayerCaptureOptionRoundEventTwoActorsInteraction.ID:
		roundActors = u.Film.GetNonMeet2Actors()
	case constant.PlayerCaptureOptionRoundEventPlayerActorInteraction.ID:
		roundActors = Actors{u.Film.Actors.GetSingleByRandom()}
	case constant.PlayerCaptureOptionRoundEventSpecialActorEvent.ID:
		specialNpc1 := SpecialNpcs[0]
		roundActors = Actors{
			{
				ID: "specialNpc1",
				//Name:       u.NPCName,
				Name:       specialNpc1.Name,
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
				Name:       u.NPCName,
				Occupation: specialNpc1.Occupation,
				Identity:   "",
				//Event:       "",
				Personality: strings.Join(specialNpc1.Personality, ","),
				Hobbies:     specialNpc1.Hobbies,
			},
			u.Film.Actors.GetSingleByRandom(),
		}
	}

	if len(roundActors) == 0 {
		global.Logger.Error("没有找到对应的演员，数据错误")
		return nil, nil
	}

	u.CurrentRoundActors = roundActors

	// 用户每一轮的随机属性
	u.CurrentRoundAttributeValues = u.Attribute.AttributeValueByRandom(4)

	var (
		randomActionTags              = constant.ActionTags.RandomNoRepeat(4)
		components, attributeValueMap = constant.PlayerCaptureOptions.SelectMessageComponent(u.CurrentRoundAttributeValues)
	)
	// 给当前回合的button属性值存起来，用户点击时出结果用
	u.CurrentRoundButtonAttributeMap = attributeValueMap

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
		"{player_infos}": u.Format(),
		"{special_info}": specialStr,
		"{location}":     location,
		"{weather}":      weather,
		"{event_theme}":  eventTheme,
		"{movie_script}": u.Film.Screenplay.Format(),
	}
	// 特殊npc不需要eventTheme
	if u.CurrentRoundEvent.ID == constant.PlayerCaptureOptionRoundEventSpecialActorEvent.ID {
		delete(variables, "{event_theme}")
	}
	return components, variables
}

func (u *User) VariablesOutcome(optionID string) map[string]string {

	u.CurrentRoundOptionOutcome = u.CurrentRoundEvent.Probability(u.CurrentRoundDifficultyValue,
		u.CurrentRoundButtonAttributeMap[optionID].Value,
		u.NPCName)

	variables := map[string]string{
		"{npc_infos}":        u.CurrentRoundActors.Format(),
		"{player_infos}":     u.Format(),
		"{event_desc}":       u.CurrentRoundGptOptions.Event,
		"{conversation}":     u.CurrentRoundGptOptions.FormatDialogue(strings.Join(u.CurrentRoundActors.Names(), "\n")),
		"{players_reaction}": u.GetCurrentRoundGptOptions(optionID),
		"{film_attr}":        u.CurrentRoundOptionOutcome.Text,
	}

	return variables
}

func (u *User) VariablesOutcomeNews(gptOutcome *constant.GptOutcome) map[string]string {
	variables := map[string]string{
		"{npc_infos}":    u.CurrentRoundActors.Format(),
		"{player_infos}": u.Format(),
		"{movie_script}": u.Film.Screenplay.FormatSimple(),
		//"{event_desc}":       u.CurrentRoundGptOptions.Event,// 使用事件版
		"{event_desc}":       u.CurrentRoundGptOptions.FormatDialogue(u.CurrentRoundActors.Names()...), // 使用对话版
		"{players_reaction}": gptOutcome.Reaction,
		"{event_outcome}":    gptOutcome.Outcome,
		"{post_headline}":    gptOutcome.Headline,
	}
	return variables
}
