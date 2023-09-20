package util

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"util/log"

	"gopkg.in/yaml.v3"
)

// ReadCfg 读取第一个运行参数
func ReadCfg(uri string, ptr any) error {
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

// LogCfg 日志的配置
type LogCfg struct {
	log.FileConfig `yaml:",inline"`
	// 标识名称
	Name string `json:"name" yaml:"name" validate:"required,max=32"`
	// 日志头格式
	HeaderFormat string `json:"headerFormat" yaml:"headerFormat" validate:"omitempty,oneof=default fileName filePath"`
	// 禁用的日志级别
	DisableLevel []string `json:"disableLevel" yaml:"disableLevel" validate:"omitempty,dive,oneof=all debug info warn error"`
}

// Init 初始化文件日志
func (c *LogCfg) Init() error {
	file, err := log.NewFile(&c.FileConfig)
	if err != nil {
		return err
	}
	logger := log.NewLogger(file, log.Header(c.HeaderFormat), c.Name)
	logger.DisableLevels(c.DisableLevel)
	log.SetLogger(logger)
	//
	return nil
}

// SerCfg 服务配置
type SerCfg struct {
	// 监听地址
	Addr string `json:"addr" yaml:"addr" validate:"required"`
	// 证书路径
	CertFile string `json:"certFile" yaml:"certFile" validate:"omitempty,filepath"`
	// 证书路径
	KeyFile string `json:"keyFile" yaml:"keyFile" validate:"omitempty,filepath"`
}
