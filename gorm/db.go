package gorm

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/glebarez/sqlite"
	"github.com/go-sql-driver/mysql"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// Open 初始化连接
func Open(dsn string) (*gorm.DB, error) {
	var cfg gorm.Config
	cfg.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
		NoLowerCase:   true,
	}
	// mysql
	_dsn := strings.TrimPrefix(dsn, "mysql://")
	if _dsn != dsn {
		return openMysql(_dsn, &cfg)
	}
	// sqlite
	return openSqlite(_dsn, &cfg)
}

// openMysql 初始化 mysql
func openMysql(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
	// 解析出 schema
	mysqlCfg, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	// 打开连接，不要数据库
	_dsn := strings.Replace(dsn, mysqlCfg.DBName, "", 1)
	db, err := sql.Open("mysql", _dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	// 如果没有就创建数据库
	_, err = db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4;", mysqlCfg.DBName))
	if err != nil {
		return nil, err
	}
	return gorm.Open(gormmysql.Open(dsn), cfg)
}

// openSqlite 初始化 sqlite
func openSqlite(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
	// sqlite
	db, err := gorm.Open(sqlite.Open(dsn), cfg)
	if err != nil {
		return nil, err
	}
	err = db.Exec("PRAGMA foreign_keys = ON;").Error
	if err != nil {
		return nil, err
	}
	//
	return db, nil
}
