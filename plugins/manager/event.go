package manager

func registerEvent() {
	// 该插件按需引入，自行修改
	//engine := control.Register("event", &control.Options[*robot.Ctx]{
	//	HideMenu: true,
	//})

	//engine.OnMessage().SetBlock(false).Handle(func(ctx *robot.Ctx) {
	//	// 监听加好友事件
	//	if ctx.IsEventFriendVerify() {
	//		f := ctx.Event.FriendVerify
	//		if strings.ToLower(f.Content) != "wxbot" {
	//			return
	//		}
	//		if err := ctx.AgreeFriendVerify(f.V3, f.V4, f.Scene); err != nil {
	//			log.Errorf("同意好友请求失败: %v", err)
	//			return
	//		}
	//		ctx.SendText(f.WxId, "你好，我是wxbot，感谢您发现并使用该项目！\n\n我即将拉您进入wxbot交流群")
	//		time.Sleep(3 * time.Second)
	//		ctx.InviteIntoGroup("39171925457@chatroom", f.WxId, 2)
	//	}
	//})
}
