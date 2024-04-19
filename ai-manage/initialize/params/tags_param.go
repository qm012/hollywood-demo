package params

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TagsParamRequired 标签列表属性
type TagsParamRequired struct {
	Tags []string `form:"tags" json:"tags" binding:"required,min=1,max=10,sliceRepeat,dive,objectId"` // 标签列表
}

// TagsParamOmitempty 标签列表属性
type TagsParamOmitempty struct {
	Tags []string `form:"tags" json:"tags" binding:"omitempty,min=1,max=10,sliceRepeat,dive,objectId"` // 标签列表
}

func (b *TagsParamOmitempty) GetTags() []primitive.ObjectID {
	if b == nil {
		return []primitive.ObjectID{}
	}

	tags := make([]primitive.ObjectID, 0, len(b.Tags)) // 最终tags数据
	tempMap := make(map[string]struct{}, len(b.Tags))  // 过滤重复tag用

	for i := 0; i < len(b.Tags); i++ {
		tempId := b.Tags[i]
		_, ok := tempMap[tempId]
		if !ok {
			id, err := primitive.ObjectIDFromHex(tempId)
			if err != nil {
				fmt.Printf("BatchCreateSystemMaterialReq TagObjectIds err:%s\n", err.Error())
				continue
			}
			tempMap[tempId] = struct{}{}
			tags = append(tags, id)
		}
	}

	return tags
}

func (b *TagsParamRequired) GetTags() []primitive.ObjectID {
	if b == nil {
		return []primitive.ObjectID{}
	}

	tags := make([]primitive.ObjectID, 0, len(b.Tags)) // 最终tags数据
	tempMap := make(map[string]struct{}, len(b.Tags))  // 过滤重复tag用

	for i := 0; i < len(b.Tags); i++ {
		tempId := b.Tags[i]
		_, ok := tempMap[tempId]
		if !ok {
			id, err := primitive.ObjectIDFromHex(tempId)
			if err != nil {
				fmt.Printf("BatchCreateSystemMaterialReq TagObjectIds err:%s\n", err.Error())
				continue
			}
			tempMap[tempId] = struct{}{}
			tags = append(tags, id)
		}
	}

	return tags
}
