package qianxun

type RobotInfoResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Wxid      string `json:"wxid"`
		WxNum     string `json:"wxNum"`
		Nick      string `json:"nick"`
		Device    string `json:"device"`
		Phone     string `json:"phone"`
		AvatarUrl string `json:"avatarUrl"`
		Country   string `json:"country"`
		Province  string `json:"province"`
		City      string `json:"city"`
		Email     string `json:"email"`
		Qq        string `json:"qq"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// ObjectInfoResp 对象可以是好友、群、公众号
type ObjectInfoResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result struct {
		Wxid                   string `json:"wxid"`
		WxNum                  string `json:"wxNum"`
		Nick                   string `json:"nick"`
		Remark                 string `json:"remark"`
		NickBrief              string `json:"nickBrief"`
		NickWhole              string `json:"nickWhole"`
		RemarkBrief            string `json:"remarkBrief"`
		RemarkWhole            string `json:"remarkWhole"`
		EnBrief                string `json:"enBrief"`
		EnWhole                string `json:"enWhole"`
		V3                     string `json:"v3"`
		V4                     string `json:"v4"`
		Sign                   string `json:"sign"`
		Country                string `json:"country"`
		Province               string `json:"province"`
		City                   string `json:"city"`
		MomentsBackgroudImgUrl string `json:"momentsBackgroudImgUrl"`
		AvatarMinUrl           string `json:"avatarMinUrl"`
		AvatarMaxUrl           string `json:"avatarMaxUrl"`
		Sex                    string `json:"sex"`
		MemberNum              int    `json:"memberNum"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// FriendsListResp 获取好友列表响应
type FriendsListResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Wxid                   string `json:"wxid"`
		WxNum                  string `json:"wxNum"`
		Nick                   string `json:"nick"`
		Remark                 string `json:"remark"`
		NickBrief              string `json:"nickBrief"`
		NickWhole              string `json:"nickWhole"`
		RemarkBrief            string `json:"remarkBrief"`
		RemarkWhole            string `json:"remarkWhole"`
		EnBrief                string `json:"enBrief"`
		EnWhole                string `json:"enWhole"`
		V3                     string `json:"v3"`
		Sign                   string `json:"sign"`
		Country                string `json:"country"`
		Province               string `json:"province"`
		City                   string `json:"city"`
		MomentsBackgroudImgUrl string `json:"momentsBackgroudImgUrl"`
		AvatarMinUrl           string `json:"avatarMinUrl"`
		AvatarMaxUrl           string `json:"avatarMaxUrl"`
		Sex                    string `json:"sex"`
		MemberNum              int    `json:"memberNum"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// GroupListResp 获取群组列表响应
type GroupListResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Wxid                   string `json:"wxid"`
		WxNum                  string `json:"wxNum"`
		Nick                   string `json:"nick"`
		Remark                 string `json:"remark"`
		NickBrief              string `json:"nickBrief"`
		NickWhole              string `json:"nickWhole"`
		RemarkBrief            string `json:"remarkBrief"`
		RemarkWhole            string `json:"remarkWhole"`
		EnBrief                string `json:"enBrief"`
		EnWhole                string `json:"enWhole"`
		V3                     string `json:"v3"`
		Sign                   string `json:"sign"`
		Country                string `json:"country"`
		Province               string `json:"province"`
		City                   string `json:"city"`
		MomentsBackgroudImgUrl string `json:"momentsBackgroudImgUrl"`
		AvatarMinUrl           string `json:"avatarMinUrl"`
		AvatarMaxUrl           string `json:"avatarMaxUrl"`
		Sex                    string `json:"sex"`
		MemberNum              int    `json:"memberNum"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// GroupMemberListResp 获取群成员列表响应
type GroupMemberListResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Wxid      string `json:"wxid"`
		GroupNick string `json:"groupNick"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}

// SubscriptionListResp 获取订阅号列表响应
type SubscriptionListResp struct {
	Code   int    `json:"code"`
	Msg    string `json:"msg"`
	Result []struct {
		Wxid                   string `json:"wxid"`
		WxNum                  string `json:"wxNum"`
		Nick                   string `json:"nick"`
		Remark                 string `json:"remark"`
		NickBrief              string `json:"nickBrief"`
		NickWhole              string `json:"nickWhole"`
		RemarkBrief            string `json:"remarkBrief"`
		RemarkWhole            string `json:"remarkWhole"`
		EnBrief                string `json:"enBrief"`
		EnWhole                string `json:"enWhole"`
		V3                     string `json:"v3"`
		Sign                   string `json:"sign"`
		Country                string `json:"country"`
		Province               string `json:"province"`
		City                   string `json:"city"`
		MomentsBackgroudImgUrl string `json:"momentsBackgroudImgUrl"`
		AvatarMinUrl           string `json:"avatarMinUrl"`
		AvatarMaxUrl           string `json:"avatarMaxUrl"`
		Sex                    string `json:"sex"`
		MemberNum              int    `json:"memberNum"`
	} `json:"result"`
	Wxid      string `json:"wxid"`
	Port      int    `json:"port"`
	Pid       int    `json:"pid"`
	Flag      string `json:"flag"`
	Timestamp string `json:"timestamp"`
}
