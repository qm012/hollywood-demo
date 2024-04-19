package params

import (
	"strconv"
)

// StatusParam 状态属性
type StatusParam struct {
	Status string `json:"status" binding:"required,oneof=true false"`
}

func (s *StatusParam) GetStatus() bool {
	status, _ := strconv.ParseBool(s.Status)
	return status
}
