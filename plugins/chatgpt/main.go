package chatgpt

import (
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
	db         sqlite.DB             // 数据库
	chatCTXMap sync.Map              // 群号/私聊:消息上下文
	chatDone   = make(chan struct{}) // 用于结束会话
)

// ApiKey apikey表，存放openai key
type ApiKey struct {
	Key string `gorm:"column:key;index"`
}

// GptModel gptmodel表，存放gpt模型相关配置参数
type GptModel struct {
	Model            string  `gorm:"column:model"`
	MaxTokens        int     `gorm:"column:max_tokens"`
	Temperature      float64 `gorm:"column:temperature"`
	TopP             float64 `gorm:"column:top_p"`
	PresencePenalty  float64 `gorm:"column:presence_penalty"`
	FrequencyPenalty float64 `gorm:"column:frequency_penalty"`
}

func init() {
	engine := control.Register("chatgpt", &control.Options{
		Alias:      "ChatGPT",
		Help:       "输入 {开始会话} => 进行ChatGPT连续会话\n输入 {提问 [问题]} => 可以单独提问，没有上下文",
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
		log.Fatalf("create apikey table failed: %v", err)
	}
	// 初始化gpt 模型参数配置
	if err := db.CreateAndFirstOrCreate("gptmodel", &GptModel{
		Model:            gpt3.TextDavinci003Engine,
		MaxTokens:        512,
		Temperature:      0.7,
		TopP:             1,
		PresencePenalty:  0,
		FrequencyPenalty: 0,
	}); err != nil {
		log.Fatalf("create gptmodel table failed: %v", err)
	}

	// 连续会话
	engine.OnFullMatch("开始会话").SetBlock(true).Handle(func(ctx *robot.Ctx) {
		// 检查是否已经在进行会话
		if _, ok := chatCTXMap.Load(ctx.Event.FromUniqueID); ok {
			ctx.ReplyTextAndAt("当前已经在会话中了")
			return
		}

		// 开始会话
		recv, cancel := ctx.EventChannel(ctx.CheckGroupSession()).Repeat()
		defer cancel()
		chatCTXMap.LoadOrStore(ctx.Event.FromUniqueID, "")
		ctx.ReplyTextAndAt("收到！已开始ChatGPT连续会话中，输入\"结束会话\"结束会话，或5分钟后自动结束，请开始吧！")
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

				// 整理问题
				question := msg + "\n"
				if c, ok := chatCTXMap.Load(ctx.Event.FromUniqueID); ok {
					question = c.(string) + question
				}
				answer, err := AskChatGpt(question, 2*time.Second)
				if err != nil {
					ctx.ReplyTextAndAt("ChatGPT出错了, err: " + err.Error())
					continue
				}
				chatCTXMap.Store(ctx.Event.FromUniqueID, question+"\n"+answer)
				if newAnswer, isNeedReply := filterAnswer(answer); isNeedReply {
					retryAnswer, err := AskChatGpt(question + "\n" + answer + newAnswer)
					if err != nil {
						ctx.ReplyTextAndAt("ChatGPT出错了, err: " + err.Error())
						continue
					}
					chatCTXMap.Store(ctx.Event.FromUniqueID, question+"\n"+answer)
					ctx.ReplyTextAndAt(retryAnswer)
				} else {
					ctx.ReplyTextAndAt(newAnswer)
				}
			}
		}
	})

	// 单独提问，没有上下文处理
	engine.OnRegex(`^提问 (.*)$`).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		question := ctx.State["regex_matched"].([]string)[1]
		answer, err := AskChatGpt(question, time.Second)
		if err != nil {
			log.Errorf("ChatGPT出错了, err: %s", err.Error())
			return
		}
		if newAnswer, isNeedRetry := filterAnswer(answer); isNeedRetry {
			retryAnswer, err := AskChatGpt(question + "\n" + answer + newAnswer)
			if err != nil {
				log.Errorf("ChatGPT出错了, err: %s", err.Error())
				return
			}
			ctx.ReplyTextAndAt(fmt.Sprintf("问：%s \n--------------------\n答：%s", question, retryAnswer))
		} else {
			ctx.ReplyTextAndAt(fmt.Sprintf("问：%s \n--------------------\n答：%s", question, newAnswer))
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
		gptClient = nil
		ctx.ReplyText("apiKey设置成功")
	})

	// 删除openai api key
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
		gptClient = nil
		ctx.ReplyText("apiKey删除成功")
	})

	// 设置gpt3模型参数
	engine.OnRegex("set chatgpt model (.*)", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		args := ctx.State["regex_matched"].([]string)[1]
		if args == "reset" {
			if err := resetGptModel(); err != nil {
				ctx.ReplyText("重置模型参数失败, err: " + err.Error())
				return
			} else {
				gptModel = nil
				ctx.ReplyText("重置模型参数成功")
				return
			}
		}

		kv := strings.Split(args, "=")
		if len(kv) != 2 {
			ctx.ReplyText("参数格式错误")
			return
		}
		k, v := kv[0], kv[1]
		updates := make(map[string]interface{})

		switch k {
		case "ModelName":
			updates["model"] = v
		case "MaxTokens":
			updates["max_tokens"] = v
		case "Temperature":
			updates["temperature"] = v
		case "TopP":
			updates["top_p"] = v
		case "FrequencyPenalty":
			updates["frequency_penalty"] = v
		case "PresencePenalty":
			updates["presence_penalty"] = v
		default:
			ctx.ReplyTextAndAt(fmt.Sprintf("配置模型没有[%s]这个参数，请核实", k))
		}

		if err := db.Orm.Table("gptmodel").Where("1=1").Updates(updates).Error; err != nil {
			ctx.ReplyTextAndAt("更新失败")
			return
		}
		gptModel = nil
		ctx.ReplyTextAndAt("更新成功")
	})

	// 获取插件配置
	engine.OnFullMatch("get chatgpt info", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		// 获取模型配置
		var gptModel GptModel
		if err := db.Orm.Table("gptmodel").Limit(1).Find(&gptModel).Error; err != nil {
			log.Errorf("[ChatGPT] 获取模型配置失败, err: %s", err.Error())
			ctx.ReplyTextAndAt("插件 - ChatGPT\n获取模型配置失败")
			return
		}

		replyMsg := ""
		replyMsg += "----------\n"
		replyMsg += "ModelName: %s\n"
		replyMsg += "MaxTokens: %d\n"
		replyMsg += "Temperature: %.2f\n"
		replyMsg += "TopP: %.2f\n"
		replyMsg += "FrequencyPenalty: %.2f\n"
		replyMsg += "PresencePenalty: %.2f\n----------\n"
		replyMsg = fmt.Sprintf(replyMsg, gptModel.Model, gptModel.MaxTokens, gptModel.Temperature, gptModel.TopP, gptModel.FrequencyPenalty, gptModel.PresencePenalty)

		// key设置
		var keys []ApiKey
		if err := db.Orm.Table("apikey").Find(&keys).Error; err != nil || len(keys) == 0 {
			log.Errorf("[ChatGPT] 获取apiKey失败, err: %s", err.Error())
			ctx.ReplyTextAndAt("插件 - ChatGPT\n获取apiKey失败")
			return
		}
		for i := range keys {
			replyMsg += fmt.Sprintf("apiKey: %s\n", keys[i].Key)
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - ChatGPT\n%s", replyMsg))
	})
}

