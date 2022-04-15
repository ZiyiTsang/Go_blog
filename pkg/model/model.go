package model

import (
	"Go_blog/pkg/logTool"
	"github.com/zalando/go-keyring"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() *gorm.DB {
	var err error
	passwd, err := keyring.Get("mysql", "root")
	logTool.CheckError(err)
	address, err := keyring.Get("mysql", "address")
	logTool.CheckError(err)
	my_config := mysql.Config{DSN: "root:" + passwd + "@tcp(" + address + ":3306)/go_blog?charset=utf8&parseTime=True"}
	config := mysql.New(my_config)
	DB, err = gorm.Open(config, &gorm.Config{})
	logTool.CheckError(err)
	return DB
}
