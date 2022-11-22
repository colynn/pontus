package app

// PaginationReq ..
type PaginationReq struct {
	PageIndex int `json:"page_index"`
	PageSize  int `json:"page_size"`
}
