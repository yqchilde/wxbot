package control

func (manager *Manager) initMessage() error {
	return manager.D.Table("__message").AutoMigrate(&MessageRecord{})
}
