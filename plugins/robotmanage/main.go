package robotmanage

import (
	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type RobotManage struct{ engine.PluginMagic }

var (
	pluginInfo = &RobotManage{
		engine.PluginMagic{
			HiddenMenu: true,
		},
	}
	_ = engine.InstallPlugin(pluginInfo)
)

func (m *RobotManage) OnRegister() {}

func (m *RobotManage) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.Content.FromWxid != robot.MyRobot.Manager {
			msg.ReplyTextAndAt("抱歉，你没有权限")
			return
		}
		if msg.Content.Msg == "刷新群数据" {
			if _, err := robot.MyRobot.GetGroupList(true); err != nil {
				msg.ReplyTextAndAt("刷新群数据失败，err: " + err.Error())
				return
			} else {
				msg.ReplyTextAndAt("刷新群数据成功")
				return
			}
		}
	}
}
