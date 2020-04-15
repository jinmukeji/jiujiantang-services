package mysqldb

import (
	"github.com/jinmukeji/plat-pkg/v2/store/mysql"
)

// DbClient 是数据访问管理器
type DbClient struct {
	mysql.MySqlStore
}

// NewDbClient 根据传入的 options 返回一个新的 DbClient
func NewDbClient(mySqlStore mysql.MySqlStore) (*DbClient, error) {
	return &DbClient{mySqlStore}, nil
}
