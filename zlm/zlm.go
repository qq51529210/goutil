package zlm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mms/cfg"
	"mms/util"
	"mms/util/log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

// 对应接口参数
const (
	True  = "1"
	False = "0"
)

// 流的协议
const (
	RTP  = "rtp"
	RTMP = "rtmp"
	RTSP = "rtsp"
	HLS  = "hls"
	TS   = "ts"
	FMP4 = "fmp4"
)

const (
	// VHost 默认的 vhost
	VHost = "__defaultVhost__"
	// CtxKeyTraceID 用于 log 追踪 id
	CtxKeyTraceID = "TraceID"
	//
	vhostQuery  = "vhost"
	secretQuery = "secret"
	//
	logTraceID = "zlm"
)

const (
	// 正确码
	codeTrue = 0
)

// Init 初始化
func Init() error {
	// 服务
	err := _servers.init()
	if err != nil {
		return err
	}
	// token
	_playToken.init()
	//
	return nil
}

// WriteSnapshotTo 将指定 app 和 stream 的截图快照写到 write
func WriteSnapshotTo(writer io.Writer, app, stream string) error {
	// 打开文件
	filePath := filepath.Join(cfg.Cfg.Media.SnapshotDir, app, stream)
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	// 读取
	_, err = io.Copy(writer, file)
	return err
}

// ReadSnapshot 返回 app 和 stream 的截图快照数据
func ReadSnapshot(buf *bytes.Buffer, app, stream string) error {
	// 文件
	file, err := os.Open(filepath.Join(cfg.Cfg.Media.SnapshotDir, app, stream))
	if err != nil {
		return err
	}
	defer file.Close()
	// 读取
	_, err = io.Copy(buf, file)
	return err
}

// httpCallRes 封装请求
func httpCallRes[ReqQuery, ResData any](ctx context.Context, ser *Server, path string, query *ReqQuery, res *ResData) error {
	return httpCall(ctx, ser, path, query, func(res *http.Response) error {
		// 必须是 200
		if res.StatusCode != http.StatusOK {
			return util.HTTPStatusError(res.StatusCode)
		}
		// 解析
		return json.NewDecoder(res.Body).Decode(res)
	})
}

// httpCall 封装请求
func httpCall[Query any](ctx context.Context, ser *Server, path string, query *Query, onRes func(res *http.Response) error) error {
	// 参数
	q := make(url.Values)
	q.Set(secretQuery, ser.Secret)
	q.Set(vhostQuery, VHost)
	if query != nil {
		q = util.HTTPQuery(query, q)
	}
	apiURL := fmt.Sprintf("%s/index/api/%s?%s", ser.APIBaseURL, path, q.Encode())
	// 请求
	old := time.Now()
	err := util.HTTP[any](ctx, http.DefaultClient, http.MethodGet, apiURL, nil, nil, onRes)
	// 日志
	now := time.Now()
	log.Debugf("api call %s cost %v", apiURL, now.Sub(old))
	if err != nil {
		return fmt.Errorf("api call %s error %s", apiURL, err.Error())
	}
	// 时间
	ser.apiCallTime = &now
	//
	return nil
}
