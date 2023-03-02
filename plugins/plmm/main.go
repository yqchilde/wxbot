package plmm

import (
	"fmt"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db   sqlite.DB
	plmm Plmm
)

func init() {
	engine := control.Register("plmm", &control.Options{
		Alias: "漂亮妹妹",
		Help: "指令:\n" +
			"* 漂亮妹妹 -> 获取漂亮妹妹",
		DataFolder: "plmm",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/plmm.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.CreateAndFirstOrCreate("plmm", &plmm); err != nil {
		log.Fatalf("create plmm table failed: %v", err)
	}

	engine.OnFullMatch("漂亮妹妹").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if plmm.AppId == "" || plmm.AppSecret == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置appId和appSecret\n指令：set plmm appId __\n指令：set plmm appSecret __\n相关秘钥申请地址：https://www.mxnzp.com/doc/detail?id=15")
			return
		}

		if len(plmmUrlStorage) > 0 {
			if err := ctx.ReplyImage(plmmUrlStorage[0]); err != nil {
				log.Errorf("[plmm] 发送图片失败: %v", err)
			}
			plmmUrlStorage = plmmUrlStorage[1:]
		} else {
			var resp PlmmApiResponse
			api := fmt.Sprintf("https://www.mxnzp.com/api/image/girl/list/random?app_id=%s&app_secret=%s", plmm.AppId, plmm.AppSecret)
			if err := req.C().SetBaseURL(api).Get().Do().Into(&resp); err != nil {
				return
			}
			if resp.Code != 1 {
				return
			}
			for _, val := range resp.Data {
				plmmUrlStorage = append(plmmUrlStorage, val.ImageUrl)
			}
			if err := ctx.ReplyImage(plmmUrlStorage[0]); err != nil {
				log.Errorf("[plmm] 发送图片失败: %v", err)
			}
			plmmUrlStorage = plmmUrlStorage[1:]
		}
	})

	// 设置appId
	engine.OnRegex("set plmm appId ([0-9a-z]{16})", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		appId := ctx.State["regex_matched"].([]string)[1]
		if err := db.Orm.Table("plmm").Where("1 = 1").Update("app_id", appId).Error; err != nil {
			ctx.ReplyText("appId设置失败")
			return
		}
		plmm.AppId = appId
		ctx.ReplyText("appId设置成功")
	})

	// 设置appSecret
	engine.OnRegex("set plmm appSecret ([0-9a-zA-Z]{32})", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		appSecret := ctx.State["regex_matched"].([]string)[1]
		if err := db.Orm.Table("plmm").Where("1 = 1").Update("app_secret", appSecret).Error; err != nil {
			ctx.ReplyText("appSecret设置失败")
			return
		}
		plmm.AppSecret = appSecret
		ctx.ReplyText("appSecret设置成功")
	})

	// 获取插件配置信息
	engine.OnFullMatch("get plmm info", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var plmm Plmm
		if err := db.Orm.Table("plmm").Limit(1).Find(&plmm).Error; err != nil {
			return
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - 漂亮妹妹\nappId: %s\nappSecret: %s", plmm.AppId, plmm.AppSecret))
	})
}
