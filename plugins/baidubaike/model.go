package baidubaike

type ApiResponse struct {
	Key        string `json:"key"`
	Desc       string `json:"desc"`
	Abstract   string `json:"abstract"`
	Copyrights string `json:"copyrights"`
}
