package controllers

// import (
// 	"database/sql"
// 	"errors"
// )

// // Global variable to hold the connection (optional, see considerations)

// var Db *sql.DB

// func Init() error {
// 	var err error
// 	Db, err = ConnectToMySQL()

// 	return err

// }
// func ConnectToMySQL() (*sql.DB, error) {
// 	// Replace with your actual connection details
// 	// connectionString example: "user:password@tcp(host:port)/databaseName"

// 	db, err := sql.Open("mysql", connectionString)
// 	if err != nil {
// 		defer db.Close()
// 		return nil, err

// 	}

// 	err = db.Ping()
// 	if err != nil {
// 		defer db.Close()

// 		return nil, err
// 	}

// 	return db, nil
// }

// func IsConnectionValid(db *sql.DB) error {
// 	// Optional: Add a function to check if the connection is still valid (e.g., ping)
// 	err := db.Ping()
// 	if err != nil {
// 		return errors.New("database connection lost")
// 	}
// 	return nil
// }
import (
	_ "os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Db *gorm.DB
var Err error
var connectionString string = "root:password@tcp(localhost:3306)/data?charset=utf8mb4&parseTime=True&loc=Local"

func DBconnect() {
	// 参考 https://github.com/go-sql-driver/mysql#dsn-data-source-name 获取详情
	// dsn := os.Getenv("DB")
	Db, Err = gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if Err != nil {
		panic(Err)
	}

}
