package robot

// Event 记录一次回调事件
type Event struct {
	Type                string               // 消息类型
	RobotWxId           string               // 机器人微信id
	IsAtMe              bool                 // 机器人是否被@了，@所有人不算
	FromUniqueID        string               // 消息来源唯一id, 私聊为发送者微信id, 群聊为群id
	FromUniqueName      string               // 消息来源唯一名称, 私聊为发送者昵称, 群聊为群名称
	FromWxId            string               // 消息来源微信id
	FromName            string               // 消息来源昵称
	FromGroup           string               // 消息来源群id
	FromGroupName       string               // 消息来源群名称
	RawMessage          string               // 原始消息
	Message             *Message             // 消息内容
	MPMessage           *Message             // 订阅公众号消息
	FriendVerifyMessage *FriendVerifyMessage // 好友验证消息
	TransferMessage     *TransferMessage     // 转账消息
	WithdrawMessage     *WithdrawMessage     // 撤回消息
	ReferenceMessage    *ReferenceMessage    // 引用消息
}

// Message 记录消息的具体内容
type Message struct {
	Id      string // 消息id
	Type    int64  // 消息类型
	Content string // 消息内容
}

// FriendVerifyMessage 记录好友验证消息的具体内容
type FriendVerifyMessage struct {
	WxId      string // 发送者微信id
	Nick      string // 发送者昵称
	V1        string // 验证V1
	V2        string // 验证V2
	V3        string // 验证V3
	V4        string // 验证V4
	AvatarUrl string // 头像url
	Content   string // 验证内容
	Scene     string // 验证场景
}

// TransferMessage 记录转账消息的具体内容
type TransferMessage struct {
	FromWxId     string // 发送者微信ID
	MsgSource    int64  // 消息来源 1:收到转账 2:对方接收转账 3:发出转账 4:自己接收转账 5:对方退还 6:自己退还
	TransferType int64  // 转账类型 1:即时到账 2:延时到账
	Money        string // 转账金额，单位元
	Memo         string // 转账备注
	TransferId   string // 转账ID
	TransferTime string // 转账时间，10位时间戳
}

// WithdrawMessage 记录撤回消息的具体内容
type WithdrawMessage struct {
	FromType  int64  // 消息来源 1:私聊 2:群聊
	FromGroup string // 消息来源群ID
	FromWxId  string // 消息来源微信ID
	MsgSource int64  // 消息来源 1:别人撤回 2:自己使用手机撤回 3:自己使用电脑撤回
	Msg       string // 消息内容
}

// ReferenceMessage 记录引用消息的具体内容
type ReferenceMessage struct {
	FromUser    string // 消息来源群ID
	ChatUser    string // 消息来源微信ID
	DisplayName string // 消息来源微信昵称
	Content     string // 消息内容
}
