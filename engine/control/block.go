package control

var blockCache = make(map[string]bool)

func (manager *Manager) initBlock() error {
	return manager.D.Table("__block").AutoMigrate(&BotBlockConfig{})
}

func (manager *Manager) DoBlock(uid string) error {
	manager.Lock()
	defer manager.Unlock()
	blockCache[uid] = true
	return manager.D.Table("__block").Create(&BotBlockConfig{UserID: uid}).Error
}

func (manager *Manager) DoUnblock(uid string) error {
	manager.Lock()
	defer manager.Unlock()
	blockCache[uid] = false
	return manager.D.Table("__block").Where("uid = ?", uid).Delete(&BotBlockConfig{}).Error
}

func (manager *Manager) IsBlocked(uid string) bool {
	manager.RLock()
	isBlock, ok := blockCache[uid]
	manager.RUnlock()
	if ok {
		return isBlock
	}
	manager.Lock()
	defer manager.Unlock()
	isBlock = manager.D.Table("__block").Where("uid = ?", uid).Find(&BotBlockConfig{}).RowsAffected > 0
	blockCache[uid] = isBlock
	return isBlock
}
