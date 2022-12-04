package memepicture

import (
	"encoding/base64"
	"strings"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("memepicture", &control.Options[*robot.Ctx]{
		Alias: "表情原图",
		Help:  "输入 {表情原图} => 30s内发送表情获取表情原图",
	})

	engine.OnFullMatch("表情原图", robot.MustMemePicture).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		url := ctx.State["image_url"].(string)
		if strings.HasPrefix(url[:20], "http") {
			url = base64.StdEncoding.EncodeToString([]byte(url))
			var resp apiResponse
			if err := req.C().Post("http://api.yqqy.top/thumbnail").
				SetFormData(map[string]string{
					"type":  "url",
					"image": url,
				}).Do().Into(&resp); err == nil {
				ctx.ReplyShareLink("快来下载你要的表情原图", "打开后长按图片可保存到本地哦", "http://api.yqqy.top/thumbnail?image="+resp.Data.Thumbnail, "http://api.yqqy.top/direct?image="+resp.Data.Original)
			}
		} else {
			var resp apiResponse
			if err := req.C().Post("http://api.yqqy.top/thumbnail").
				SetFormData(map[string]string{
					"type":  "base64",
					"image": url,
				}).Do().Into(&resp); err == nil {
				ctx.ReplyShareLink("快来下载你要的表情原图", "打开后长按图片可保存到本地哦", "http://api.yqqy.top/thumbnail?image="+resp.Data.Thumbnail, "http://api.yqqy.top/direct?image="+resp.Data.Original)
			}
		}
	})
}

type apiResponse struct {
	Code int `json:"code"`
	Data struct {
		Thumbnail string `json:"thumbnail"`
		Original  string `json:"original"`
	} `json:"data"`
}
