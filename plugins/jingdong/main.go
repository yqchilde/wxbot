package jingdong

import (
	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("jdbean", &control.Options{
		Alias:            "京豆上车",
		Help:             "输入 {京东上车} => 快上车和我一起挂京豆",
		DisableOnDefault: true,
	})
	engine.OnFullMatch("京豆上车").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if err := ctx.ReplyImage("C:\\Users\\Administrator\\Pictures\\jd\\qrcode.png"); err != nil {
			ctx.ReplyText(err.Error())
		}
	})
}
