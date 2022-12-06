package chatgpt

import (
	"context"
	"fmt"
	"gorm.io/gorm"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/glebarez/sqlite"
	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/robot"
)

func init() {
	engine := control.Register("chatgpt", &control.Options[*robot.Ctx]{
		Alias:      "ChatGPT",
		Help:       "输入 {gpt 内容} => 获取ChatGPT回复",
		DataFolder: "chatgpt",
	})

	db, err := gorm.Open(sqlite.Open(engine.GetDataFolder() + "/chatgpt.db"))
	if err != nil {
		log.Fatal(err)
	}
	db.Table("chatgpt").AutoMigrate(&ChatGPT{})

	engine.OnPrefix("gpt").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var chatGPT ChatGPT
		dbRet := db.Table("chatGPT").FirstOrCreate(&chatGPT)
		if err := dbRet.Error; err != nil {
			return
		}
		if chatGPT.Token == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置token\n指令：set chatgpt token __\n相关秘钥申请地址：https://openai.com/api")
			return
		}

		args, content := ctx.State["args"].(string), ""
		client := gpt3.NewClient(chatGPT.Token, gpt3.WithTimeout(time.Minute))
		_ = client.CompletionStreamWithEngine(context.Background(), gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
			Prompt:    []string{args},
			MaxTokens: gpt3.IntPtr(512),
			Echo:      true,
		}, func(resp *gpt3.CompletionResponse) {
			content += resp.Choices[0].Text
		})
		ctx.ReplyText(content)
	})

	engine.OnRegex("set chatgpt token ([0-9a-zA-Z-]{51})", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		token := ctx.State["regex_matched"].([]string)[1]
		db.Table("chatgpt").Where("1 = 1").Update("token", token)
		ctx.ReplyText("token设置成功")
	})

	engine.OnFullMatch("get chatgpt info", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var chatGPT ChatGPT
		if err := db.Table("chatgpt").Limit(1).Find(&chatGPT).Error; err != nil {
			return
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - ChatGPT\ntoken: %s", chatGPT.Token))
	})
}

type ChatGPT struct {
	Token string `gorm:"column:token"`
}
