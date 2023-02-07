package ghmonitor

import (
	"encoding/xml"
	"strings"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db      sqlite.DB
	monitor Monitor
)

type Monitor struct {
	Mode     int    `gorm:"column:mode"`      // 模式
	GhWxId   string `gorm:"column:gh_wxid"`   // 监控公众号wxId
	PushWxId string `gorm:"column:push_wxid"` // 推送转发wxId
}

func init() {
	engine := control.Register("ghmonitor", &control.Options[*robot.Ctx]{
		Alias:      "公众号监控",
		DataFolder: "ghmonitor",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/monitor.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.Create("monitor", &monitor); err != nil {
		log.Fatalf("create monitor table failed: %v", err)
	}

	engine.OnRegex(`monitor (gh_.*) push (.*)`, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		ghAccount := ctx.State["regex_matched"].([]string)[1]
		pushAccount := strings.Split(ctx.State["regex_matched"].([]string)[2], ";")

		for i := range pushAccount {
			db.Orm.Table("monitor").FirstOrCreate(&Monitor{
				Mode:     1,
				GhWxId:   ghAccount,
				PushWxId: pushAccount[i],
			})
		}

		ctx.ReplyText("设置成功，查看请输入 monitor get")
	})

	// 监控模式
	// 1. 指定公众号的所有发布消息都转发
	// 2. 标题或描述有关键字的发布消息转发
	// 3. 文章中内容有关键字的发布消息转发

	engine.OnMessage().SetBlock(false).Handle(func(ctx *robot.Ctx) {
		if ctx.IsEventSubscription() {
			var monitorList []Monitor
			if err := db.Orm.Table("monitor").Find(&monitorList).Error; err != nil {
				return
			}

			for _, data := range monitorList {
				switch data.Mode {
				case 1: // 模式1实现
					content := ctx.Event.SubscriptionMessage.Content
					var msgModel SubscriptionMsgModel
					if err := xml.Unmarshal([]byte(content), &msgModel); err != nil {
						return
					}
					msgModel.Fromusername = robot.WxBot.BotConfig.BotWxId
					if newXml, err := xml.Marshal(msgModel); err == nil {
						ctx.SendXML(data.PushWxId, string(newXml))
					}
				}
			}
		}
	})
}
