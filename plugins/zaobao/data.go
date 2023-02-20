package zaobao

import (
	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

type zaoBaoResp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Date      string   `json:"date"`
		News      []string `json:"news"`
		Weiyu     string   `json:"weiyu"`
		Image     string   `json:"image"`
		HeadImage string   `json:"head_image"`
	} `json:"data"`
	Time  int    `json:"time"`
	LogId string `json:"log_id"`
}

func getZaoBao(token string) error {
	var data zaoBaoResp
	if err := req.C().Get("https://v2.alapi.cn/api/zaobao").
		SetQueryParams(map[string]string{
			"format": "json",
			"token":  token,
		}).Do().Into(&data); err != nil {
		log.Errorf("Zaobao获取失败: %v", err)
		return err
	}

	if err := db.Orm.Table("zaobao").Where("1=1").Updates(map[string]interface{}{
		"date":  data.Data.Date,
		"image": data.Data.Image,
	}).Error; err == nil {
		zaoBao.Date = data.Data.Date
		zaoBao.Image = data.Data.Image
	}
	return nil
}
