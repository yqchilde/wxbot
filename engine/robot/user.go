package robot

// User 抽象用户对象，对象可以是好友、群组、公众号
type User struct {
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

	self *Self
}

// Self 包装了关于bot、好友、群组、公众号的操作
type Self struct {
	bot     *Bot
	friends Friends
	groups  Groups
	mps     MPs
}

func (s *Self) CheckUserObjNil() bool {
	return s.friends == nil && s.groups == nil && s.mps == nil
}

// Init 初始化获取好友、群、公众号列表
func (s *Self) Init() error {
	if _, err := s.Friends(true); err != nil {
		return err
	}
	if _, err := s.Groups(true); err != nil {
		return err
	}
	if _, err := s.MPs(true); err != nil {
		return err
	}
	return nil
}

// Friends 获取所有的好友
func (s *Self) Friends(update ...bool) (Friends, error) {
	if (len(update) > 0 && update[0]) || s.CheckUserObjNil() {
		if friendsList, err := s.bot.framework.GetFriends(true); err != nil {
			return nil, err
		} else {
			s.friends = make(Friends, 0)
			for _, friend := range friendsList {
				friend.self = s
				s.friends = append(s.friends, &Friend{User: friend})
			}
		}
	}
	return s.friends, nil
}

// Groups 获取所有的群组
func (s *Self) Groups(update ...bool) (Groups, error) {
	if (len(update) > 0 && update[0]) || s.CheckUserObjNil() {
		if groupsList, err := s.bot.framework.GetGroups(true); err != nil {
			return nil, err
		} else {
			s.groups = make(Groups, 0)
			for _, group := range groupsList {
				group.self = s
				s.groups = append(s.groups, &Group{User: group})
			}
		}
	}
	return s.groups, nil
}

// MPs 获取所有的公众号
func (s *Self) MPs(update ...bool) (MPs, error) {
	if (len(update) > 0 && update[0]) || s.CheckUserObjNil() {
		if mpsList, err := s.bot.framework.GetMPs(true); err != nil {
			return nil, err
		} else {
			s.mps = make(MPs, 0)
			for _, mp := range mpsList {
				mp.self = s
				s.mps = append(s.mps, &MP{User: mp})
			}
		}
	}
	return s.mps, nil
}

// sendText 发送文本消息
func (s *Self) sendText(user *User, text string) error {
	return s.bot.framework.SendText(user.WxId, text)
}

// sendImage 发送图片消息
func (s *Self) sendImage(user *User, path string) error {
	return s.bot.framework.SendImage(user.WxId, path)
}

// sendShareLink 发送分享链接消息
func (s *Self) sendShareLink(user *User, title, desc, imageUrl, jumpUrl string) error {
	return s.bot.framework.SendShareLink(user.WxId, title, desc, imageUrl, jumpUrl)
}

// sendFile 发送文件消息
func (s *Self) sendFile(user *User, path string) error {
	return s.bot.framework.SendFile(user.WxId, path)
}

// sendVideo 发送视频消息
func (s *Self) sendVideo(user *User, path string) error {
	return s.bot.framework.SendVideo(user.WxId, path)
}

// sendEmoji 发送表情消息
func (s *Self) sendEmoji(user *User, path string) error {
	return s.bot.framework.SendEmoji(user.WxId, path)
}

// sendMusic 发送音乐消息
func (s *Self) sendMusic(user *User, name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return s.bot.framework.SendMusic(user.WxId, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// sendMiniProgram 发送小程序消息
func (s *Self) sendMiniProgram(user *User, ghId, title, content, imagePath, jumpPath string) error {
	return s.bot.framework.SendMiniProgram(user.WxId, ghId, title, content, imagePath, jumpPath)
}

// sendMessageRecord 发送消息记录
func (s *Self) sendMessageRecord(user *User, title string, dataList []map[string]interface{}) error {
	return s.bot.framework.SendMessageRecord(user.WxId, title, dataList)
}

// sendMessageRecordXML 发送消息记录(支持xml)
func (s *Self) sendMessageRecordXML(user *User, xmlStr string) error {
	return s.bot.framework.SendMessageRecordXML(user.WxId, xmlStr)
}

// sendFavorites 发送收藏消息
func (s *Self) sendFavorites(user *User, favoritesId string) error {
	return s.bot.framework.SendFavorites(user.WxId, favoritesId)
}

// sendXML 发送xml消息
func (s *Self) sendXML(user *User, xmlStr string) error {
	return s.bot.framework.SendXML(user.WxId, xmlStr)
}

// sendBusinessCard 发送名片消息
func (s *Self) sendBusinessCard(user *User, targetWxId string) error {
	return s.bot.framework.SendBusinessCard(user.WxId, targetWxId)
}
