package zaobao

import (
	"fmt"
	"sync"
	"time"

	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db     sqlite.DB
	zaoBao ZaoBao
	sendMu sync.Mutex
)

type ZaoBao struct {
	Token string `gorm:"column:token"`
	Date  string `gorm:"column:date"`
	Image string `gorm:"column:image"`
}

func init() {
	engine := control.Register("zaobao", &control.Options[*robot.Ctx]{
		Alias:      "每日早报",
		Help:       "输入 {每日早报|早报} => 获取每天60s读懂世界",
		DataFolder: "zaobao",
		OnEnable: func(ctx *robot.Ctx) {
			// todo 启动将定时任务加入到定时任务列表
			ctx.ReplyText("启用成功")
		},
		OnDisable: func(ctx *robot.Ctx) {
			// todo 停止将定时任务从定时任务列表移除
			ctx.ReplyText("禁用成功")
		},
		OnCronjob: func(ctx *robot.Ctx) {
			wxId := ctx.Event.FromUniqueID
			if zaoBao.Token == "" {
				log.Debugf("[cronjob] 早报token为空，wxId: %s", wxId)
				return
			}

			go func() {
				if zaoBao.Date == time.Now().Format("2006-01-02") {
					sendMu.Lock()
					ctx.SendImage(wxId, zaoBao.Image)
					sendMu.Unlock()
					return
				}
				ticker := time.NewTicker(1 * time.Minute)
				for range ticker.C {
					if zaoBao.Date != time.Now().Format("2006-01-02") {
						log.Debugf("[cronjob] 早报数据未更新，wxId: %s, 当前时间: %s，早报时间: %s", wxId, time.Now().Format("2006-01-02"), zaoBao.Date)
						continue
					}
					ticker.Stop()
					sendMu.Lock()
					ctx.SendImage(wxId, zaoBao.Image)
					sendMu.Unlock()
					break
				}
			}()
		},
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/zaobao.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.CreateAndFirstOrCreate("zaobao", &zaoBao); err != nil {
		log.Fatalf("create weather table failed: %v", err)
	}

	go func() {
		if zaoBao.Token == "" {
			log.Debug("早报token为空，请先设置token")
			return
		}
		ticker := time.NewTicker(30 * time.Minute)
		for range ticker.C {
			if zaoBao.Image != "" && zaoBao.Date == time.Now().Format("2006-01-02") {
				continue
			}
			if err := getZaoBao(zaoBao.Token); err != nil {
				log.Errorf("获取早报失败: %v", err)
			}
		}
	}()

	engine.OnFullMatchGroup([]string{"早报", "每日早报"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if zaoBao.Token == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置token\n指令：set zaobao token __\n相关秘钥申请地址：https://admin.alapi.cn")
			return
		}
		if zaoBao.Image == "" {
			if err := getZaoBao(zaoBao.Token); err != nil {
				log.Errorf("获取早报失败: %v", err)
				return
			}
		}
		if zaoBao.Date != time.Now().Format("2006-01-02") {
			log.Errorf("早报数据未更新，当前时间: %s, 早报时间: %s", time.Now().Format("2006-01-02"), zaoBao.Date)
			ctx.ReplyTextAndAt("早报数据未更新，请稍后再试")
			return
		}
		ctx.ReplyImage(zaoBao.Image)
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

func getZaoBao(token string) error {
	var data zaoBaoJson
	if err := req.C().Get("https://v2.alapi.cn/api/zaobao").
		SetQueryParams(map[string]string{
			"format": "json",
			"token":  token,
		}).Do().Into(&data); err != nil {
		log.Errorf("Zaobao获取失败: %v", err)
		return err
	}

	if err := db.Orm.Table("zaobao").Where("1=1").Updates(map[string]interface{}{
		"date":  data.Data.Date,
		"image": data.Data.Image,
	}).Error; err == nil {
		zaoBao.Date = data.Data.Date
		zaoBao.Image = data.Data.Image
	}
	return nil
}

type zaoBaoJson struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data struct {
		Date      string   `json:"date"`
		News      []string `json:"news"`
		Weiyu     string   `json:"weiyu"`
		Image     string   `json:"image"`
		HeadImage string   `json:"head_image"`
	} `json:"data"`
	Time  int    `json:"time"`
	LogId string `json:"log_id"`
}
