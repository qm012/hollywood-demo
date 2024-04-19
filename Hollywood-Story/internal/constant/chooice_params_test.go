package constant

import (
	"testing"
	"time"
)

func TestLocation_Random(t *testing.T) {

	for i := 0; i < 10; i++ {
		time.Sleep(time.Microsecond * 100)
		t.Log("地点：", Location.Random(), "事件主题：", EventTopic.Random(), "行事标签：", ActionTags.Random(), "天气：", Weathers.Random())
	}

}
