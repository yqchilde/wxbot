package robot

type BotList struct {
	Code       int    `json:"Code"`
	Result     string `json:"Result"`
	ReturnJson struct {
		Number int   `json:"Number"`
		Data   []Bot `json:"data"`
	} `json:"ReturnJson"`
}

type Bot struct {
	Pid                      int    `json:"pid"`
	Username                 string `json:"username"`
	Wxid                     string `json:"wxid"`
	WxNum                    string `json:"wx_num"`
	WxHeadimgurl             string `json:"wx_headimgurl"`
	EnterpriseWechat         int    `json:"Enterprise wechat"`
	EnterpriseWechatClientId int    `json:"Enterprise wechat clientId"`
}

type GroupList struct {
	Code       int     `json:"Code"`
	Result     string  `json:"Result"`
	ReturnJson []Group `json:"ReturnJson"`
}

type Group struct {
	Avatar      string `json:"avatar"`
	IsManager   int    `json:"is_manager"`
	ManagerWxid string `json:"manager_wxid"`
	Nickname    string `json:"nickname"`
	TotalMember int    `json:"total_member"`
	Wxid        string `json:"wxid"`
}
