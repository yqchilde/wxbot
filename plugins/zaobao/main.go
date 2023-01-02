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
	cronMu sync.Mutex
)

type ZaoBao struct {
	Token string `gorm:"column:token"`
	Image string `gorm:"column:image"`
}

func init() {
	engine := control.Register("zaobao", &control.Options[*robot.Ctx]{
		Alias:      "每日早报",
		Help:       "输入 {每日早报|早报} => 获取每天60s读懂世界",
		DataFolder: "zaobao",
		OnCronjob: func(ctx *robot.Ctx) {
			wxId := ctx.Event.FromUniqueID
			cronMu.Lock()
			defer cronMu.Unlock()
			if zaoBao.Image == "" {
				log.Println("没有")
				zaoBao.Image = getZaoBaoImageUrl(zaoBao.Token)
			}
			ctx.SendImage(wxId, zaoBao.Image)
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
			return
		}
		if zaoBao.Image == "" {
			zaoBao.Image = getZaoBaoImageUrl(zaoBao.Token)
			db.Orm.Table("zaobao").Where("1=1").Update("image", zaoBao.Image)
		}
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			zaoBao.Image = getZaoBaoImageUrl(zaoBao.Token)
		}
	}()

	engine.OnFullMatchGroup([]string{"早报", "每日早报"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if zaoBao.Token == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置token\n指令：set zaobao token __\n相关秘钥申请地址：https://admin.alapi.cn")
			return
		}
		if zaoBao.Image == "" {
			zaoBao.Image = getZaoBaoImageUrl(zaoBao.Token)
			db.Orm.Table("zaobao").Where("1=1").Update("image", zaoBao.Image)
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

func getZaoBaoImageUrl(token string) string {
	var data zaoBaoJson
	if err := req.C().Get("https://v2.alapi.cn/api/zaobao").
		SetQueryParams(map[string]string{
			"format": "json",
			"token":  token,
		}).Do().Into(&data); err != nil {
		log.Errorf("Zaobao获取失败: %v", err)
		return ""
	}

	return data.Data.Image
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
