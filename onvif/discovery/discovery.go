package discovery

import (
	"encoding/xml"
	"time"
)

const (
	// 用于读取消息的大小，10k 足够了
	readBufLen = 1024 * 10
	// MulticastAddr 默认的使用的多播地址
	MulticastAddr = "239.255.255.250:3702"
)

const (
	// 发送的消息，就差 uuid 了
	msgFmt = `<?xml version="1.0" encoding="UTF-8"?>
<Envelope xmlns="http://www.w3.org/2003/05/soap-envelope" xmlns:tds="http://www.onvif.org/ver10/device/wsdl">
  <Header xmlns:a="http://schemas.xmlsoap.org/ws/2004/08/addressing">
    <a:MessageID>uuid:%s</a:MessageID>
    <a:To>urn:schemas-xmlsoap-org:ws:2005:04:discovery</a:To>
    <a:Action>http://schemas.xmlsoap.org/ws/2005/04/discovery/Probe</a:Action>
  </Header>
  <Body>
    <Probe xmlns="http://schemas.xmlsoap.org/ws/2005/04/discovery" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance">
		<Types>tds:Device</Types>
    </Probe>
  </Body>
</Envelope>`
)

// result 用于解析地址
type result struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		ProbeMatches struct {
			ProbeMatch struct {
				XAddrs string
			}
		}
	}
}

// Discover 启动协程在后台运行，通过 callback 回调函数通知结果或者错误
// iface 是网卡接口名称
// addr 是监听的多播地址
// callback 是回调，返回 true 表示结束运行哦
// 返回初始化发生的错误
func Discover(ifaceName, mutilAddr string, timeout time.Duration, callback Callback) error {
	// d := new(discoverer)
	// d.callback = callback
	// // 初始化
	// if err := d.init(ifaceName, mutilAddr, timeout); err != nil {
	// 	return err
	// }
	// // 启动读写
	// go d.readRoutine()
	// go d.writeRoutine()
	//
	return nil
}
