package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/PullRequestInc/go-gpt3"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db         sqlite.DB
	apiKeys    []ApiKey
	gpt3Client gpt3.Client
	chatCTXMap sync.Map // 群号/私聊:消息上下文
	chatDone   = make(chan struct{})
)

// ApiKey api_key表，存放api_key
type ApiKey struct {
	Key string `gorm:"column:key;index"`
}

func init() {
	engine := control.Register("chatgpt", &control.Options{
		Alias:      "ChatGPT",
		Help:       "输入 {开始会话} => 进行ChatGPT连续会话",
		DataFolder: "chatgpt",
		OnDisable: func(ctx *robot.Ctx) {
			ctx.ReplyText("禁用成功")
			chatDone <- struct{}{}
		},
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/chatgpt.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.Create("apikey", &ApiKey{}); err != nil {
		log.Fatalf("create chatgpt table failed: %v", err)
	}

	engine.OnFullMatch("开始会话").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if err := db.Orm.Table("apikey").Find(&apiKeys).Error; err != nil {
			log.Errorf("开始ChatGPT会话失败，error:%s", err.Error())
			ctx.ReplyTextAndAt("开启失败")
			return
		}
		if len(apiKeys) == 0 {
			ctx.ReplyTextAndAt("请先私聊机器人配置apiKey\n指令：set chatgpt apiKey __(多个key用;符号隔开)\napiKey获取请到https://beta.openai.com获取")
			return
		}
		gpt3Client = gpt3.NewClient(apiKeys[0].Key, gpt3.WithTimeout(time.Minute))
		if _, ok := chatCTXMap.Load(ctx.Event.FromUniqueID); ok {
			ctx.ReplyTextAndAt("当前已经在进行ChatGPT会话了")
			return
		}

		recv, cancel := ctx.EventChannel(ctx.CheckGroupSession()).Repeat()
		defer cancel()
		ctx.ReplyTextAndAt("收到！已开始ChatGPT会话，输入\"结束会话\"结束会话，或5分钟后自动结束，请开始吧！")
		chatCTXMap.LoadOrStore(ctx.Event.FromUniqueID, "")
		for {
			select {
			case <-time.After(time.Minute * 5):
				chatCTXMap.LoadAndDelete(ctx.Event.FromUniqueID)
				ctx.ReplyTextAndAt("😊检测到您已有5分钟不再提问，那我先主动结束会话咯")
				return
			case <-chatDone:
				chatCTXMap.LoadAndDelete(ctx.Event.FromUniqueID)
				ctx.ReplyTextAndAt("已退出ChatGPT")
				return
			case ctx := <-recv:
				msg := ctx.MessageString()
				if msg == "" {
					continue
				} else if msg == "结束会话" {
					chatCTXMap.LoadAndDelete(ctx.Event.FromUniqueID)
					ctx.ReplyTextAndAt("已结束聊天的上下文语境，您可以重新发起提问")
					return
				} else if msg == "清空会话" {
					chatCTXMap.Store(ctx.Event.FromUniqueID, "")
					ctx.ReplyTextAndAt("已清空会话，您可以继续提问新的问题")
				}
				question, answer := msg+"\n", ""
				if question == "" {
					continue
				}
				if c, ok := chatCTXMap.Load(ctx.Event.FromUniqueID); ok {
					question = c.(string) + question
				}
				time.Sleep(2 * time.Second)
				answer, err := askChatGPT(question)
				if err != nil {
					ctx.ReplyTextAndAt("ChatGPT出错了, err: " + err.Error())
					continue
				}
				chatCTXMap.Store(ctx.Event.FromUniqueID, question+"\n"+answer)
				if r, need := filterReply(answer); need {
					answer, err := askChatGPT(question + "\n" + answer + r)
					if err != nil {
						ctx.ReplyTextAndAt("ChatGPT出错了, err: " + err.Error())
						continue
					}
					chatCTXMap.Store(ctx.Event.FromUniqueID, question+"\n"+answer)
					ctx.ReplyTextAndAt(answer)
				} else {
					ctx.ReplyTextAndAt(r)
				}
			}
		}
	})

	// 设置openai api key
	engine.OnRegex("set chatgpt api[K|k]ey (.*)", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		keys := strings.Split(ctx.State["regex_matched"].([]string)[1], ";")
		failedKeys := make([]string, 0)
		for i := range keys {
			data := ApiKey{Key: keys[i]}
			if err := db.Orm.Table("apikey").Where(&data).FirstOrCreate(&data).Error; err != nil {
				failedKeys = append(failedKeys, keys[i])
				continue
			}
		}
		if len(failedKeys) > 0 {
			ctx.ReplyText(fmt.Sprintf("以下apiKey设置失败: %v", failedKeys))
			return
		}
		ctx.ReplyText("apiKey设置成功")
	})

	// 获取插件配置
	engine.OnFullMatch("get chatgpt info", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var apiKeys []ApiKey
		if err := db.Orm.Table("apikey").Find(&apiKeys).Error; err != nil {
			return
		}
		if len(apiKeys) == 0 {
			ctx.ReplyTextAndAt("插件 - ChatGPT\napiKey: 未设置")
			return
		}
		var apiKeyMsg string
		for i := range apiKeys {
			log.Println(apiKeys[i])
			apiKeyMsg += fmt.Sprintf("apiKey: %s\n", apiKeys[i].Key)
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - ChatGPT\n%s", apiKeyMsg))
	})

	engine.OnRegex("del chatgpt api[K|k]ey (.*)", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		keys := strings.Split(ctx.State["regex_matched"].([]string)[1], ";")
		failedKeys := make([]string, 0)
		for i := range keys {
			if err := db.Orm.Table("apikey").Where("key = ?", keys[i]).Delete(&ApiKey{}).Error; err != nil {
				failedKeys = append(failedKeys, keys[i])
				continue
			}
		}
		if len(failedKeys) > 0 {
			ctx.ReplyText(fmt.Sprintf("以下apiKey删除失败: %v", failedKeys))
			return
		}
		ctx.ReplyText("apiKey删除成功")
	})
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
		if strings.Contains(err.Error(), "You exceeded your current quota") {
			log.Printf("当前apiKey[%s]配额已用完, 将删除并切换到下一个", apiKeys[0].Key)
			db.Orm.Table("apikey").Where("key = ?", apiKeys[0].Key).Delete(&ApiKey{})
			apiKeys = apiKeys[1:]
			gpt3Client = gpt3.NewClient(apiKeys[0].Key, gpt3.WithTimeout(time.Minute))
			return askChatGPT(question)
		} else if strings.Contains(err.Error(), "The server had an error while processing your request") {
			log.Println("OpenAi服务出现问题，将重试")
			return askChatGPT(question)
		} else if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
			log.Println("OpenAi服务请求超时，将重试")
			return askChatGPT(question)
		} else if strings.Contains(err.Error(), "Please reduce your prompt") {
			return "", errors.New("OpenAi免费上下文长度限制为4097个词组，您的上下文长度已超出限制，请发送\"清空会话\"以清空上下文")
		} else {
			return "", err
		}
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
