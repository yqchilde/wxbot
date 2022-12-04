package qianxun

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/yqchilde/pkgs/log"

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

func (f *Framework) Callback(handler func(*robot.Event, robot.APICaller)) {
	http.HandleFunc("/wxbot/callback", func(w http.ResponseWriter, r *http.Request) {
		recv, err := io.ReadAll(r.Body)
		if err != nil {
			log.Errorf("[千寻] 接收回调错误, error: %v", err)
			return
		}
		body := string(recv)
		event := robot.Event{
			RobotWxId:     gjson.Get(body, "wxid").String(),
			IsAtMe:        gjson.Get(body, "event").String() == strconv.Itoa(eventPrivateChat),
			IsPrivateChat: gjson.Get(body, "event").String() == strconv.Itoa(eventPrivateChat),
			IsGroupChat:   gjson.Get(body, "event").String() == strconv.Itoa(eventGroupChat),
			Message: robot.Message{
				Msg:      gjson.Get(body, "data.data.msg").String(),
				MsgId:    "",
				MsgType:  gjson.Get(body, "data.data.msgType").Int(),
				FromWxId: gjson.Get(body, "data.data.fromWxid").String(),
				FromName: "",
			},
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
	if f.ServePort == 0 {
		f.ServePort = 9528
	}
	log.Printf("[千寻] 回调地址, http://%s:%d/wxbot/callback", "127.0.0.1", f.ServePort)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", f.ServePort), nil); err != nil {
		log.Fatalf("[千寻] 回调服务启动失败, error: %v", err)
	}
}
