package model

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
	"math/rand"
	"strconv"
	"time"
	"vland.live/app/global"
)

type occupations map[string][]string
type personalities []string

var (

	// Occupations 职业
	Occupations = occupations{
		"actor":        []string{"double", "cameo appearance", "supporting", "extra"},
		"screenwriter": []string{"screenwriter"},
		"director":     []string{"assistant director"},
		"grip":         []string{"gaffer", "cameraman", "make-up artist", "costume designer", "set dresser", "property master"},
		"consultant":   []string{"literary consultant", "historical consultant", "dance consultant"},
		"producer":     []string{"producer", "producer agent"},
	}
	OccupationsV1 = occupations{
		"演员": []string{"替身演员", "客串演员", "配角", "群众演员"},
		"编剧": []string{"编剧"},
		"导演": []string{"副导演"},
		"场务": []string{"灯光", "摄影", "化妆", "造型", "布景", "道具"},
		"顾问": []string{"文学顾问", "历史顾问", "舞蹈顾问"},
		"制片": []string{"制片人", "制片方代理"},
	}
	// identities 身份
	//identities = []string{"double", "cameo appearance", "supporting", "extra", "leading", "screenwriter", "director", "assistant director", "gaffer", "cameraman", "make-up artist", "costume designer", "set dresser",
	//	"property master", "literary consultant", "historical consultant", "dance consultant", "producer", "producer agent"}
	// actorNames 名字
	actorNames = []string{"John", "Charles", "Mike", "Mary", "Marilyn", "James", "Scarlett", "Dan", "Gloria", "Lee", "Michael", "Laura", "Jackson", "Kevin", "Elliot", "Emily", "Ashley", "Sasha", "Martin", "Martina"}
	// events 事件
	//events = []string{"闷闷不乐", "口角", "争执", "大打出手", "结仇", "聊天", "丢失", "结识", "遇到", "意外", "突然", "讨论", "发现", "思考", "研究", "决定", "希望/想要", "拒绝", "接受", "顿悟"}

	// Personalities 性格
	//Personalities = []string{"作恶", "喜欢批评他人", "喜欢风险行为", "行善", "保持中立", "热情", "固执", "喜欢艺术", "喜欢社交", "喜欢和别人争吵", "喜欢教导别人", "喜欢进行不切实际的想象", "喜欢同情他人", "对别人的隐私感兴趣", "不易惹怒别人", "喜欢讲笑话", "安静", "喜欢倾听别人", "喜欢调整自己的行为策略", "喜欢做事前进行思考",
	//	"奸诈", "愤世嫉俗", "冒险精神", "善良感性", "理性中立", "热情", "坚韧不拔", "艺术气质", "社交狂", "忠诚", "好为人师", "梦想家", "富有爱心", "寻求挑战", "圆滑", "幽默", "灵性主义", "倾听者", "灵活变通者", "深思熟虑"}
	Personalities = personalities{
		"do evil", "like to criticize others", "like risky behavior", "do good", "remain neutral", "be passionate", "be stubborn", "seek approval from others", "enjoy provoking others", "enjoy self-sacrifice",
		"enjoy imagining crises", "enjoy sympathizing with others", "enjoy prying into others' privacy", "look down on others", "enjoy surprising others", "enjoy telling lies", "enjoy listening to others",
		"enjoy speaking ill of others", "act independently; go one's own way.", "sex maniac"}
	// hobbies 爱好
	//hobbies = []string{"天文", "阅读", "徒步", "游泳", "绘画", "钓鱼", "瑜伽", "跑步", "烹饪", "摄影", "音乐", "篮球", "旅行", "舞蹈", "观影", "动物", "手工艺", "植物养护", "手机摄影", "编程", "电子游戏", "桌游"}
)

// DeleteByID 删除一个演员
func sliceDeleteByIndex(slice []string, index string) []string {

	var newSlice = make([]string, 0, len(slice))
	for _, value := range slice {
		if value != index {
			newSlice = append(newSlice, value)
		}
	}
	return newSlice
}

