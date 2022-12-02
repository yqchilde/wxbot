package robot

// Event 记录一次回调事件
type Event struct {
	RobotWxId     string  // 机器人微信id
	IsAtMe        bool    // 机器人是否被@了，@所有人不算
	IsPrivateChat bool    // 是否是私聊消息
	IsGroupChat   bool    // 是否是群聊消息
	Message       Message // 消息内容
}

// Message 记录消息的具体内容
type Message struct {
	FromWxId      string
	FromName      string
	FromGroup     string
	FromGroupName string
	Msg           string
	MsgId         string
	MsgType       int
}
