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
func Do[reqData, resData any](ctx context.Context, url string, rqd *reqData, rsd *resData) error {
	// 格式化
	var body bytes.Buffer
	// err := xml.NewEncoder(&body).Encode(rqd)
	// for test
	enc := xml.NewEncoder(&body)
	enc.Indent("", " ")
	err := enc.Encode(rqd)
	if err != nil {
		return fmt.Errorf("encode xml %v", err)
	}
	// for test
	fmt.Println(string(body.Bytes()))
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
	// 状态码
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error status code %d", res.StatusCode)
	}
	// 解析
	// err = xml.NewDecoder(res.Body).Decode(rsd)
	// for test
	io.Copy(&body, res.Body)
	fmt.Println(string(body.Bytes()))
	err = xml.NewDecoder(&body).Decode(rsd)
	if err != nil {
		return fmt.Errorf("decode response %v", err)
	}
	//
	return nil
}
