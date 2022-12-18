package zaobao

import (
	"github.com/imroc/req/v3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db     sqlite.DB
	zaobao ZaoBao
)

type ZaoBao struct {
	Token string `gorm:"column:token"`
}

func init() {
	engine := control.Register("zaobao", &control.Options[*robot.Ctx]{
		Alias:      "每日早报",
		Help:       "输入 {每日早报|早报} => 获取每天60s读懂世界",
		DataFolder: "zaobao",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/zaobao.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.CreateAndFirstOrCreate("zaobao", &zaobao); err != nil {
		log.Fatalf("create weather table failed: %v", err)
	}

	engine.OnFullMatchGroup([]string{"早报", "每日早报"}).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if zaobao.Token == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置token\n指令：set zaobao token __\n相关秘钥申请地址：https://admin.alapi.cn")
			return
		}

		resp := getZaoBaoImageUrl(zaobao.Token)
		ctx.ReplyImage(resp)
	})

	engine.OnRegex("set zaobao token ([0-9a-zA-Z]{16})", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		token := ctx.State["regex_matched"].([]string)[1]
		if err := db.Orm.Table("zaobao").Where("1 = 1").Update("token", token).Error; err != nil {
			ctx.ReplyTextAndAt("token配置失败")
			return
		}
		zaobao.Token = token
		ctx.ReplyText("token设置成功")
	})
}

func getZaoBaoImageUrl(token string) string {
	var data zaoBaoJson
	if err := req.C().Get("https://v2.alapi.cn/api/zaobao").
		SetQueryParams(map[string]string{
			"format": "json",
			"token":  token,
		}).
		Do().Into(&data); err != nil {
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
