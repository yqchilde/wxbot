package control

import (
	"fmt"

	"github.com/yqchilde/pkgs/log"
)

var banCache = make(map[string]struct{})

// Ban 禁止某人在某群使用本插件
func (m *Control[CTX]) Ban(uid, gid string) {
	var err error
	if gid != "" {
		label := fmt.Sprintf("%s_%s_%s", m.Service, uid, gid)
		m.Manager.Lock()
		err = m.Manager.D.Table(m.Service + "ban").Create(&PluginBanConfig{Label: label, UserID: uid, GroupID: gid}).Error
		banCache[label] = struct{}{}
		m.Manager.Unlock()
		if err == nil {
			log.Debugf("[control] plugin %s is banned in group %d for user %d.", m.Service, gid, uid)
			return
		}
	}
	// 所有群
	label := fmt.Sprintf("%s_%s_%s", m.Service, uid, "all")
	m.Manager.Lock()
	err = m.Manager.D.Table(m.Service + "ban").Create(&PluginBanConfig{Label: label, UserID: uid, GroupID: "all"}).Error
	banCache[label] = struct{}{}
	m.Manager.Unlock()
	if err == nil {
		log.Debugf("[control] plugin %s is banned in all group for user %d.", m.Service, uid)
	}
}

// Permit 允许某人在某群使用本插件
func (m *Control[CTX]) Permit(uid, gid string) {
	if gid != "" {
		label := fmt.Sprintf("%s_%s_%s", m.Service, uid, gid)
		m.Manager.Lock()
		m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).Delete(&PluginBanConfig{})
		delete(banCache, label)
		m.Manager.Unlock()
		log.Debugf("[control] plugin %s is permitted in group %d for user %d.", m.Service, gid, uid)
		return
	}
	// 所有群
	label := fmt.Sprintf("%s_%s_%s", m.Service, uid, "all")
	m.Manager.Lock()
	m.Manager.D.Table(m.Service+"ban").Where("label = ?", label).Delete(&PluginBanConfig{})
	delete(banCache, label)
	m.Manager.Unlock()
	log.Debugf("[control] plugin %s is permitted in all group for user %d.", m.Service, uid)
}

// IsBannedIn 某人是否在某群被ban
func (m *Control[CTX]) IsBannedIn(uid, gid string) bool {
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
			log.Debugf("[control] plugin %s is banned in group %d for user %d.", m.Service, b.GroupID, b.UserID)
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
		log.Debugf("[control] plugin %s is banned in all group for user %d.", m.Service, b.UserID)
		m.Manager.Lock()
		banCache[label] = struct{}{}
		m.Manager.Unlock()
		return true
	}
	return false
}
