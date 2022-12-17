package control

import (
	"os"
	"strings"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/yqchilde/wxbot/engine/pkg/log"
	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
)

type Manager[CTX any] struct {
	sync.RWMutex
	M map[string]*Control[CTX]
	D *gorm.DB
}

func NewManager[CTX any](dbpath string) (m Manager[CTX]) {
	switch {
	case dbpath == "":
		dbpath = "plugins.db"
	case strings.HasSuffix(dbpath, "/"):
		if err := os.MkdirAll(dbpath, 0755); err != nil {
			log.Fatal(err)
		}
		dbpath += "plugins.db"
	default:
		i := strings.LastIndex(dbpath, "/")
		if i > 0 {
			if err := os.MkdirAll(dbpath[:i], 0755); err != nil {
				log.Fatal(err)
			}
		}
	}
	var db sqlite.DB
	if err := sqlite.Open(dbpath, &db, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}); err != nil {
		log.Fatal("open plugins database failed: ", err)
	}
	m = Manager[CTX]{
		M: map[string]*Control[CTX]{},
		D: db.Orm,
	}
	if err := m.initBlock(); err != nil {
		log.Fatal("init block failed: ", err)
	}
	if err := m.initResponse(); err != nil {
		log.Fatal("init response failed: ", err)
	}
	return
}

// Lookup 查找插件管理器
func (manager *Manager[CTX]) Lookup(service string) (*Control[CTX], bool) {
	manager.RLock()
	m, ok := manager.M[service]
	manager.RUnlock()
	return m, ok
}

// LookupAll 查找全部插件管理器
func (manager *Manager[CTX]) LookupAll() map[string]*Control[CTX] {
	return manager.M
}
