package Config

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var (
	DB *gorm.DB
)

func InitConnect()  {
	dsn:="root:root@(127.0.0.1:3306)/test?charset=utf8&parseTime=true&loc=Local"
	DB,_=gorm.Open("mysql",dsn)
	//if err != nil {
	//	panic(err)
	//}
	// 强制限制表明为自己定义的模型名单数形式
	DB.SingularTable(true)

	//student.CreateTable()
	//AutoMigrate()
}

func CloseDB() {
	_ = DB.Close()
}
