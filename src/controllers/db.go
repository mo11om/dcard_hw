package controllers

import (
	_ "os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Err error
var connectionString string = "dcard:i_love_dcard@tcp(db:3306)/data?charset=utf8mb4&parseTime=True"

func DBconnect() error {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := os.Getenv("DB")
	Db, Err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if Err != nil {
		return Err
	}
	return nil
}
