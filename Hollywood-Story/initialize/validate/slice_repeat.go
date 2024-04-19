package validate

import (
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"vland.live/app/global"
)

func SliceRepeat(fl validator.FieldLevel) bool {
	tempMap := make(map[interface{}]int, 20)
	switch ids := fl.Field().Interface().(type) {
	case []string:
		for i := 0; i < len(ids); i++ {
			name := ids[i]
			count, ok := tempMap[name]
			count = 1
			if ok {
				count = count + 1
			}
			tempMap[name] = count
		}
	case []int:
		for i := 0; i < len(ids); i++ {
			name := ids[i]
			count, ok := tempMap[name]
			count = 1
			if ok {
				count = count + 1
			}
			tempMap[name] = count
		}
	case []int64:
		for i := 0; i < len(ids); i++ {
			name := ids[i]
			count, ok := tempMap[name]
			count = 1
			if ok {
				count = count + 1
			}
			tempMap[name] = count
		}
	default:
		global.Logger.Error("validator sliceRepeat 不支持的类型")
		return false
	}
	// 验证
	for k, v := range tempMap {
		if v > 1 {
			global.Logger.Error("validator sliceRepeat 验证出现重复值", zap.Any("value", k), zap.Any("count", v))
			return false
		}
	}
	return true
}
