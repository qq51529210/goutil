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
	traceID string
}

// NewGORMLog 用于接收 gorm 的日志
func NewGORMLog(traceID string) logger.Interface {
	return &gormLog{traceID: traceID}
}

func (g *gormLog) LogMode(logger.LogLevel) logger.Interface {
	return g
}

func (g *gormLog) Info(ctx context.Context, str string, args ...interface{}) {
	log.InfoTrace(g.traceID, str)
}

func (g *gormLog) Warn(ctx context.Context, str string, args ...interface{}) {
	log.WarnTrace(g.traceID, str)
}

func (g *gormLog) Error(ctx context.Context, str string, args ...interface{}) {
	log.ErrorTrace(g.traceID, str)
}

func (g *gormLog) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	log.DebugfTrace(g.traceID, "%s cost %v", sql, time.Since(begin))
	//
	if err != nil {
		log.ErrorTrace(g.traceID, err)
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
//	  H *int8 `gq:"null"` if H==0/1 db.Where("`H` IS NULL/IS NOT NULL")
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
		p = strings.TrimPrefix(tn, "null=")
		if p != tn {
			a, ok := fv.Interface().(int8)
			if ok {
				if a == 0 {
					db = db.Where(fmt.Sprintf("`%s` IS NULL", p))
				} else {
					db = db.Where(fmt.Sprintf("`%s` IS NOT NULL", p))
				}
			}
			continue
		}
		if tn == "null" {
			a, ok := fv.Interface().(int8)
			if ok {
				if a == 0 {
					db = db.Where(fmt.Sprintf("`%s` IS NULL", ft.Name))
				} else {
					db = db.Where(fmt.Sprintf("`%s` IS NOT NULL", ft.Name))
				}
			}
			continue
		}
	}
	//
	return db
}

// GORMTime 创建和更新时间
type GORMTime struct {
	// 创建时间戳，单位秒
	CreatedAt int64 `json:"createdAt" gorm:""`
	// 更新时间戳，单位秒
	UpdatedAt int64 `json:"updatedAt" gorm:""`
}

// GORMPageQuery 分页查询参数
type GORMPageQuery struct {
	// 偏移，小于 0 不匹配
	Offset *int `json:"offset,omitempty" form:"offset" binding:"omitempty,min=0"`
	// 条数，小于 1 不匹配
	Count *int `json:"count,omitempty" form:"count" binding:"omitempty,min=1"`
	// 排序，"column [desc]"
	Order string `json:"order,omitempty" form:"order"`
	// 是否需要返回总数
	Total int8 `json:"total,omitempty" form:"order" binding:"omitempty,oneof=0 1"`
}

// GORMPageResult 是 GORMPage 的返回值
type GORMPageResult[M any] struct {
	// 总数
	Total int64 `json:"total,omitempty"`
	// 列表
	Data []M `json:"data"`
}

// GORMPage 用于分页查询
func GORMPage[M any](db *gorm.DB, page *GORMPageQuery, res *GORMPageResult[M]) (err error) {
	if page != nil {
		// 总数
		if page.Total == 1 {
			err = db.Count(&res.Total).Error
			if err != nil {
				return err
			}
		}
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
	// 分页
	err = db.Find(&res.Data).Error
	if err != nil {
		return err
	}
	//
	return nil
}

// GORMAll 用于查询全部
func GORMAll[M any](db *gorm.DB, query any) (ms []M, err error) {
	// 查询条件
	if query != nil {
		db = InitGORMQuery(db, query)
	}
	// 查询
	err = db.Scan(&ms).Error
	return
}

// GORMDA 封装访问
type GORMDA[K, M any] struct {
	DB *gorm.DB
	M  M
}

// Add 添加
func (g *GORMDA[K, M]) Add(ctx context.Context, m M) (int64, error) {
	db := g.DB.WithContext(ctx).Create(m)
	return db.RowsAffected, db.Error
}

// Update 更新
func (g *GORMDA[K, M]) Update(ctx context.Context, m M) (int64, error) {
	db := g.DB.WithContext(ctx).Updates(m)
	return db.RowsAffected, db.Error
}

// Get 获取
func (g *GORMDA[K, M]) Get(ctx context.Context, m M) error {
	return g.DB.WithContext(ctx).Scan(m).Error
}

// Delete 删除
func (g *GORMDA[K, M]) Delete(ctx context.Context, k K) (int64, error) {
	db := g.DB.WithContext(ctx).Where("`ID`=?", k).Delete(g.M)
	return db.RowsAffected, db.Error
}

// BatchDelete 删除
func (g *GORMDA[K, M]) BatchDelete(ctx context.Context, k []K) (int64, error) {
	db := g.DB.WithContext(ctx).Where("`ID` IN ?", k).Delete(g.M)
	return db.RowsAffected, db.Error
}
