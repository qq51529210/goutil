package gorm

import (
	"context"
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

// IsDataNotFound 是否没有数据
func IsDataNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}

// IsDuplicateError 判断是否重复 key
func IsDuplicateError(err error) bool {
	if e, ok := err.(*mysql.MySQLError); ok {
		return e.Number == 1062
	}
	return false
}

type DB[V any] struct {
	D *gorm.DB
	M V
}

// Get 单个
func (m *DB[V]) Get(ctx context.Context, v V, fields ...string) error {
	db := m.D.WithContext(ctx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	return db.Take(v).Error
}

// First 第一个
func (m *DB[V]) First(ctx context.Context, v V, q any, fields ...string) error {
	db := m.D.WithContext(ctx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	if q != nil {
		db = InitQuery(db, q)
	}
	return db.First(v).Error
}

// Add 添加
func (m *DB[V]) Add(ctx context.Context, v V) error {
	return m.D.WithContext(ctx).Create(v).Error
}

// BatchAdd 批量添加
func (m *DB[V]) BatchAdd(ctx context.Context, vs []V) error {
	return m.D.WithContext(ctx).Create(vs).Error
}

// Update 更新
func (m *DB[V]) Update(ctx context.Context, v V, fields ...string) (int64, error) {
	db := m.D.WithContext(ctx)
	if len(fields) > 0 {
		db = db.Select(fields)
	}
	db = db.Updates(v)
	return db.RowsAffected, db.Error
}

// BatchUpdate 批量更新
func (m *DB[V]) BatchUpdate(ctx context.Context, vs []V, fields ...string) (int64, error) {
	var row int64
	return row, m.D.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, v := range vs {
			db := tx
			if len(fields) > 0 {
				db = db.Select(fields)
			}
			db = db.Updates(v)
			if db.Error != nil {
				return db.Error
			}
			row += db.RowsAffected
		}
		return nil
	})
}

// Delete 删除
func (m *DB[V]) Delete(ctx context.Context, v V) (int64, error) {
	db := m.D.WithContext(ctx).Delete(v)
	return db.RowsAffected, db.Error
}

// BatchDelete 批量删除，v 做为 table name 使用
func (m *DB[V]) BatchDelete(ctx context.Context, query any) (int64, error) {
	db := InitQuery(m.D.WithContext(ctx), query).Delete(m.M)
	return db.RowsAffected, db.Error
}

// Page 分页
func (m *DB[V]) Page(ctx context.Context, page *PageQuery, query any, res *PageResult[V]) error {
	db := m.D.WithContext(ctx)
	if query != nil {
		db = InitQuery(db, query)
	}
	return Page(db, page, res)
}

// All 所有
func (m *DB[V]) All(ctx context.Context, query any) ([]V, error) {
	return All[V](m.D.WithContext(ctx), query)
}
