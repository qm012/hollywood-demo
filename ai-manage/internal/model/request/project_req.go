package request

import (
	"github.com/qm012/dun"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"vland.live/app/initialize/params"
)

type CreateAdminProjectReq struct {
	Name string `json:"name" binding:"required,max=50"` // 项目名称
}

type UpdateAdminProjectReq struct {
	params.IdParamOmit
	Name string `json:"name" binding:"required,max=50"` // 项目名称
}

type DeleteAdminProjectReq struct {
	params.IdParamOmit
}

type SearchAdminProjectReq struct {
	Name string `form:"name" binding:"omitempty,max=50"` // 项目名称
	dun.PageSearch
}

func (s *SearchAdminProjectReq) Filter() bson.D {
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
