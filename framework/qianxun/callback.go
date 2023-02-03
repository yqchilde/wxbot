package qianxun

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/yqchilde/pkgs/net"
	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/robot"
)

const (
	eventAccountChange   = 10014 // 账号变动事件
	eventGroupChat       = 10008 // 群聊消息事件
	eventPrivateChat     = 10009 // 私聊消息事件
	eventSelfMessage     = 10010 // 自己消息事件
	eventTransfer        = 10006 // 收到转账事件
	eventMessageWithdraw = 10013 // 消息撤回事件
	eventFriendVerify    = 10011 // 好友请求事件
	eventPayment         = 10007 // 收到支付事件
)

type Framework struct {
	BotWxId   string // 机器人微信ID
	ApiUrl    string // http api地址
	ApiToken  string // http api鉴权token
	ServePort uint   // 本地服务端口，用于接收回调
}

func New(botWxId, apiUrl, apiToken string, servePort uint) *Framework {
	return &Framework{
		BotWxId:   botWxId,
		ApiUrl:    apiUrl,
		ApiToken:  apiToken,
		ServePort: servePort,
	}
}

func (f *Framework) Callback(handler func(*robot.Event, robot.IFramework)) {
	robot.Framework.Store(f)
	http.HandleFunc("/wxbot/callback", func(w http.ResponseWriter, r *http.Request) {
		recv, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[千寻] 接收回调错误, error: %v", err)
			return
		}
		resp := string(recv)
		event := buildEvent(resp)
		handler(event, f)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	})
	if f.ServePort == 0 {
		f.ServePort = 9528
	}

	if ip, err := net.GetIPWithLocal(); err != nil {
		log.Printf("[千寻] WxBot回调地址: http://%s:%d/wxbot/callback", "127.0.0.1", f.ServePort)
	} else {
		log.Printf("[千寻] WxBot回调地址: http://%s:%d/wxbot/callback", ip, f.ServePort)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServePort), nil); err != nil {
		log.Fatalf("[千寻] WxBot回调服务启动失败, error: %v", err)
	}
}

func buildEvent(resp string) *robot.Event {
	event := robot.Event{RobotWxId: gjson.Get(resp, "wxid").String()}
	switch gjson.Get(resp, "event").Int() {
	case eventAccountChange:
		// todo
	case eventGroupChat:
		switch gjson.Get(resp, "data.data.msgType").Int() {
		case 10000:
			event.Type = robot.EventSystem
			event.Message = &robot.Message{
				Content: gjson.Get(resp, "data.data.msg").String(),
			}
		default:
			event.Type = robot.EventGroupChat
			event.FromUniqueID = gjson.Get(resp, "data.data.fromWxid").String()
			event.FromGroup = gjson.Get(resp, "data.data.fromWxid").String()
			event.FromGroupName = ""
			event.FromWxId = gjson.Get(resp, "data.data.finalFromWxid").String()
			event.FromName = ""
			event.Message = &robot.Message{
				Type:    gjson.Get(resp, "data.data.msgType").Int(),
				Content: gjson.Get(resp, "data.data.msg").String(),
			}
			if gjson.Get(resp, fmt.Sprintf("data.data.atWxidList.#(==%s)", event.RobotWxId)).Exists() {
				if !strings.Contains(event.Message.Content, "@所有人") {
					event.IsAtMe = true
				}
			}
		}
	case eventPrivateChat:
		event.Type = robot.EventPrivateChat
		event.FromUniqueID = gjson.Get(resp, "data.data.fromWxid").String()
		event.FromWxId = gjson.Get(resp, "data.data.fromWxid").String()
		event.FromName = ""
		event.IsAtMe = true
		event.Message = &robot.Message{
			Type:    gjson.Get(resp, "data.data.msgType").Int(),
			Content: gjson.Get(resp, "data.data.msg").String(),
		}
	case eventSelfMessage:
		event.Type = robot.EventSelfMessage
		event.Message = &robot.Message{
			Type:    gjson.Get(resp, "data.data.msgType").Int(),
			Content: gjson.Get(resp, "data.data.msg").String(),
		}
	case eventTransfer:
		event.Type = robot.EventTransfer
		event.Transfer = &robot.Transfer{
			FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
			MsgSource:    gjson.Get(resp, "data.data.msgSource").Int(),
			TransferType: gjson.Get(resp, "data.data.transType").Int(),
			Money:        gjson.Get(resp, "data.data.money").String(),
			Memo:         gjson.Get(resp, "data.data.memo").String(),
			TransferId:   gjson.Get(resp, "data.data.transferid").String(),
			TransferTime: gjson.Get(resp, "data.data.invalidtime").String(),
		}
	case eventMessageWithdraw:
		event.Type = robot.EventMessageWithdraw
		event.Withdraw = &robot.Withdraw{
			FromType:  gjson.Get(resp, "data.data.fromType").Int(),
			MsgSource: gjson.Get(resp, "data.data.msgSource").Int(),
			Msg:       gjson.Get(resp, "data.data.msg").String(),
		}
		if event.Withdraw.FromType == 1 {
			event.Withdraw.FromWxId = gjson.Get(resp, "data.data.fromWxid").String()
		} else if event.Withdraw.FromType == 2 {
			event.Withdraw.FromGroup = gjson.Get(resp, "data.data.fromWxid").String()
			event.Withdraw.FromWxId = gjson.Get(resp, "data.data.finalFromWxid").String()
		}
	case eventFriendVerify:
		event.Type = robot.EventFriendVerify
		event.FriendVerify = &robot.FriendVerify{
			WxId:      gjson.Get(resp, "data.data.wxid").String(),
			Nick:      gjson.Get(resp, "data.data.nick").String(),
			V3:        gjson.Get(resp, "data.data.v3").String(),
			V4:        gjson.Get(resp, "data.data.v4").String(),
			AvatarUrl: gjson.Get(resp, "data.data.avatarMinUrl").String(),
			Content:   gjson.Get(resp, "data.data.content").String(),
			Scene:     gjson.Get(resp, "data.data.scene").String(),
		}
	case eventPayment:
		// todo
	}
	return &event
}
