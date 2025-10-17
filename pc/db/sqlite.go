package db

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/wallet/pc/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	_db *gorm.DB
	om  sync.Once
)

func Sqlite() *gorm.DB {
	return _db
}

const (
	maxIdleConns       = 10
	maxOpenConns       = 500
	connMaxIdleTimeMin = 10
	connMaxLifetimeMin = 10
)

var (
	dbPathEmpty = errors.New("dbPath is empty")
)

func InitSqlite() error {
	var err error
	om.Do(func() {

		c := config.Conf()
		dbPath := c.Db.Path

		if dbPath == "" {
			err = dbPathEmpty
			return
		}

		dbDir := filepath.Dir(dbPath)
		if _, _err := os.Stat(dbDir); os.IsNotExist(_err) {
			_err = os.MkdirAll(dbDir, 0644)
			if _err != nil {
				err = errors.Wrap(_err, "os.MkdirAll")
				return
			}
		}

		dsn := fmt.Sprintf("file:%s?cache=shared&_fk=1", dbPath)

		gd, _err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
		if _err != nil {
			err = errors.Wrap(_err, "gorm.Open")
			return
		}
		sd, _err := gd.DB()
		if _err != nil {
			err = errors.Wrap(_err, "gd.DB()")
			return
		}

		sd.SetMaxIdleConns(maxIdleConns)
		sd.SetMaxOpenConns(maxOpenConns)
		sd.SetConnMaxIdleTime(connMaxIdleTimeMin * time.Minute)
		sd.SetConnMaxLifetime(connMaxLifetimeMin * time.Minute)
		_db = gd
	})
	return err
}
