package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ReadCfg 读取配置，如果没有入参，则读取与 app 同目录下的 cfg.yaml
func ReadCfg(ptr any) error {
	// 默认是 app.cfg
	var uri string
	if len(os.Args) > 1 {
		uri = os.Args[1]
	} else {
		p, err := filepath.Abs(os.Args[0])
		if err != nil {
			return err
		}
		uri = filepath.Join(filepath.Dir(p), "cfg.yaml")
	}
	// 不同方式加载
	_u, err := url.Parse(uri)
	if err != nil {
		return err
	}
	switch _u.Scheme {
	case "http", "https":
		return ReadHTTPCfg(uri, ptr)
	default:
		ext := filepath.Ext(uri)
		switch ext {
		case ".json":
			return ReadJSONCfg(uri, ptr)
		case ".yaml", ".yml":
			return ReadYAMLCfg(uri, ptr)
		default:
			return fmt.Errorf("unsupported config type %s", ext)
		}
	}
}

// ReadHTTPCfg 读取 http 的 json 或者 yaml 数据并解析到 ptr
func ReadHTTPCfg(uri string, ptr any) error {
	// 下载
	res, err := http.Get(uri)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	// 看看是什么
	contentType := res.Header.Get("Content-Type")
	if strings.Contains(contentType, "application/json") {
		// 解析
		return json.NewDecoder(res.Body).Decode(ptr)
	}
	if strings.Contains(contentType, "application/yaml") ||
		strings.Contains(contentType, "application/x-yaml") {
		// 解析
		return yaml.NewDecoder(res.Body).Decode(ptr)
	}
	//
	return fmt.Errorf("unsupported content type %s", contentType)
}

// ReadJSONCfg 读取 json 格式的文件并解析到 ptr
func ReadJSONCfg(path string, ptr any) error {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// 解析
	return json.NewDecoder(file).Decode(ptr)
}

// ReadYAMLCfg 读取 yaml 格式的文件并解析到 ptr
func ReadYAMLCfg(path string, ptr any) error {
	// 打开文件
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()
	// 解析
	return yaml.NewDecoder(file).Decode(ptr)
}

var (
	// ReadEnvTag 是 ReadEnv 解析 tag 的名称
	ReadEnvTag = "env"
)

// ReadEnv 读取环境变量
func ReadEnv(ptr any) {
	readEnv(reflect.ValueOf(ptr), ReadEnvTag)
}

// ReadEnvWithTag 读取环境变量
func ReadEnvWithTag(ptr any, tag string) {
	readEnv(reflect.ValueOf(ptr), tag)
}

func readEnv(v reflect.Value, tagName string) bool {
	// 检查，结构体指针
	vk := v.Kind()
	if vk != reflect.Pointer {
		panic("input must be pointer")
	}
	v = v.Elem()
	vk = v.Kind()
	if vk != reflect.Struct {
		panic("input must be struct type")
	}
	hasValue := false
	// 字段
	vt := v.Type()
	for i := 0; i < vt.NumField(); i++ {
		// 字段类型
		ft := vt.Field(i)
		// 不可导出的
		if !ft.IsExported() {
			continue
		}
		// 字段值
		fv := v.Field(i)
		fk := fv.Kind()
		// 是指针
		if fk == reflect.Pointer {
			// 空指针
			if fv.IsNil() {
				// 是结构指针
				if fv.Elem().Kind() == reflect.Struct {
					// 先 new 一个
					_fv := reflect.New(fv.Type().Elem())
					if readEnv(_fv, tagName) {
						fv.Set(_fv)
					}
					continue
				} else {
					fv = fv.Elem()
					// 往下走
				}
			} else {
				if fv.Elem().Kind() == reflect.Struct {
					readEnv(fv, tagName)
					continue
				} else {
					fv = fv.Elem()
					// 往下走
				}
			}
		}
		// 是结构
		if fk == reflect.Struct {
			readEnv(fv.Addr(), tagName)
			continue
		}
		// 标签
		tag := ft.Tag.Get(tagName)
		if tag == "" {
			continue
		}
		// 环境变量的值
		val := os.Getenv(tag)
		if val == "" {
			continue
		}
		switch fk {
		case reflect.String:
			fv.SetString(val)
			hasValue = true
		case reflect.Bool:
			fv.SetBool(val == "true")
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			n, err := strconv.ParseInt(val, 10, 64)
			if err == nil {
				fv.SetInt(n)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			n, err := strconv.ParseUint(val, 10, 64)
			if err == nil {
				fv.SetUint(n)
			}
		case reflect.Float32, reflect.Float64:
			n, err := strconv.ParseFloat(val, 64)
			if err == nil {
				fv.SetFloat(n)
			}
		}
	}
	return hasValue
}
