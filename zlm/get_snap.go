package zlm

import (
	"context"
	gh "goutil/http"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// GetSnapReq 是 GetSnap 的参数
type GetSnapReq struct {
	// http://localhost:8080
	BaseURL string
	// 访问密钥
	Secret string `query:"secret"`
	// 虚拟主机，例如 __defaultVhost__
	VHost string `query:"vhost"`
	// 需要截图的url，可以是本机的，也可以是远程主机的
	URL string `query:"url"`
	// 截图失败超时时间，防止FFmpeg一直等待截图
	TimeoutSec string `query:"timeout_sec"`
	// 截图的过期时间，该时间内产生的截图都会作为缓存返回
	ExpireSec string `query:"expire_sec"`
}

const (
	apiGetSnap = "getSnap"
)

// GetSnap 调用 /index/api/getSnap
// 获取截图或生成实时截图并返回，jpeg格式的图片，可以在浏览器直接打开
func GetSnap(ctx context.Context, req *GetSnapReq, out io.Writer) error {
	url := requestURL(req.BaseURL, apiGetSnap, req)
	// 请求
	old := time.Now()
	err := gh.Request(ctx, http.DefaultClient, http.MethodGet, url, nil, nil,
		func(res *http.Response) error {
			// 必须是 200
			if res.StatusCode != http.StatusOK {
				return gh.StatusError(res.StatusCode)
			}
			// 写入
			_, err := io.Copy(out, res.Body)
			return err
		})
	Logger.DebugfDepth(1, "[%v] %s", time.Since(old), url)
	//
	return err
}

// ReadSnapshot 返回 app 和 stream 的截图快照数据
func ReadSnapshot(dir, app, stream string) ([]byte, error) {
	return os.ReadFile(filepath.Join(dir, app, stream))
}
