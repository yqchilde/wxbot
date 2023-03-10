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
	Mode      int    `gorm:"column:mode"`       // 模式
	GhWxId    string `gorm:"column:gh_wxid"`    // 监控公众号wxId
	Keys      string `gorm:"column:keys"`       // 监控关键字
	PushWxIds string `gorm:"column:push_wxids"` // 推送转发wxId,多个用英文逗号分隔
}

func init() {
	engine := control.Register("ghmonitor", &control.Options{
		Alias: "公众号监控",
		Help: "权限:\n" +
			"仅限机器人管理员\n\n" +
			"指令:\n" +
			"* 监控公众号 (gh_.*) 转发到 (.*)\n" +
			"* 监控公众号关键词 (.*) 转发到 (.*)",
		DataFolder: "ghmonitor",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/monitor.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.Create("monitor", &monitor); err != nil {
		log.Fatalf("create monitor table failed: %v", err)
	}

	// 监控模式
	// 1. 指定公众号的所有发布消息都转发
	// 2. 标题或描述有关键字的发布消息转发
	// 3. 文章中内容有关键字的发布消息转发

	// 配置监控模式1 ->指定公众号的所有发布消息都转发
	engine.OnRegex(`监控公众号 (gh_.*) 转发到 (.*)`, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		ghAccount := ctx.State["regex_matched"].([]string)[1]
		pushAccount := ctx.State["regex_matched"].([]string)[2]

		var data *Monitor
		dbRet := db.Orm.Table("monitor").Where("mode = 1 AND gh_wxid = ?", ghAccount).Limit(1).Find(&data)
		if err := dbRet.Error; err != nil {
			log.Errorf("设置失败，error: %v", err)
			ctx.ReplyTextAndAt("设置失败，请查看日志")
			return
		}

		if dbRet.RowsAffected == 0 {
			if err := db.Orm.Table("monitor").Create(&Monitor{
				Mode:      1,
				GhWxId:    ghAccount,
				PushWxIds: pushAccount,
			}).Error; err != nil {
				log.Errorf("设置失败，error: %v", err)
				ctx.ReplyTextAndAt("设置失败，请查看日志")
				return
			}
		} else {
			newPushWxIds := SliceUnion(strings.Split(pushAccount, ","), strings.Split(data.PushWxIds, ","))
			pushWxIds := strings.Join(newPushWxIds, ",")
			if err := db.Orm.Table("monitor").Where("1=1").Update("push_wxids", pushWxIds).Error; err != nil {
				log.Errorf("设置失败，error: %v", err)
				ctx.ReplyTextAndAt("设置失败，请查看日志")
				return
			}
		}
		ctx.ReplyText("设置成功")
	})

	// 配置监控模式2 ->标题或描述有关键字的发布消息转发
	engine.OnRegex(`监控公众号关键词 (.*) 转发到 (.*)`, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		keys := ctx.State["regex_matched"].([]string)[1]
		pushAccount := ctx.State["regex_matched"].([]string)[2]

		var data *Monitor
		dbRet := db.Orm.Table("monitor").Where("mode = 2").Limit(1).Find(&data)
		if err := dbRet.Error; err != nil {
			log.Errorf("设置失败，error: %v", err)
			ctx.ReplyTextAndAt("设置失败，请查看日志")
		}

		if dbRet.RowsAffected == 0 {
			if err := db.Orm.Table("monitor").Create(&Monitor{
				Mode:      2,
				Keys:      keys,
				PushWxIds: pushAccount,
			}).Error; err != nil {
				log.Errorf("设置失败，error: %v", err)
				ctx.ReplyTextAndAt("设置失败，请查看日志")
				return
			}
		} else {
			newKeys := SliceUnion(strings.Split(keys, ","), strings.Split(data.Keys, ","))
			newPushWxIds := SliceUnion(strings.Split(pushAccount, ","), strings.Split(data.PushWxIds, ","))
			keys := strings.Join(newKeys, ",")
			pushWxIds := strings.Join(newPushWxIds, ",")
			if err := db.Orm.Table("monitor").Where("1=1").Updates(map[string]interface{}{
				"keys":       keys,
				"push_wxids": pushWxIds,
			}).Error; err != nil {
				log.Errorf("设置失败，error: %v", err)
				ctx.ReplyTextAndAt("设置失败，请查看日志")
				return
			}
		}
		ctx.ReplyText("设置成功")
	})

	engine.OnMessage().SetBlock(false).Handle(func(ctx *robot.Ctx) {
		if ctx.IsEventSubscription() {
			var monitorList []Monitor
			if err := db.Orm.Table("monitor").Find(&monitorList).Error; err != nil {
				return
			}

			for _, data := range monitorList {
				switch data.Mode {
				case 1: // 模式1实现
					if data.GhWxId == ctx.Event.FromWxId {
						content := ctx.Event.MPMessage.Content
						var msgModel SubscriptionMsgModel
						if err := xml.Unmarshal([]byte(content), &msgModel); err != nil {
							return
						}
						msgModel.Fromusername = ctx.Bot.GetConfig().BotWxId
						if newXml, err := xml.Marshal(msgModel); err == nil {
							for _, wxId := range strings.Split(data.PushWxIds, ",") {
								ctx.SendXML(wxId, string(newXml))
							}
						}
					}
				case 2: // 模式2实现
					content := ctx.Event.MPMessage.Content
					var msgModel SubscriptionMsgModel
					if err := xml.Unmarshal([]byte(content), &msgModel); err != nil {
						return
					}
					for _, key := range strings.Split(data.Keys, ",") {
						if strings.Contains(msgModel.Appmsg.Title, key) || strings.Contains(msgModel.Appmsg.Des, key) {
							msgModel.Fromusername = ctx.Bot.GetConfig().BotWxId
							if newXml, err := xml.Marshal(msgModel); err == nil {
								for _, wxId := range strings.Split(data.PushWxIds, ",") {
									ctx.SendXML(wxId, string(newXml))
								}
							}
						}
					}
				}
			}
		}
	})
}
