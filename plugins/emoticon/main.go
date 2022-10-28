package emoticon

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
)

type Emoticon struct {
	engine.PluginMagic
	Enable bool   `yaml:"enable"`
	Dir    string `yaml:"dir"`
}

var (
	pluginInfo = &Emoticon{
		PluginMagic: engine.PluginMagic{
			Desc:     "🚀 输入 {表情原图} => 10s内发送表情获取表情原图",
			Commands: []string{"表情原图"},
		},
	}
	plugin      = engine.InstallPlugin(pluginInfo)
	userCommand = make(map[string]string) // 用户指令 key:username val:command
	waitCommand = make(chan *robot.Message)
	mutex       sync.Mutex
)

func (e *Emoticon) OnRegister() {}

func (e *Emoticon) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchTextCommand(pluginInfo.Commands) {
			if addCommand(msg.Content.FromWxid, msg.Content.Msg) {
				return
			}

			if msg.IsSendByPrivateChat() {
				msg.ReplyText("请在10s内发送表情获取表情原图")
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				go waitEmoticon(ctx, cancel, msg)
			} else if msg.IsSendByGroupChat() {
				msg.ReplyTextAndAt("请在10s内发送表情获取表情原图")
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				go waitEmoticon(ctx, cancel, msg)
			}

		}

		if msg.IsEmoticon() {
			for i := range userCommand {
				for j := range pluginInfo.Commands {
					if userCommand[i] == pluginInfo.Commands[j] {
						waitCommand <- msg
						break
					}
				}
			}
		}
	}
}

// 添加用户指令
func addCommand(sender, command string) bool {
	mutex.Lock()
	defer mutex.Unlock()

	if val, ok := userCommand[sender]; ok && val == command {
		return true
	} else {
		userCommand[sender] = command
		return false
	}
}

// 移除用户指令
func removeCommand(sender string) {
	mutex.Lock()
	defer mutex.Unlock()

	delete(userCommand, sender)
}

func waitEmoticon(ctx context.Context, cancel context.CancelFunc, msg *robot.Message) {
	defer func() {
		cancel()
		removeCommand(msg.Content.FromWxid)
	}()

	for {
		select {
		case <-ctx.Done():
			fmt.Println("waitEmoticon timeout")
			if msg.IsSendByPrivateChat() {
				msg.ReplyText("10s内未发送表情，获取表情原图失败")
			} else if msg.IsSendByGroupChat() {
				msg.ReplyTextAndAt("10s内未发送表情，获取表情原图失败")
			}
			return
		case msg := <-waitCommand:
			emoticonUrl := msg.Content.Msg[5 : len(msg.Content.Msg)-1]
			if err := msg.ReplyImage(emoticonUrl); err != nil {
				msg.ReplyText(err.Error())
			}
			return
		}
	}
}
