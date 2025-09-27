package v1

import (
	"github.com/gin-gonic/gin"
	"trip-map/internal/model/common/response"
	"trip-map/internal/model/request"
	"trip-map/internal/service"
)

func GetNoteListByCIdAndPage(c *gin.Context) {
	// 初始化请求参数结构体
	var req request.NotePageRequest

	// 从请求参数中绑定数据（支持query参数、表单数据等）
	if err := c.ShouldBind(&req); err != nil {
		// 绑定失败时返回错误信息
		response.FailWithMessage("参数解析失败: "+err.Error(), c)
		return
	}

	// 验证必要参数
	if req.CId == "" {
		response.FailWithMessage("CId不能为空", c)
		return
	}

	// 处理分页默认值（如果前端未传递）
	if req.Offset <= 0 {
		req.Offset = 1 // 默认第一页
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 8
	}

	res, err := service.GetNoteListByCIdAndPage(&req)
	if err != nil {
		response.FailWithMessage("查询失败: "+err.Error(), c)
		return
	}

	// 返回分页数据
	response.OkWithData(res, c)
}
