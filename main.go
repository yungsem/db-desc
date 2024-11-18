package main

import (
	"github.com/yungsem/db-desc/cnf"
	_ "github.com/yungsem/db-desc/cnf"
	"github.com/yungsem/db-desc/database"
	"github.com/yungsem/db-desc/excel"
	"github.com/yungsem/db-desc/schema"
	"log"
)

func main() {
	// 初始化配置
	conf, err := cnf.NewConf()
	if err != nil {
		log.Fatal(err.Error())
	}

	// 初始化 DB
	db, err := database.NewDB(conf.DB.Type,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Schema)
	if err != nil {
		log.Fatal(err.Error())
	}

	// 创建 Describer
	desc, err := schema.NewDescriber(db)
	if err != nil {
		log.Fatal(err.Error())
	}

	// 获取表信息
	tableInfos := desc.Describe()

	// 输出到 excel
	err = excel.GenerateExcel(tableInfos)
	if err != nil {
		log.Fatal(err.Error())
	}
}
