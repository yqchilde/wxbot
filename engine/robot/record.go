package robot

import (
	"fmt"
	"time"

	"github.com/yqchilde/wxbot/engine/pkg/sqlite"
	"github.com/yqchilde/wxbot/engine/pkg/utils"
)

// MessageRecord 消息记录表
type MessageRecord struct {
	ID         uint      `gorm:"primarykey"`         // 主键
	Type       string    `gorm:"column:type"`        // 消息类型
	FromWxId   string    `gorm:"column:from_wxid"`   // 消息来源wxid，群消息为群wxid，私聊消息为发送者wxid
	FromNick   string    `gorm:"column:from_nick"`   // 消息来源昵称，群消息为群昵称，私聊消息为发送者昵称
	SenderWxId string    `gorm:"column:sender_wxid"` // 消息具体发送者wxid
	SenderNick string    `gorm:"column:sender_nick"` // 消息具体发送者昵称
	Content    string    `gorm:"column:content"`     // 消息内容
	CreatedAt  time.Time `gorm:"column:created_at"`  // 创建时间
}

var db sqlite.DB

func initMessageRecordDB() error {
	if db.Orm != nil {
		return nil
	}

	dbPath := "data/manager/manager.db"
	if !utils.CheckPathExists(dbPath) {
		return fmt.Errorf("db file not found: %s", dbPath)
	}
	if err := sqlite.Open(dbPath, &db); err != nil {
		return err
	}
	return nil
}

// RecordConditions 历史记录搜索条件
type RecordConditions struct {
	FromWxId   string // 消息来源wxid，群消息为群wxid，私聊消息为对方wxid
	SenderWxId string // 消息具体发送者wxid
	CreatedAt  string // 消息创建时间，格式为yyyy-mm-dd
}

// GetRecordHistory 获取消息记录
func (ctx *Ctx) GetRecordHistory(cond *RecordConditions) ([]MessageRecord, error) {
	if err := initMessageRecordDB(); err != nil {
		return nil, err
	}
	recordDB := db.Orm.Table("__message")
	if cond != nil && cond.FromWxId != "" {
		recordDB.Where("from_wxid = ?", cond.FromWxId)
	}
	if cond != nil && cond.SenderWxId != "" {
		recordDB.Where("sender_wxid = ?", cond.SenderWxId)
	}
	if cond != nil && cond.CreatedAt != "" {
		recordDB.Where("STRFTIME('%Y-%m-%d', created_at, 'localtime') = ?", cond.CreatedAt)
	}

	var msgRecord []MessageRecord
	if err := recordDB.Find(&msgRecord).Error; err != nil {
		return nil, err
	}
	return msgRecord, nil
}
