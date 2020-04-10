package mysqldb

import (
	"fmt"

	// import mysql driver fo gorm
	"context"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinmukeji/go-pkg/mysqldb"
	"github.com/jinmukeji/plat-pkg/store"
)

// DbClient 是数据访问管理器
type DbClient struct {
	*store.MySqlStore
}

// NewDbClient 根据传入的 options 返回一个新的 DbClient
func NewDbClient(opts ...mysqldb.Option) (*DbClient, error) {
	options := mysqldb.NewOptions(opts...)

	db, err := mysqldb.NewDbClientFromOptions(options)
	if err != nil {
		return nil, err
	}
	return &DbClient{store.NewMySqlStore(db)}, nil
}

// SafeCloseDB  安全的关闭数据库连接
func (db *DbClient) SafeCloseDB(ctx context.Context) {
	err := db.DB(ctx).Close()
	if err != nil {
		fmt.Println("Closing DB error:", err)
	}
}

// 实现 Datastore 行为
var _ Datastore = (*DbClient)(nil)