package params

import "go.mongodb.org/mongo-driver/bson/primitive"

// IdParam 主键属性
type IdParam struct {
	Id string `json:"id" uri:"id" form:"id" binding:"required,objectId"`
}

func (i *IdParam) GetId() primitive.ObjectID {
	objectId, _ := primitive.ObjectIDFromHex(i.Id)
	return objectId
}

// IdParamOmit 主键属性
type IdParamOmit struct {
	Id string `json:"id" uri:"id" form:"id" binding:"omitempty,objectId"`
}

func (i *IdParamOmit) GetId() primitive.ObjectID {
	if i.Id == "" {
		return primitive.NilObjectID
	}
	objectId, _ := primitive.ObjectIDFromHex(i.Id)
	return objectId
}
