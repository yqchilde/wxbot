package memepicture

import (
	"encoding/base64"

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
		ctx.ReplyShareLink("快来下载你要的表情原图", "打开后长按图片可保存到本地哦", "http://api.yqqy.top/img?url="+base64.StdEncoding.EncodeToString([]byte(url)), "http://api.yqqy.top/direct?url="+base64.StdEncoding.EncodeToString([]byte(url)))
	})
}
