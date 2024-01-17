package media

import (
	"context"
	"encoding/xml"
	"goutil/onvif/wsdl/ver10/schema"
	"goutil/soap"
)

// StreamProtocol 流协议枚举
type StreamProtocol string

// StreamProtocol 枚举
const (
	StreamProtocolUDP  StreamProtocol = "UDP"
	StreamProtocolTCP  StreamProtocol = "TCP"
	StreamProtocolRTSP StreamProtocol = "RTSP"
	StreamProtocolHTTP StreamProtocol = "HTTP"
)

// StreamType 流类型枚举
type StreamType string

// StreamType 枚举
const (
	StreamTypeRTPUnicast   StreamType = "RTP-Unicast"
	StreamTypeRTPMulticast StreamType = "RTP-Multicast"
)

// GetStreamURL 查询媒体流
func GetStreamURL(ctx context.Context, url string, security *soap.Security, profileToken string,
	streamProtocol StreamProtocol, streamType StreamType) (*schema.MediaURL, error) {
	// 请求体
	var req soap.Envelope[*soap.Security, struct {
		XMLName     xml.Name `xml:"trt:GetStreamUri"`
		StreamSetup struct {
			XMLName   xml.Name   `xml:"trt:StreamSetup"`
			Stream    StreamType `xml:"tt:Stream"`
			Transport struct {
				Protocol StreamProtocol `xml:"tt:Protocol"`
			} `xml:"tt:Transport"`
		}
		ProfileToken string `xml:"trt:ProfileToken"`
	}]
	req.SetSoapTag()
	req.Attr = append(req.Attr, envelopeAttr...)
	req.Attr = append(req.Attr, schema.NamespaceAttr)
	req.Attr = append(req.Attr, soap.NamespaceAttr)
	req.Header.Data = security
	req.Body.Data.ProfileToken = profileToken
	req.Body.Data.StreamSetup.Stream = streamType
	req.Body.Data.StreamSetup.Transport.Protocol = streamProtocol
	// 响应体
	var res soap.Envelope[any, struct {
		XMLName  xml.Name         `xml:"GetStreamUriResponse"`
		MediaURL *schema.MediaURL `xml:"MediaUri"`
	}]
	// 发送
	err := soap.Do(ctx, url, &req, &res)
	if err != nil {
		return nil, err
	}
	// 成功
	return res.Body.Data.MediaURL, nil
}
