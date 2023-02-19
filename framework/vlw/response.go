package vlw

type RobotInfoResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson struct {
		Number int `json:"Number"`
		Data   []struct {
			Pid                      int    `json:"pid"`
			Username                 string `json:"username"`
			Wxid                     string `json:"wxid"`
			WxNum                    string `json:"wx_num"`
			WxHeadimgurl             string `json:"wx_headimgurl"`
			EnterpriseWechat         int    `json:"Enterprise wechat"`
			EnterpriseWechatClientId int    `json:"Enterprise wechat clientId"`
		} `json:"data"`
	} `json:"ReturnJson"`
}

// ObjectInfoResp 对象可以是好友、群、公众号
type ObjectInfoResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson struct {
		Data struct {
			Account     string `json:"account"`
			Avatar      string `json:"avatar"`
			City        string `json:"city"`
			Country     string `json:"country"`
			Nickname    string `json:"nickname"`
			Province    string `json:"province"`
			Remark      string `json:"remark"`
			Sex         int    `json:"sex"`
			Signature   string `json:"signature"`
			SmallAvatar string `json:"small_avatar"`
			SnsPic      string `json:"sns_pic"`
			SourceType  int    `json:"source_type"`
			Status      int    `json:"status"`
			V1          string `json:"v1"`
			V2          string `json:"v2"`
			Wxid        string `json:"wxid"`
		} `json:"data"`
		Type int `json:"type"`
	} `json:"ReturnJson"`
}

// FriendsListResp 获取好友列表响应
type FriendsListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson []struct {
		WxNum    string `json:"wx_num"`
		Avatar   string `json:"avatar"`
		City     string `json:"city"`
		Country  string `json:"country"`
		Nickname string `json:"nickname"`
		Province string `json:"province"`
		Note     string `json:"note"`
		Sex      int    `json:"sex"`
		Wxid     string `json:"wxid"`
	} `json:"ReturnJson"`
}

// GroupListResp 获取群组列表响应
type GroupListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson []struct {
		Avatar      string `json:"avatar"`
		IsManager   int    `json:"is_manager"`
		ManagerWxid string `json:"manager_wxid"`
		Nickname    string `json:"nickname"`
		TotalMember int    `json:"total_member"`
		Wxid        string `json:"wxid"`
	} `json:"ReturnJson"`
}

// GroupMemberListResp 获取群成员列表响应
type GroupMemberListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson struct {
		GroupWxid     string `json:"group_wxid"`
		GroupName     string `json:"group_name"`
		Count         int    `json:"count"`
		OwnerWxid     string `json:"owner_wxid"`
		OwnerNickname string `json:"owner_nickname"`
		MemberList    []struct {
			WxNum         string `json:"wx_num"`
			Avatar        string `json:"avatar"`
			City          string `json:"city"`
			Country       string `json:"country"`
			GroupNickname string `json:"group_nickname"`
			Nickname      string `json:"nickname"`
			Province      string `json:"province"`
			Remark        string `json:"remark"`
			Sex           int    `json:"sex"`
			Wxid          string `json:"wxid"`
		} `json:"member_list"`
	} `json:"ReturnJson"`
}

// SubscriptionListResp 获取订阅号列表响应
type SubscriptionListResp struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson []struct {
		Avatar   string `json:"avatar"`
		Nickname string `json:"nickname"`
		Wxid     string `json:"wxid"`
	} `json:"ReturnJson"`
}
