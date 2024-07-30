package zlm

import (
	"context"
	ghttp "goutil/http"
	"io"
	"net/http"
	"os"
)

// GetSnapReq 是 GetSnap 的参数
type GetSnapReq struct {
	// 媒体流源地址
	URL string `query:"url"`
	// 超时时间，单位秒
	Timeout string `query:"timeout_sec"`
	// 缓存过期时间，单位秒
	Expire string `query:"expire_sec"`
}

const (
	GetSnapPath = apiPathPrefix + "/getSnap"
)

// GetSnap 调用 /index/api/getSnap ，返回 jpeg 格式的图片
func GetSnap(ctx context.Context, ser Server, req *GetSnapReq) ([]byte, error) {
	var data []byte
	return data, ghttp.JSONRequest(ctx, http.DefaultClient, http.MethodGet,
		ser.BaseURL()+GetSnapPath, ghttp.Query(req, NewRequestQuery(ser)), nil,
		func(res *http.Response) error {
			// 必须是 200
			if res.StatusCode != http.StatusOK {
				return ghttp.StatusError(res.StatusCode)
			}
			// 读取数据
			var err error
			data, err = io.ReadAll(res.Body)
			return err
		})
}

// SaveSnap 保存截图到磁盘
func SaveSnap(ctx context.Context, ser Server, req *GetSnapReq, path string) error {
	// 获取
	d, err := GetSnap(ctx, ser, req)
	if err != nil {
		return err
	}
	// 写入文件
	return os.WriteFile(path, d, os.ModePerm)
}
