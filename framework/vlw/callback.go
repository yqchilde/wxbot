package vlw

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
	eventGroupChat           = "EventGroupChat"           // 群聊消息事件
	eventPrivateChat         = "EventPrivateChat"         // 私聊消息事件
	eventDeviceCallback      = "EventDeviceCallback"      // 设备回调事件
	eventFriendVerify        = "EventFrieneVerify"        // 好友请求事件
	eventGroupNameChange     = "EventGroupNameChange"     // 群名称变动事件
	eventGroupMemberAdd      = "EventGroupMemberAdd"      // 群成员增加事件
	eventGroupMemberDecrease = "EventGroupMemberDecrease" // 群成员减少事件
	eventInvitedInGroup      = "EventInvitedInGroup"      // 被邀请入群事件
	eventQRCodePayment       = "EventQRcodePayment"       // 面对面收款事件
	eventDownloadFile        = "EventDownloadFile"        // 文件下载结束事件
	eventGroupEstablish      = "EventGroupEstablish"      // 创建新的群聊事件
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
			log.Errorf("[VLW] 接收回调错误, error: %v", err)
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
		log.Printf("[VLW] WxBot回调地址: http://%s:%d/wxbot/callback", "127.0.0.1", f.ServerPort)
	} else {
		log.Printf("[VLW] WxBot回调地址: http://%s:%d/wxbot/callback", ip, f.ServerPort)
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServerPort), nil)
	if err != nil {
		log.Fatalf("[VLW] WxBot回调服务启动失败, error: %v", err)
	}
}

func buildEvent(resp string, f *Framework) *robot.Event {
	var event robot.Event
	switch gjson.Get(resp, "Event").String() {
	case eventGroupChat:
		contentType := gjson.Get(resp, "content.type").Int()
		if contentType == 10000 {
			event = robot.Event{
				Type: robot.EventSystem,
				Message: &robot.Message{
					Content: gjson.Get(resp, "content.msg").String(),
				},
			}
		} else {
			event = robot.Event{
				Type:          robot.EventGroupChat,
				FromUniqueID:  gjson.Get(resp, "content.from_group").String(),
				FromGroup:     gjson.Get(resp, "content.from_group").String(),
				FromGroupName: gjson.Get(resp, "content.from_group_name").String(),
				FromWxId:      gjson.Get(resp, "content.from_wxid").String(),
				FromName:      gjson.Get(resp, "content.from_name").String(),
				Message: &robot.Message{
					Id:      gjson.Get(resp, "content.msg_id").String(),
					Type:    gjson.Get(resp, "content.type").Int(),
					Content: gjson.Get(resp, "content.msg").String(),
				},
			}
			if gjson.Get(resp, fmt.Sprintf("content.msg_source.atuserlist.#(wxid==%s)", event.RobotWxId)).Exists() {
				if !gjson.Get(resp, "content.msg_source.atuserlist.#(nickname==@所有人)").Exists() {
					event.IsAtMe = true
				}
			}
		}
	case eventPrivateChat:
		contentType := gjson.Get(resp, "content.type").Int()
		if contentType == 49 { // 分享卡片
			// 公众号处理 gh_开头
			if strings.HasPrefix(gjson.Get(resp, "content.from_wxid").String(), "gh_") {
				event = robot.Event{
					Type:         robot.EventMPChat,
					FromUniqueID: gjson.Get(resp, "content.from_wxid").String(),
					FromWxId:     gjson.Get(resp, "content.from_wxid").String(),
					FromName:     gjson.Get(resp, "content.from_name").String(),
					SubscriptionMessage: &robot.Message{
						Id:      gjson.Get(resp, "content.msg_id").String(),
						Type:    gjson.Get(resp, "content.type").Int(),
						Content: gjson.Get(resp, "content.msg").String(),
					},
				}
				for _, data := range robot.GetBot().MPsFromCache() {
					if data.WxId == event.FromWxId {
						event.FromName = data.Nick
						break
					}
				}
			}
		} else if contentType == 2000 { // 转账
			event = robot.Event{
				Type:         robot.EventTransfer,
				FromUniqueID: gjson.Get(resp, "content.from_wxid").String(),
				FromWxId:     gjson.Get(resp, "content.from_wxid").String(),
				FromName:     gjson.Get(resp, "content.from_name").String(),
				Transfer: &robot.Transfer{
					FromWxId:   gjson.Get(resp, "content.from_wxid").String(),
					MsgSource:  gjson.Get(gjson.Get(resp, "content.msg").String(), "paysubtype").Int(),
					Money:      gjson.Get(gjson.Get(resp, "content.msg").String(), "money").String(),
					Memo:       gjson.Get(gjson.Get(resp, "content.msg").String(), "pay_momo").String(),
					TransferId: gjson.Get(gjson.Get(resp, "content.msg").String(), "payer_pay_id").String(),
				},
			}
		} else { // 私聊
			event = robot.Event{
				Type:         robot.EventPrivateChat,
				IsAtMe:       true,
				FromUniqueID: gjson.Get(resp, "content.from_wxid").String(),
				FromWxId:     gjson.Get(resp, "content.from_wxid").String(),
				FromName:     gjson.Get(resp, "content.from_name").String(),
				Message: &robot.Message{
					Id:      gjson.Get(resp, "content.msg_id").String(),
					Type:    gjson.Get(resp, "content.type").Int(),
					Content: gjson.Get(resp, "content.msg").String(),
				},
			}
		}
	case eventDeviceCallback:
		contentType := gjson.Get(resp, "content.type").Int()
		if contentType == 1 {
			event = robot.Event{
				Type: robot.EventSelfMessage, // 可能不准确，待反馈
				Message: &robot.Message{
					Id:      gjson.Get(resp, "content.msg_id").String(),
					Type:    gjson.Get(resp, "content.type").Int(),
					Content: gjson.Get(resp, "content.msg").String(),
				},
			}
		}
	case eventFriendVerify:
	case eventGroupNameChange:
		// vlw框架有bug，将通过其他方式实现通知
	case eventGroupMemberAdd:
		event.Type = robot.EventGroupMemberIncrease
		// todo
	case eventGroupMemberDecrease:
		event.Type = robot.EventGroupMemberDecrease
		// todo
	case eventInvitedInGroup:
	case eventQRCodePayment:
	case eventDownloadFile:
	case eventGroupEstablish:
	}

	event.RobotWxId = gjson.Get(resp, "content.robot_wxid").String()
	event.RawMessage = resp
	return &event
}
