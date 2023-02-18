package robot

import (
	"strings"
	"time"
)

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

// AsUsers 将好友列表转换为User列表
func (f Friends) AsUsers() []*User {
	var users []*User
	for _, friend := range f {
		users = append(users, friend.User)
	}
	return users
}

// GetByWxId 根据微信ID获取好友
func (f Friends) GetByWxId(wxId string) *Friend {
	for _, friend := range f {
		if friend.WxId == wxId {
			return friend
		}
	}
	return nil
}

// GetByWxNum 根据微信号获取好友
func (f Friends) GetByWxNum(wxNum string) *Friend {
	for _, friend := range f {
		if friend.WxNum == wxNum {
			return friend
		}
	}
	return nil
}

// GetByNick 根据昵称获取好友
func (f Friends) GetByNick(nick string) *Friend {
	for _, friend := range f {
		if friend.Nick == nick {
			return friend
		}
	}
	return nil
}

// GetByRemark 根据备注获取好友
func (f Friends) GetByRemark(remark string) *Friend {
	for _, friend := range f {
		if friend.Remark == remark {
			return friend
		}
	}
	return nil
}

// GetByRemarkOrNick 根据备注或昵称获取好友
func (f Friends) GetByRemarkOrNick(remarkOrNick string) *Friend {
	for _, friend := range f {
		if friend.Remark == remarkOrNick || friend.Nick == remarkOrNick {
			return friend
		}
	}
	return nil
}

// GetByWxIds 根据微信ID列表获取好友列表
func (f Friends) GetByWxIds(wxIds []string) Friends {
	var result Friends
	for _, wxId := range wxIds {
		if friend := f.GetByWxId(wxId); friend != nil {
			result = append(result, friend)
		}
	}
	return result
}

// GetByWxNums 根据微信号列表获取好友列表
func (f Friends) GetByWxNums(wxNums []string) Friends {
	var result Friends
	for _, wxNum := range wxNums {
		if friend := f.GetByWxNum(wxNum); friend != nil {
			result = append(result, friend)
		}
	}
	return result
}

// GetByNicks 根据昵称列表获取好友列表
func (f Friends) GetByNicks(nicks []string) Friends {
	var result Friends
	for _, nick := range nicks {
		if friend := f.GetByNick(nick); friend != nil {
			result = append(result, friend)
		}
	}
	return result
}

// GetByRemarks 根据备注列表获取好友列表
func (f Friends) GetByRemarks(remarks []string) Friends {
	var result Friends
	for _, remark := range remarks {
		if friend := f.GetByRemark(remark); friend != nil {
			result = append(result, friend)
		}
	}
	return result
}

// GetByRemarkOrNicks 根据备注或昵称列表获取好友列表
func (f Friends) GetByRemarkOrNicks(remarkOrNicks []string) Friends {
	var result Friends
	for _, remarkOrNick := range remarkOrNicks {
		if friend := f.GetByRemarkOrNick(remarkOrNick); friend != nil {
			result = append(result, friend)
		}
	}
	return result
}

