package model

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const (
	ActorsSeparator = "and "
)

type Actors []*Actor

type Actor struct {
	ID          string   `json:"id" bson:"id"`                   // 主键
	Name        string   `json:"name" bson:"name"`               // 名字
	Occupation  string   `json:"occupation" bson:"occupation"`   // 职业
	Identity    string   `json:"identity" bson:"identity"`       // 身份
	Event       string   `json:"event_v1" bson:"event"`          // 事件
	Personality string   `json:"personality" bson:"personality"` // 性格
	Hobbies     []string `json:"hobbies" bson:"hobbies"`         // 爱好
}

func (a *Actor) Title() string {
	return fmt.Sprintf("%s（%s）", a.Name, a.Occupation)
}

// Copy 拷贝一份演员信息
func (a *Actor) Copy() *Actor {
	return &Actor{
		ID:          a.ID,
		Name:        a.Name,
		Occupation:  a.Occupation,
		Identity:    a.Identity,
		Event:       a.Event,
		Personality: a.Personality,
		Hobbies:     a.Hobbies,
	}
}

func (a Actors) Titles() []string {
	var titles = make([]string, 0, len(a))
	for _, actor := range a {
		titles = append(titles, actor.Title())
	}
	return titles
}

func (a Actors) Names() []string {
	var names = make([]string, 0, len(a))
	for _, actor := range a {
		names = append(names, actor.Name)
	}
	return names
}

// Copy 深拷贝一份演员列表
func (a Actors) Copy() Actors {
	// 需要将数据深拷贝，防止数据变更时产生一些连锁反应：bug
	var copiedActors = make(Actors, 0, len(a))
	for _, originalActor := range a {
		copiedActors = append(copiedActors, originalActor.Copy())
	}
	return copiedActors
}

func (a Actors) PrettyFormat() string {
	//var buffer bytes.Buffer
	//var data = make([][]string, 0, len(a))
	//for i, actor := range a {
	//	i++
	//	event := actor.Event
	//	if event == "" {
	//		event = "None"
	//	}
	//	data = append(data, []string{strconv.Itoa(i), actor.Name, actor.Occupation, actor.Identity, actor.Personality})
	//}
	//table := tablewriter.NewWriter(&buffer)
	//table.SetHeader([]string{"Sort", "Name", "Occupation", "Identity", "Personality"})
	//table.SetRowLine(true)
	//
	//table.SetCenterSeparator("*")
	//table.SetColumnSeparator("│")
	//table.SetRowSeparator("─")
	//table.SetColWidth(10)
	//table.AppendBulk(data)
	//table.Render()
	//
	//format := buffer.String()
	//format = fmt.Sprintf("```%s```", format)
	//
	//return format

	var builder strings.Builder

	// 构建表头
	builder.WriteString("## Sort | Name | Occupation | Identity | Personality\n")
	for i, actor := range a {
		i++
		event := actor.Event
		if event == "" {
			event = "None"
		}
		builder.WriteString(fmt.Sprintf("%v. | %-8v | %-12v | %-21v | %-v\n",
			i, actor.Name, actor.Occupation, actor.Identity, actor.Personality))
	}
	return builder.String()
}

func (a Actors) Format() string {
	var builder strings.Builder
	for _, actor := range a {
		builder.WriteString(fmt.Sprintf("name:%s, occupation:%s, personality:%s\n\t",
			actor.Name, actor.Occupation, actor.Personality))
	}
	return builder.String()
}

// DeleteByID 删除一个演员
func (a Actors) DeleteByID(ID string) Actors {
	if a == nil {
		return nil
	}
	var newActors = make(Actors, 0, len(a))
	for _, actor := range a {
		if actor.ID != ID {
			newActors = append(newActors, actor)
		}
	}
	return newActors
}

func (a Actors) GetByID(ID string) *Actor {
	if ID == "" {
		return nil
	}
	for _, actor := range a {
		if ID == actor.ID {
			return actor.Copy()
		}
	}
	return nil
}

// GetMultiByRandom 随机获取 num 个演员数据
//
//	场景1. 群聊
func (a Actors) GetMultiByRandom(num int) Actors {
	if num <= 0 || len(a) <= 0 {
		return Actors{}
	}

	copiedActors := a.Copy()

	// 打乱演员列表
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	r.Shuffle(len(copiedActors), func(i, j int) {
		copiedActors[i], copiedActors[j] = copiedActors[j], copiedActors[i]
	})

	// 打乱后在返回，不要提前返回，否则返回的数据是按顺序的。
	if num > len(a) {
		return copiedActors
	}

	return copiedActors[0:num]
}

func (a Actors) GetSingleByRandom() *Actor {
	copiedActors := (a).GetMultiByRandom(1)
	if len(copiedActors) == 0 {
		return nil
	}
	return copiedActors[0]
}
