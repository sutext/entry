package sqlite

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"sutext.github.io/entry/model"
)

type sqliteDriver struct {
	dsn     string
	options []gorm.Option
}

func New(dsn string, opts ...gorm.Option) model.Driver {
	return &sqliteDriver{dsn: dsn, options: opts}
}
func Named(name string, opts ...gorm.Option) model.Driver {
	return &sqliteDriver{dsn: "file:" + name + "?cache=shared&mode=rwc", options: opts}
}
func (d *sqliteDriver) Open() (*gorm.DB, error) {
	di := sqlite.Open(d.dsn)
	return gorm.Open(di, d.options...)
}
