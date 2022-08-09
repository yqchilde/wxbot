package message

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/eatmoreapple/openwechat"
	"github.com/yqchilde/pkgs/log"

	"wxBot/internal/config"
	"wxBot/internal/model"
	"wxBot/internal/pkg/download"
	"wxBot/internal/pkg/holiday"
	"wxBot/internal/service"
)

type Message struct {
	sync.Mutex

	// 用户指令 key:username val:command
	users map[string]string
}

var (
	messageObj  *Message
	waitCommand = make(chan *openwechat.Message)
)

func NewMessage() {
	messageObj = &Message{users: make(map[string]string)}
}

func HandleMessage(msg *openwechat.Message) { messageObj.handleMessage(msg) }
func (m *Message) handleMessage(msg *openwechat.Message) {
	// 忽略自己的消息
	if msg.IsSendBySelf() {
		return
	}

	if msg.IsText() {
		// 分析指令
		command := msg.Content
		command = strings.ReplaceAll(command, "\n", "")
		command = strings.TrimLeft(command, " ")
		command = strings.TrimRight(command, " ")
		if !strings.HasPrefix(command, "/") {
			return
		}

		// 群聊
		if msg.IsSendByGroup() {
			// 根据用户存指令
			sender, err := msg.SenderInGroup()
			if err != nil {
				log.Errorf("handleMessage get sender error: %v", err)
				return
			} else {
				if m.addCommand(sender.UserName, command) {
					return
				}
			}

			log.Printf("listen groupChat command: %s", command)
			switch command {
			case "/img":
				msg.ReplyText(m.getAtMessage(sender.NickName, "请在10秒内发送表情图"))
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				go m.waitEmoticon(ctx, cancel, msg, sender)
			case "/plmm":
				service.GetPlmmPhoto(msg)
				m.removeCommand(sender.UserName)
			case "/sj":
				service.GetShaoJiPhoto(msg)
				m.removeCommand(sender.UserName)
			case "/myb":
				if notes, err := holiday.DailyLifeNotes(); err == nil {
					msg.ReplyText(notes)
				}
			case "/menu":
				reply := showGroupChatMenu()
				msg.ReplyText(reply)
			}
		}

		// 单聊
		if msg.IsSendByFriend() {
			sender, err := msg.Sender()
			if err != nil {
				log.Errorf("handleMessage get sender error: %v", err)
				return
			}

			log.Printf("listen singleChat command: %s", command)
			switch command {
			case "/img":
				msg.ReplyText(m.getAtMessage(sender.NickName, "请在10秒内发送表情图"))
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				go m.waitEmoticon(ctx, cancel, msg, sender)
			case "/myb":
				if notes, err := holiday.DailyLifeNotes(); err == nil {
					msg.ReplyText(notes)
				}
			case "/menu":
				reply := showSingleChatMenu()
				msg.ReplyText(reply)
			}
		}
	}

	if msg.IsEmoticon() {
		for _, command := range m.users {
			if command == "/img" {
				waitCommand <- msg
			}
		}
	}
}

// 添加用户指令
func (m *Message) addCommand(sender, command string) bool {
	m.Lock()
	defer m.Unlock()

	if val, ok := m.users[sender]; ok && val == command {
		return true
	} else {
		m.users[sender] = command
		return false
	}
}

// 移除用户指令
func (m *Message) removeCommand(sender string) {
	m.Lock()
	defer m.Unlock()

	delete(m.users, sender)
}

// 打印at消息内容
func (m *Message) getAtMessage(nickName, content string) string {
	return fmt.Sprintf("@%s\u2005%s", nickName, content)
}

// 等待收到emoticon
func (m *Message) waitEmoticon(ctx context.Context, cancel context.CancelFunc, msg *openwechat.Message, sender *openwechat.User) {
	defer func() {
		cancel()
		m.removeCommand(sender.UserName)
	}()

	for {
		select {
		case <-ctx.Done():
			msg.ReplyText(m.getAtMessage(sender.NickName, "操作超时，请重新输入命令"))
			return
		case msg := <-waitCommand:
			emoticon, err := UnMarshalForEmoticon(msg.Content)
			if err != nil {
				log.Errorf("waitEmoticon UnMarshalForEmoticon error: %v", err)
				return
			}
			emoticonUrl := emoticon.Emoji.Cdnurl
			msg.ReplyText(m.getAtMessage(sender.NickName, "表情包原图如下"))
			fileName := fmt.Sprintf("%s/%s", config.GetEmoticonConf().Dir, time.Now().Format("20060102150405"))
			fileName, err = download.SingleDownload(model.ImgInfo{Url: emoticonUrl, Name: fileName})
			if err != nil {
				log.Errorf("Failed to download emoticon, err: %v", err)
				return
			}

			emoticonFile, err := os.Open(fileName)
			if err != nil {
				log.Error(err)
				return
			}
			if filepath.Ext(fileName) == ".gif" {
				msg.ReplyFile(emoticonFile)
			} else {
				msg.ReplyImage(emoticonFile)
			}
			emoticonFile.Close()
			os.Remove(fileName)
			return
		}
	}
}

// 单聊菜单
func showSingleChatMenu() string {
	command := `Bug Bot🤖
				🚀 输入 /img => 10s内发送表情可收货表情原图
				🚀 输入 /myb => 获取摸鱼办消息
				🚀 单聊

				- - - - - - - - - - - - - - - - - - - - - 
				👴🏻?? 可收货漂亮妹妹`
	command = strings.ReplaceAll(command, "\t", "")
	return command
}

// 群聊菜单
func showGroupChatMenu() string {
	command := `Bug Bot🤖
				🚀 输入 /img => 10s内发送表情可收货表情原图
				🚀 输入 /plmm => 可收货漂亮妹妹
				🚀 输入 /myb => 获取摸鱼办消息
				🚀 输入 /sj => 可收货🔥🐔

				- - - - - - - - - - - - - - - - - - - - - 
				👴🏻?? 可收货漂亮妹妹`
	command = strings.ReplaceAll(command, "\t", "")
	return command
}
