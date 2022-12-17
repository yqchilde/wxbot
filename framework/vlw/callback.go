package vlw

import (
	"fmt"
	"io"
	"net/http"

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
	http.HandleFunc("/wxbot/callback", func(w http.ResponseWriter, r *http.Request) {
		recv, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[VLW] 接收回调错误, error: %v", err)
			return
		}
		body := string(recv)
		event := robot.Event{RobotWxId: gjson.Get(body, "content.robot_wxid").String()}
		switch gjson.Get(body, "Event").String() {
		case eventPrivateChat:
			event.Type = robot.EventPrivateChat
			event.FromUniqueID = gjson.Get(body, "content.from_wxid").String()
			event.FromWxId = gjson.Get(body, "content.from_wxid").String()
			event.FromName = gjson.Get(body, "content.from_name").String()
			event.IsAtMe = true
			event.Message = &robot.Message{
				Id:      gjson.Get(body, "content.msg_id").String(),
				Type:    gjson.Get(body, "content.type").Int(),
				Content: gjson.Get(body, "content.msg").String(),
			}
		case eventGroupChat:
			event.Type = robot.EventGroupChat
			event.FromUniqueID = gjson.Get(body, "content.from_group").String()
			event.FromGroup = gjson.Get(body, "content.from_group").String()
			event.FromGroupName = gjson.Get(body, "content.from_group_name").String()
			event.FromWxId = gjson.Get(body, "content.from_wxid").String()
			event.FromName = gjson.Get(body, "content.from_name").String()
			if gjson.Get(body, "content.msg_source.atuserlist").Exists() {
				gjson.Get(body, "content.msg_source.atuserlist").ForEach(func(key, val gjson.Result) bool {
					if gjson.Get(val.String(), "wxid").String() == event.RobotWxId &&
						gjson.Get(val.String(), "nickname").String() != "@所有人" {
						event.IsAtMe = true
					}
					return true
				})
			}
			event.Message = &robot.Message{
				Id:      gjson.Get(body, "content.msg_id").String(),
				Type:    gjson.Get(body, "content.type").Int(),
				Content: gjson.Get(body, "content.msg").String(),
			}
		}
		handler(&event, f)
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	})
	log.Printf("[VLW] 回调地址: http://%s:%d/wxbot/callback", "127.0.0.1", f.ServePort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServePort), nil)
	if err != nil {
		log.Fatalf("[VLW] 回调服务启动失败, error: %v", err)
	}
}
