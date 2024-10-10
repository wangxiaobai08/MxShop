package main

import (
	"crypto/sha512"
	"fmt"
	"github.com/anaskhan96/go-password-encoder"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"user_srv/model"
)

var DB *gorm.DB

func main() {

	//连接数据库
	var _ error
	dsn := "root:123456@tcp(127.0.0.1:3306)/mxshop_user_service?charset=utf8mb4&parseTime=True&loc=Local"
	newLogger := logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
		SlowThreshold: time.Second,   // Slow SQL threshold
		LogLevel:      logger.Silent, // Log level
		Colorful:      false,         // Disable color
	})
	//Gorm 有一个 默认 logger 实现，默认情况下，它会打印慢 SQL 和错误
	//Logger 接受的选项不多，您可以在初始化时自定义它
	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic(err)
	}

	//md5加盐密码加密
	options := &password.Options{
		SaltLen:      16,
		Iterations:   100,
		KeyLen:       32,
		HashFunction: sha512.New,
	}
	salt, encodedPwd := password.Encode("admin", options)
	newPassword := fmt.Sprintf("$pbkdf2-sha512$%s$%s", salt, encodedPwd)

	for i := 0; i < 10; i++ {
		user := model.User{
			NickName: fmt.Sprintf("nickname:%d", i),
			Mobile:   fmt.Sprintf("13655035209%d", i),
			Password: newPassword,
		}
		DB.Save(&user)
	}
}
