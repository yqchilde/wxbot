package control

import (
	"errors"
	"fmt"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

var banCache = make(map[string]struct{})

// Ban 禁止某人在某群使用本插件
func (m *Control) Ban(uid, gid string) error {
	if gid != "" {
		label := fmt.Sprintf("%s_%s_%s", m.Service, uid, gid)
		m.Manager.Lock()
		err := m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).FirstOrCreate(&PluginBanConfig{Label: label, UserID: uid, GroupID: gid}).Error
		banCache[label] = struct{}{}
		m.Manager.Unlock()
		if err != nil {
			log.Errorf("(plugin) %s banned in group %s for user %s, failed: %v", m.Service, gid, uid, err)
			return errors.New("ban失败")
		}
		return nil
	}
	// 所有群
	label := fmt.Sprintf("%s_%s_%s", m.Service, uid, "all")
	m.Manager.Lock()
	err := m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).Create(&PluginBanConfig{Label: label, UserID: uid, GroupID: "all"}).Error
	banCache[label] = struct{}{}
	m.Manager.Unlock()
	if err != nil {
		log.Errorf("(plugin) %s banned in all group for user %s, failed: %v", m.Service, gid, uid, err)
		return errors.New("ban失败")
	}
	return nil
}

// UnBan 允许某人在某群使用本插件
func (m *Control) UnBan(uid, gid string) error {
	if gid != "" {
		label := fmt.Sprintf("%s_%s_%s", m.Service, uid, gid)
		m.Manager.Lock()
		err := m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).Delete(&PluginBanConfig{}).Error
		delete(banCache, label)
		m.Manager.Unlock()
		if err != nil {
			log.Errorf("(plugin) %s unbanned in group %s for user %s, failed: %v", m.Service, gid, uid, err)
			return errors.New("unban失败")
		}
		return nil
	}
	// 所有群
	label := fmt.Sprintf("%s_%s_%s", m.Service, uid, "all")
	m.Manager.Lock()
	err := m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).Delete(&PluginBanConfig{}).Error
	delete(banCache, label)
	m.Manager.Unlock()
	if err != nil {
		log.Errorf("(plugin) %s unbanned in all group for user %s, failed: %v", m.Service, gid, uid, err)
		return errors.New("unban失败")
	}
	return nil
}

// IsBannedIn 某人是否在某群被ban
func (m *Control) IsBannedIn(uid, gid string) bool {
	var b PluginBanConfig
	var err error
	if len(gid) != 0 {
		label := fmt.Sprintf("%s_%s_%s", m.Service, uid, gid)
		m.Manager.RLock()
		if _, ok := banCache[label]; ok {
			m.Manager.RUnlock()
			return true
		}
		err = m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).Find(&b).Error
		m.Manager.RUnlock()
		if err == nil && gid == b.GroupID && uid == b.UserID {
			log.Debugf("[control] plugin %s is banned in group %s for user %s", m.Service, b.GroupID, b.UserID)
			m.Manager.Lock()
			banCache[label] = struct{}{}
			m.Manager.Unlock()
			return true
		}
	}
	label := fmt.Sprintf("%s_%s_%s", m.Service, uid, "all")
	m.Manager.RLock()
	if _, ok := banCache[label]; ok {
		m.Manager.RUnlock()
		return true
	}
	err = m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).Find(&b).Error
	m.Manager.RUnlock()
	if err == nil && b.GroupID == "all" && uid == b.UserID {
		log.Debugf("[control] plugin %s is banned in all group for user %s", m.Service, b.UserID)
		m.Manager.Lock()
		banCache[label] = struct{}{}
		m.Manager.Unlock()
		return true
	}
	return false
}
