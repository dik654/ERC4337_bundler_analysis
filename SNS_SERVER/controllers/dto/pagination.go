package dto

type PaginationRequest struct {
	PageNum  int64 `json:"page_num"`
	PageSize int64 `json:"page_size"`
}
