package dbutils

import (
	"fmt"
	"github.com/cengsin/oracle"
	"github.com/sunmfei/mfus/common/MFei"
	"github.com/sunmfei/mfus/config/conf"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"strconv"
	"time"
)

// db连接

// Setup 初始化连接
func Setup(typ string) *gorm.DB {
	// sun = newConnection()&parseTime=True&loc=Local
	var dialector gorm.Dialector
	databaseConnConf := conf.NewDatabaseConnConf(typ)
	dbURI := databaseConnConf.Url

	var maxLifetime time.Duration
	maxIdle := 10
	maxOpen := 100
	maxLifetime = 6000 * 1

	var errConv error
	errConv = nil
	if len(databaseConnConf.ConnMaxLifetime) != 0 && databaseConnConf.ConnMaxLifetime != "" {
		var num int
		num, errConv = strconv.Atoi(databaseConnConf.ConnMaxLifetime)
		maxLifetime = time.Duration(num * 1)

	} else if len(databaseConnConf.MaxOpenConn) != 0 && databaseConnConf.MaxOpenConn != "" {
		maxOpen, errConv = strconv.Atoi(databaseConnConf.ConnMaxLifetime)
	} else if len(databaseConnConf.MaxIdleConn) != 0 && databaseConnConf.MaxIdleConn != "" {
		maxIdle, errConv = strconv.Atoi(databaseConnConf.ConnMaxLifetime)
	}
	if errConv != nil {
		MFei.LOGGER.Error("配置数据库错误!", errConv)
		return nil
	}

	MFei.LOGGER.Info("databaseConnConf:", databaseConnConf)

	if databaseConnConf.Type == "oracle" {
		dialector = oracle.Open(dbURI)

	} else if databaseConnConf.Type == "mysql" {
		dialector = mysql.New(mysql.Config{
			DSN:                       dbURI, // data source name
			DefaultStringSize:         256,   // default size for string fields
			DisableDatetimePrecision:  true,  // disable datetime precision, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false, // auto configure based on currently MySQL version
		})
	} else if databaseConnConf.Type == "postgres" {
		dialector = postgres.New(postgres.Config{
			DSN:                  "user=gorm password=gorm dbname=gorm port=9920 sslmode=disable TimeZone=Asia/Shanghai",
			PreferSimpleProtocol: true, // disables implicit prepared statement usage
		})
	} else { // sqlite3
		dbURI = fmt.Sprintf("test.sun")
		dialector = sqlite.Open("test.sun")
	}

	//conn, err := gorm.Open(dialector, &gorm.Config{})
	conn, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{
			SlowThreshold: 1 * time.Millisecond,
			LogLevel:      logger.Info,
			Colorful:      false,
		}),
	})
	if err != nil {
		MFei.LOGGER.Error("get sun server failed.", err)
		return nil
	}
	sqlDB, err := conn.DB()
	if err != nil {
		MFei.LOGGER.Error("connect sun server failed.", err)
		return nil
	} else if err := sqlDB.Ping(); err != nil {
		MFei.LOGGER.Error("connect sun server failed.", err)
	}

	sqlDB.SetMaxIdleConns(maxIdle)        // SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxOpenConns(maxOpen)        // SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetConnMaxLifetime(maxLifetime) // SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	return conn
}
