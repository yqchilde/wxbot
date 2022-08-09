package model

type ShaoJiApiResponse struct {
	Success bool   `json:"success"`
	Imgurl  string `json:"imgurl"`
	Info    struct {
		Width  int    `json:"width"`
		Height int    `json:"height"`
		Type   string `json:"type"`
	} `json:"info"`
}
