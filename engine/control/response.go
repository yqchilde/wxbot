package control

import (
	"fmt"
)

var respCache = make(map[string]bool)

func (manager *Manager[CTX]) initResponse() error {
	return manager.D.Table("__resp").AutoMigrate(&BotResponseConfig{})
}

func (manager *Manager[CTX]) Response(gid string) error {
	if manager.CanResponse(gid) {
		return fmt.Errorf("group-%s is already response", gid)
	}
	manager.Lock()
	defer manager.Unlock()
	respCache[gid] = true
	return manager.D.Table("__resp").Create(&BotResponseConfig{GroupID: gid}).Error
}

func (manager *Manager[CTX]) Silence(gid string) error {
	if !manager.CanResponse(gid) {
		return fmt.Errorf("group-%s is already silence", gid)
	}
	manager.Lock()
	defer manager.Unlock()
	respCache[gid] = false
	return manager.D.Table("__resp").Where("gid = ?", gid).Delete(&BotResponseConfig{}).Error
}

func (manager *Manager[CTX]) CanResponse(gid string) bool {
	manager.RLock()
	isResp, ok := respCache["all"]
	manager.RUnlock()
	if !ok {
		var r BotResponseConfig
		if err := manager.D.Table("__resp").Where("gid = ?", "all").Find(&r).Error; err != nil {
			manager.Lock()
			respCache["all"] = r.Status
			manager.Unlock()
			isResp, ok = r.Status, true
		}
	}
	if ok && isResp {
		return true
	}

	manager.RLock()
	isResp, ok = respCache[gid]
	manager.RUnlock()
	if !ok {
		var r BotResponseConfig
		if err := manager.D.Table("__resp").Where("gid = ?", gid).Find(&r).Error; err != nil {
			manager.Lock()
			respCache[gid] = r.Status
			manager.Unlock()
			isResp, ok = r.Status, true
		}
	}
	return ok && isResp
}
