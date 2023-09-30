package zlm

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mms/util"
	"net/http"
	"os"
	"path/filepath"
)

// GetSnapReq 是 GetSnap 的参数
type GetSnapReq struct {
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
func (s *Server) GetSnap(ctx context.Context, req *GetSnapReq, out io.Writer) error {
	err := httpCall[any](ctx, s, apiGetSnap, nil, func(res *http.Response) error {
		// 必须是 200
		if res.StatusCode != http.StatusOK {
			return util.HTTPStatusError(res.StatusCode)
		}
		// 写入
		_, err := io.Copy(out, res.Body)
		return err
	})
	if err != nil {
		return err
	}
	return nil
}

// SaveSnap 保存截图
func (s *Server) SaveSnap(app, stream string, timeout string) error {
	// 创建目录
	dir := filepath.Join(s.SnapshotDir, s.ID, app)
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return err
	}
	// 请求
	var data bytes.Buffer
	err = s.GetSnap(context.Background(), &GetSnapReq{
		TimeoutSec: timeout,
		ExpireSec:  timeout,
		URL:        fmt.Sprintf("rtmp://%s:%s/%s/%s", s.PrivateIP, s.RTMPPort, app, stream),
	}, &data)
	if err != nil {
		return err
	}
	// 保存文件
	file, err := os.OpenFile(filepath.Join(dir, stream), os.O_TRUNC|os.O_CREATE|os.O_SYNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer file.Close()
	// 写入文件
	_, err = io.Copy(file, &data)
	return err
}
