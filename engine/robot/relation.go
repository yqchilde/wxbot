package robot

import "time"

// Friend 好友对象
type Friend struct{ *User }

// AsUser 将当前对象转换为User对象
func (f *Friend) AsUser() *User {
	return f.User
}

// SendText  发送文本消息
func (f *Friend) SendText(content string) error {
	return f.self.sendText(f.User, content)
}

// SendImage 发送图片消息
func (f *Friend) SendImage(path string) error {
	return f.self.sendImage(f.User, path)
}

// SendShareLink 发送分享链接消息
func (f *Friend) SendShareLink(title, desc, imageUrl, jumpUrl string) error {
	return f.self.sendShareLink(f.User, title, desc, imageUrl, jumpUrl)
}

// SendFile 发送文件消息
func (f *Friend) SendFile(path string) error {
	return f.self.sendFile(f.User, path)
}

// SendVideo 发送视频消息
func (f *Friend) SendVideo(path string) error {
	return f.self.sendVideo(f.User, path)
}

// SendEmoji 发送表情消息
func (f *Friend) SendEmoji(path string) error {
	return f.self.sendEmoji(f.User, path)
}

// SendMusic 发送音乐消息
func (f *Friend) SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return f.self.sendMusic(f.User, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// SendMiniProgram 发送小程序消息
func (f *Friend) SendMiniProgram(ghId, title, content, imagePath, jumpPath string) error {
	return f.self.sendMiniProgram(f.User, ghId, title, content, imagePath, jumpPath)
}

// SendMessageRecord 发送消息记录
func (f *Friend) SendMessageRecord(title string, dataList []map[string]interface{}) error {
	return f.self.sendMessageRecord(f.User, title, dataList)
}

// SendMessageRecordXML 发送消息记录(支持xml)
func (f *Friend) SendMessageRecordXML(xmlStr string) error {
	return f.self.sendMessageRecordXML(f.User, xmlStr)
}

// SendFavorites 发送收藏消息
func (f *Friend) SendFavorites(favoritesId string) error {
	return f.self.sendFavorites(f.User, favoritesId)
}

// SendXML 发送xml消息
func (f *Friend) SendXML(xmlStr string) error {
	return f.self.sendXML(f.User, xmlStr)
}

// SendBusinessCard 发送名片消息
func (f *Friend) SendBusinessCard(targetWxId string) error {
	return f.self.sendBusinessCard(f.User, targetWxId)
}

// Friends 好友列表
type Friends []*Friend

// Count 获取好友的数量
func (f Friends) Count() int {
	return len(f)
}

// SendText 依次发送文本消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendText(content string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendText(content); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendImage 依次发送图片消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendImage(path string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendImage(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendShareLink 依次发送分享链接消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendShareLink(title, desc, imageUrl, jumpUrl string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendShareLink(title, desc, imageUrl, jumpUrl); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendFile 依次发送文件消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendFile(path string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendFile(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendVideo 依次发送视频消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendVideo(path string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendVideo(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendEmoji 依次发送表情消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendEmoji(path string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendEmoji(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMusic 依次发送音乐消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMiniProgram 依次发送小程序消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendMiniProgram(ghId, title, content, imagePath, jumpPath string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendMiniProgram(ghId, title, content, imagePath, jumpPath); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMessageRecord 依次发送消息记录, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendMessageRecord(title string, dataList []map[string]interface{}, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendMessageRecord(title, dataList); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMessageRecordXML 依次发送消息记录, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendMessageRecordXML(xmlStr string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendMessageRecordXML(xmlStr); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendFavorites 依次发送收藏消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendFavorites(favoritesId string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendFavorites(favoritesId); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendXML 依次发送XML消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendXML(xmlStr string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendXML(xmlStr); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendBusinessCard 依次发送名片消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (f Friends) SendBusinessCard(targetWxId string, delays ...time.Duration) error {
	for _, friend := range f {
		if err := friend.SendBusinessCard(targetWxId); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// Group 群组对象
type Group struct{ *User }

// AsUser 将当前对象转换为User对象
func (g *Group) AsUser() *User {
	return g.User
}

// SendText  发送文本消息
func (g *Group) SendText(content string) error {
	return g.self.sendText(g.User, content)
}

// SendImage 发送图片消息
func (g *Group) SendImage(path string) error {
	return g.self.sendImage(g.User, path)
}

// SendShareLink 发送分享链接消息
func (g *Group) SendShareLink(title, desc, imageUrl, jumpUrl string) error {
	return g.self.sendShareLink(g.User, title, desc, imageUrl, jumpUrl)
}

// SendFile 发送文件消息
func (g *Group) SendFile(path string) error {
	return g.self.sendFile(g.User, path)
}

// SendVideo 发送视频消息
func (g *Group) SendVideo(path string) error {
	return g.self.sendVideo(g.User, path)
}

// SendEmoji 发送表情消息
func (g *Group) SendEmoji(path string) error {
	return g.self.sendEmoji(g.User, path)
}

// SendMusic 发送音乐消息
func (g *Group) SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return g.self.sendMusic(g.User, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// SendMiniProgram 发送小程序消息
func (g *Group) SendMiniProgram(ghId, title, content, imagePath, jumpPath string) error {
	return g.self.sendMiniProgram(g.User, ghId, title, content, imagePath, jumpPath)
}

// SendMessageRecord 发送消息记录
func (g *Group) SendMessageRecord(title string, dataList []map[string]interface{}) error {
	return g.self.sendMessageRecord(g.User, title, dataList)
}

// SendMessageRecordXML 发送消息记录(支持xml)
func (g *Group) SendMessageRecordXML(xmlStr string) error {
	return g.self.sendMessageRecordXML(g.User, xmlStr)
}

// SendFavorites 发送收藏消息
func (g *Group) SendFavorites(favoritesId string) error {
	return g.self.sendFavorites(g.User, favoritesId)
}

// SendXML 发送xml消息
func (g *Group) SendXML(xmlStr string) error {
	return g.self.sendXML(g.User, xmlStr)
}

// SendBusinessCard 发送名片消息
func (g *Group) SendBusinessCard(targetWxId string) error {
	return g.self.sendBusinessCard(g.User, targetWxId)
}

// Groups 群组列表
type Groups []*Group

// Count 获取群组的数量
func (g Groups) Count() int {
	return len(g)
}

// SendText 依次发送文本消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendText(content string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendText(content); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendImage 依次发送图片消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendImage(path string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendImage(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendShareLink 依次发送分享链接消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendShareLink(title, desc, imageUrl, jumpUrl string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendShareLink(title, desc, imageUrl, jumpUrl); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendFile 依次发送文件消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendFile(path string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendFile(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendVideo 依次发送视频消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendVideo(path string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendVideo(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendEmoji 依次发送表情消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendEmoji(path string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendEmoji(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMusic 依次发送音乐消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMiniProgram 依次发送小程序消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendMiniProgram(ghId, title, content, imagePath, jumpPath string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendMiniProgram(ghId, title, content, imagePath, jumpPath); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMessageRecord 依次发送消息记录, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendMessageRecord(title string, dataList []map[string]interface{}, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendMessageRecord(title, dataList); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMessageRecordXML 依次发送消息记录, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendMessageRecordXML(xmlStr string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendMessageRecordXML(xmlStr); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendFavorites 依次发送收藏夹消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendFavorites(favoritesId string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendFavorites(favoritesId); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendXML 依次发送XML消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendXML(xmlStr string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendXML(xmlStr); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendBusinessCard 依次发送名片消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (g Groups) SendBusinessCard(targetWxId string, delays ...time.Duration) error {
	for _, group := range g {
		if err := group.SendBusinessCard(targetWxId); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil

}

// GroupMembers 群成员列表对象
type GroupMembers []*User

// Count 获取群成员数量
func (g GroupMembers) Count() int {
	return len(g)
}

// MP 公众号对象
type MP struct{ *User }

// AsUser 将当前对象转换为User对象
func (m *MP) AsUser() *User {
	return m.User
}

// SendText  发送文本消息
func (m *MP) SendText(content string) error {
	return m.self.sendText(m.User, content)
}

// SendImage 发送图片消息
func (m *MP) SendImage(path string) error {
	return m.self.sendImage(m.User, path)
}

// SendShareLink 发送分享链接消息
func (m *MP) SendShareLink(title, desc, imageUrl, jumpUrl string) error {
	return m.self.sendShareLink(m.User, title, desc, imageUrl, jumpUrl)
}

// SendFile 发送文件消息
func (m *MP) SendFile(path string) error {
	return m.self.sendFile(m.User, path)
}

// SendVideo 发送视频消息
func (m *MP) SendVideo(path string) error {
	return m.self.sendVideo(m.User, path)
}

// SendEmoji 发送表情消息
func (m *MP) SendEmoji(path string) error {
	return m.self.sendEmoji(m.User, path)
}

// SendMusic 发送音乐消息
func (m *MP) SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl string) error {
	return m.self.sendMusic(m.User, name, author, app, jumpUrl, musicUrl, coverUrl)
}

// SendMiniProgram 发送小程序消息
func (m *MP) SendMiniProgram(ghId, title, content, imagePath, jumpPath string) error {
	return m.self.sendMiniProgram(m.User, ghId, title, content, imagePath, jumpPath)
}

// SendMessageRecord 发送消息记录
func (m *MP) SendMessageRecord(title string, dataList []map[string]interface{}) error {
	return m.self.sendMessageRecord(m.User, title, dataList)
}

// SendMessageRecordXML 发送消息记录(支持xml)
func (m *MP) SendMessageRecordXML(xmlStr string) error {
	return m.self.sendMessageRecordXML(m.User, xmlStr)
}

// SendFavorites 发送收藏消息
func (m *MP) SendFavorites(favoritesId string) error {
	return m.self.sendFavorites(m.User, favoritesId)
}

// SendXML 发送xml消息
func (m *MP) SendXML(xmlStr string) error {
	return m.self.sendXML(m.User, xmlStr)
}

// SendBusinessCard 发送名片消息
func (m *MP) SendBusinessCard(targetWxId string) error {
	return m.self.sendBusinessCard(m.User, targetWxId)
}

// MPs 公众号列表
type MPs []*MP

// Count 获取公众号的数量
func (m MPs) Count() int {
	return len(m)
}

// SendText 依次发送文本消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendText(content string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendText(content); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendImage 依次发送图片消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendImage(path string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendImage(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendShareLink 依次发送分享链接消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendShareLink(title, desc, imageUrl, jumpUrl string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendShareLink(title, desc, imageUrl, jumpUrl); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendFile 依次发送文件消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendFile(path string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendFile(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendVideo 依次发送视频消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendVideo(path string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendVideo(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendEmoji 依次发送表情消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendEmoji(path string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendEmoji(path); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMusic 依次发送音乐消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendMusic(name, author, app, jumpUrl, musicUrl, coverUrl); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMiniProgram 依次发送小程序消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendMiniProgram(ghId, title, content, imagePath, jumpPath string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendMiniProgram(ghId, title, content, imagePath, jumpPath); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMessageRecord 依次发送消息记录, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendMessageRecord(title string, dataList []map[string]interface{}, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendMessageRecord(title, dataList); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendMessageRecordXML 依次发送消息记录, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendMessageRecordXML(xmlStr string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendMessageRecordXML(xmlStr); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendFavorites 依次发送收藏夹消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendFavorites(favoritesId string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendFavorites(favoritesId); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendXML 依次发送XML消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendXML(xmlStr string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendXML(xmlStr); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}

// SendBusinessCard 依次发送名片消息, delays为每个好友发送消息的间隔时间, 默认为1秒, 为0时不间隔
func (m MPs) SendBusinessCard(targetWxId string, delays ...time.Duration) error {
	for _, mp := range m {
		if err := mp.SendBusinessCard(targetWxId); err != nil {
			return err
		}
		if len(delays) > 0 {
			time.Sleep(delays[0])
		} else {
			time.Sleep(1 * time.Second)
		}
	}
	return nil
}
