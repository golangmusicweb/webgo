package entity

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"webgo/setting"
)

func GetDbEngine(name string) *xorm.Engine {
	var config setting.Config
	config.LoadConfig()

	db, ok := config.Datasource[name]
	if ok == false {
		fmt.Printf("The dbname '%s' does not exists!", db)
		panic(ok)
	}
	engine, err := xorm.NewEngine(db["driveName"], db["dataSourceName"])
	if err != nil {
		panic(err)
	}

	return engine
}
