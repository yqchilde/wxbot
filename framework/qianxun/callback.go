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
	ServerPort uint   // 本地服务端口，用于接收回调
	BotWxId    string // 机器人微信ID
	ApiUrl     string // http api地址
	ApiToken   string // http api鉴权token
}

func New(serverPort uint, botWxId, apiUrl, apiToken string) *Framework {
	return &Framework{
		ServerPort: serverPort,
		BotWxId:    botWxId,
		ApiUrl:     apiUrl,
		ApiToken:   apiToken,
	}
}

func (f *Framework) Callback(handler func(*robot.Event, robot.IFramework)) {
	http.HandleFunc("/wxbot/callback", func(w http.ResponseWriter, r *http.Request) {
		recv, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[千寻] 接收回调错误, error: %v", err)
			return
		}
		resp := string(recv)
		event := buildEvent(resp, f)
		handler(event, f)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	})
	if f.ServerPort == 0 {
		f.ServerPort = 9528
	}

	if ip, err := net.GetIPWithLocal(); err != nil {
		log.Printf("[千寻] WxBot回调地址: http://%s:%d/wxbot/callback", "127.0.0.1", f.ServerPort)
	} else {
		log.Printf("[千寻] WxBot回调地址: http://%s:%d/wxbot/callback", ip, f.ServerPort)
	}
	if err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServerPort), nil); err != nil {
		log.Fatalf("[千寻] WxBot回调服务启动失败, error: %v", err)
	}
}

func buildEvent(resp string, f *Framework) *robot.Event {
	var event robot.Event
	switch gjson.Get(resp, "event").Int() {
	case eventAccountChange:
		// todo
	case eventGroupChat:
		msgType := gjson.Get(resp, "data.data.msgType").Int()
		if msgType == 10000 { // 系统消息
			event = robot.Event{
				Type: robot.EventSystem,
				Message: &robot.Message{
					Content: gjson.Get(resp, "data.data.msg").String(),
				},
			}
		} else { // 群聊
			event = robot.Event{
				Type:          robot.EventGroupChat,
				FromUniqueID:  gjson.Get(resp, "data.data.fromWxid").String(),
				FromGroup:     gjson.Get(resp, "data.data.fromWxid").String(),
				FromGroupName: "",
				FromWxId:      gjson.Get(resp, "data.data.finalFromWxid").String(),
				FromName:      "",
				Message: &robot.Message{
					Type:    gjson.Get(resp, "data.data.msgType").Int(),
					Content: gjson.Get(resp, "data.data.msg").String(),
				},
			}
			if gjson.Get(resp, fmt.Sprintf("data.data.atWxidList.#(==%s)", event.RobotWxId)).Exists() {
				if !strings.Contains(event.Message.Content, "@所有人") {
					event.IsAtMe = true
				}
			}
			for _, data := range robot.GetBot().Groups() {
				if data.WxId == event.FromGroup {
					event.FromGroupName = data.Nick
					event.FromUniqueName = data.Nick
					break
				}
			}
		}
	case eventPrivateChat:
		fromType := gjson.Get(resp, "data.data.fromType").Int()
		if fromType == 3 { // 公众号
			event = robot.Event{
				Type:         robot.EventMPChat,
				FromUniqueID: gjson.Get(resp, "data.data.fromWxid").String(),
				FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
				FromName:     "",
				SubscriptionMessage: &robot.Message{
					Type:    gjson.Get(resp, "data.data.msgType").Int(),
					Content: gjson.Get(resp, "data.data.msg").String(),
				},
			}
			for _, data := range robot.GetBot().MPs() {
				if data.WxId == event.FromWxId {
					event.FromName = data.Nick
					event.FromUniqueName = data.Nick
					break
				}
			}
		} else { // 私聊
			event = robot.Event{
				Type:         robot.EventPrivateChat,
				FromUniqueID: gjson.Get(resp, "data.data.fromWxid").String(),
				FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
				FromName:     "",
				IsAtMe:       true,
				Message: &robot.Message{
					Type:    gjson.Get(resp, "data.data.msgType").Int(),
					Content: gjson.Get(resp, "data.data.msg").String(),
				},
			}
			for _, data := range robot.GetBot().Friends() {
				if data.WxId == event.FromWxId {
					event.FromName = data.Nick
					event.FromUniqueName = data.Nick
					break
				}
			}
		}
	case eventSelfMessage:
		event = robot.Event{
			Type: robot.EventSelfMessage,
			Message: &robot.Message{
				Type:    gjson.Get(resp, "data.data.msgType").Int(),
				Content: gjson.Get(resp, "data.data.msg").String(),
			},
		}
	case eventTransfer:
		event = robot.Event{
			Type: robot.EventTransfer,
			Transfer: &robot.Transfer{
				FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
				MsgSource:    gjson.Get(resp, "data.data.msgSource").Int(),
				TransferType: gjson.Get(resp, "data.data.transType").Int(),
				Money:        gjson.Get(resp, "data.data.money").String(),
				Memo:         gjson.Get(resp, "data.data.memo").String(),
				TransferId:   gjson.Get(resp, "data.data.transferid").String(),
				TransferTime: gjson.Get(resp, "data.data.invalidtime").String(),
			},
		}
	case eventMessageWithdraw:
		fromType := gjson.Get(resp, "data.data.fromType").Int()
		if fromType == 1 {
			event = robot.Event{
				Type: robot.EventMessageWithdraw,
				Withdraw: &robot.Withdraw{
					FromType:  fromType,
					FromWxId:  gjson.Get(resp, "data.data.fromWxid").String(),
					MsgSource: gjson.Get(resp, "data.data.msgSource").Int(),
					Msg:       gjson.Get(resp, "data.data.msg").String(),
				},
			}
		} else if fromType == 2 {
			event = robot.Event{
				Type: robot.EventMessageWithdraw,
				Withdraw: &robot.Withdraw{
					FromType:  fromType,
					FromGroup: gjson.Get(resp, "data.data.fromWxid").String(),
					FromWxId:  gjson.Get(resp, "data.data.finalFromWxid").String(),
					MsgSource: gjson.Get(resp, "data.data.msgSource").Int(),
					Msg:       gjson.Get(resp, "data.data.msg").String(),
				},
			}
		}
	case eventFriendVerify:
		event = robot.Event{
			Type: robot.EventFriendVerify,
			FriendVerify: &robot.FriendVerify{
				WxId:      gjson.Get(resp, "data.data.wxid").String(),
				Nick:      gjson.Get(resp, "data.data.nick").String(),
				V3:        gjson.Get(resp, "data.data.v3").String(),
				V4:        gjson.Get(resp, "data.data.v4").String(),
				AvatarUrl: gjson.Get(resp, "data.data.avatarMinUrl").String(),
				Content:   gjson.Get(resp, "data.data.content").String(),
				Scene:     gjson.Get(resp, "data.data.scene").String(),
			},
		}
	case eventPayment:
		// todo
	}

	event.RobotWxId = gjson.Get(resp, "wxid").String()
	event.RawMessage = resp
	return &event
}
