package vlw

import (
	"fmt"
	"io"
	"net/http"

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
			log.Errorf("[VLW] 接收回调错误, error: %v", err)
			return
		}
		resp := string(recv)
		event := buildEvent(resp)
		handler(event, f)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	})
	if ip, err := net.GetIPWithLocal(); err != nil {
		log.Printf("[VLW] WxBot回调地址: http://%s:%d/wxbot/callback", "127.0.0.1", f.ServePort)
	} else {
		log.Printf("[VLW] WxBot回调地址: http://%s:%d/wxbot/callback", ip, f.ServePort)
	}
	err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServePort), nil)
	if err != nil {
		log.Fatalf("[VLW] WxBot回调服务启动失败, error: %v", err)
	}
}

func buildEvent(resp string) *robot.Event {
	event := robot.Event{RobotWxId: gjson.Get(resp, "content.robot_wxid").String()}
	switch gjson.Get(resp, "Event").String() {
	case eventGroupChat:
		// 根据消息类型细分多种事件
		switch gjson.Get(resp, "content.type").Int() {
		case 10000:
			event.Type = robot.EventSystem
			event.Message = &robot.Message{
				Content: gjson.Get(resp, "content.msg").String(),
			}
		default:
			event.Type = robot.EventGroupChat
			event.FromUniqueID = gjson.Get(resp, "content.from_group").String()
			event.FromGroup = gjson.Get(resp, "content.from_group").String()
			event.FromGroupName = gjson.Get(resp, "content.from_group_name").String()
			event.FromWxId = gjson.Get(resp, "content.from_wxid").String()
			event.FromName = gjson.Get(resp, "content.from_name").String()
			event.Message = &robot.Message{
				Id:      gjson.Get(resp, "content.msg_id").String(),
				Type:    gjson.Get(resp, "content.type").Int(),
				Content: gjson.Get(resp, "content.msg").String(),
			}
			if gjson.Get(resp, fmt.Sprintf("content.msg_source.atuserlist.#(wxid==%s)", event.RobotWxId)).Exists() {
				if !gjson.Get(resp, "content.msg_source.atuserlist.#(nickname==@所有人)").Exists() {
					event.IsAtMe = true
				}
			}
		}
	case eventPrivateChat:
		event.FromUniqueID = gjson.Get(resp, "content.from_wxid").String()
		event.FromWxId = gjson.Get(resp, "content.from_wxid").String()
		event.FromName = gjson.Get(resp, "content.from_name").String()

		// 根据消息类型细分多种事件
		switch gjson.Get(resp, "content.type").Int() {
		case 2000:
			event.Type = robot.EventTransfer
			event.Transfer = &robot.Transfer{
				FromWxId:   gjson.Get(resp, "content.from_wxid").String(),
				MsgSource:  gjson.Get(gjson.Get(resp, "content.msg").String(), "paysubtype").Int(),
				Money:      gjson.Get(gjson.Get(resp, "content.msg").String(), "money").String(),
				Memo:       gjson.Get(gjson.Get(resp, "content.msg").String(), "pay_momo").String(),
				TransferId: gjson.Get(gjson.Get(resp, "content.msg").String(), "payer_pay_id").String(),
			}
		default:
			event.Type = robot.EventPrivateChat
			event.IsAtMe = true
			event.Message = &robot.Message{
				Id:      gjson.Get(resp, "content.msg_id").String(),
				Type:    gjson.Get(resp, "content.type").Int(),
				Content: gjson.Get(resp, "content.msg").String(),
			}
		}
	case eventDeviceCallback:
		switch gjson.Get(resp, "content.type").Int() {
		case 1:
			event.Type = robot.EventSelfMessage // 可能不准确，待反馈
			event.Message = &robot.Message{
				Id:      gjson.Get(resp, "content.msg_id").String(),
				Type:    gjson.Get(resp, "content.type").Int(),
				Content: gjson.Get(resp, "content.msg").String(),
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
	return &event
}
