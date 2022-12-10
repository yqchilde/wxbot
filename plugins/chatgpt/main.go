package chatgpt

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/yqchilde/pkgs/log"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db         sqlite.DB
	chatGPT    ChatGPT
	gpt3Client gpt3.Client
	chatCTXMap sync.Map // 群号/私聊:消息上下文
)

type chatCTX struct {
	prompt  string
	created time.Time
}

func init() {
	engine := control.Register("chatgpt", &control.Options[*robot.Ctx]{
		Alias:      "ChatGPT",
		Help:       "输入 {# 问题} => 获取ChatGPT回复",
		DataFolder: "chatgpt",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/chatgpt.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.CreateAndFirstOrCreate("chatgpt", &chatGPT); err != nil {
		log.Fatalf("create chatgpt table failed: %v", err)
	}

	gpt3Client = gpt3.NewClient(chatGPT.ApiKey, gpt3.WithTimeout(time.Minute))
	engine.OnPrefix("#").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		question, answer := ctx.State["args"].(string)+"\n", ""
		if question == "" {
			return
		}
		if chatGPT.ApiKey == "" {
			ctx.ReplyTextAndAt("请先私聊机器人配置apiKey\n指令：set chatgpt apiKey __\napiKey获取请到https://beta.openai.com获取")
			return
		}
		chatClear := []string{"清除上下文", "换个话题", "换个问题"}
		for i := range chatClear {
			if strings.Contains(question, chatClear[i]) {
				chatCTXMap.Delete(ctx.Event.FromUniqueID)
				ctx.ReplyText("😎我已结束聊天的上下文语境，您可以重新发起提问")
				return
			}
		}
		if c, ok := chatCTXMap.Load(ctx.Event.FromUniqueID); ok {
			if time.Now().Sub(c.(chatCTX).created) > time.Minute*5 {
				chatCTXMap.Delete(ctx.Event.FromUniqueID)
				ctx.ReplyTextAndAt("😊收到您的问题了，由于距离上一次提问已超过5分钟，我在重新构建上下文，马上就好~")
			} else {
				question = c.(chatCTX).prompt + question
			}
		} else {
			ctx.ReplyTextAndAt("😊收到您的问题了，正在构建上下文中，由于训练我的工程师们将我放在了大陆另一端，所以回复可能会有点慢哦~")
		}
		time.Sleep(5 * time.Second)
		answer, err := askChatGPT(question)
		if err != nil {
			ctx.ReplyTextAndAt("ChatGPT出错了, err: " + err.Error())
			return
		}
		chatCTXMap.Store(ctx.Event.FromUniqueID, chatCTX{prompt: question + "\n" + answer, created: time.Now()})
		if r, need := filterReply(answer); need {
			answer, err := askChatGPT(question + "\n" + answer + r)
			if err != nil {
				ctx.ReplyTextAndAt("ChatGPT出错了, err: " + err.Error())
				return
			}
			chatCTXMap.Store(ctx.Event.FromUniqueID, chatCTX{prompt: question + "\n" + answer, created: time.Now()})
			ctx.ReplyTextAndAt(answer)
		} else {
			ctx.ReplyTextAndAt(r)
		}
	})

	// 设置openai api key
	engine.OnRegex("set chatgpt apiKey (.*)", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		apiKey := ctx.State["regex_matched"].([]string)[1]
		if err := db.Orm.Table("chatgpt").Where("1 = 1").Update("api_key", apiKey).Error; err != nil {
			ctx.ReplyTextAndAt("设置apiKey失败")
			return
		}
		chatGPT.ApiKey = apiKey
		gpt3Client = gpt3.NewClient(chatGPT.ApiKey, gpt3.WithTimeout(time.Minute))
		ctx.ReplyText("apiKey设置成功")
	})

	// 获取插件配置
	engine.OnFullMatch("get chatgpt info", robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var chatGPT ChatGPT
		if err := db.Orm.Table("chatgpt").Limit(1).Find(&chatGPT).Error; err != nil {
			return
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - ChatGPT\napiKey: %s", chatGPT.ApiKey))
	})
}

type ChatGPT struct {
	ApiKey string `gorm:"column:api_key"`
}

func askChatGPT(question string) (string, error) {
	resp, err := gpt3Client.CompletionWithEngine(context.Background(), gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt:           []string{question},
		MaxTokens:        gpt3.IntPtr(512),
		Temperature:      gpt3.Float32Ptr(0.7),
		TopP:             gpt3.Float32Ptr(1),
		Echo:             false,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	})
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Text, nil
}

func filterReply(msg string) (string, bool) {
	punctuation := ",，!！?？"
	msg = strings.TrimSpace(msg)
	if len(msg) == 1 {
		return msg, true
	}
	if len(msg) == 3 && strings.ContainsAny(msg, punctuation) {
		return msg, true
	}
	msg = strings.TrimLeftFunc(msg, func(r rune) bool {
		if strings.ContainsAny(string(r), punctuation) {
			return true
		}
		return false
	})
	return msg, false
}
