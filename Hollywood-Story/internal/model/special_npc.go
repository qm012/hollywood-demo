package model

type specialNpcs []*specialNpc

type specialNpc struct {
	ID              string            // 主键
	Name            string            // 名字
	Gender          string            // 性别
	Age             int               // 年龄
	Occupation      string            // 职业
	Personality     []string          // 性格
	Hobbies         []string          // 爱好
	Birth           string            // 出身
	Dream           string            // 梦想
	Relationship    map[string]string // 和其他NPC的关系
	ReputationValue string            // 名望值
	Memories        []string          // 记忆流
}

var (
	SpecialNpcs = specialNpcs{
		{
			ID:          "specialNpc1",
			Name:        "Noah Lee",
			Gender:      "Male",
			Age:         32,
			Occupation:  "actor",
			Personality: []string{"passionate about acting", "dreamer", "romantic", "witty and clever"},
			Hobbies:     []string{"cooking", "jazz"},
			Birth:       "a middle-class family in Michigan",
			Dream:       "find a partner to share his Hollywood dream",
			Relationship: map[string]string{
				"Mia Jefferson": "Ex-girlfriend",
				"Player":        "has feelings",
			},
			ReputationValue: "moderately well-known actor",
			Memories: []string{
				"His ex-girlfriend spilled coffee in the cafe, damaging the protagonist's phone.",
				"To compensate, he gifted the female lead a new phone.",
				"During their conversation, he discovered that the female lead aspired to be an actress, so he introduced her to audition for a movie.",
				"The female lead passed the audition, and he was happy for her, but this incident sparked jealousy in his ex-girlfriend, Mia Jefferson.",
			},
		},
		{
			ID:          "specialNpc2",
			Name:        "Mia Jefferson",
			Gender:      "Female",
			Age:         28,
			Occupation:  "actor",
			Personality: []string{"envious", "realistic", "mean", "glass heart"},
			Hobbies:     []string{"skating", "shopping", "surfing"},
			Birth:       "a celebrity family",
			Dream:       "demonstrate her competence to be a Hollywood star without the support from her family",
			Relationship: map[string]string{
				"Noah Lee": "Ex-boyfriend",
				"Player":   "jealous mockery",
			},
			ReputationValue: "a famous actress",
			Memories: []string{
				"When Noah Lee broke up with her, she spilled coffee, damaging the player's phone.",
				"She left without apologizing, but later heard that Noah Lee introduced the player to filmmaking.",
				"She was extremely jealous that the player gained Noah Lee's attention and hoped for the player's failure.",
			},
		},
	}
)
