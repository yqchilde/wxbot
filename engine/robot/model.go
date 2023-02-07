package robot

// ObjectInfo 对象信息
// 对象可以是好友、群、公众号
type ObjectInfo struct {
	WxId                    string `json:"wxId"`                    // 微信ID
	WxNum                   string `json:"wxNum"`                   // 微信号
	Nick                    string `json:"nick"`                    // 昵称
	Remark                  string `json:"remark"`                  // 备注
	NickBrief               string `json:"nickBrief"`               // 昵称简拼
	NickWhole               string `json:"nickWhole"`               // 昵称全拼
	RemarkBrief             string `json:"remarkBrief"`             // 备注简拼
	RemarkWhole             string `json:"remarkWhole"`             // 备注全拼
	EnBrief                 string `json:"enBrief"`                 // 英文简拼
	EnWhole                 string `json:"enWhole"`                 // 英文全拼
	V3                      string `json:"v3"`                      // v3数据，同意好友验证时需要
	V4                      string `json:"v4"`                      // v4数据，同意好友验证时需要
	Sign                    string `json:"sign"`                    // 签名，需要在会话列表中
	Country                 string `json:"country"`                 // 国家，需要在会话列表中
	Province                string `json:"province"`                // 省份，需要在会话列表中
	City                    string `json:"city"`                    // 城市，需要在会话列表中
	MomentsBackgroundImgUrl string `json:"momentsBackgroundImgUrl"` // 朋友圈背景图，需要在朋友圈中
	AvatarMinUrl            string `json:"avatarMinUrl"`            // 头像小图，需要在会话列表中
	AvatarMaxUrl            string `json:"avatarMaxUrl"`            // 头像大图，需要在会话列表中
	Sex                     string `json:"sex"`                     // 性别，1男，2女，0未知
	MemberNum               int    `json:"memberNum"`               // 群成员数量，仅当对象是群聊时有效
}

// FriendInfo 好友信息
type FriendInfo struct {
	WxId                    string `json:"wxId"`                    // 微信ID
	WxNum                   string `json:"wxNum"`                   // 微信号
	Nick                    string `json:"nick"`                    // 昵称
	Remark                  string `json:"remark"`                  // 备注
	NickBrief               string `json:"nickBrief"`               // 昵称简拼
	NickWhole               string `json:"nickWhole"`               // 昵称全拼
	RemarkBrief             string `json:"remarkBrief"`             // 备注简拼
	RemarkWhole             string `json:"remarkWhole"`             // 备注全拼
	EnBrief                 string `json:"enBrief"`                 // 英文简拼
	EnWhole                 string `json:"enWhole"`                 // 英文全拼
	V3                      string `json:"v3"`                      // v3数据，同意好友验证时需要
	Sign                    string `json:"sign"`                    // 签名，需要在会话列表中
	Country                 string `json:"country"`                 // 国家，需要在会话列表中
	Province                string `json:"province"`                // 省份，需要在会话列表中
	City                    string `json:"city"`                    // 城市，需要在会话列表中
	MomentsBackgroundImgUrl string `json:"momentsBackgroundImgUrl"` // 朋友圈背景图，需要在朋友圈中
	AvatarMinUrl            string `json:"avatarMinUrl"`            // 头像小图，需要在会话列表中
	AvatarMaxUrl            string `json:"avatarMaxUrl"`            // 头像大图，需要在会话列表中
	Sex                     string `json:"sex"`                     // 性别，1男，2女，0未知
	MemberNum               int    `json:"memberNum"`               // 群成员数量，仅当对象是群聊时有效
}

// GroupInfo 群信息
type GroupInfo struct {
	WxId         string `json:"wxId"`         // 微信ID
	WxNum        string `json:"wxNum"`        // 微信号
	Nick         string `json:"nick"`         // 昵称
	Remark       string `json:"remark"`       // 备注
	NickBrief    string `json:"nickBrief"`    // 昵称简拼
	NickWhole    string `json:"nickWhole"`    // 昵称全拼
	RemarkBrief  string `json:"remarkBrief"`  // 备注简拼
	RemarkWhole  string `json:"remarkWhole"`  // 备注全拼
	EnBrief      string `json:"enBrief"`      // 英文简拼
	EnWhole      string `json:"enWhole"`      // 英文全拼
	MemberNum    int    `json:"memberNum"`    // 群成员数量，仅当对象是群聊时有效
	AvatarMinUrl string `json:"avatarMinUrl"` // 头像小图，需要在会话列表中
	AvatarMaxUrl string `json:"avatarMaxUrl"` // 头像大图，需要在会话列表中
}

// SubscriptionInfo 订阅公众号信息
type SubscriptionInfo struct {
	WxId                    string `json:"wxId"`                    // 微信ID
	WxNum                   string `json:"wxNum"`                   // 微信号
	Nick                    string `json:"nick"`                    // 昵称
	Remark                  string `json:"remark"`                  // 备注
	NickBrief               string `json:"nickBrief"`               // 昵称简拼
	NickWhole               string `json:"nickWhole"`               // 昵称全拼
	RemarkBrief             string `json:"remarkBrief"`             // 备注简拼
	RemarkWhole             string `json:"remarkWhole"`             // 备注全拼
	EnBrief                 string `json:"enBrief"`                 // 英文简拼
	EnWhole                 string `json:"enWhole"`                 // 英文全拼
	V3                      string `json:"v3"`                      // v3数据，同意好友验证时需要
	Sign                    string `json:"sign"`                    // 签名，需要在会话列表中
	Country                 string `json:"country"`                 // 国家，需要在会话列表中
	Province                string `json:"province"`                // 省份，需要在会话列表中
	City                    string `json:"city"`                    // 城市，需要在会话列表中
	MomentsBackgroundImgUrl string `json:"momentsBackgroundImgUrl"` // 朋友圈背景图，需要在朋友圈中
	AvatarMinUrl            string `json:"avatarMinUrl"`            // 头像小图，需要在会话列表中
	AvatarMaxUrl            string `json:"avatarMaxUrl"`            // 头像大图，需要在会话列表中
}
