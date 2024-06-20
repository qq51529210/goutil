package xml

import "encoding/xml"

// Result 用于解析响应中有 Result 的
type Result struct {
	// 基本
	XMLName  xml.Name
	CmdType  string
	SN       string
	DeviceID string
	// 执行结果标志
	Result string `xml:",omitempty"`
}
