package zaobao

import (
	"fmt"
	"path/filepath"
	"time"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/pkg/utils"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db     sqlite.DB
	zaoBao ZaoBao
)

type ZaoBao struct {
	Token string `gorm:"column:token"`
	Date  string `gorm:"column:date"`
	Image string `gorm:"column:image"`
}

func init() {
	engine := control.Register("zaobao", &control.Options{
		Alias: "每日早报",
		Help: "指令:\n" +
			"* 早报 -> 获取每天60s读懂世界\n" +
			"* 每日早报 -> 获取每天60s读懂世界",
		DataFolder: "zaobao",
		OnEnable: func(ctx *robot.Ctx) {
			// todo 启动将定时任务加入到定时任务列表
			ctx.ReplyText("启用成功")
		},
		OnDisable: func(ctx *robot.Ctx) {
			// todo 停止将定时任务从定时任务列表移除
			ctx.ReplyText("禁用成功")
		},
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/zaobao.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.CreateAndFirstOrCreate("zaobao", &zaoBao); err != nil {
		log.Fatalf("create weather table failed: %v", err)
	}

	go pollingTask()

	engine.OnFullMatchGroup([]string{"早报", "每日早报"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if zaoBao.Token == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置token\n指令：set zaobao token __\n相关秘钥申请地址：https://admin.alapi.cn")
			return
		}
		if time.Now().Hour() < 5 {
			ctx.ReplyTextAndAt("早报数据每天5点后更新，当前时间不可用")
			return
		}
		imgCache := filepath.Join(engine.GetCacheFolder(), time.Now().Local().Format("20060102")+".jpg")
		if !utils.IsImageFile(imgCache) {
			if err := flushZaoBao(zaoBao.Token, imgCache); err != nil {
				ctx.ReplyTextAndAt("获取早报失败")
				return
			}
		}
		ctx.ReplyImage("local://" + imgCache)
	})

	engine.OnRegex("set zaobao token ([0-9a-zA-Z]{16})", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		token := ctx.State["regex_matched"].([]string)[1]
		if err := db.Orm.Table("zaobao").Where("1 = 1").Update("token", token).Error; err != nil {
			ctx.ReplyTextAndAt("token配置失败")
			return
		}
		zaoBao.Token = token
		ctx.ReplyText("token设置成功")
	})

	engine.OnFullMatch("get zaobao info", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var data ZaoBao
		if err := db.Orm.Table("zaobao").Limit(1).Find(&data).Error; err != nil {
			return
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - 每日早报\ntoken: %s", data.Token))
	})
}

func pollingTask() {
	// 计算下一个整点
	now := time.Now().Local()
	next := now.Add(10 * time.Minute).Truncate(10 * time.Minute)
	diff := next.Sub(now)
	timer := time.NewTimer(diff)
	<-timer.C
	timer.Stop()

	// 轮询任务
	ticker := time.NewTicker(10 * time.Minute)
	for range ticker.C {
		// 避开0点-5点(应该不会有人设置这个时间吧)
		if time.Now().Hour() < 5 {
			continue
		}

		// 早报token为空
		if zaoBao.Token == "" {
			continue
		}

		// 早报未更新
		imgCache := filepath.Join("./data/plugins/zaobao/cache", time.Now().Local().Format("20060102")+".jpg")
		if !utils.IsImageFile(imgCache) {
			if err := flushZaoBao(zaoBao.Token, imgCache); err != nil {
				continue
			}
		}
	}
}
