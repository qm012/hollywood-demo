package constant

import "github.com/qm012/dun"

var (
	ErrRefreshAttribute = dun.NewStatusCode(1001, "只能在开始回合之前刷新修改")
	ErrEventID          = dun.NewStatusCode(1002, "事件ID错误")
	ErrOptionID         = dun.NewStatusCode(1003, "选项ID错误")
	ErrNeedStartGame    = dun.NewStatusCode(1004, "流程错误，需要先开始/进入下一轮")
	ErrNeedOption       = dun.NewStatusCode(1005, "流程错误，需要先选择一个选项")
)
