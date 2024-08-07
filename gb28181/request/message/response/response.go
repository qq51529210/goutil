package response

import (
	"context"
	"goutil/gb28181/request"
	"goutil/gb28181/xml"
	"goutil/sip"
)

// SendResult 封装代码，向级联发送应答请求
func SendResult(ctx context.Context, ser *sip.Server, cascade request.Request, data *xml.Message, result string) error {
	// body
	var body xml.Message
	body.XMLName.Local = xml.TypeResponse
	body.CmdType = data.CmdType
	body.SN = data.SN
	body.DeviceID = data.DeviceID
	body.Result = result
	// 发送
	return request.SendMessage(ctx, ser, cascade, &body, nil)
}
