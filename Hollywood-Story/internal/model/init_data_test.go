package model

import (
	"testing"
	"vland.live/app/internal/constant"
)

func TestGetActorsByNum(t *testing.T) {
	actors := GetActorsByNum(constant.GlobalFilmActorNum)
	for i := 0; i < 1000; i++ {
		actor := actors.GetSingleByRandom()
		if actor == nil {
			t.Error("生成了错误的数据")
		}
	}

	t.Log("\n", actors.PrettyFormat())
}
