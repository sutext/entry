package pgsql

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sutext.github.io/entry/model"
)

type pgsqlDriver struct {
	dsn     string
	options []gorm.Option
}

func New(dsn string, opts ...gorm.Option) model.Driver {
	return &pgsqlDriver{dsn: dsn, options: opts}
}
func DSN(host, port, username, password, database string) model.Driver {
	return &pgsqlDriver{
		dsn: "host=" + host + " user=" + username + " password=" + password + " dbname=" + database + " port=" + port + " sslmode=disable",
	}
}
func (d *pgsqlDriver) Open() (*gorm.DB, error) {
	di := postgres.Open(d.dsn)
	return gorm.Open(di, d.options...)
}
