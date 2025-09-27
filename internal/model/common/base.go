package common

type PageInfo struct {
	Offset   int `json:"offset" form:"offset"`     // 页码
	PageSize int `json:"pageSize" form:"pageSize"` // 每页大小
}

type Point struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}
