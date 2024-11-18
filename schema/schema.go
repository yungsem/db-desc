package schema

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	"github.com/yungsem/db-desc/database"
)

// TableInfo 表示表信息
type TableInfo struct {
	TableName    string         `db:"TABLE_NAME"`
	TableComment sql.NullString `db:"TABLE_COMMENT"`
	ColumnInfos  []ColumnInfo
}

// ColumnInfo 表示列信息
type ColumnInfo struct {
	TableName    string         `db:"TABLE_NAME"`
	Name         string         `db:"NAME"`
	Kind         string         `db:"KIND"`
	Length       sql.NullString `db:"LENGTH"`
	Precision    sql.NullString `db:"PRECISION"`
	NullFlag     string         `db:"NULL_FLAG"`
	DefaultValue sql.NullString `db:"DEFAULT_VALUE"`
	Comment      sql.NullString `db:"COMMENTS"`
	PkFlag       string         `db:"PK_FLAG"`
}

type Schema struct {
	db          *database.DB
	tableInfos  []TableInfo
	columnInfos []ColumnInfo
}

// convertColumnToMap 对 columnInfos 按 TableName 字段分组，分组后的结果是一个 map ，
// key 为 TableName ，value 为 []ColumnInfo
func (s *Schema) convertColumnToMap(columnInfos []ColumnInfo) map[string][]ColumnInfo {
	result := make(map[string][]ColumnInfo)
	for _, info := range columnInfos {
		result[info.TableName] = append(result[info.TableName], info)
	}
	return result
}

// Describe 实现了 Table 接口的 TableInfos 方法
func (s *Schema) Describe() []TableInfo {
	// 转换成 map
	columnInfoMap := s.convertColumnToMap(s.columnInfos)
	// 组装 tis
	var tis []TableInfo
	for _, ti := range s.tableInfos {
		ti.ColumnInfos = columnInfoMap[ti.TableName]
		tis = append(tis, ti)
	}

	return tis
}

// listAllTable 返回 schema 中用户空间所有的表
func listAllTable(db *sqlx.DB, schemaName string, sql string) ([]TableInfo, error) {
	var tableInfos []TableInfo

	err := db.Select(&tableInfos, sql, schemaName)
	if err != nil {
		return nil, err
	}

	return tableInfos, nil
}

// listAllColumn 返回 schema 中用户空间所有表的所有列
func listAllColumn(db *sqlx.DB, schemaName string, sql string) ([]ColumnInfo, error) {
	var columnInfos []ColumnInfo

	err := db.Select(&columnInfos, sql, schemaName)
	if err != nil {
		return nil, err
	}

	return columnInfos, nil
}
