package service

import (
	"log"
	"trip-map/global"
	"trip-map/internal/model"
	"trip-map/internal/model/common/response"
	"trip-map/internal/model/request"
)

func BatchSaveNodeList(nodeList *[]model.Note) error {
	return global.Db.Save(nodeList).Error
}

func GetNoteListByCIdAndPage(req *request.NotePageRequest) (*response.PageResponseData[model.Note], error) {
	var notes []model.Note
	var total int64

	err := global.Db.
		Where("is_delete = 0").
		Where("c_poi_id = ?", req.CId).
		Model(&model.Note{}).
		Count(&total).
		Error
	if err != nil {
		log.Printf("查询分页查询景点帖子总数失败%v\n", err)
		return response.PageResponseWithEmpty[model.Note](&req.PageInfo), err
	}

	offset := (req.Offset - 1) * req.PageSize

	err = global.Db.
		Where("is_delete = 0").
		Where("c_poi_id = ?", req.CId).
		Offset(offset).
		Limit(req.PageSize).
		Find(&notes).
		Error
	if err != nil {
		log.Printf("查询分页查询景点帖子列表失败%v\n", err)
		return response.PageResponseWithEmpty[model.Note](&req.PageInfo), err
	}

	return response.PageResponseWithData[model.Note](notes, total, &req.PageInfo), nil
}
