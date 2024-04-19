package params

import (
	"errors"
	"strings"
)

// PhoneParam 手机号属性
type PhoneParam struct {
	AreaCode string `json:"areaCode" uri:"areaCode" form:"areaCode" binding:"required,max=6"`
	Phone    string `json:"phone" uri:"phone" form:"phone" binding:"required,max=20"`
}

func (p *PhoneParam) Verify() error {
	// 验证前缀
	if strings.HasPrefix(p.AreaCode, "+") {
		return errors.New("区号不需要+号")
	}

	// 验证中国区号的手机号长度
	if p.AreaCode == "86" || p.AreaCode == "+86" {
		if len(p.Phone) != 11 {
			return errors.New("手机号格式错误，中国区内手机号长度为11位")
		}
	}

	return nil
}

// PhoneParamOmit 手机号属性
type PhoneParamOmit struct {
	AreaCode string `json:"areaCode" uri:"areaCode" form:"areaCode" binding:"omitempty,max=6"`
	Phone    string `json:"phone" uri:"phone" form:"phone" binding:"omitempty,max=20"`
}

func (p *PhoneParamOmit) Verify() error {
	// 验证前缀
	if strings.HasPrefix(p.AreaCode, "+") {
		return errors.New("区号不需要+号")
	}

	// 验证中国区号的手机号长度
	if p.AreaCode == "86" || p.AreaCode == "+86" {
		if len(p.Phone) != 11 {
			return errors.New("手机号格式错误，中国区内手机号长度为11位")
		}
	}

	return nil
}
