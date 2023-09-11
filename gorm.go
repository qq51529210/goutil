package util

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"
	"time"
	"util/log"

	"github.com/go-sql-driver/mysql"

	"github.com/glebarez/sqlite"
	gormmysql "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// InitGORM 初始化连接
func InitGORM(dsn string) (*gorm.DB, error) {
	var cfg gorm.Config
	cfg.NamingStrategy = schema.NamingStrategy{
		SingularTable: true,
		NoLowerCase:   true,
	}
	// mysql
	_dsn := strings.TrimPrefix(dsn, "mysql://")
	if _dsn != dsn {
		return gormMysql(_dsn, &cfg)
	}
	// sqlite
	return gormSqlite(_dsn, &cfg)
}

// gormMysql 初始化 mysql
func gormMysql(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
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

// gormSqlite 初始化 sqlite
func gormSqlite(dsn string, cfg *gorm.Config) (*gorm.DB, error) {
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

type gormLog struct {
}

// NewGORMLog 用于接收 gorm 的日志
func NewGORMLog() logger.Interface {
	return new(gormLog)
}

func (lg *gormLog) LogMode(logger.LogLevel) logger.Interface {
	return lg
}

func (lg *gormLog) Info(ctx context.Context, str string, args ...interface{}) {
	log.Info(str)
}

func (lg *gormLog) Warn(ctx context.Context, str string, args ...interface{}) {
	log.Warn(str)
}

func (lg *gormLog) Error(ctx context.Context, str string, args ...interface{}) {
	log.Error(str)
}

func (lg *gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	log.Debugf("%s cost %v", sql, time.Since(begin))
	//
	if err != nil {
		log.Error(err)
		return
	}
}

var (
	// InitGORMQueryTag 是 InitGORMQuery 解析 tag 的名称
	InitGORMQueryTag = "gq"
)

// InitGORMQuery 将 v 格式化到 where ，全部是 AND ，略过空值
//
//	type query struct {
//	  A *int64 `gq:"eq"` db.Where("`A` = ?", A)
//	  B string `gq:"like"` db.Where("`B` LiKE %%s%%", B)
//	  C *int64 `gq:"gt=A"` db.Where("`A` < ?", C)
//	  D *int64 `gq:"gte=A"` db.Where("`A` <= ?", D)
//	  E *int64 `gq:"lt=A"` db.Where("`A` > ?", E)
//	  F *int64 `gq:"let=A"` db.Where("`A` >= ?", F)
//	  G *int64 `gq:"neq"` db.Where("`G` != ?", G)
//	}
//
// 先这样，以后遇到再加
func InitGORMQuery(db *gorm.DB, q any) *gorm.DB {
	v := reflect.ValueOf(q)
	vk := v.Kind()
	if vk == reflect.Pointer {
		v = v.Elem()
		vk = v.Kind()
	}
	if vk != reflect.Struct {
		panic("v must be struct or struct ptr")
	}
	return initGORMQuery(db, v)
}

func initGORMQuery(db *gorm.DB, v reflect.Value) *gorm.DB {
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		fv := v.Field(i)
		if !fv.IsValid() {
			continue
		}
		fvk := fv.Kind()
		if fvk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				continue
			}
			fv = fv.Elem()
			fvk = fv.Kind()
		}
		// 结构
		if fvk == reflect.Struct {
			initGORMQuery(db, fv)
			continue
		}
		if fvk == reflect.String {
			// 空值
			if fv.IsZero() {
				continue
			}
		}
		ft := vt.Field(i)
		tn := ft.Tag.Get(InitGORMQueryTag)
		p := strings.TrimPrefix(tn, "eq=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` = ?", p), fv.Interface())
			continue
		}
		if tn == "eq" {
			db = db.Where(fmt.Sprintf("`%s` = ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "neq=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` != ?", p), fv.Interface())
			continue
		}
		if tn == "neq" {
			db = db.Where(fmt.Sprintf("`%s` != ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "like=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", p), fmt.Sprintf("%%%v%%", fv.Interface()))
			continue
		}
		if tn == "like" {
			db = db.Where(fmt.Sprintf("`%s` LIKE ?", ft.Name), fmt.Sprintf("%%%v%%", fv.Interface()))
			continue
		}
		p = strings.TrimPrefix(tn, "gt=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` < ?", p), fv.Interface())
			continue
		}
		if tn == "gt" {
			db = db.Where(fmt.Sprintf("`%s` < ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "gte=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` <= ?", p), fv.Interface())
			continue
		}
		if tn == "gte" {
			db = db.Where(fmt.Sprintf("`%s` <= ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "lt=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` > ?", p), fv.Interface())
			continue
		}
		if tn == "lt" {
			db = db.Where(fmt.Sprintf("`%s` > ?", ft.Name), fv.Interface())
			continue
		}
		p = strings.TrimPrefix(tn, "lte=")
		if p != tn {
			db = db.Where(fmt.Sprintf("`%s` >= ?", p), fv.Interface())
			continue
		}
		if tn == "lte" {
			db = db.Where(fmt.Sprintf("`%s` >= ?", ft.Name), fv.Interface())
			continue
		}
	}
	//
	return db
}

// GORMTime 创建和更新时间
type GORMTime struct {
	// 数据库的创建时间，时间戳
	CreatedAt int64 `json:"createdAt" gorm:""`
	// 数据库的更新时间，时间戳
	UpdatedAt int64 `json:"updatedAt" gorm:""`
}

// GORMPageQuery 分页查询参数
type GORMPageQuery struct {
	// 偏移，小于 0 不匹配
	Offset *int `form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count *int `form:"count" binding:"omitempty,min=1"`
	// 排序，"column [desc]"
	Order string `form:"order"`
}

// GORMQuery 是 All 函数格式化查询参数的接口
type GORMQuery interface {
	Init(*gorm.DB) *gorm.DB
}

// GORMPageResult 是 GORMPage 的返回值
type GORMPageResult[M any] struct {
	// 总数
	Total int64 `json:"total"`
	// 列表
	Data []M `json:"data"`
}

// GORMPage 用于分页查询
func GORMPage[M any](db *gorm.DB, page *GORMPageQuery, res *GORMPageResult[M]) error {
	// 总数
	err := db.Count(&res.Total).Error
	if err != nil {
		return err
	}
	if page != nil {
		// 分页
		if page.Offset != nil {
			db = db.Offset(*page.Offset)
		}
		if page.Count != nil {
			db = db.Limit(*page.Count)
		}
		// 排序
		if page.Order != "" {
			db = db.Order(page.Order)
		}
	}
	err = db.Find(&res.Data).Error
	if err != nil {
		return err
	}
	//
	return nil
}
