package robot

const (
	MsgTypeText           = 1     // 文本消息
	MsgTypeImage          = 3     // 图片消息
	MsgTypeVoice          = 34    // 语音消息
	MsgTypeAuthentication = 37    // 认证消息
	MsgTypePossibleFriend = 40    // 好友推荐消息
	MsgTypeShareCard      = 42    // 名片消息
	MsgTypeVideo          = 43    // 视频消息
	MsgTypeMemePicture    = 47    // 表情消息
	MsgTypeLocation       = 48    // 地理位置消息
	MsgTypeApp            = 49    // APP消息
	MsgTypeMicroVideo     = 62    // 小视频消息
	MsgTypeSystem         = 10000 // 系统消息
	MsgTypeRecalled       = 10002 // 消息撤回
	MsgTypeReference      = 10003 // 消息引用
)
