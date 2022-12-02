package control

import (
	"github.com/yqchilde/pkgs/log"
)

// Control 插件控制器
type Control[CTX any] struct {
	Service string          // Service 插件服务名
	Cache   map[string]bool // Cache 缓存
	Options Options[CTX]    // Options 插件配置
	Manager *Manager[CTX]   // Manager 插件管理器
}

func (manager *Manager[CTX]) NewControl(service string, o *Options[CTX]) *Control[CTX] {
	m := &Control[CTX]{
		Service: service,
		Cache:   make(map[string]bool, 16),
		Options: func() Options[CTX] {
			if o == nil {
				return Options[CTX]{}
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
	if err := manager.D.Table(service).Where("gid = ?", "all").Find(&c).Error; err != nil {
		m.Options.DisableOnDefault = c.Enable
	}
	return m
}

// Handler 返回预处理器
func (m *Control[CTX]) Handler(gid, uid string) bool {
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
func (m *Control[CTX]) Enable(groupID string) error {
	c := PluginConfig{GroupID: groupID, Enable: true}
	tx := m.Manager.D.Begin()
	if err := tx.Table(m.Service).Where("gid = ?", groupID).Delete(&PluginConfig{}).Error; err != nil {
		log.Errorf("(plugin) %s enable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return err
	}
	if err := tx.Table(m.Service).Create(&c).Error; err != nil {
		log.Errorf("(plugin) %s enable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	m.Manager.Lock()
	m.Cache[groupID] = true
	m.Manager.Unlock()
	return nil
}

// Disable 使插件在某个群中禁用
func (m *Control[CTX]) Disable(groupID string) error {
	c := PluginConfig{GroupID: groupID, Enable: false}
	tx := m.Manager.D.Begin()
	if err := tx.Table(m.Service).Where("gid = ?", groupID).Delete(&PluginConfig{}).Error; err != nil {
		log.Errorf("(plugin) %s disable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return err
	}
	if err := tx.Table(m.Service).Create(&c).Error; err != nil {
		log.Errorf("(plugin) %s disable in %s failed: %v", m.Service, groupID, err)
		tx.Rollback()
		return err
	}
	tx.Commit()
	m.Manager.Lock()
	m.Cache[groupID] = false
	m.Manager.Unlock()
	return nil
}

// IsEnabledIn 查询开启群组
func (m *Control[CTX]) IsEnabledIn(gid string) bool {
	var c PluginConfig

	{
		// 查询是否全局开启
		m.Manager.RLock()
		yes, ok := m.Cache["all"]
		m.Manager.RUnlock()
		if !ok {
			err := m.Manager.D.Table(m.Service).Where("gid = ?", "all").Find(&c).Error
			if err == nil && c.Enable {
				yes, ok = c.Enable, true
				m.Manager.Lock()
				m.Cache["all"] = yes
				m.Manager.Unlock()
			}
		}
		if ok && yes {
			return true
		}
	}

	// 查询是否单独开启
	{
		m.Manager.RLock()
		yes, ok := m.Cache[gid]
		m.Manager.RUnlock()
		if !ok {
			err := m.Manager.D.Table(m.Service).Where("gid = ?", gid).Find(&c).Error
			if err == nil && c.Enable {
				yes, ok = c.Enable, true
				m.Manager.Lock()
				m.Cache[gid] = yes
				m.Manager.Unlock()
			}
		}
		if ok && yes {
			return true
		}
	}
	return false
}
