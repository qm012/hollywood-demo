package response

type SearchAdminProjectResp struct {
	ID             string `json:"id"`              // 主键
	Name           string `json:"name"`            // 项目名称
	PromptQuantity int    `json:"prompt_quantity"` // 使用数量
	ModifiedAt     int64  `json:"modified_at"`     // 更新时间
	CreatedAt      int64  `json:"created_at"`      // 创建时间
}
