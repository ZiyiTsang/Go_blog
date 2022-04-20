package DBTool

import (
	"Go_blog/pkg/logTool"
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/zalando/go-keyring"
	"time"
)

var DB *sql.DB

func Initialize() {
	initDB()
}
func initDB() {
	var err error
	mysqlPasswd, err := keyring.Get("mysql", "root")
	logTool.CheckError(err)
	mysqlAddress, err := keyring.Get("mysql", "address")
	logTool.CheckError(err)
	config := mysql.Config{
		User:                 "root",
		Passwd:               mysqlPasswd,
		Addr:                 mysqlAddress,
		Net:                  "tcp",
		DBName:               "go_blog",
		AllowNativePasswords: true,
		Timeout:              time.Hour * 2,
		CheckConnLiveness:    true,
	}
	DB, err = sql.Open("mysql", config.FormatDSN())
	logTool.CheckError(err)
	//my mySQL "wait_timeout" shows "7200"(s)=2hour,I set same as it did.
	DB.SetConnMaxLifetime(2 * time.Hour)
	//my mySQL "max_connections" shows 2520,so I set 2000 here.
	DB.SetMaxOpenConns(1000)
	//I think it is ok for more than 10.
	DB.SetMaxIdleConns(40)
	err = DB.Ping()
	logTool.CheckError(err)
	fmt.Println("init DB successful")
}
func createTables() {
	createArticlesSQL := `CREATE TABLE IF NOT EXISTS articles(
    id bigint(20) PRIMARY KEY AUTO_INCREMENT NOT NULL,
    title varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL,
    body longtext COLLATE utf8mb4_unicode_ci
); `

	_, err := DB.Exec(createArticlesSQL)
	logTool.CheckError(err)
}
