package control

import (
	"errors"

	"github.com/yqchilde/wxbot/engine/pkg/log"
)

// Control 插件控制器
type Control struct {
	Service string          // Service 插件服务名
	Cache   map[string]bool // Cache 缓存
	Options Options         // Options 插件配置
	Manager *Manager        // Manager 插件管理器
}

func (manager *Manager) NewControl(service string, o *Options) *Control {
	m := &Control{
		Service: service,
		Cache:   make(map[string]bool),
		Options: func() Options {
			if o == nil {
				return Options{}
			}
			return *o
		}(),
		Manager: manager,
	}
	manager.Lock()
	defer manager.Unlock()
	manager.M[service] = m
	if err := manager.D.Table(service).AutoMigrate(&PluginConfig{}); err != nil {
		log.Fatal(err)
	}
	if err := manager.D.Table(service + "ban").AutoMigrate(&PluginBanConfig{}); err != nil {
		log.Fatal(err)
	}
	var c PluginConfig
	if err := manager.D.Table(service).Where("gid = ?", "all").First(&c).Error; err == nil {
		m.Options.DisableOnDefault = c.Enable
	}
	return m
}

// Handler 返回预处理器
func (m *Control) Handler(gid, uid string) bool {
	if m.Manager.IsBlocked(uid) {
		return false
	}
	if gid == "" {
		gid = uid
	}
	if !m.Manager.CanResponse(gid) || m.IsBannedIn(uid, gid) {
		return false
	}
	return m.IsEnabledIn(gid)
}

// Enable 使插件在某个群中启用
func (m *Control) Enable(groupID string) error {
	if groupID != "all" {
		if isEnable, ok := m.IsEnabledAll(true); ok {
			if isEnable {
				return errors.New("该插件已在全局启用")
			}
			return errors.New("该插件已全局禁用，如需启用请先在关闭全局禁用")
		}
	}

	c := PluginConfig{GroupID: groupID, Enable: true}
	tx := m.Manager.D.Begin()
	if err := tx.Table(m.Service).Delete(&PluginConfig{}, "gid = ?", groupID).Error; err != nil {
		log.Errorf("(plugin) %s enable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return errors.New("启用失败")
	}
	if err := tx.Table(m.Service).Create(&c).Error; err != nil {
		log.Errorf("(plugin) %s enable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return errors.New("启用失败")
	}
	tx.Commit()
	m.Manager.Lock()
	m.Cache[groupID] = true
	m.Manager.Unlock()
	return nil
}

// Disable 使插件在某个群中禁用
func (m *Control) Disable(groupID string) error {
	if groupID != "all" {
		if isEnable, ok := m.IsEnabledAll(false); ok {
			if isEnable {
				return errors.New("该插件已在全局禁用")
			}
			return errors.New("该插件已全局启用，如需启用请先在关闭全局启用")
		}
	}

	c := PluginConfig{GroupID: groupID, Enable: false}
	tx := m.Manager.D.Begin()
	if err := tx.Table(m.Service).Delete(&PluginConfig{}, "gid = ?", groupID).Error; err != nil {
		log.Errorf("(plugin) %s disable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return errors.New("禁用失败")
	}
	if err := tx.Table(m.Service).Create(&c).Error; err != nil {
		log.Errorf("(plugin) %s disable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return errors.New("禁用失败")
	}
	tx.Commit()
	m.Manager.Lock()
	m.Cache[groupID] = false
	m.Manager.Unlock()
	return nil
}

// CloseGlobalMode 关闭全局模式
func (m *Control) CloseGlobalMode() error {
	if err := m.Manager.D.Table(m.Service).Delete(&PluginConfig{}, "gid = ?", "all").Error; err != nil {
		log.Errorf("(plugin) %s close global failed: %v", m.Service, err)
		return errors.New("关闭失败")
	}
	m.Manager.Lock()
	delete(m.Cache, "all")
	m.Manager.Unlock()
	return nil
}

// IsEnabledIn 查询开启群组
func (m *Control) IsEnabledIn(gid string) bool {
	m.Manager.RLock()
	isEnable, ok := m.Cache["all"]
	m.Manager.RUnlock()
	if ok {
		return isEnable
	}
	m.Manager.RLock()
	isEnable, ok = m.Cache[gid]
	m.Manager.RUnlock()
	if ok {
		return isEnable
	}
	m.Manager.Lock()
	defer m.Manager.Unlock()
	var c PluginConfig
	if m.Manager.D.Table(m.Service).First(&c, "gid = ?", "all").Error == nil {
		m.Cache["all"] = c.Enable
		return c.Enable
	}
	if m.Manager.D.Table(m.Service).First(&c, "gid = ?", gid).Error == nil {
		m.Cache[gid] = c.Enable
		return c.Enable
	}
	return !m.Options.DisableOnDefault
}

// IsEnabledAll 查询是否全局开启
func (m *Control) IsEnabledAll(enable bool) (isEnable bool, ok bool) {
	m.Manager.RLock()
	isEnable, ok = m.Cache["all"]
	m.Manager.RUnlock()
	if ok {
		return isEnable == enable, ok
	}
	m.Manager.Lock()
	defer m.Manager.Unlock()
	var c PluginConfig
	if m.Manager.D.Table(m.Service).First(&c, "gid = ?", "all").Error == nil {
		m.Cache["all"] = c.Enable
		return c.Enable == enable, ok
	}
	return false, ok
}
