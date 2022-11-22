package db

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/colynn/pontus/config"

	_ "github.com/go-sql-driver/mysql" //加载mysql
	"github.com/jinzhu/gorm"
)

// Eloquent ..
var Eloquent *gorm.DB

// Init ..
func Init() {
	c := config.GetConfig()
	host := c.GetString("database.host")
	port := c.GetInt("database.port")
	dbType := c.GetString("database.dbtype")
	database := c.GetString("database.database")
	username := c.GetString("database.username")
	password := c.GetString("database.password")

	if dbType != "mysql" {
		fmt.Println("db type unknown")
	}
	var err error

	// TODO: [the way of get conn] maybe database have multiple driver, should be change to interface
	conn := GetMysqlConn(username, password, host, database, port)
	var db Database
	if dbType == "mysql" {
		db = new(Mysql)
		Eloquent, err = db.Open(dbType, conn)
	} else {
		panic("db type unknown")
	}

	if err != nil {
		log.Printf("%s connect error %v", dbType, err)
		os.Exit(1)
	}
	log.Printf("%s connect success!", dbType)
	if Eloquent.Error != nil {
		log.Printf("database error %v", Eloquent.Error)
		os.Exit(1)
	}
	if os.Getenv("ENV") == "local" {
		Eloquent.LogMode(false)
	}
}

// GetMysqlConn ..
func GetMysqlConn(username, passwd, host, database string, port int) string {
	var conn bytes.Buffer
	conn.WriteString(username)
	conn.WriteString(":")
	conn.WriteString(passwd)
	conn.WriteString("@tcp(")
	conn.WriteString(host)
	conn.WriteString(":")
	conn.WriteString(strconv.Itoa(port))
	conn.WriteString(")")
	conn.WriteString("/")
	conn.WriteString(database)
	// TODO: `loc` get from config
	conn.WriteString("?charset=utf8&parseTime=True&loc=Asia%2FShanghai&timeout=1000ms")
	return conn.String()
}

// Database ..
type Database interface {
	Open(dbType string, conn string) (db *gorm.DB, err error)
}

// Mysql ..
type Mysql struct {
}

// Open ..
func (*Mysql) Open(dbType string, conn string) (db *gorm.DB, err error) {
	eloquent, err := gorm.Open(dbType, conn)
	return eloquent, err
}
