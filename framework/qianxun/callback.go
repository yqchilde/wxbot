package qianxun

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

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
	BotWxId  string // 机器人微信ID
	ApiUrl   string // http api地址
	ApiToken string // http api鉴权token
}

func New(botWxId, apiUrl, apiToken string) *Framework {
	return &Framework{
		BotWxId:  botWxId,
		ApiUrl:   apiUrl,
		ApiToken: apiToken,
	}
}

func (f *Framework) Callback(ctx *gin.Context, handler func(*robot.Event, robot.IFramework)) {
	recv, err := ctx.GetRawData()
	if err != nil {
		log.Errorf("[千寻] 接收回调错误, error: %v", err)
		return
	}
	handler(buildEvent(string(recv)), f)
	ctx.JSON(http.StatusOK, gin.H{"code": 0})
}

func buildEvent(resp string) *robot.Event {
	var event robot.Event
	switch gjson.Get(resp, "event").Int() {
	case eventAccountChange:
		// todo
	case eventGroupChat:
		switch gjson.Get(resp, "data.data.msgType").Int() {
		case 10000: // 系统消息
			event = robot.Event{
				Type: robot.EventSystem,
				Message: &robot.Message{
					Content: gjson.Get(resp, "data.data.msg").String(),
				},
			}
		case 49: // 群聊发app应用消息
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

			var refer ReferenceXml
			if err := xml.Unmarshal([]byte(gjson.Get(resp, "data.data.msg").String()), &refer); err == nil {
				if refer.Appmsg.Refermsg != nil { // 引用消息
					event.Message.Type = robot.MsgTypeText // 方便匹配
					event.Message.Content = refer.Appmsg.Title
					event.ReferenceMessage = &robot.ReferenceMessage{
						FromUser:    refer.Appmsg.Refermsg.Fromusr,
						ChatUser:    refer.Appmsg.Refermsg.Chatusr,
						DisplayName: refer.Appmsg.Refermsg.Displayname,
						Content:     refer.Appmsg.Refermsg.Content,
					}
				}
			}
		default:
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
			if gjson.Get(resp, fmt.Sprintf("data.data.atWxidList.#(==%s)", gjson.Get(resp, "wxid").String())).Exists() {
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
		switch gjson.Get(resp, "data.data.fromType").Int() {
		case 3: // 公众号
			event = robot.Event{
				Type:         robot.EventMPChat,
				FromUniqueID: gjson.Get(resp, "data.data.fromWxid").String(),
				FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
				FromName:     "",
				MPMessage: &robot.Message{
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
		default: // 私聊
			switch gjson.Get(resp, "data.data.msgType").Int() {
			case 49: // 私聊发app应用消息
				event = robot.Event{
					Type:         robot.EventPrivateChat,
					FromUniqueID: gjson.Get(resp, "data.data.fromWxid").String(),
					FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
					FromName:     "",
					Message: &robot.Message{
						Type:    gjson.Get(resp, "data.data.msgType").Int(),
						Content: gjson.Get(resp, "data.data.msg").String(),
					},
				}

				var refer ReferenceXml
				if err := xml.Unmarshal([]byte(gjson.Get(resp, "data.data.msg").String()), &refer); err == nil {
					if refer.Appmsg.Refermsg != nil { // 引用消息
						event.Message.Type = robot.MsgTypeText // 方便匹配
						event.Message.Content = refer.Appmsg.Title
						event.ReferenceMessage = &robot.ReferenceMessage{
							FromUser:    refer.Appmsg.Refermsg.Fromusr,
							ChatUser:    refer.Appmsg.Refermsg.Chatusr,
							DisplayName: refer.Appmsg.Refermsg.Displayname,
							Content:     refer.Appmsg.Refermsg.Content,
						}
					}
				}
			default:
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
		}
	case eventSelfMessage:
		event = robot.Event{
			Type:         robot.EventSelfMessage,
			FromUniqueID: gjson.Get(resp, "data.data.fromWxid").String(),
			FromWxId:     gjson.Get(resp, "data.data.fromWxid").String(),
			Message: &robot.Message{
				Type:    gjson.Get(resp, "data.data.msgType").Int(),
				Content: gjson.Get(resp, "data.data.msg").String(),
			},
		}
	case eventTransfer:
		event = robot.Event{
			Type: robot.EventTransfer,
			TransferMessage: &robot.TransferMessage{
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
				WithdrawMessage: &robot.WithdrawMessage{
					FromType:  fromType,
					FromWxId:  gjson.Get(resp, "data.data.fromWxid").String(),
					MsgSource: gjson.Get(resp, "data.data.msgSource").Int(),
					Msg:       gjson.Get(resp, "data.data.msg").String(),
				},
			}
		} else if fromType == 2 {
			event = robot.Event{
				Type: robot.EventMessageWithdraw,
				WithdrawMessage: &robot.WithdrawMessage{
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
			FriendVerifyMessage: &robot.FriendVerifyMessage{
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
