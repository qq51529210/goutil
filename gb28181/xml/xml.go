package xml

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"goutil/sip"
	"io"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

// 根标签的值
const (
	TypeQuery    = "Query"
	TypeNotify   = "Notify"
	TypeControl  = "Control"
	TypeResponse = "Response"
)

// CmdType 标签的值
const (
	CmdAlarm          = "Alarm"
	CmdCatalog        = "Catalog"
	CmdKeepalive      = "Keepalive"
	CmdDeviceInfo     = "DeviceInfo"
	CmdRecordInfo     = "RecordInfo"
	CmdDeviceStatus   = "DeviceStatus"
	CmdDeviceControl  = "DeviceControl"
	CmdDeviceConfig   = "DeviceConfig"
	CmdBroadcast      = "Broadcast"
	CmdMobilePosition = "MobilePosition"
	CmdMediaStatus    = "MediaStatus"
	CmdConfigDownload = "ConfigDownload"
	CmdPresetQuery    = "PresetQuery"
)

// 编码
var (
	EncodingGBK    = "GBK"
	EncodingGB2312 = "GB2312"
	EncodingUTF8   = "UTF-8"
)

// 公共错误
var (
	ErrFormat  = sip.NewResponseError(sip.StatusBadRequest, "Error XML Format", "")
	ErrType    = sip.NewResponseError(sip.StatusBadRequest, "Error XML Type", "")
	ErrCmdType = sip.NewResponseError(sip.StatusBadRequest, "Error XML CmdType", "")
	ErrData    = sip.NewResponseError(sip.StatusBadRequest, "Error XML Data", "")
)

// Decode 从 r 解析数据到 v
// 自动分辨字符编码
func Decode(r io.Reader, v any) error {
	d := xml.NewDecoder(r)
	d.CharsetReader = charsetReader
	return d.Decode(v)
}

// charsetReader 字符编码
func charsetReader(charset string, input io.Reader) (io.Reader, error) {
	charset = strings.ToUpper(charset)
	switch charset {
	case EncodingGBK:
		return simplifiedchinese.GBK.NewDecoder().Reader(input), nil
	case EncodingGB2312:
		return transform.NewReader(input, simplifiedchinese.GB18030.NewDecoder()), nil
	}
	return input, nil
}

// Encode 使用 encoding 编码格式化 v 到 buf
func Encode(buf *bytes.Buffer, encoding string, v any) {
	fmt.Fprintf(buf, "<?xml version=\"1.0\" encoding=\"%s\"?>\n", encoding)
	var enc *xml.Encoder
	switch encoding {
	case EncodingGB2312:
		enc = xml.NewEncoder(transform.NewWriter(buf, simplifiedchinese.GB18030.NewEncoder()))
	case EncodingGBK:
		enc = xml.NewEncoder(transform.NewWriter(buf, simplifiedchinese.GBK.NewEncoder()))
	default:
		enc = xml.NewEncoder(buf)
	}
	enc.Indent("", " ")
	enc.Encode(v)
	buf.WriteByte('\n')
}
