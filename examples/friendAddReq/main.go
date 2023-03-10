package friendAddReq

import (
	"strings"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("friendAddReq", &control.Options{
		Alias: "自动通过好友添加请求",
	})

	// SetBlock(false) 代表不阻断后面的匹配
	engine.OnMessage().SetBlock(false).Handle(func(ctx *robot.Ctx) {
		// 监听加好友事件
		if ctx.IsEventFriendVerify() {
			f := ctx.Event.FriendVerifyMessage
			// 判断一下好友验证消息是否为"wxbot"
			if strings.ToLower(f.Content) != "wxbot" {
				return
			}
			if err := ctx.AgreeFriendVerify(f.V3, f.V4, f.Scene); err != nil {
				log.Errorf("同意好友请求失败: %v", err)
				return
			}
			ctx.SendText(f.WxId, "你好，我是wxbot，感谢您发现并使用该项目！")
		}
	})
}