// GetActorsByNum 随机获取 num 个演员
func GetActorsByNum(num int) Actors {

	var (
		staffMinSize                = 13 // 班底最低人数
		staffActorMinSizeMultiplier = 2  // 班底演员的最低主演人数的倍数
		staffLandingSize            = 2  // 班底主演人数
	)

	if num < staffMinSize {
		num = staffMinSize
	}
	if num > len(actorNames) {
		global.Logger.Error("演员数量不足", zap.Any("actorNames", actorNames), zap.Any("len(actorNames)", len(actorNames)), zap.Any("num", num))
		num = len(actorNames)
	}

	//var (
	//	occupationActor        = "演员"
	//	occupationScreenwriter = "编剧"
	//	occupationDirector     = "导演"
	//	occupationGrip         = "场务"
	//	occupationConsultant   = "顾问"
	//	occupationProducer     = "制片"
	//)

	//var (
	//	identityMainDirector = "主导演"
	//	identityMainLeading  = "主要演员"
	//)

	var (
		occupationActor        = "actor"
		occupationScreenwriter = "screenwriter"
		occupationDirector     = "director"
		occupationGrip         = "grip"
		occupationConsultant   = "consultant"
		occupationProducer     = "producer"
	)

	var (
		identityMainDirector = "director"
		identityMainLeading  = "leading"
	)

	var (
		remainingActorNum = num
	)

	var (
		finalOccupations = Occupations
	)

	//- NPC职业目前有6类，分别为：制作人（不超过NPC数量的1/10）、导演（不超过NPC数量的1/5）、编剧（不超过NPC数量的1/5）、演员（最低为主演人数的二倍）、顾问（不超过NPC数量的1/10）、场务（工作组）
	var (
		r                   = rand.New(rand.NewSource(time.Now().UnixNano()))
		actors              = make(Actors, 0, num)
		copiedActorNames    = make([]string, len(actorNames))
		copiedPersonalities = make([]string, len(Personalities))
	)
	copy(copiedActorNames, actorNames)
	copy(copiedPersonalities, Personalities)

	// 主导演
	{
		// 演员名称：出现过的信息组合就固定下来。不再参与随机。
		name := copiedActorNames[r.Intn(len(copiedActorNames))]
		copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

		personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
		copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

		actors = append(actors, &Actor{
			ID:         strconv.Itoa(remainingActorNum),
			Name:       name,
			Occupation: occupationDirector,
			Identity:   identityMainDirector,
			//Event:       "",
			Personality: personality,
			//Hobbies:     nil,
		})
	}

	// 副导演
	{
		count := num/5 - 1 // 减去已经计算的主导演
		for i := 0; i < count; i++ {
			remainingActorNum -= 1
			name := copiedActorNames[r.Intn(len(copiedActorNames))]
			copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

			personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
			copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

			identities := finalOccupations[occupationDirector]

			actors = append(actors, &Actor{
				ID:         strconv.Itoa(remainingActorNum),
				Name:       name,
				Occupation: occupationDirector,
				Identity:   identities[r.Intn(len(identities))],
				//Event:       "",
				Personality: personality,
				//Hobbies:     nil,
			})
		}
	}

	// 编剧
	{
		count := num / 5
		for i := 0; i < count; i++ {
			remainingActorNum -= 1
			name := copiedActorNames[r.Intn(len(copiedActorNames))]
			copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

			personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
			copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

			identities := finalOccupations[occupationScreenwriter]

			actors = append(actors, &Actor{
				ID:         strconv.Itoa(remainingActorNum),
				Name:       name,
				Occupation: occupationScreenwriter,
				Identity:   identities[r.Intn(len(identities))],
				//Event:       "",
				Personality: personality,
				//Hobbies:     nil,
			})
		}
	}

	// 制作人
	{
		count := num / 10
		for i := 0; i < count; i++ {
			remainingActorNum -= 1
			name := copiedActorNames[r.Intn(len(copiedActorNames))]
			copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

			personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
			copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

			identities := finalOccupations[occupationProducer]

			actors = append(actors, &Actor{
				ID:         strconv.Itoa(remainingActorNum),
				Name:       name,
				Occupation: occupationProducer,
				Identity:   identities[r.Intn(len(identities))],
				//Event:       "",
				Personality: personality,
				//Hobbies:     nil,
			})
		}
	}

	// 顾问
	{
		count := num / 10
		for i := 0; i < count; i++ {
			remainingActorNum -= 1
			name := copiedActorNames[r.Intn(len(copiedActorNames))]
			copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

			personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
			copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

			identities := finalOccupations[occupationConsultant]

			actors = append(actors, &Actor{
				ID:         strconv.Itoa(remainingActorNum),
				Name:       name,
				Occupation: occupationConsultant,
				Identity:   identities[r.Intn(len(identities))],
				//Event:       "",
				Personality: personality,
				//Hobbies:     nil,
			})
		}
	}

	// 演员
	{
		// 固定两个主要演员
		for i := 0; i < staffLandingSize; i++ {
			remainingActorNum -= 1
			name := copiedActorNames[r.Intn(len(copiedActorNames))]
			copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

			personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
			copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

			actors = append(actors, &Actor{
				ID:         strconv.Itoa(remainingActorNum),
				Name:       name,
				Occupation: occupationActor,
				Identity:   identityMainLeading,
				//Event:       "",
				Personality: personality,
				//Hobbies:     nil,
			})
		}
		// 随机的其他演员
		count := staffLandingSize * staffActorMinSizeMultiplier
		for i := 0; i < count; i++ {
			remainingActorNum -= 1
			name := copiedActorNames[r.Intn(len(copiedActorNames))]
			copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

			personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
			copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

			identities := finalOccupations[occupationActor]

			actors = append(actors, &Actor{
				ID:         strconv.Itoa(remainingActorNum),
				Name:       name,
				Occupation: occupationActor,
				Identity:   identities[r.Intn(len(identities))],
				//Event:       "",
				Personality: personality,
				//Hobbies:     nil,
			})
		}
	}
	// 剩余的场务
	{
		for i := 0; i < remainingActorNum; i++ {
			remainingActorNum -= 1
			name := copiedActorNames[r.Intn(len(copiedActorNames))]
			copiedActorNames = sliceDeleteByIndex(copiedActorNames, name)

			personality := copiedPersonalities[r.Intn(len(copiedPersonalities))]
			copiedPersonalities = sliceDeleteByIndex(copiedPersonalities, personality)

			identities := finalOccupations[occupationGrip]

			actors = append(actors, &Actor{
				ID:         strconv.Itoa(remainingActorNum),
				Name:       name,
				Occupation: occupationGrip,
				Identity:   identities[r.Intn(len(identities))],
				//Event:       "",
				Personality: personality,
				//Hobbies:     nil,
			})
		}
	}
	return actors
}

// GetSelectMenOptions 获取性格的下拉选择
func (p personalities) GetSelectMenOptions() []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, 0, 20)
	for _, personality := range p {
		options = append(options, discordgo.SelectMenuOption{
			Label:       personality,
			Value:       personality,
			Description: "",
			Default:     false,
		})
	}

	return options
}

// GetSelectMenOptions 获取职业的下拉选择
func (o occupations) GetSelectMenOptions() []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, 0, 20)
	for key := range o {
		options = append(options, discordgo.SelectMenuOption{
			Label:       key,
			Value:       key,
			Description: "",
			Default:     false,
		})
	}

	return options
}
