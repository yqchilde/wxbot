package plmm

type PlmmApiResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		ImageUrl        string `json:"imageUrl"`
		ImageSize       string `json:"imageSize"`
		ImageFileLength int    `json:"imageFileLength"`
	} `json:"data"`
}
