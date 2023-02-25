package chatgpt

import (
	"context"
	"errors"
	"strings"
	"time"

	gogpt "github.com/sashabaranov/go-gpt3"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

var (
	gptClient *gogpt.Client
	gptModel  *GptModel
)

func AskChatGpt(prompt string, delay ...time.Duration) (answer string, err error) {
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

	// 请求gpt3
	resp, err := gptClient.CreateCompletion(context.Background(), gogpt.CompletionRequest{
		Model:            gptModel.Model,
		Prompt:           prompt,
		MaxTokens:        gptModel.MaxTokens,
		Temperature:      float32(gptModel.Temperature),
		TopP:             float32(gptModel.TopP),
		PresencePenalty:  float32(gptModel.PresencePenalty),
		FrequencyPenalty: float32(gptModel.FrequencyPenalty),
		Echo:             false,
		Stop:             []string{"Human:", "AI:"},
	})

	// 处理响应回来的错误
	if err != nil {
		if strings.Contains(err.Error(), "You exceeded your current quota") {
			log.Printf("当前apiKey[%s]配额已用完, 将删除并切换到下一个", apiKeys[0].Key)
			db.Orm.Table("apikey").Where("key = ?", apiKeys[0].Key).Delete(&ApiKey{})
			if len(apiKeys) == 1 {
				return "", errors.New("OpenAi配额已用完，请联系管理员")
			}
			apiKeys = apiKeys[1:]
			gptClient = gogpt.NewClient(apiKeys[0].Key)
			return AskChatGpt(prompt)
		}
		if strings.Contains(err.Error(), "The server had an error while processing your request") {
			log.Println("OpenAi服务出现问题，将重试")
			return AskChatGpt(prompt)
		}
		if strings.Contains(err.Error(), "Client.Timeout exceeded while awaiting headers") {
			log.Println("OpenAi服务请求超时，将重试")
			return AskChatGpt(prompt)
		}
		if strings.Contains(err.Error(), "Please reduce your prompt") {
			return "", errors.New("OpenAi免费上下文长度限制为4097个词组，您的上下文长度已超出限制，请发送\"清空会话\"以清空上下文")
		}
		if strings.Contains(err.Error(), "Incorrect API key") {
			return "", errors.New("OpenAi ApiKey错误，请联系管理员")
		}
		return "", err
	}
	return resp.Choices[0].Text + "\n", nil
}

// filterAnswer 过滤答案，处理一些符号问题
// return 新的答案，是否需要重试
func filterAnswer(answer string) (newAnswer string, isNeedRetry bool) {
	punctuation := ",，!！?？"
	answer = strings.TrimSpace(answer)
	if len(answer) == 1 {
		return answer, true
	}
	if len(answer) == 3 && strings.ContainsAny(answer, punctuation) {
		return answer, true
	}
	answer = strings.TrimLeftFunc(answer, func(r rune) bool {
		if strings.ContainsAny(string(r), punctuation) {
			return true
		}
		return false
	})
	return answer, false
}

func AskChatGptWithImage(prompt string, delay ...time.Duration) (b64 string, err error) {
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

	resp, err := gptClient.CreateImage(context.Background(), gogpt.ImageRequest{
		Prompt:         prompt,
		Size:           gptModel.ImageSize,
		ResponseFormat: gogpt.CreateImageResponseFormatB64JSON,
	})
	if err != nil {
		return "", err
	}
	return resp.Data[0].B64JSON, nil
}
