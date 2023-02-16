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

type Manager struct {
	sync.RWMutex
	M map[string]*Control
	D *gorm.DB
}

func NewManager(dbpath string) (m Manager) {
	i := strings.LastIndex(dbpath, "/")
	if i > 0 {
		if err := os.MkdirAll(dbpath[:i], 0755); err != nil {
			log.Fatal(err)
		}
	}
	dbpath = dbpath + "?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)"
	var db sqlite.DB
	if err := sqlite.Open(dbpath, &db, &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)}); err != nil {
		log.Fatal("open plugins database failed: ", err)
	}
	m = Manager{
		M: map[string]*Control{},
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
func (manager *Manager) Lookup(service string) (*Control, bool) {
	manager.RLock()
	m, ok := manager.M[service]
	manager.RUnlock()
	return m, ok
}

// LookupAll 查找全部插件管理器
func (manager *Manager) LookupAll() map[string]*Control {
	return manager.M
}
