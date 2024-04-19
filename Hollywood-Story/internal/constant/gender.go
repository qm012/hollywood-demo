package constant

import "github.com/bwmarrin/discordgo"

type genders []*gender

type gender struct {
	ID    string
	Name  string
	Emoji string
}

var (
	Genders = genders{GenderMale, GenderFemale, GenderQueer}
)

var (
	GenderMale = &gender{
		ID:    "Male",
		Name:  "Male",
		Emoji: "♂️",
	}
	GenderFemale = &gender{
		ID:    "Female",
		Name:  "Female",
		Emoji: "♀️",
	}
	// GenderQueer “酷儿”是所有不符合主流性与性别规范的性少数群体所使用的身份、政治和学术用语。
	//		它既是身份标签（性别酷儿），也是一种政治策略（性别酷儿/酷儿身份），同时也是一种文化分析概念（酷儿理论）。
	//		该词来源于英文的“Queer”，原意是“怪诞、奇怪、（性）变态.
	GenderQueer = &gender{
		ID:    "Queer",
		Name:  "Queer",
		Emoji: "🏳️‍⚧️", // "🏳️‍⚧️" （"⚧️" is Invalid emoji）
	}
)

func (g genders) GetSelectMenuOptions() []discordgo.SelectMenuOption {
	selectMenuOptions := make([]discordgo.SelectMenuOption, 0, len(g))
	for _, value := range g {
		selectMenuOptions = append(selectMenuOptions, discordgo.SelectMenuOption{
			Label:       value.ID,
			Value:       value.ID,
			Description: value.Name,
			Emoji: &discordgo.ComponentEmoji{
				Name: value.Emoji,
			},
			Default: false,
		})
	}
	return selectMenuOptions
}
