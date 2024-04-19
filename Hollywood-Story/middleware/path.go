package middleware

import "net/http"

var (

	// 管理员才可以访问的url
	adminPath = map[string]struct{}{
		// ai
		http.MethodPost + "/ai-manage/v1/chat/completions": {},

		// 项目
		http.MethodPost + "/ai-manage/projects":       {},
		http.MethodPut + "/ai-manage/projects/:id":    {},
		http.MethodDelete + "/ai-manage/projects/:id": {},
		http.MethodGet + "/ai-manage/projects":        {},

		// prompt
		http.MethodPost + "/ai-manage/prompts":                                         {},
		http.MethodPut + "/ai-manage/prompts/:id":                                      {},
		http.MethodPut + "/ai-manage/prompts/:id/locked":                               {},
		http.MethodDelete + "/ai-manage/prompts/:id":                                   {},
		http.MethodGet + "/ai-manage/prompts":                                          {},
		http.MethodGet + "/ai-manage/prompts/:id":                                      {},
		http.MethodPost + "/ai-manage/prompts/:id/versions/save":                       {},
		http.MethodPost + "/ai-manage/prompts/:id/versions":                            {},
		http.MethodPatch + "/ai-manage/prompts/:id/versions/:version_id/is_production": {},
		http.MethodPatch + "/ai-manage/prompts/:id/versions/:version_id/name":          {},
		http.MethodDelete + "/ai-manage/prompts/:id/versions/:version_id":              {},
	}
)
