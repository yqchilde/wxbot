package emoticon

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/yqchilde/wxbot/engine"
	"github.com/yqchilde/wxbot/engine/robot"
	"github.com/yqchilde/wxbot/engine/util"
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

func (e *Emoticon) OnRegister() {
	err := os.MkdirAll(plugin.RawConfig.Get("dir").(string), os.ModePerm)
	if err != nil {
		panic("init img dir error: " + err.Error())
	}
}

func (e *Emoticon) OnEvent(msg *robot.Message) {
	if msg != nil {
		if msg.MatchTextCommand(pluginInfo.Commands) {
			if msg.IsSendByFriend() {
				sender, err := msg.Sender()
				if err != nil {
					log.Printf("handleMessage get sender error: %v", err)
					return
				}
				if addCommand(sender.UserName, msg.Content) {
					return
				}

				msg.ReplyText(getAtMessage(sender.NickName, "请在10秒内发送表情图"))
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				go waitEmoticon(ctx, cancel, msg, sender)

			} else {
				sender, err := msg.SenderInGroup()
				if err != nil {
					log.Printf("handleMessage get sender error: %v", err)
					return
				} else {
					if addCommand(sender.UserName, msg.Content) {
						return
					}
				}

				msg.ReplyText(getAtMessage(sender.NickName, "请在10秒内发送表情图"))
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				go waitEmoticon(ctx, cancel, msg, sender)
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

func getAtMessage(nickName, content string) string {
	return fmt.Sprintf("@%s\u2005%s", nickName, content)
}

func waitEmoticon(ctx context.Context, cancel context.CancelFunc, msg *robot.Message, sender *robot.User) {
	defer func() {
		cancel()
		removeCommand(sender.UserName)
	}()

	for {
		select {
		case <-ctx.Done():
			_, err := msg.ReplyText(getAtMessage(sender.NickName, "操作超时，请重新输入命令"))
			if err != nil {
				log.Printf("handleMessage reply error: %v", err)
			}
			return
		case msg := <-waitCommand:
			emoticon, err := robot.UnMarshalForEmoticon(msg.Content)
			if err != nil {
				log.Printf("waitEmoticon UnMarshalForEmoticon error: %v", err)
				return
			}
			emoticonUrl := emoticon.Emoji.Cdnurl
			msg.ReplyText(getAtMessage(sender.NickName, "表情包原图如下"))
			fileName := fmt.Sprintf("%s/%s", plugin.RawConfig.Get("dir"), time.Now().Format("20060102150405"))
			fileName, err = util.SingleDownload(util.ImgInfo{Url: emoticonUrl, Name: fileName})
			if err != nil {
				log.Printf("Failed to download emoticon, err: %v", err)
				return
			}

			emoticonFile, err := os.Open(fileName)
			if err != nil {
				log.Println(err)
				return
			}
			if filepath.Ext(fileName) == ".gif" {
				if _, err := msg.ReplyFile(emoticonFile); err != nil {
					if strings.Contains(err.Error(), "operate too often") {
						msg.ReplyText("Warn: 被微信ban了，请稍后再试")
					} else {
						log.Printf("msg.ReplyImage reply image error: %v", err)
					}
				}
			} else {
				if _, err := msg.ReplyImage(emoticonFile); err != nil {
					if strings.Contains(err.Error(), "operate too often") {
						msg.ReplyText("Warn: 被微信ban了，请稍后再试")
					} else {
						log.Printf("msg.ReplyImage reply image error: %v", err)
					}
				}
			}
			emoticonFile.Close()
			os.Remove(fileName)
			return
		}
	}
}
