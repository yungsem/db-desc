package db

import (
	"github.com/jmoiron/sqlx"
)

// Mysql 表示 mysql 数据库
type Mysql struct {
	db *sqlx.DB
}

// listAllTable 返回 db 中用户空间所有的表
func (r *Mysql) listAllTable() ([]TableInfo, error) {
	sql := `
		SELECT 
			table_name,
			comments
		FROM 
			user_tab_comments
		WHERE 
			table_type = 'TABLE' -- 仅选择表
			AND table_name NOT LIKE 'BIN$%' -- 排除回收站中的表
			AND table_name NOT LIKE 'APEX%' -- 排除 Oracle Application Express 表
			AND table_name NOT LIKE 'MLOG$%' -- 排除物化视图日志表
			AND table_name NOT LIKE 'RUPD$%' -- 排除物化视图日志表
			AND table_name NOT LIKE 'RIMP$%' -- 排除物化视图日志表
			AND table_name NOT LIKE 'REDO%' -- 排除重做日志表
			AND table_name NOT LIKE 'C_OBJ#%' -- 排除系统表
			AND table_name NOT LIKE 'OBJ$%' -- 排除系统表
			AND table_name NOT LIKE 'COL$%' -- 排除系统表
			AND table_name NOT LIKE 'CON$%' -- 排除系统表
			AND table_name NOT LIKE 'DF%' -- 排除系统表
			AND table_name NOT LIKE 'ICOL$%' -- 排除系统表
			AND table_name NOT LIKE 'I_OBJ#%' -- 排除系统表
			AND table_name NOT LIKE 'I_USER#%' -- 排除系统表
			AND table_name NOT LIKE 'TRIGGER$%' -- 排除系统表
			AND table_name NOT LIKE 'LOB$%' -- 排除系统表
			AND table_name NOT LIKE 'NEVER%' -- 排除系统表
			AND table_name NOT LIKE 'RECYCLEBIN%' -- 排除回收站表
			AND table_name NOT LIKE 'RM_$%' -- 排除系统表
			AND table_name NOT LIKE 'DBMS%' -- 排除系统表
			AND table_name NOT LIKE 'PLAN_TABLE' -- 排除系统表
			AND table_name NOT LIKE 'ORA$%' -- 排除系统表
			AND table_name NOT LIKE 'TAB$%' -- 排除系统表
			AND table_name NOT LIKE 'USER$%' -- 排除系统表
			AND table_name NOT LIKE 'TMP$%' -- 排除临时表
			AND table_name NOT LIKE 'XDS%' -- 排除系统表
			AND table_name NOT LIKE 'XS%' -- 排除系统表
			AND table_name NOT LIKE 'WRI$_%' -- 排除系统表
			AND table_name NOT LIKE 'WRH$_%' -- 排除系统表
			AND table_name NOT LIKE 'AWR%' -- 排除系统表
			AND table_name NOT LIKE 'SQLPLUS%' -- 排除系统表
			AND table_name NOT LIKE 'DBA%' -- 排除系统表
			AND table_name NOT LIKE 'DUAL' -- 排除系统表
			AND table_name NOT LIKE 'DUMMY' -- 排除系统表
		ORDER BY 
			table_name;
	`

	var tableInfos []TableInfo

	err := r.db.Select(&tableInfos, sql)
	if err != nil {
		return nil, err
	}

	return tableInfos, nil
}

// listAllColumn 返回 db 中用户空间所有表的所有列
func (r *Mysql) listAllColumn() ([]ColumnInfo, error) {
	sql := `
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

	var columnInfos []ColumnInfo

	err := r.db.Select(&columnInfos, sql)
	if err != nil {
		return nil, err
	}

	return columnInfos, nil
}

// DescribeTable 实现了 Table 接口的 TableInfos 方法
func (r *Mysql) DescribeTable() ([]TableInfo, error) {
	tableInfos, err := r.listAllTable()
	if err != nil {
		return nil, err
	}

	columnInfos, err := r.listAllColumn()
	if err != nil {
		return nil, err
	}

	return makeTableInfo(tableInfos, columnInfos), nil
}