// FuzzyGetByRemarkOrNick 根据备注或昵称模糊匹配好友列表
func (f Friends) FuzzyGetByRemarkOrNick(remarkOrNick string) Friends {
	var result Friends
	for _, friend := range f {
		if strings.Contains(friend.Remark, remarkOrNick) || strings.Contains(friend.Nick, remarkOrNick) {
			result = append(result, friend)
		}
	}
	return result
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

// AsUsers 将群组列表转换为User列表
func (g Groups) AsUsers() []*User {
	users := make([]*User, len(g))
	for i, group := range g {
		users[i] = group.AsUser()
	}
	return users
}

// GetByWxId 根据微信ID获取群组
func (g Groups) GetByWxId(wxId string) *Group {
	for _, group := range g {
		if group.WxId == wxId {
			return group
		}
	}
	return nil
}

// GetByWxNum 根据微信号获取群组
func (g Groups) GetByWxNum(wxNum string) *Group {
	for _, group := range g {
		if group.WxNum == wxNum {
			return group
		}
	}
	return nil
}

// GetByNick 根据昵称获取群组
func (g Groups) GetByNick(nick string) *Group {
	for _, group := range g {
		if group.Nick == nick {
			return group
		}
	}
	return nil
}

// GetByRemark 根据备注获取群组
func (g Groups) GetByRemark(remark string) *Group {
	for _, group := range g {
		if group.Remark == remark {
			return group
		}
	}
	return nil
}

// GetByRemarkOrNick 根据备注或昵称获取群组
func (g Groups) GetByRemarkOrNick(remarkOrNick string) *Group {
	for _, group := range g {
		if group.Remark == remarkOrNick || group.Nick == remarkOrNick {
			return group
		}
	}
	return nil
}

// GetByWxIds 根据微信ID列表获取群组列表
func (g Groups) GetByWxIds(wxIds []string) Groups {
	var result Groups
	for _, wxId := range wxIds {
		if group := g.GetByWxId(wxId); group != nil {
			result = append(result, group)
		}
	}
	return result
}

// GetByWxNums 根据微信号列表获取群组列表
func (g Groups) GetByWxNums(wxNums []string) Groups {
	var result Groups
	for _, wxNum := range wxNums {
		if group := g.GetByWxNum(wxNum); group != nil {
			result = append(result, group)
		}
	}
	return result
}

// GetByNicks 根据昵称列表获取群组列表
func (g Groups) GetByNicks(nicks []string) Groups {
	var result Groups
	for _, nick := range nicks {
		if group := g.GetByNick(nick); group != nil {
			result = append(result, group)
		}
	}
	return result
}

// GetByRemarks 根据备注列表获取群组列表
func (g Groups) GetByRemarks(remarks []string) Groups {
	var result Groups
	for _, remark := range remarks {
		if group := g.GetByRemark(remark); group != nil {
			result = append(result, group)
		}
	}
	return result
}

// GetByRemarkOrNicks 根据备注或昵称列表获取群组列表
func (g Groups) GetByRemarkOrNicks(remarkOrNicks []string) Groups {
	var result Groups
	for _, remarkOrNick := range remarkOrNicks {
		if group := g.GetByRemarkOrNick(remarkOrNick); group != nil {
			result = append(result, group)
		}
	}
	return result
}

// FuzzyGetByRemarkOrNick 根据备注或昵称模糊匹配好友列表
func (g Groups) FuzzyGetByRemarkOrNick(remarkOrNick string) Groups {
	var result Groups
	for _, group := range g {
		if strings.Contains(group.Remark, remarkOrNick) || strings.Contains(group.Nick, remarkOrNick) {
			result = append(result, group)
		}
	}
	return result
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

// AsUsers 将当前对象转换为User对象
func (m MPs) AsUsers() []*User {
	var users []*User
	for _, mp := range m {
		users = append(users, mp.AsUser())
	}
	return users
}

// GetByWxId 根据微信ID获取公众号
func (m MPs) GetByWxId(wxId string) *MP {
	for _, mp := range m {
		if mp.WxId == wxId {
			return mp
		}
	}
	return nil
}

// GetByWxNum 根据微信号获取公众号
func (m MPs) GetByWxNum(wxNum string) *MP {
	for _, mp := range m {
		if mp.WxNum == wxNum {
			return mp
		}
	}
	return nil
}

// GetByNick 根据昵称获取公众号
func (m MPs) GetByNick(nick string) *MP {
	for _, mp := range m {
		if mp.Nick == nick {
			return mp
		}
	}
	return nil
}

// GetByRemark 根据备注获取公众号
func (m MPs) GetByRemark(remark string) *MP {
	for _, mp := range m {
		if mp.Remark == remark {
			return mp
		}
	}
	return nil
}

// GetByRemarkOrNick 根据备注或昵称获取公众号
func (m MPs) GetByRemarkOrNick(remarkOrNick string) *MP {
	for _, mp := range m {
		if mp.Remark == remarkOrNick || mp.Nick == remarkOrNick {
			return mp
		}
	}
	return nil
}

// GetByWxIds 根据微信ID列表获取公众号列表
func (m MPs) GetByWxIds(wxIds []string) MPs {
	var result MPs
	for _, wxId := range wxIds {
		if mp := m.GetByWxId(wxId); mp != nil {
			result = append(result, mp)
		}
	}
	return result
}

// GetByWxNums 根据微信号列表获取公众号列表
func (m MPs) GetByWxNums(wxNums []string) MPs {
	var result MPs
	for _, wxNum := range wxNums {
		if mp := m.GetByWxNum(wxNum); mp != nil {
			result = append(result, mp)
		}
	}
	return result
}

// GetByNicks 根据昵称列表获取公众号列表
func (m MPs) GetByNicks(nicks []string) MPs {
	var result MPs
	for _, nick := range nicks {
		if mp := m.GetByNick(nick); mp != nil {
			result = append(result, mp)
		}
	}
	return result
}

// GetByRemarks 根据备注列表获取公众号列表
func (m MPs) GetByRemarks(remarks []string) MPs {
	var result MPs
	for _, remark := range remarks {
		if mp := m.GetByRemark(remark); mp != nil {
			result = append(result, mp)
		}
	}
	return result
}

// GetByRemarkOrNicks 根据备注或昵称列表获取公众号列表
func (m MPs) GetByRemarkOrNicks(remarkOrNicks []string) MPs {
	var result MPs
	for _, remarkOrNick := range remarkOrNicks {
		if mp := m.GetByRemarkOrNick(remarkOrNick); mp != nil {
			result = append(result, mp)
		}
	}
	return result
}

// FuzzyGetByRemarkOrNick 根据备注或昵称模糊匹配好友列表
func (m MPs) FuzzyGetByRemarkOrNick(remarkOrNick string) MPs {
	var result MPs
	for _, mp := range m {
		if strings.Contains(mp.Remark, remarkOrNick) || strings.Contains(mp.Nick, remarkOrNick) {
			result = append(result, mp)
		}
	}
	return result
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
