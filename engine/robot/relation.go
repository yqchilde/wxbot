package robot

// Friend 好友对象
type Friend struct{ *User }

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

// Group 群组对象
type Group struct{ *User }

// Groups 群组列表
type Groups []*Group

// MP 公众号对象
type MP struct{ *User }

// MPs 公众号列表
type MPs []*MP
