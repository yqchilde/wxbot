package robot

// Event 记录一次回调事件
type Event struct {
	RobotWxId     string  // 机器人微信id
	IsAtMe        bool    // 机器人是否被@了，@所有人不算
	IsPrivateChat bool    // 是否是私聊消息
	IsGroupChat   bool    // 是否是群聊消息
	FromUniqueID  string  // 消息来源唯一id, 私聊为发送者微信id, 群聊为群id
	FromWxId      string  // 消息来源微信id
	FromName      string  // 消息来源昵称
	FromGroup     string  // 消息来源群id
	FromGroupName string  // 消息来源群名称
	Message       Message // 消息内容
}

// Message 记录消息的具体内容
type Message struct {
	Msg     string
	MsgId   string
	MsgType int64
}
