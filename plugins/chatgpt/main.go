package chatgpt

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/yqchilde/wxbot/engine/control"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/pkg/utils"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	db          sqlite.DB // 数据库
	chatRoomCtx sync.Map  // 聊天室消息上下文
)

// ChatRoom chatRoomCtx -> ChatRoom => 维系每个人的上下文
type ChatRoom struct {
	chatId   string                         // 聊天室ID, 格式为: 聊天室ID_发送人ID
	chatTime time.Time                      // 聊天时间
	role     string                         // 角色
	content  []openai.ChatCompletionMessage // 聊天上下文内容
}

// ApiKey 表名:apikey，存放openai key
type ApiKey struct {
	Key string `gorm:"column:key;index"`
}

// ApiProxy 表名:apiproxy，存放openai 代理url地址
type ApiProxy struct {
	Id  uint   `gorm:"column:id;index"`
	Url string `gorm:"column:url;"`
}

// GptModel 表名:gptmodel，存放gpt模型相关配置参数
type GptModel struct {
	Model            string  `gorm:"column:model"`
	MaxTokens        int     `gorm:"column:max_tokens"`
	Temperature      float64 `gorm:"column:temperature"`
	TopP             float64 `gorm:"column:top_p"`
	PresencePenalty  float64 `gorm:"column:presence_penalty"`
	FrequencyPenalty float64 `gorm:"column:frequency_penalty"`
	ImageSize        string  `gorm:"column:image_size"`
}

// SystemRoles 表名:roles，存放系统角色
type SystemRoles struct {
	Role string `gorm:"column:role"`
	Desc string `gorm:"column:desc"`
}

var defaultGptModel = GptModel{
	Model:            "gpt-3.5-turbo",
	MaxTokens:        4096,
	Temperature:      0.8,
	TopP:             1.0,
	PresencePenalty:  0.0,
	FrequencyPenalty: 0.6,
	ImageSize:        "512x512",
}

