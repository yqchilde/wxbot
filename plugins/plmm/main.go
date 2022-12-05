package plmm

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/imroc/req/v3"
	"github.com/yqchilde/pkgs/log"
	"gorm.io/gorm"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("plmm", &control.Options[*robot.Ctx]{
		Alias:            "漂亮妹妹",
		Help:             "输入 {漂亮妹妹} => 获取漂亮妹妹",
		DataFolder:       "plmm",
		DisableOnDefault: true,
	})

	db, err := gorm.Open(sqlite.Open(engine.GetDataFolder() + "/plmm.db"))
	if err != nil {
		log.Fatal(err)
	}
	db.Table("plmm").AutoMigrate(&Plmm{})

	engine.OnFullMatch("漂亮妹妹").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var plmm Plmm
		dbRet := db.Table("plmm").FirstOrCreate(&plmm)
		if err := dbRet.Error; err != nil {
			log.Println(err)
			return
		}
		if plmm.AppId == "" || plmm.AppSecret == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置appId和appSecret\n指令：set plmm appId __\n指令：set plmm appSecret __\n相关秘钥申请地址：https://www.mxnzp.com/doc/detail?id=15")
			return
		}

		if len(plmmUrlStorage) > 50 {
			if err := ctx.ReplyImage(plmmUrlStorage[0]); err != nil {
				ctx.ReplyText(err.Error())
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
				ctx.ReplyText(err.Error())
			}
			plmmUrlStorage = plmmUrlStorage[1:]
		}
	})

	// 设置appId
	engine.OnRegex("set plmm appId ([0-9a-z]{16})", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		appId := ctx.State["regex_matched"].([]string)[1]
		db.Table("plmm").Where("1 = 1").Update("app_id", appId)
		ctx.ReplyText("appId设置成功")
	})

	// 设置appSecret
	engine.OnRegex("set plmm appSecret ([0-9a-zA-Z]{32})", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		appSecret := ctx.State["regex_matched"].([]string)[1]
		db.Table("plmm").Where("1 = 1").Update("app_secret", appSecret)
		ctx.ReplyText("appSecret设置成功")
	})

	// 获取插件配置信息
	engine.OnFullMatch("get plmm info", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var plmm Plmm
		if err := db.Table("plmm").Limit(1).Find(&plmm).Error; err != nil {
			return
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - 漂亮妹妹\nappId: %s\nappSecret: %s", plmm.AppId, plmm.AppSecret))
	})
}