// apikey缓存
var apiKeys []ApiKey

// 获取gpt3客户端
func getGptClient() (gpt3.Client, error) {
	var keys []ApiKey
	if err := db.Orm.Table("apikey").Find(&keys).Error; err != nil {
		log.Errorf("[ChatGPT] 获取apikey失败, error:%s", err.Error())
		return nil, errors.New("获取apiKey失败")
	}
	if len(keys) == 0 {
		log.Errorf("[ChatGPT] 未设置apikey")
		return nil, fmt.Errorf("请先私聊机器人配置apiKey\n指令：set chatgpt apiKey __(多个key用;符号隔开)\napiKey获取请到https://beta.openai.com获取")
	}
	apiKeys = keys
	return gpt3.NewClient(keys[0].Key, gpt3.WithTimeout(time.Minute)), nil
}

// 获取gpt3模型配置
func getGptModel() (*GptModel, error) {
	var gptModel GptModel
	if err := db.Orm.Table("gptmodel").Limit(1).Find(&gptModel).Error; err != nil {
		log.Errorf("[ChatGPT] 获取模型配置失败, err: %s", err.Error())
		return nil, errors.New("获取模型配置失败")
	}
	return &gptModel, nil
}

// 重置gpt3模型配置
func resetGptModel() error {
	updates := map[string]interface{}{
		"model":             "text-davinci-003",
		"max_tokens":        512,
		"temperature":       0.7,
		"top_p":             1,
		"frequency_penalty": 0,
		"presence_penalty":  0,
	}
	if err := db.Orm.Table("gptmodel").Where("1=1").Updates(updates).Error; err != nil {
		log.Errorf("[ChatGPT] 重置模型配置失败, err: %s", err.Error())
		return err
	}
	return nil
}
