package util

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"

	"github.com/glebarez/sqlite"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// NewGORMConfig 返回配置
func NewGORMConfig() *gorm.Config {
	var config gorm.Config
	config.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
		NoLowerCase:   true,
	}
	return &config
}

// InitGORM 初始化并返回连接
func InitGORM(uri string, cfg *gorm.Config) (*gorm.DB, error) {
	// mysql
	_uri := strings.TrimPrefix(uri, "mysql://")
	if _uri != uri {
		return gormMysql(_uri, cfg)
	}
	// sqlite
	return gormSqlite(_uri, cfg)
}

func gormMysql(uri string, cfg *gorm.Config) (*gorm.DB, error) {
	// 解析出 schema
	_cfg, err := mysql.ParseDSN(uri)
	if err != nil {
		return nil, err
	}
	_uri := strings.Replace(uri, _cfg.DBName, "", 1)
	_db, err := sql.Open("mysql", _uri)
	if err != nil {
		return nil, err
	}
	// 创建
	_, err = _db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS `%s` DEFAULT CHARACTER SET utf8mb4;", _cfg.DBName))
	if err != nil {
		return nil, err
	}
	_, err = _db.Exec(fmt.Sprintf("USE `%s`;", _cfg.DBName))
	if err != nil {
		return nil, err
	}
	return gorm.Open(gormmysql.New(gormmysql.Config{
		DSNConfig: _cfg,
		Conn:      _db,
	}), cfg)
}

func gormSqlite(uri string, cfg *gorm.Config) (*gorm.DB, error) {
	// sqlite
	db, err := gorm.Open(sqlite.Open(uri), cfg)
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
