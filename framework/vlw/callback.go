package vlw

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/tidwall/gjson"
	"github.com/yqchilde/pkgs/log"

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

func (f *Framework) Callback(handler func(*robot.Event, robot.APICaller)) {
	http.HandleFunc("/wxbot/callback", func(w http.ResponseWriter, r *http.Request) {
		recv, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[VLW] 接收回调错误, error: %v", err)
			return
		}
		body := string(recv)
		os.WriteFile("callback.json", recv, 0644)
		event := robot.Event{
			RobotWxId:     gjson.Get(body, "content.robot_wxid").String(),
			IsAtMe:        gjson.Get(body, "Event").String() == eventPrivateChat,
			IsPrivateChat: gjson.Get(body, "Event").String() == eventPrivateChat,
			IsGroupChat:   gjson.Get(body, "Event").String() == eventGroupChat,
			Message: robot.Message{
				Msg:      gjson.Get(body, "content.msg").String(),
				MsgId:    gjson.Get(body, "content.msg_id").String(),
				MsgType:  gjson.Get(body, "content.type").Int(),
				FromWxId: gjson.Get(body, "content.from_wxid").String(),
				FromName: gjson.Get(body, "content.from_name").String(),
			},
		}
		if event.IsGroupChat {
			event.Message.FromGroup = gjson.Get(body, "content.from_group").String()
			event.Message.FromGroupName = gjson.Get(body, "content.from_group_name").String()
			if gjson.Get(body, "content.msg_source.atuserlist").Exists() {
				gjson.Get(body, "content.msg_source.atuserlist").ForEach(func(key, val gjson.Result) bool {
					if gjson.Get(val.String(), "wxid").String() == event.RobotWxId &&
						gjson.Get(val.String(), "nickname").String() != "@所有人" {
						event.IsAtMe = true
					}
					return true
				})
			}
		}
		handler(&event, f)

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	})
	log.Printf("[VLW] 回调地址, http://%s:%d/wxbot/callback", "127.0.0.1", f.ServePort)
	err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServePort), nil)
	if err != nil {
		log.Fatalf("[VLW] 回调服务启动失败, error: %v", err)
	}
}
