package util

import (
	"util/log"
)

// LogCfg 日志的配置
type LogCfg struct {
	log.FileConfig `yaml:",inline"`
	// 日志头格式
	HeaderFormat string `json:"headerFormat" yaml:"headerFormat" validate:"omitempty,dive,oneof=default fileName filePath"`
	// 禁用的日志级别
	DisableLevel []string `json:"disableLevel" yaml:"disableLevel" validate:"omitempty,dive,oneof=all debug info warn error"`
}

// InitLog 初始化文件日志
func InitLog(cfg *LogCfg, name string) error {
	file, err := log.NewFile(&cfg.FileConfig)
	if err != nil {
		return err
	}
	logger := log.NewLogger(file, log.Header(cfg.HeaderFormat), name)
	logger.DisableLevels(cfg.DisableLevel)
	log.SetLogger(logger)
	//
	return nil
}
