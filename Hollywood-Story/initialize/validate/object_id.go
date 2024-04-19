package validate

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
	"vland.live/app/global"
)

// ObjectId 验证mongo的objectId是否符合规范
func ObjectId(fl validator.FieldLevel) bool {
	id, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}

	_, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		global.Logger.Error("validator ObjectId 错误", zap.Error(err),
			zap.Any("id", id),
			zap.Any("FieldName", fl.FieldName()),
			zap.Any("Top", fl.Top()),
			zap.Any("Parent", fl.Parent()),
			zap.Any("Param", fl.Param()),
		)
		return false
	}

	return true
}
