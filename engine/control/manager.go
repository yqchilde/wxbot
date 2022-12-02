package control

import (
	"os"
	"strings"
	"sync"

	"github.com/glebarez/sqlite"
	"github.com/yqchilde/pkgs/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	db, err := gorm.Open(sqlite.Open(dbpath), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		log.Fatal("open plugins database failed: ", err)
	}
	m = Manager[CTX]{
		M: map[string]*Control[CTX]{},
		D: db,
	}
	err = m.initBlock()
	if err != nil {
		panic(err)
	}
	err = m.initResponse()
	if err != nil {
		panic(err)
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
