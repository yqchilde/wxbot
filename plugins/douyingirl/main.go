package douyingirl

import (
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("dygirl", &control.Options[*robot.Ctx]{
		Alias:            "抖音小姐姐",
		Help:             "输入 {抖音小姐姐} => 获取抖音小姐姐视频",
		DisableOnDefault: true,
	})
	engine.OnFullMatch("抖音小姐姐").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		ctx.ReplyShareLink("抖音小姐姐", "每次点进来都不一样呦", "https://imgbed.link/file/8708", "http://api.yqqy.top/douyin")
	})
}
