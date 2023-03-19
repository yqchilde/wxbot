package chatgpt

import (
	"errors"
	"fmt"
	"time"

	"github.com/sashabaranov/go-openai"

	"github.com/yqchilde/wxbot/engine/robot"
)

// 设置单次提问指令
func setSingleCommand(ctx *robot.Ctx, msg string, command string) {
	switch command {
	case "提问":
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
		answer = replaceSensitiveWords(answer)
		ctx.ReplyTextAndAt(fmt.Sprintf("问：%s \n--------------------\n答：%s", msg, answer))
	}
}
