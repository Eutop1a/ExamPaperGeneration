package mapper

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type Mysql struct {
	Username string
	Password string
	Host     string
	Port     string
	Database string
}

func init() {
	c := Mysql{
		Username: "root",
		Password: "123456",
		Host:     "127.0.0.1",
		Port:     "3306",
		Database: "graduation",
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		panic("failed to connect database")
	}
	DB = db
}
