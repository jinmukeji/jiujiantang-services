package mysqldb

import "strings"

const (
	// maxQuerySize 最大正整数
	maxQuerySize = 2147483647
)

// sqlEscape 转移 SQL 语句
func sqlEscape(s string) string {
	s = strings.Replace(s, "%", "\\%", -1)
	s = strings.Replace(s, "_", "\\_", -1)
	s = strings.Replace(s, "\\", "\\\\", -1)
	return s
}
