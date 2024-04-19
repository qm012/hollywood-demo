package request

import (
	"github.com/qm012/dun"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vland.live/app/initialize/params"
)

type CreateAdminScreenplayReq struct {
	Key            string   `json:"key" binding:"required"`              // 键
	Name           string   `json:"name" binding:"required"`             // 名称
	Labels         []string `json:"labels" binding:"required"`           // 标签
	Synopsis       string   `json:"synopsis" binding:"required"`         // 概要
	WelcomeMessage string   `json:"welcome_message" binding:"omitempty"` // 欢迎语
}

type UpdateAdminScreenplayReq struct {
	params.IdParamOmit
	Key            string   `json:"key" binding:"required"`              // 键
	Name           string   `json:"name" binding:"required"`             // 名称
	Labels         []string `json:"labels" binding:"required"`           // 标签
	Synopsis       string   `json:"synopsis" binding:"required"`         // 概要
	WelcomeMessage string   `json:"welcome_message" binding:"omitempty"` // 欢迎语
}

type SearchAdminScreenplayReq struct {
	dun.PageSearch
	Name string `form:"name" binding:"omitempty,max=50"` // 剧本名称
}

func (s *SearchAdminScreenplayReq) Filter() bson.D {
	filter := make(bson.D, 0, 2)
	if len(s.Name) != 0 {

		tempFilter := make(bson.A, 0, 2)

		tempFilter = append(tempFilter, bson.M{
			"name": bson.M{"$regex": primitive.Regex{
				Pattern: s.Name,
				Options: "i", // 不区分大小写 i
			}},
		})
		filter = append(filter, bson.E{Key: "$or", Value: tempFilter})
	}

	return filter
}

type SearchScreenplayReq struct {
	dun.PageSearch
	Name     string `form:"name" binding:"omitempty,max=50"`             // 剧本名称
	ActorNum int    `form:"actor_num" binding:"required,min=13,max=100"` // 演员列表数量
}

func (s *SearchScreenplayReq) Filter() bson.D {
	filter := make(bson.D, 0, 2)
	if len(s.Name) != 0 {

		tempFilter := make(bson.A, 0, 2)

		tempFilter = append(tempFilter, bson.M{
			"name": bson.M{"$regex": primitive.Regex{
				Pattern: s.Name,
				Options: "i", // 不区分大小写 i
			}},
		})
		filter = append(filter, bson.E{Key: "$or", Value: tempFilter})
	}

	return filter
}

type DeleteAdminScreenplayReq struct {
	params.IdParamOmit
}
