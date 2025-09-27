package request

import "trip-map/internal/model/common"

type NotePageRequest struct {
	common.PageInfo
	CId string `json:"cId" form:"cId"`
}
