package mysql

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sutext.github.io/entry/model"
)

type mysqlDriver struct {
	dsn     string
	options []gorm.Option
}

func New(dsn string, opts ...gorm.Option) model.Driver {
	return &mysqlDriver{dsn: dsn, options: opts}
}
func DSN(host, port, username, password, database string) model.Driver {
	return &mysqlDriver{
		dsn: username + ":" + password + "@tcp(" + host + ":" + port + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local",
	}
}
func (d *mysqlDriver) Open() (*gorm.DB, error) {
	di := mysql.Open(d.dsn)
	return gorm.Open(di, d.options...)
}
