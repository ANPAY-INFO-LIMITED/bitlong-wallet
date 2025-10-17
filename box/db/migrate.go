package db

import "github.com/wallet/box/models"

var (
	mds = []any{
		&models.Lnt{},
		&models.Info{},
		&models.Key{},
		&models.Cpa{},
	}
)

func Migrate() error {

	if err := Sqlite().AutoMigrate(mds...); err != nil {
		return err
	}
	return nil
}
