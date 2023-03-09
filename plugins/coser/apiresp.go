package coser

type APIResp struct {
	Code int    `json:"code,string"`
	Text string `json:"text,omitempty"`
	Data struct {
		Title string `json:"Title,omitempty"`
		Data  []string
	} `json:"data"`
}
