package control

import (
	"crypto/rand"
	"io"

	"github.com/yqchilde/wxbot/engine/robot"
)

func (manager *Manager) initBase() error {
	// 初始化文件服务秘钥
	key := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, key); err != nil {
		return err
	}

	if err := manager.D.Table("__base").AutoMigrate(&BaseConfig{}); err != nil {
		return err
	}
	baseConfig := BaseConfig{FileSecret: key}
	if err := manager.D.Table("__base").FirstOrCreate(&baseConfig).Error; err != nil {
		return err
	}
	robot.SetFileSecret(baseConfig.FileSecret)
	return nil
}
