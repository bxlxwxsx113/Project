package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	dbhost := beego.AppConfig.String("db.host")
	dbuser := beego.AppConfig.String("db.user")
	dbpassword := beego.AppConfig.String("db.password")
	dbport := beego.AppConfig.String("db.port")
	dbname := beego.AppConfig.String("db.name")
	if dbport == "" {
		dbport = "3306"
	}
	//"root:111111@tcp(127.0.0.1:3306)/HelloBeego?charset=utf8"
	dsn := dbuser + ":" + dbpassword + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"

	fmt.Println("dsn = ", dsn)

	orm.RegisterDataBase("default", "mysql", dsn, 30)

	// register mode
	orm.RegisterModel(new(Admin), new(Auth), new(InfoClass), new(InfoList))
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
}

func TableName(name string) string {
	return beego.AppConfig.String("db.prefix") + name
}
