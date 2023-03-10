package vlw

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/tidwall/gjson"

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
		log.Errorf("[VLW] 接收回调错误, error: %v", err)
		return
	}
	handler(buildEvent(string(recv)), f)
	ctx.JSON(http.StatusOK, gin.H{"code": 0})
}

func buildEvent(resp string) *robot.Event {
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
				Type:           robot.EventGroupChat,
				FromUniqueID:   gjson.Get(resp, "content.from_group").String(),
				FromUniqueName: gjson.Get(resp, "content.from_group_name").String(),
				FromGroup:      gjson.Get(resp, "content.from_group").String(),
				FromGroupName:  gjson.Get(resp, "content.from_group_name").String(),
				FromWxId:       gjson.Get(resp, "content.from_wxid").String(),
				FromName:       gjson.Get(resp, "content.from_name").String(),
				Message: &robot.Message{
					Id:      gjson.Get(resp, "content.msg_id").String(),
					Type:    gjson.Get(resp, "content.type").Int(),
					Content: gjson.Get(resp, "content.msg").String(),
				},
			}
			if gjson.Get(resp, fmt.Sprintf("content.msg_source.atuserlist.#(wxid==%s)", gjson.Get(resp, "content.robot_wxid").String())).Exists() {
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
					Type:           robot.EventMPChat,
					FromUniqueID:   gjson.Get(resp, "content.from_wxid").String(),
					FromUniqueName: gjson.Get(resp, "content.from_name").String(),
					FromWxId:       gjson.Get(resp, "content.from_wxid").String(),
					FromName:       gjson.Get(resp, "content.from_name").String(),
					MPMessage: &robot.Message{
						Id:      gjson.Get(resp, "content.msg_id").String(),
						Type:    gjson.Get(resp, "content.type").Int(),
						Content: gjson.Get(resp, "content.msg").String(),
					},
				}
				for _, data := range robot.GetBot().MPs() {
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
				TransferMessage: &robot.TransferMessage{
					FromWxId:   gjson.Get(resp, "content.from_wxid").String(),
					MsgSource:  gjson.Get(gjson.Get(resp, "content.msg").String(), "paysubtype").Int(),
					Money:      gjson.Get(gjson.Get(resp, "content.msg").String(), "money").String(),
					Memo:       gjson.Get(gjson.Get(resp, "content.msg").String(), "pay_momo").String(),
					TransferId: gjson.Get(gjson.Get(resp, "content.msg").String(), "payer_pay_id").String(),
				},
			}
		} else { // 私聊
			event = robot.Event{
				Type:           robot.EventPrivateChat,
				IsAtMe:         true,
				FromUniqueID:   gjson.Get(resp, "content.from_wxid").String(),
				FromUniqueName: gjson.Get(resp, "content.from_name").String(),
				FromWxId:       gjson.Get(resp, "content.from_wxid").String(),
				FromName:       gjson.Get(resp, "content.from_name").String(),
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
				Type:         robot.EventSelfMessage, // 可能不准确，待反馈
				FromUniqueID: gjson.Get(resp, "content.from_wxid").String(),
				FromWxId:     gjson.Get(resp, "content.from_wxid").String(),
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
