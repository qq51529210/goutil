package zlm

import "context"

// OnServerStarted 处理 zlm 的 on_server_started 回调
func OnServerStarted(ctx context.Context, ip string, cfg map[string]string) {
	// 获取实例
	ser := GetServer(cfg["general.mediaServerId"])
	if ser == nil {
		return
	}
	// 更新
	ser.updateConfig(cfg)
}
