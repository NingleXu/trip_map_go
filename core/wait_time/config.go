package wait_time

import "fmt"

type OpenState int

const (
	OpenStateUnknown  OpenState = iota
	OpenStateOpen               // 开放
	OpenStateClosed             // 关闭
	OpenStateMaintain           // 维护
)

type ScenicSpotWaitTime struct {
	ScenicSpotId string
	WaitMinute   int64 `json:"wait_minute"`
	StartTime    string
	EndTime      string
	OpenState    OpenState `json:"open_state"`
}

type WaitTimeHandler interface {
	GetWaitTime() (map[string]ScenicSpotWaitTime, error)
	ConvertID(originalID string) (customID string, err error)
}

var ScenicAreaWaitTimeHandlerMap = map[string]WaitTimeHandler{
	"7469417": NewZhuHaiChimelongOceanKingdom(),
}

type BaseHandler struct {
	IDMapping map[string]string
}

func (h *BaseHandler) ConvertID(originalID string) (string, error) {
	if customID, ok := h.IDMapping[originalID]; ok {
		return customID, nil
	}
	return "", fmt.Errorf("ID未映射: %s", originalID)
}
