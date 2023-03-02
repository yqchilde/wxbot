package zaobao

import (
	"fmt"
	"time"

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

func flushZaoBao(token, dstFile string) error {
	var data zaoBaoResp
	err := req.C().Get("https://v2.alapi.cn/api/zaobao").
		SetQueryParams(map[string]string{"format": "json", "token": token}).
		Do().Into(&data)
	if err != nil {
		log.Errorf("[zaoBao]获取数据失败: %v", err)
		return fmt.Errorf("获取数据失败")
	}

	if data.Data.Date != time.Now().Local().Format("2006-01-02") {
		return fmt.Errorf("获取数据失败: %v", "数据未更新，请稍后再试")
	}

	updates := map[string]interface{}{
		"date":  data.Data.Date,
		"image": data.Data.Image,
	}
	if err := db.Orm.Table("zaobao").Where("1=1").Updates(updates).Error; err != nil {
		log.Errorf("[zaoBao]更新数据库失败: %v", err)
		return fmt.Errorf("获取数据失败")
	}
	zaoBao.Date = data.Data.Date
	zaoBao.Image = data.Data.Image

	// 下载图片
	if err := req.C().Get(data.Data.Image).SetOutputFile(dstFile).Do().Err; err != nil {
		log.Errorf("[zaoBao]下载图片失败: %v", err)
		return fmt.Errorf("获取数据失败")
	}
	return nil
}
