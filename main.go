package main

import (
	"github.com/donnie4w/go-logger/logger"
	"github.com/yungsem/db-desc/cnf"
	_ "github.com/yungsem/db-desc/cnf"
	"github.com/yungsem/db-desc/database"
	"github.com/yungsem/db-desc/excel"
	"github.com/yungsem/db-desc/schema"
)

func main() {
	// 初始化日志
	log := logger.NewLogger()
	log.SetFormat(
		logger.FORMAT_LEVELFLAG |
			logger.FORMAT_LONGFILENAME |
			logger.FORMAT_MICROSECONDS |
			logger.FORMAT_FUNC)
	log.SetFormatter("{time} - {level} - {file} - {message} \n")
	log.SetOption(&logger.Option{Console: true, Stacktrace: logger.LEVEL_WARN})

	// 初始化配置
	conf, err := cnf.NewConf()
	if err != nil {
		log.Errorf("init conf error: %v", err)
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
