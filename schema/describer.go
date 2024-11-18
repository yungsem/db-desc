package schema

import (
	"github.com/yungsem/db-desc/database"
)

// Describer 接口表示 schema 描述器
// 定义了描述 schema 信息的行为
type Describer interface {
	// Describe 返回所有的表信息
	Describe() []TableInfo
}

// NewDescriber 创建一个 Describer
// 如果 dbType = mysql ，返回的是 Mysql 实例
// 如果 dbType = sqlserver ，返回的是 Sqlserver 实例
// 如果 dbType = oracle ，返回的是 Oracle 实例
// 默认返回 Mysql 实例
func NewDescriber(db *database.DB) (Describer, error) {
	switch db.DBType {
	case database.TypeMysql:
		return NewMysql(db)
	case database.TypeSqlserver:
		return NewSqlserver(db)
	case database.TypeOracle:
		return NewOracle(db)
	default:
		return NewMysql(db)
	}
}
