package model

import (
	"strings"
	"testing"
	"vland.live/app/internal/constant"
)

func TestUser_GetNonMeet2Actors(t *testing.T) {
	var user = NewUser("NO:00001", "Qiming Tang")
	user.Nickname = "agoni"
	user.Film = &UserFilm{
		Screenplay: constant.Screenplays.GetByID("shiguangzhimen"),
		Actors:     GetActorsByNum(constant.GlobalFilmActorNum),
	}
	user.Film.Init()

	for i := 1; i <= constant.GlobalTotalRoundNum; i++ {
		t.Log(user.RoundText())
		roundActors := user.Film.GetNonMeet2Actors()

		t.Log("演员：", strings.Join(roundActors.Titles(), "，"))
		user.Film.DeleteMeetingActors(roundActors)
		user.CurrentRound++
	}

	actors := user.Film.AppointActors(5)
	t.Logf("\n%v\n", actors)
}
