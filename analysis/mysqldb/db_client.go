package mysqldb

import (
	"fmt"

	// import mysql driver fo gorm
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinmukeji/go-pkg/log/gormlogger"
	"github.com/jinzhu/gorm"
)

// DbClient 是数据访问管理器
type DbClient struct {
	*gorm.DB
	opts Options
}

// NewDbClient 根据传入的 options 返回一个新的 DbClient
func NewDbClient(opts ...Option) (*DbClient, error) {
	options := newOptions(opts...)
	// mysql 连接字符串格式:
	// 	`username:password@tcp(localhost:3306)/db_name?charset=utf8mb4&parseTime=True&loc=utc`
	connStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=%s&parseTime=%t&loc=%s",
		options.Username,
		options.Password,
		options.Address,
		options.Database,
		options.Charset,
		options.ParseTime,
		options.Locale)

	db, err := gorm.Open("mysql", connStr)
	if err != nil {
		return nil, err
	}

	// gorm setting
	db.SingularTable(true)
	db.DB().SetMaxOpenConns(options.MaxConnections)
	db.SetLogger(gormlogger.New(options.Address, options.Database))
	db.LogMode(options.EnableLog)

	return &DbClient{db, options}, nil
}

// Options 返回 DbClient 的 Options.
func (c *DbClient) Options() Options {
	return c.opts
}
