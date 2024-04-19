package constant

import "testing"

func TestGptNews_Format(t *testing.T) {
	news := GptNews{
		Post: "🌟📽️🌟 Just when we thought things couldn't get any more dramatic on the set of 'Endless Pursuit', a clash between our talented actress Agoni and the stubborn screenwriter Emily has caused chaos! 🤯📝 Who will emerge victorious in this battle of creative minds? Stay tuned for more behind-the-scenes drama! #EndlessPursuit #HollywoodDrama #FilmChaos",
		Comments: []string{
			"I can't believe the nerve of Emily! It's Agoni's character, she should have some say in it too! 🙄",
			"I'm living for this drama! Can't wait to see the sparks fly on screen! 🔥",
			"Emily needs to chill and let the actress do her thing. It's all about teamwork, right? 🎬",
			"This is the kind of drama I live for in Hollywood! 🍿 #TeamAgoni",
			"I bet this tension will make the movie even more intense and captivating! 🤩",
			"I hope they can resolve their differences for the sake of the movie. 🤞",
		},
	}
	t.Log(news.Format())
}
