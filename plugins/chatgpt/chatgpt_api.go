package chatgpt

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

var (
	gptClient *openai.Client
	gptModel  *GptModel
)

var (
	ErrNoKey              = fmt.Errorf("请先私聊机器人配置apiKey\n指令：set chatgpt apikey __(多个key用;符号隔开)\napiKey获取请到https://beta.openai.com获取")
	ErrMaxTokens          = errors.New("OpenAi免费上下文长度限制为4097个词组，您的上下文长度已超出限制")
	ErrExceededQuota      = errors.New("OpenAi配额已用完，请联系管理员")
	ErrIncorrectKey       = errors.New("OpenAi ApiKey错误，请联系管理员")
	ErrServiceUnavailable = errors.New("ChatGPT服务异常，请稍后再试")
)

// apikey缓存
var apiKeys []ApiKey

// 获取gpt3客户端
func getGptClient() (*openai.Client, error) {
	var keys []ApiKey
	if err := db.Orm.Table("apikey").Find(&keys).Error; err != nil {
		log.Errorf("[ChatGPT] 获取apikey失败, error:%s", err.Error())
		return nil, errors.New("获取apiKey失败")
	}
	if len(keys) == 0 {
		return nil, ErrNoKey
	}
	apiKeys = keys

	var proxy []ApiProxy
	if err := db.Orm.Table("apiproxy").Find(&proxy).Error; err != nil {
		log.Errorf("[ChatGPT] 获取apiProxy失败, error:%s", err.Error())
		return nil, errors.New("获取apiProxy失败")
	}

	config := openai.DefaultConfig(keys[0].Key)
	for i := range proxy {
		if proxy[i].Id == 1 {
			config.BaseURL = proxy[i].Url
		}
		if proxy[i].Id == 2 {
			proxyUrl, err := url.Parse(proxy[i].Url)
			if err != nil {
				log.Errorf("[ChatGPT] 解析http_proxy失败, error:%s", err.Error())
				return nil, errors.New("解析http_proxy失败")
			}
			transport := &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			}
			config.HTTPClient = &http.Client{
				Transport: transport,
				Timeout:   time.Minute * 2,
			}
		}
	}

	return openai.NewClientWithConfig(config), nil
}

// 获取gpt3模型配置
func getGptModel() (*GptModel, error) {
	var model GptModel
	if err := db.Orm.Table("gptmodel").Limit(1).Find(&model).Error; err != nil {
		log.Errorf("[ChatGPT] 获取模型配置失败, err: %s", err.Error())
		return nil, errors.New("获取模型配置失败")
	}
	if model.ImageSize == "" {
		model.ImageSize = openai.CreateImageSize512x512
	}
	return &model, nil
}

// 重置gpt3模型配置
func resetGptModel() error {
	if err := db.Orm.Table("gptmodel").Where("1=1").Updates(&defaultGptModel).Error; err != nil {
		log.Errorf("[ChatGPT] 重置模型配置失败, err: %s", err.Error())
		return err
	}
	return nil
}

// AskChatGpt 向ChatGPT请求回复
func AskChatGpt(ctx *robot.Ctx, messages []openai.ChatCompletionMessage, delay ...time.Duration) (answer string, err error) {
	// 获取客户端
	if gptClient == nil {
		gptClient, err = getGptClient()
		if err != nil {
			return "", err
		}
	}

	// 获取模型
	if gptModel == nil {
		gptModel, err = getGptModel()
		if err != nil {
			return "", err
		}
	}

	// 延迟请求
	if len(delay) > 0 {
		time.Sleep(delay[0])
	}

	// 处理用户role
	var role string
	if val, ok := chatRoomCtx.Load(ctx.Event.FromUniqueID + "_" + ctx.Event.FromWxId); ok {
		role = val.(ChatRoom).role
	}
	if role == "" {
		role = "默认"
	}

	var chatMessages []openai.ChatCompletionMessage
	if strings.Contains(SystemRole.MustGet(role).(string), "%s") {
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    "system",
			Content: fmt.Sprintf(SystemRole.MustGet(role).(string), robot.GetBot().GetConfig().BotNickname),
		})
	} else {
		chatMessages = append(chatMessages, openai.ChatCompletionMessage{
			Role:    "system",
			Content: SystemRole.MustGet(role).(string),
		})
	}
	chatMessages = append(chatMessages, messages...)

	resp, err := gptClient.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    gptModel.Model,
		Messages: chatMessages,
	})
	// 处理响应回来的错误
	if err != nil {
		if strings.Contains(err.Error(), "You exceeded your current quota") {
			log.Printf("当前apiKey[%s]配额已用完, 将删除并切换到下一个", apiKeys[0].Key)
			db.Orm.Table("apikey").Where("key = ?", apiKeys[0].Key).Delete(&ApiKey{})
			if len(apiKeys) == 1 {
				return "", ErrExceededQuota
			}
			apiKeys = apiKeys[1:]
			gptClient = openai.NewClient(apiKeys[0].Key)
			return AskChatGpt(ctx, messages)
		}
		if strings.Contains(err.Error(), "Please reduce your prompt") || strings.Contains(err.Error(), "Please reduce the length of the messages") {
			return "", ErrMaxTokens
		}
		if strings.Contains(err.Error(), "Incorrect API key") {
			return "", ErrIncorrectKey
		}
		if strings.Contains(err.Error(), "invalid character") {
			return "", ErrServiceUnavailable
		}
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", ErrServiceUnavailable
	}
	return resp.Choices[0].Message.Content, nil
}

// AskChatGptWithImage 向ChatGPT请求回复图片
func AskChatGptWithImage(ctx *robot.Ctx, prompt string, delay ...time.Duration) (b64 string, err error) {
	// 获取客户端
	if gptClient == nil {
		gptClient, err = getGptClient()
		if err != nil {
			return "", err
		}
	}

	// 获取模型
	if gptModel == nil {
		gptModel, err = getGptModel()
		if err != nil {
			return "", err
		}
	}

	// 延迟请求
	if len(delay) > 0 {
		time.Sleep(delay[0])
	}

	resp, err := gptClient.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         prompt,
		Size:           gptModel.ImageSize,
		ResponseFormat: openai.CreateImageResponseFormatB64JSON,
	})
	// 处理响应回来的错误
	if err != nil {
		if strings.Contains(err.Error(), "You exceeded your current quota") {
			log.Printf("当前apiKey[%s]配额已用完, 将删除并切换到下一个", apiKeys[0].Key)
			db.Orm.Table("apikey").Where("key = ?", apiKeys[0].Key).Delete(&ApiKey{})
			if len(apiKeys) == 1 {
				return "", ErrExceededQuota
			}
			apiKeys = apiKeys[1:]
			gptClient = openai.NewClient(apiKeys[0].Key)
			return AskChatGptWithImage(ctx, prompt, delay...)
		}
		if strings.Contains(err.Error(), "Please reduce your prompt") || strings.Contains(err.Error(), "Please reduce the length of the messages") {
			return "", ErrMaxTokens
		}
		if strings.Contains(err.Error(), "Incorrect API key") {
			return "", ErrIncorrectKey
		}
		if strings.Contains(err.Error(), "invalid character") {
			return "", ErrServiceUnavailable
		}
		return "", err
	}
	if len(resp.Data) == 0 {
		return "", ErrServiceUnavailable
	}
	return resp.Data[0].B64JSON, nil
}
