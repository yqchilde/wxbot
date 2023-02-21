package sqlite

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DB struct {
	Orm *gorm.DB
}

func CreateDBFile(dbPath string) error {
	if _, err := os.Stat(filepath.Dir(dbPath)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			return errors.New("创建数据库文件夹失败: " + err.Error())
		}
	}
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		if _, err := os.Create(dbPath); err != nil {
			return errors.New("创建数据库文件失败: " + err.Error())
		}
	}
	return nil
}

// Open 创建数据库连接
func Open(dbPath string, db *DB, opts ...gorm.Option) error {
	if err := CreateDBFile(dbPath); err != nil {
		return err
	}
	_db, err := gorm.Open(sqlite.Open(dbPath), opts...)
	if err != nil {
		return err
	}
	sqlDB, err := _db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxOpenConns(1)
	db.Orm = _db
	return nil
}

// Create 创建数据表
func (d *DB) Create(table string, dst ...interface{}) error {
	return d.Orm.Table(table).AutoMigrate(dst...)
}

// CreateAndFirstOrCreate 创建数据表并创建第一条数据，如果已有一条数据则不创建
func (d *DB) CreateAndFirstOrCreate(table string, dest interface{}, conds ...interface{}) error {
	if err := d.Orm.Table(table).AutoMigrate(dest); err != nil {
		return err
	}
	return d.Orm.Table(table).FirstOrCreate(dest, conds...).Error
}