func init() {
	engine := control.Register("chatgpt", &control.Options{
		Alias: "ChatGPT",
		Help: "指令:\n" +
			"* @机器人 [内容] -> 进行AI对话，计入上下文\n" +
			"* @机器人 提问 [问题] -> 单独提问，不计入上下文\n" +
			"* @机器人 作画 [描述] -> 进行AI作画\n" +
			"* @机器人 清空会话 -> 可清空与您的上下文\n" +
			"* @机器人 角色列表 -> 获取可切换的AI角色\n" +
			"* @机器人 当前角色 -> 获取当前用户的AI角色\n" +
			"* @机器人 创建角色 [角色名] [角色描述]\n" +
			"* @机器人 删除角色 [角色名]\n" +
			"* @机器人 切换角色 [角色名]",
		DataFolder: "chatgpt",
	})

	if err := sqlite.Open(engine.GetDataFolder()+"/chatgpt.db", &db); err != nil {
		log.Fatalf("open sqlite db failed: %v", err)
	}
	if err := db.Create("apikey", &ApiKey{}); err != nil {
		log.Fatalf("create apikey table failed: %v", err)
	}
	if err := db.Create("apiproxy", &ApiProxy{}); err != nil {
		log.Fatalf("create apiproxy table failed: %v", err)
	}
	if err := db.Create("roles", &SystemRoles{}); err != nil {
		log.Fatalf("create roles table failed: %v", err)
	}
	// 初始化gpt 模型参数配置
	initGptModel := defaultGptModel
	if err := db.CreateAndFirstOrCreate("gptmodel", &initGptModel); err != nil {
		log.Fatalf("create gptmodel table failed: %v", err)
	}
	// 初始化系统角色
	initRole()

	// 群聊并且艾特机器人
	engine.OnMessage(robot.OnlyAtMe).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		var (
			now = time.Now().Local()
			msg = ctx.MessageString()

			chatRoom = ChatRoom{
				chatId:   fmt.Sprintf("%s_%s", ctx.Event.FromWxId, ctx.Event.FromWxId),
				chatTime: time.Now().Local(),
				content:  []openai.ChatCompletionMessage{},
			}
		)

		// 预判断
		switch {
		case strings.TrimSpace(msg) == "菜单" || strings.TrimSpace(msg) == "帮助":
			ctx.ReplyTextAndAt("请发送菜单查看我还有哪些功能，无需@我哦")
			return
		case strings.TrimSpace(msg) == "清空会话":
			chatRoomCtx.Store(chatRoom.chatId, chatRoom)
			ctx.ReplyTextAndAt("已清空和您的上下文会话")
			return
		case strings.HasPrefix(msg, "提问"):
			messages := []openai.ChatCompletionMessage{{Role: "user", Content: msg}}
			answer, err := AskChatGpt(ctx, messages, time.Second)
			if err != nil {
				if errors.Is(err, ErrNoKey) {
					ctx.ReplyTextAndAt(err.Error())
				} else {
					ctx.ReplyTextAndAt("ChatGPT出错了，Err：" + err.Error())
				}
				return
			}
			ctx.ReplyTextAndAt(fmt.Sprintf("问：%s \n--------------------\n答：%s", msg, answer))
			return
		case strings.HasPrefix(msg, "作画"):
			b64, err := AskChatGptWithImage(ctx, msg, time.Second)
			if err != nil {
				log.Errorf("ChatGPT出错了，Err：%s", err.Error())
				ctx.ReplyTextAndAt("ChatGPT出错了，Err：" + err.Error())
				return
			}
			filename := fmt.Sprintf("%s/%s.png", engine.GetCacheFolder(), msg)
			if err := utils.Base64ToImage(b64, filename); err != nil {
				log.Errorf("作画失败，Err: %s", err.Error())
				ctx.ReplyTextAndAt("作画失败，请重试")
				return
			}
			ctx.ReplyImage("local://" + filename)
			return
		case strings.TrimSpace(msg) == "角色列表":
			replyMsg := "角色列表:\n"
			SystemRole.Each(func(key string, value interface{}) {
				replyMsg += fmt.Sprintf("%s\n", key)
			})
			ctx.ReplyTextAndAt(replyMsg)
			return
		case strings.TrimSpace(msg) == "当前角色":
			var role string
			if val, ok := chatRoomCtx.Load(ctx.Event.FromUniqueID + "_" + ctx.Event.FromWxId); ok {
				role = val.(ChatRoom).role
			}
			if role == "" {
				ctx.ReplyTextAndAt("当前角色为: 默认")
			} else {
				ctx.ReplyTextAndAt("当前角色为: " + role)
			}
			return
		case strings.HasPrefix(msg, "创建角色"):
			matched := regexp.MustCompile(`创建角色\s*(\S+)\s*(\S+)`).FindStringSubmatch(msg)
			role := matched[1]
			if _, ok := SystemRole.Get(role); ok {
				ctx.ReplyTextAndAt(fmt.Sprintf("角色[%s]已存在", role))
				return
			}
			desc := matched[2]
			if err := db.Orm.Table("roles").Create(&SystemRoles{Role: role, Desc: desc}).Error; err != nil {
				ctx.ReplyTextAndAt("创建角色失败")
				return
			}
			SystemRole.Set(role, desc)
			ctx.ReplyTextAndAt("创建角色成功")
			return
		case strings.HasPrefix(msg, "删除角色"):
			matched := regexp.MustCompile(`删除角色\s*(\S+)`).FindStringSubmatch(msg)
			role := matched[1]
			if _, ok := SystemRole.Get(role); !ok {
				ctx.ReplyTextAndAt(fmt.Sprintf("角色[%s]不存在", role))
				return
			}
			if err := db.Orm.Table("roles").Where("role = ?", role).Delete(&SystemRoles{}).Error; err != nil {
				ctx.ReplyTextAndAt("删除角色失败")
				return
			}
			SystemRole.Delete(role)
			ctx.ReplyTextAndAt("删除角色成功")
			return
		case strings.HasPrefix(msg, "切换角色"):
			matched := regexp.MustCompile(`切换角色\s*(\S+)`).FindStringSubmatch(msg)
			role := matched[1]
			if _, ok := SystemRole.Get(role); !ok {
				ctx.ReplyTextAndAt(fmt.Sprintf("角色[%s]不存在", role))
				return
			}

			var chatRoom = ChatRoom{
				chatId:   fmt.Sprintf("%s_%s", ctx.Event.FromUniqueID, ctx.Event.FromWxId),
				chatTime: time.Now().Local(),
				role:     role,
				content:  []openai.ChatCompletionMessage{},
			}
			chatRoomCtx.Store(chatRoom.chatId, chatRoom)
			ctx.ReplyTextAndAt("切换角色成功")
			return
		}

		// 正式处理
		if c, ok := chatRoomCtx.Load(chatRoom.chatId); ok {
			// 判断距离上次聊天是否超过10分钟了
			if now.Sub(c.(ChatRoom).chatTime) > 10*time.Minute {
				chatRoomCtx.Store(chatRoom.chatId, chatRoom)
				chatRoom.content = []openai.ChatCompletionMessage{{Role: "user", Content: msg}}
			} else {
				chatRoom.content = append(c.(ChatRoom).content, openai.ChatCompletionMessage{Role: "user", Content: msg})
			}
		} else {
			chatRoom.content = []openai.ChatCompletionMessage{{Role: "user", Content: msg}}
		}

		answer, err := AskChatGpt(ctx, chatRoom.content, time.Second)
		if err != nil {
			switch {
			case errors.Is(err, ErrNoKey):
				ctx.ReplyTextAndAt(err.Error())
			case errors.Is(err, ErrMaxTokens):
				ctx.ReplyTextAndAt("和你的聊天上下文内容太多啦，我的记忆好像在消退.. 糟糕，我忘记了..，请重新问我吧")
				chatRoomCtx.Store(chatRoom.chatId, chatRoom)
			default:
				ctx.ReplyTextAndAt("ChatGPT出错了，Err：" + err.Error())
			}
			return
		}
		chatRoom.content = append(chatRoom.content, openai.ChatCompletionMessage{Role: "assistant", Content: answer})
		chatRoomCtx.Store(chatRoom.chatId, chatRoom)
		ctx.ReplyTextAndAt(answer)
	})

	// 设置openai api 代理
	engine.OnRegex("set chatgpt proxy (.*)", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		url := ctx.State["regex_matched"].([]string)[1]
		data := ApiProxy{
			Id:  1,
			Url: url,
		}
		if err := db.Orm.Table("apiproxy").Save(&data).Error; err != nil {
			ctx.ReplyText(fmt.Sprintf("设置api代理地址失败: %v", url))
			return
		}
		gptClient = nil
		ctx.ReplyText("api代理设置成功")
		return
	})

	// 删除openai api 代理
	engine.OnRegex("del chatgpt proxy", robot.OnlyPrivate, robot.AdminPermission).SetBlock(true).Handle(func(ctx *robot.Ctx) {
		if err := db.Orm.Table("apiproxy").Where("id = 1").Delete(&ApiProxy{}).Error; err != nil {
			ctx.ReplyText(fmt.Sprintf("删除api代理地址失败: %v", err.Error()))
			return
		}
		gptClient = nil
		ctx.ReplyText("api代理删除成功")
		return
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
		case "ImageSize":
			updates["image_size"] = v
		default:
			ctx.ReplyTextAndAt(fmt.Sprintf("配置模型没有[%s]这个参数，请核实", k))
			return
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
		replyMsg += "ImageSize: %s\n----------\n"
		replyMsg = fmt.Sprintf(replyMsg, gptModel.Model, gptModel.ImageSize)

		// key设置
		var keys []ApiKey
		if err := db.Orm.Table("apikey").Find(&keys).Error; err != nil {
			log.Errorf("[ChatGPT] 获取apiKey失败, err: %s", err.Error())
			ctx.ReplyTextAndAt("插件 - ChatGPT\n获取apiKey失败")
			return
		}
		for i := range keys {
			replyMsg += fmt.Sprintf("apiKey: %s\n", keys[i].Key)
		}
		// Proxy查询
		var proxy ApiProxy
		if err := db.Orm.Table("apiproxy").Find(&proxy).Error; err != nil {
			log.Errorf("[ChatGPT] 获取apiproxy失败, err: %s", err.Error())
			ctx.ReplyTextAndAt("插件 - ChatGPT\n获取apiProxy失败")
			return
		}
		if len(proxy.Url) > 0 {
			replyMsg += fmt.Sprintf("apiProxy: %s\n", proxy.Url)
		}
		ctx.ReplyTextAndAt(fmt.Sprintf("插件 - ChatGPT\n%s", replyMsg))
	})
}
