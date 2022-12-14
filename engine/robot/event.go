package robot

// Event 记录一次回调事件
type Event struct {
	Type          string        // 消息类型
	RobotWxId     string        // 机器人微信id
	IsAtMe        bool          // 机器人是否被@了，@所有人不算
	FromUniqueID  string        // 消息来源唯一id, 私聊为发送者微信id, 群聊为群id
	FromWxId      string        // 消息来源微信id
	FromName      string        // 消息来源昵称
	FromGroup     string        // 消息来源群id
	FromGroupName string        // 消息来源群名称
	Message       *Message      // 消息内容
	FriendVerify  *FriendVerify // 好友验证消息
}

// Message 记录消息的具体内容
type Message struct {
	Id      string // 消息id
	Type    int64  // 消息类型
	Content string // 消息内容
}

// FriendVerify 记录好友验证消息的具体内容
type FriendVerify struct {
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
