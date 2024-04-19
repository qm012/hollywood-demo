package params

import "go.mongodb.org/mongo-driver/bson/primitive"

// IdsParam 主键列表属性
type IdsParam struct {
	Ids []string `form:"ids" json:"ids" binding:"required,min=1,max=200,sliceRepeat,dive,objectId"` // 批量操作的id列表
}

func (ip *IdsParam) GetIds() []primitive.ObjectID {
	length := len(ip.Ids)
	ids := make([]primitive.ObjectID, 0, length)

	for i := 0; i < length; i++ {
		objectId, _ := primitive.ObjectIDFromHex(ip.Ids[i])
		ids = append(ids, objectId)
	}

	return ids
}

// IdsParamOmit 主键列表属性
type IdsParamOmit struct {
	Ids []string `form:"ids" json:"ids" binding:"omitempty,min=1,max=200,sliceRepeat,dive,objectId"` // 批量操作的id列表
}

func (ipo *IdsParamOmit) GetIds() []primitive.ObjectID {
	length := len(ipo.Ids)
	ids := make([]primitive.ObjectID, 0, length)

	for i := 0; i < length; i++ {
		objectId, _ := primitive.ObjectIDFromHex(ipo.Ids[i])
		ids = append(ids, objectId)
	}

	return ids
}
