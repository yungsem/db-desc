package schema

import (
	"github.com/yungsem/db-desc/database"
)

const (
	listTableSqlOracle = `
		SELECT 
			table_name AS TABLE_NAME,
			comments AS TABLE_COMMENT
		FROM 
			user_tab_comments
		WHERE 
			table_type = 'TABLE'
			AND table_name NOT LIKE 'BIN$%'
			AND table_name NOT LIKE 'APEX%'
			AND table_name NOT LIKE 'MLOG$%'
			AND table_name NOT LIKE 'RUPD$%'
			AND table_name NOT LIKE 'RIMP$%'
			AND table_name NOT LIKE 'REDO%'
			AND table_name NOT LIKE 'C_OBJ#%'
			AND table_name NOT LIKE 'OBJ$%'
			AND table_name NOT LIKE 'COL$%'
			AND table_name NOT LIKE 'CON$%'
			AND table_name NOT LIKE 'DF%'
			AND table_name NOT LIKE 'ICOL$%'
			AND table_name NOT LIKE 'I_OBJ#%'
			AND table_name NOT LIKE 'I_USER#%'
			AND table_name NOT LIKE 'TRIGGER$%'
			AND table_name NOT LIKE 'LOB$%'
			AND table_name NOT LIKE 'NEVER%'
			AND table_name NOT LIKE 'RECYCLEBIN%'
			AND table_name NOT LIKE 'RM_$%'
			AND table_name NOT LIKE 'DBMS%'
			AND table_name NOT LIKE 'PLAN_TABLE'
			AND table_name NOT LIKE 'ORA$%'
			AND table_name NOT LIKE 'TAB$%'
			AND table_name NOT LIKE 'USER$%'
			AND table_name NOT LIKE 'TMP$%'
			AND table_name NOT LIKE 'XDS%'
			AND table_name NOT LIKE 'XS%'
			AND table_name NOT LIKE 'WRI$_%'
			AND table_name NOT LIKE 'WRH$_%'
			AND table_name NOT LIKE 'AWR%'
			AND table_name NOT LIKE 'SQLPLUS%'
			AND table_name NOT LIKE 'DBA%'
			AND table_name NOT LIKE 'DUAL'
			AND table_name NOT LIKE 'DUMMY'
		ORDER BY 
			table_name
	`

	listColumnSqlOracle = `
		SELECT 
			tc.TABLE_NAME AS TABLE_NAME, 
			tc.COLUMN_NAME AS NAME, 
			tc.DATA_TYPE AS KIND,
			CASE WHEN tc.DATA_PRECISION IS NOT NULL THEN tc.DATA_PRECISION ELSE tc.DATA_LENGTH END AS LENGTH,
			tc.DATA_SCALE AS PRECISION,
			CASE WHEN tc.NULLABLE = 'N' THEN '否' ELSE '是' END AS NULL_FLAG,
			tc.DATA_DEFAULT AS DEFAULT_VALUE,
			(CASE WHEN tc.COLUMN_NAME = 'ID' THEN '主键ID'
			ELSE cc.COMMENTS END) AS COMMENTS,
			(CASE WHEN tc.COLUMN_NAME = 'ID' THEN '是'
			ELSE '否' END) AS PK_FLAG
		FROM user_tab_columns tc
		LEFT JOIN user_col_comments cc ON tc.TABLE_NAME = cc.TABLE_NAME AND tc.COLUMN_NAME = cc.COLUMN_NAME
		ORDER BY tc.column_id
	`
)

// Oracle 表示 oracle 数据库
type Oracle struct {
	Schema
}

// NewOracle 创建 oracle 实例
func NewOracle(db *database.DB) (*Oracle, error) {
	tableInfos, err := listAllTable(db.DB, db.SchemaName, listTableSqlOracle)
	if err != nil {
		return nil, err
	}

	columnInfos, err := listAllColumn(db.DB, db.SchemaName, listColumnSqlOracle)
	if err != nil {
		return nil, err
	}

	return &Oracle{
		Schema: Schema{
			db:          db,
			tableInfos:  tableInfos,
			columnInfos: columnInfos,
		},
	}, nil
}
