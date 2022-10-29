package engine

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine/config"
	"github.com/yqchilde/wxbot/engine/robot"
)

func InitRobot(conf *config.Config) error {
	// 检查配置
	var bot robot.BotConf
	conf.GetChild("robot").Unmarshal(&bot)
	if bot.Server == "" || bot.Token == "" {
		return errors.New("robot config error")
	}
	robot.MyRobot = bot
	bot.GetRobotInfo()
	log.Println("success to start robot")

	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.POST("/wxbot/callback", func(c *gin.Context) {
		var msg robot.Message
		if err := c.ShouldBindJSON(&msg); err != nil {
			c.JSON(http.StatusOK, gin.H{"Code": "-1"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"Code": "0"})

		// 菜单
		menuItems := "YY Bot🤖\n"
		for _, plugin := range Plugins {
			if plugin.RawConfig["enable"] != false {
				plugin.Config.OnEvent(&msg)
			}
			if !plugin.HiddenMenu {
				menuItems += plugin.Desc + "\n"
			}
		}

		if msg.IsAt() {
			msg.ReplyText("您可以发送menu | 菜单解锁更多功能😎")
		}
		if msg.MatchTextCommand([]string{"menu", "菜单", "/menu"}) {
			msg.ReplyText(menuItems)
		}
		if msg.IsSendByPrivateChat() {
			if msg.IsText() {
				log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", msg.Content.FromName, msg.Content.Msg))
			} else {
				log.Println(fmt.Sprintf("收到私聊(%s)消息 ==> %v", msg.Content.FromName, msg.Content.Msg))
			}
		}
		if msg.IsSendByGroupChat() {
			if msg.IsText() {
				log.Println(fmt.Sprintf("收到群聊(%s[%s])消息 ==> %v", msg.Content.FromGroupName, msg.Content.FromName, msg.Content.Msg))
			} else {
				log.Println(fmt.Sprintf("收到群聊(%s[%s])消息 ==> %v", msg.Content.FromGroupName, msg.Content.FromName, msg.Content.Msg))
			}
		}
	})
	r.Run(":9528")
	return nil
}
