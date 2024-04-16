package main

import (
	"github.com/yungsem/db-desc/cnf"
	_ "github.com/yungsem/db-desc/cnf"
	"github.com/yungsem/db-desc/db"
	"github.com/yungsem/db-desc/excel"
)

func main() {

	conf, err := cnf.NewConf()
	if err != nil {
		panic(err)
	}

	describer, err := db.NewTableDescriber(conf.DB.Type,
		conf.DB.Host,
		conf.DB.Port,
		conf.DB.Username,
		conf.DB.Password,
		conf.DB.Schema)
	if err != nil {
		panic(err)
	}

	tableInfos, err := describer.DescribeTable()
	if err != nil {
		panic(err)
	}

	err = excel.GenerateExcel(tableInfos)
	if err != nil {
		panic(err)
	}
}
