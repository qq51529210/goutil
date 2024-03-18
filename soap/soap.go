package soap

import (
	"bytes"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

// Do 发送请求, 格式化 req , 判断 status code 200 , 然后解析到 res
func Do[rqBody, rsHeader, rsBody any](ctx context.Context, url string, rqb rqBody, rsb *Envelope[rsHeader, rsBody]) error {
	// 格式化
	var body bytes.Buffer
	err := xml.NewEncoder(&body).Encode(rqb)
	if err != nil {
		return fmt.Errorf("encode xml %v", err)
	}
	// 请求
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &body)
	if err != nil {
		return fmt.Errorf("create request %v", err)
	}
	req.Header.Add("Content-Type", "application/soap+xml; charset=utf-8;")
	// 发送
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("do request %v", err)
	}
	defer res.Body.Close()
	// 先读取
	io.Copy(&body, res.Body)
	if body.Len() > 0 {
		// 解析
		err = xml.NewDecoder(&body).Decode(rsb)
		if err != nil {
			return fmt.Errorf("decode response %v", err)
		}
		if rsb.Body.Fault != nil {
			return rsb.Body.Fault
		}
	}
	// 状态码
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code %d", res.StatusCode)
	}
	//
	return nil
}
