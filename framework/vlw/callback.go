package vlw

import (
	"fmt"
	"io"
	"net/http"
	"strings"

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
	ApiUrl    string // http api地址
	ApiToken  string // http api鉴权token
	ServePort uint   // 本地服务端口，用于接收回调
}

func New(apiUrl, apiToken string, servePort uint) *Framework {
	return &Framework{
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
		event := robot.Event{
			RobotWxId:     gjson.Get(body, "wxid").String(),
			IsPrivateChat: gjson.Get(body, "event").String() == eventPrivateChat,
			IsGroupChat:   gjson.Get(body, "event").String() == eventGroupChat,
			Message: robot.Message{
				Msg:      gjson.Get(body, "data.data.msg").String(),
				MsgId:    "",
				MsgType:  int(gjson.Get(body, "data.data.msgType").Int()),
				FromWxId: gjson.Get(body, "data.data.fromWxid").String(),
				FromName: "",
			},
		}
		if event.IsPrivateChat {
			event.IsAtMe = true
		}
		if event.IsGroupChat {
			event.Message.FromGroup = gjson.Get(body, "data.data.fromWxid").String()
			event.Message.FromGroupName = ""
			event.Message.FromWxId = gjson.Get(body, "data.data.finalFromWxid").String()
			gjson.Get(body, "data.data.atWxidList").ForEach(func(key, val gjson.Result) bool {
				if val.String() == event.RobotWxId && !strings.Contains(event.Message.Msg, "@所有人") {
					event.IsAtMe = true
				}
				return true
			})
		}
		handler(&event, f)

		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(`{"code":0}`))
	})
	err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServePort), nil)
	if err != nil {
		log.Fatalf("[VLW] 回调服务启动失败, error: %v", err)
	}
}
