package request

type CaptureCityScenicSpotReq struct {
	CaptureCityScenicSpotList []CaptureCityScenicSpot `json:"captureCityScenicSpotList"`
}
type CaptureCityScenicSpot struct {
	CityName string `json:"cityName"`
	DestId   string `json:"destId"`
}
