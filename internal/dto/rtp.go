package dto

type RTPsResponse struct {
	RTPs []RTP `json:"RTPs"`
}
type RTP struct {
	GameCode string  `json:"GameCode"`
	RTP      float64 `json:"RTP"`
}

type ClearResponse struct {
	GameInfos []ClearInfo `json:"GameInfos"`
}
type ClearInfo struct {
	GameCode string `json:"GameCode"`
	Mesage   string `json:"Message"`
}
