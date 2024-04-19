package constant

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

type screenplays []*Screenplay

type Screenplay struct {
	ID             string
	Name           string
	Labels         []string
	Synopsis       string
	WelcomeMessage string
}

// Format 格式化剧本的信息
func (s *Screenplay) Format() string {
	format := fmt.Sprintf(`Title：%s
	Label: %s
	Synopsis:%s`, s.Name, strings.Join(s.Labels, ","), s.Synopsis)
	return format
}

// FormatSimple 格式化剧本的信息:简易版
func (s *Screenplay) FormatSimple() string {
	format := fmt.Sprintf(`Title：%s
	Label: %s`, s.Name, strings.Join(s.Labels, ","))
	return format
}

var Screenplays = screenplays{
	{
		ID:             "shiguangzhimen",
		Name:           "Gateway of Time",
		Labels:         []string{"Fantasy", "Business"},
		Synopsis:       "On the edge of a modern metropolis, a mysterious gateway of time suddenly emerges, connecting an ordinary corporate building to a magical and ancient fantasy kingdom. The protagonist, a successful but weary entrepreneur, accidentally traverses the gateway of time, finding himself in a world filled with magic and adventure.",
		WelcomeMessage: "```Gateway of Time\n%s：Hey %s! 🚀\nGuess what's making waves in Hollywood? \"Gateway of Time,\" and the star? Yep, it's you! 🚀💫\nSuper stoked for your fantasy adventure! Embrace those magical vibes, sprinkle some charm, and own that enchanted set. If it feels like a maze of mystery, just remember, you've got the magic touch!\nCan't wait to see you shine and turn \"Gateway of Time\" into the coolest thing since unicorns. 🦄✨ Break a leg, fantasy queen! 🎬😊```\n",
	},
	{
		ID:             "xingjiliefeng",
		Name:           "Interstellar Fissure",
		Labels:         []string{"Science Fiction", "Business"},
		Synopsis:       "In the future Earth, scientists have discovered a mysterious energy fissure leading to the depths of the universe. This fissure is believed to be a gateway to unknown galaxies, holding infinite resources and technological wonders. However, as human exploration teams traverse the fissure, they uncover a planet more mysterious and dangerous than they had anticipated.",
		WelcomeMessage: "```Interstellar Fissure\n%s：Hey %s! 🚀\nGuess what's officially taking off? \"Interstellar Fissure,\" and you're the cosmic star of the show! 🌌✨\nReady to sprinkle some stardust and make this sci-fi flick legendary? Embrace the weird, dive into the interstellar chaos, and most importantly, have a blast! 🎬\nIf the set ever feels like navigating a UFO with no manual, I got your back. We're all learning as we go, right?\nBreak a leg, space superstar! Can't wait to see you shine brighter than a supernova. 🌠👽🎥```\n",
	},
	{
		ID:             "moluzhuizong",
		Name:           "Endless Pursuit",
		Labels:         []string{"Crime", "Independent"},
		Synopsis:       "In a decaying city, Emo Havan, a former detective seeking redemption, is pulled back into the criminal world he vowed to leave behind when his partner, Max Davis, is mysteriously kidnapped. Emo discovers a city-wide conspiracy and must navigate the criminal underworld to rescue Max. As he unravels political corruption and power struggles, Emo, reunited with his journalist ex-girlfriend, faces a challenging climax. The film explores the fine line between justice and sin, with Emo's choices determining the city's destiny.",
		WelcomeMessage: "```Endless Pursuit\n%s：Hey %s! 🚀\nGuess what's causing a stir in Tinseltown? \"Endless Pursuit,\" starring the incredible you! 🕵️♀️💃\nMassive congrats on diving into the crime scene. Get ready to rock those mysterious vibes and steal the show. If it ever feels like a Hollywood whodunit, just remember, you're the heartthrob detective we've all been waiting for.\nEnjoy every scene, soak up the thrill, and let the charm flow naturally. Can't wait to see you shine and make \"Endless Pursuit\" the hottest crime story in town. Break a leg, superstar! 🎬🌟💕```\n",
	},
	{
		ID:             "chenshiderouqing",
		Name:           "City Serenade",
		Labels:         []string{"Life", "Artistic"},
		Synopsis:       "In a bustling city, artist Sophia, violinist Aidan, and retired psychologist Emily's lives unexpectedly intertwine, weaving a touching narrative. Through chance encounters and community activities, \"City Serenade\" explores the precious connections and inner peace found in the urban hustle.",
		WelcomeMessage: "```City Serenade\n%s：Hey %s! 🚀\nSo, \"City Serenade\" has officially hit the lights, and guess who's stealing the show? Yep, it's you! 🎬💖\nJust wanted to send a little note to say how awesome it is to see you rock the art film scene. Your charm is like the perfect melody in this city serenade. Enjoy every moment, soak up the magic, and let the city lights illuminate your brilliance.\nIf it ever feels like you're in a scene straight out of a dream, know that you've got this dreamboat leading the way. Can't wait to see you shine and make \"City Serenade\" the most enchanting film ever.\nBreak a leg, captivating star! 🌟😊🎥```\n",
	},
	{
		ID:             "aizaishuguang",
		Name:           "Love at Dawn",
		Labels:         []string{"Romance", "Business"},
		Synopsis:       "In bustling New York City, curator Allison and aspiring musician Mark's hearts collide in a chance meeting, sparking a romantic journey. Despite their differing worlds—Allison in the art scene, Mark chasing his musical dreams—they find a profound connection. As real-life pressures test their love, the film explores how love, with courage and compromise, can bridge seemingly disparate lives. \"Love at Dawn\" is a poignant tale of the unifying power of love.",
		WelcomeMessage: "```Love at Dawn\n%s：Hey %s! 🚀\nSo, \"Love at Dawn\" is officially in action, and guess who's stealing all the love vibes? Yep, it's you! 🎥💖\nJust wanted to drop a note and say how awesome it is to see you rock the romance scene. Your charm is like the perfect sunrise in this love story. Enjoy every moment, feel the butterflies, and let the romance blossom.\nIf it ever feels like you're in a scene straight out of a dream, just know that you've got this dreamboat leading the way. Can't wait to see you make \"Love at Dawn\" the most heart-melting film ever.\nBreak a leg, enchanting star! 🌄😊💫```\n",
	},
}

func (s screenplays) GetByID(ID string) *Screenplay {
	if ID == "" {
		return nil
	}
	for _, i2 := range s {
		if i2.ID == ID {
			return i2
		}
	}
	return nil
}

func (s screenplays) GetSelectMenuOptions() []discordgo.SelectMenuOption {
	options := make([]discordgo.SelectMenuOption, 0, 20)
	for _, screenplay := range Screenplays {
		options = append(options, discordgo.SelectMenuOption{
			Label:       screenplay.Name,
			Value:       screenplay.ID,
			Description: strings.Join(screenplay.Labels, "/"),
			Default:     false,
		})
	}

	return options
}
