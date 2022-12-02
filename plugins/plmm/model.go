package plmm

type Plmm struct {
	AppId     string `json:"appId" gorm:"column:app_id"`
	AppSecret string `json:"appSecret" gorm:"column:app_secret"`
}

var plmmUrlStorage []string

type PlmmApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ImageUrl        string `json:"imageUrl"`
		ImageSize       string `json:"imageSize"`
		ImageFileLength int    `json:"imageFileLength"`
	} `json:"data"`
}
