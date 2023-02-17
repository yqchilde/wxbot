package control

import (
	"fmt"
)

var respCache = make(map[string]bool)

func (manager *Manager) initResponse() error {
	return manager.D.Table("__resp").AutoMigrate(&BotResponseConfig{})
}

func (manager *Manager) Response(gid string) error {
	if manager.CanResponse(gid) {
		return fmt.Errorf("wxid-[%s] is already response", gid)
	}
	manager.Lock()
	defer manager.Unlock()
	respCache[gid] = true
	return manager.D.Table("__resp").Where("gid = ?", gid).Delete(&BotResponseConfig{}).Error
}

func (manager *Manager) Silence(gid string) error {
	if !manager.CanResponse(gid) {
		return fmt.Errorf("wxid-[%s] is already silence", gid)
	}
	manager.Lock()
	defer manager.Unlock()
	respCache[gid] = false
	return manager.D.Table("__resp").Create(&BotResponseConfig{GroupID: gid, Status: false}).Error
}

func (manager *Manager) CanResponse(gid string) bool {
	manager.RLock()
	isResp, ok := respCache["all"]
	manager.RUnlock()
	if ok {
		return isResp
	}
	manager.RLock()
	isResp, ok = respCache[gid]
	manager.RUnlock()
	if ok {
		return isResp
	}
	manager.Lock()
	defer manager.Unlock()
	var r BotResponseConfig
	if manager.D.Table("__resp").Where("gid = ?", "all").First(&r).Error == nil {
		respCache["all"] = r.Status
		return r.Status
	}
	if manager.D.Table("__resp").Where("gid = ?", gid).First(&r).Error == nil {
		respCache[gid] = r.Status
		return r.Status
	}
	return true
}
