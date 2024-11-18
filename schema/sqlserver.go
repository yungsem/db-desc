package schema

import (
	"github.com/yungsem/db-desc/database"
)

const (
	listTableSqlSqlserver = `
		SELECT DISTINCT
			d.name AS TABLE_NAME,
			f.value AS TABLE_COMMENT 
		FROM
			syscolumns a
			LEFT JOIN systypes b ON a.xusertype= b.xusertype
			INNER JOIN sysobjects d ON a.id= d.id 
			AND d.xtype= 'U' 
			AND d.name != 'dtproperties'
			LEFT JOIN syscomments e ON a.cdefault= e.id
			LEFT JOIN sys.extended_properties g ON a.id= G.major_id 
			AND a.colid= g.minor_id
			LEFT JOIN sys.extended_properties f ON d.id= f.major_id 
			AND f.minor_id= 0
	`
	listColumnSqlSqlServer = `
		SELECT 
			t.name AS TABLE_NAME,
			c.name AS NAME,
			ty.name AS KIND,
			c.max_length AS LENGTH,
			c.precision AS PRECISION,
				CASE WHEN c.is_nullable  = 1 THEN '是' ELSE '否' END AS NULL_FLAG,
			isnull(dc.definition, '' ) AS DEFAULT_VALUE,
			ep.value AS COMMENTS,
			CASE WHEN ic.column_id IS NULL THEN '否' ELSE '是' END AS PK_FLAG
		FROM 
			sys.tables t
		INNER JOIN 
			sys.columns c ON t.object_id = c.object_id
		INNER JOIN 
			sys.types ty ON c.system_type_id = ty.system_type_id
		LEFT JOIN 
			sys.default_constraints dc ON c.default_object_id = dc.object_id
		LEFT JOIN 
			sys.extended_properties ep ON ep.major_id = c.object_id AND ep.minor_id = c.column_id AND ep.class = 1 AND ep.name = 'MS_Description'
		LEFT JOIN 
			sys.indexes i ON t.object_id = i.object_id AND i.is_primary_key = 1
		LEFT JOIN 
			sys.index_columns ic ON i.object_id = ic.object_id AND c.column_id = ic.column_id AND i.index_id = ic.index_id
		WHERE 
    		ty.name <> 'sysname'
		ORDER BY 
			TABLE_NAME,
			c.column_id
	`
)

// Sqlserver 表示 Sqlserver 数据库
type Sqlserver struct {
	Schema
}

// NewSqlserver 创建 sqlserver 实例
func NewSqlserver(db *database.DB) (*Sqlserver, error) {
	tableInfos, err := listAllTable(db.DB, db.SchemaName, listTableSqlSqlserver)
	if err != nil {
		return nil, err
	}

	columnInfos, err := listAllColumn(db.DB, db.SchemaName, listColumnSqlSqlServer)
	if err != nil {
		return nil, err
	}

	return &Sqlserver{
		Schema: Schema{
			db:          db,
			tableInfos:  tableInfos,
			columnInfos: columnInfos,
		},
	}, nil
}
