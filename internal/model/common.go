package model

type Page struct {
	Total    int64  `json:"total" form:"total"`
	PageSize int    `json:"page_size" form:"page_size"`
	PageNum  int    `json:"page_num" form:"page_num"`
	Data     any    `json:"data" form:"data"`
	Order    string `json:"order" form:"order"`
	Sort     string `json:"sort" form:"sort"`
}
