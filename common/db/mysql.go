package db

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"github.com/wangghyz/polestar/common/model"
	"github.com/wangghyz/polestar/common/util"
)

var (
	mysqlDB *gorm.DB
)

func init() {
	db := NewMySQLConnectionInstance()
	db.AutoMigrate(&model.SysClient{})
	db.AutoMigrate(&model.SysUser{})
	db.AutoMigrate(&model.SysRole{})
	db.AutoMigrate(&model.SysPermission{})
	db.AutoMigrate(&model.SysUserRole{})
	db.AutoMigrate(&model.SysRolePermission{})
}

// NewMySQLConnectionInstance 获得 Mysql 数据库链接
func NewMySQLConnectionInstance() *gorm.DB {
	if mysqlDB != nil {
		return mysqlDB
	}

	appConfig := util.ApplicationConfig()

	var err error
	mysqlDB, err = gorm.Open("mysql", appConfig.Mysql.Url)
	if err != nil {
		log.Fatalf("数据库打开失败！\n\t%s", err)
	}

	mysqlDB.DB().SetMaxIdleConns(appConfig.Mysql.MaxIdleConns)
	mysqlDB.SingularTable(true)

	// log mode
	mysqlDB.LogMode(appConfig.Mysql.LogMode)

	return mysqlDB
}
