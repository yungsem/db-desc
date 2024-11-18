package schema

import (
	"github.com/yungsem/db-desc/database"
)

const (
	listTableSqlMysql = `
		SELECT
			TABLE_NAME,
			TABLE_COMMENT 
		FROM
			information_schema.TABLES 
		WHERE
			TABLE_SCHEMA = ?
			AND table_type = 'BASE TABLE'
	`

	listColumnSQlMysql = `
		SELECT
			TABLE_NAME,
			COLUMN_NAME AS 'NAME',
			COLUMN_TYPE AS 'KIND',
			IF(CHARACTER_MAXIMUM_LENGTH IS NULL, NUMERIC_PRECISION, CHARACTER_MAXIMUM_LENGTH) AS 'LENGTH',
			NUMERIC_SCALE AS 'PRECISION',
			IF(IS_NULLABLE = 'YES', '是', '否') AS 'NULL_FLAG',
			COLUMN_DEFAULT AS 'DEFAULT_VALUE',
			IF(COLUMN_NAME = 'ID', '主键ID', COLUMN_COMMENT) AS 'COMMENTS',
			IF(COLUMN_KEY = 'PRI', '是', '否') AS 'PK_FLAG'
		FROM
			information_schema.COLUMNS
		WHERE
			TABLE_SCHEMA = ?
		ORDER BY
			ORDINAL_POSITION
	`
)

// Mysql 表示 mysql 数据库
type Mysql struct {
	Schema
}

// NewMysql 创建 mysql 实例
func NewMysql(db *database.DB) (*Mysql, error) {
	tableInfos, err := listAllTable(db.DB, db.SchemaName, listTableSqlMysql)
	if err != nil {
		return nil, err
	}

	columnInfos, err := listAllColumn(db.DB, db.SchemaName, listColumnSQlMysql)
	if err != nil {
		return nil, err
	}

	return &Mysql{
		Schema: Schema{
			db:          db,
			tableInfos:  tableInfos,
			columnInfos: columnInfos,
		},
	}, nil
}
