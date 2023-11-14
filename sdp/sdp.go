package sdp

// // 一些常量
// const (
// 	NetTypeIN   = "IN"
// 	AddrTypeIP4 = "IP4"
// 	AddrTypeIP6 = "IP6"
// )

// const (
// 	zero = "0"
// )

// var (
// 	// ErrSDPFormat 表示 sdp 格式错误
// 	ErrSDPFormat = errors.New("error sdp format")
// 	// ErrSDPRTPMapFormat 表示 rtpmap 格式错误
// 	ErrSDPRTPMapFormat = errors.New("error sdp rtpmap format")
// )

// // SDP 表示 sdp 的字段
// type SDP struct {
// 	// 版本号
// 	Version string
// 	// 会话创建者信息
// 	Origin *Origin
// 	// 会话名称
// 	SessionName string
// 	// URI
// 	URI string
// 	// 连接信息
// 	Connection *Connection
// 	// 时间描述
// 	Time *Time
// 	// 媒体流描述
// 	Media *Media
// 	// rtpmap
// 	RTPMap map[string]string
// 	// sendonly / recvonly / sendrecv
// 	SendRecv string
// 	// 属性
// 	Attributes map[string]string
// 	// 国标 y=
// 	SSRC string
// 	// 国标 f=
// 	Codec string
// }

// // Init 初始化
// func (sdp *SDP) Init() {
// 	sdp.Origin = new(Origin)
// 	sdp.Origin.SessionID = zero
// 	sdp.Origin.SessionVersion = zero
// 	sdp.Origin.NetType = NetTypeIN
// 	sdp.Origin.AddrType = AddrTypeIP4
// 	sdp.Connection = new(Connection)
// 	sdp.Connection.NetType = NetTypeIN
// 	sdp.Connection.AddrType = AddrTypeIP4
// 	sdp.Time = new(Time)
// 	sdp.Media = new(Media)
// 	sdp.Media.Proto = ProtoUDP
// 	sdp.RTPMap = make(map[string]string)
// 	sdp.Attributes = make(map[string]string)
// }

// // ParseFrom 从 reader 中解析 sdp
// func (sdp *SDP) ParseFrom(reader io.Reader) error {
// 	sdp.RTPMap = make(map[string]string)
// 	sdp.Attributes = make(map[string]string)
// 	//
// 	scaner := bufio.NewScanner(reader)
// 	var value string
// 	for scaner.Scan() {
// 		line := scaner.Text()
// 		if line == "" {
// 			continue
// 		}
// 		// v=
// 		value = strings.TrimPrefix(line, "v=")
// 		if value != line {
// 			sdp.Version = value
// 			continue
// 		}
// 		// o=
// 		value = strings.TrimPrefix(line, "o=")
// 		if value != line {
// 			sdp.Origin = new(Origin)
// 			if err := sdp.Origin.Parse(value); err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		// s=
// 		value = strings.TrimPrefix(line, "s=")
// 		if value != line {
// 			sdp.SessionName = value
// 			continue
// 		}
// 		// u=
// 		value = strings.TrimPrefix(line, "u=")
// 		if value != line {
// 			sdp.URI = value
// 			continue
// 		}
// 		// c=
// 		value = strings.TrimPrefix(line, "c=")
// 		if value != line {
// 			sdp.Connection = new(Connection)
// 			if err := sdp.Connection.Parse(value); err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		// t=
// 		value = strings.TrimPrefix(line, "t=")
// 		if value != line {
// 			sdp.Time = new(Time)
// 			if err := sdp.Time.Parse(value); err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		// m=
// 		value = strings.TrimPrefix(line, "m=")
// 		if value != line {
// 			sdp.Media = new(Media)
// 			if err := sdp.Media.Parse(value); err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		// a=
// 		value = strings.TrimPrefix(line, "a=")
// 		if value != line {
// 			if err := sdp.parseA(value); err != nil {
// 				return err
// 			}
// 			continue
// 		}
// 		// y=
// 		value = strings.TrimPrefix(line, "y=")
// 		if value != line {
// 			sdp.SSRC = value
// 			continue
// 		}
// 		// f=
// 		value = strings.TrimPrefix(line, "f=")
// 		if value != line {
// 			sdp.Codec = value
// 			continue
// 		}
// 	}
// 	//
// 	if err := scaner.Err(); err != nil {
// 		return err
// 	}
// 	if sdp.Origin == nil || sdp.Connection == nil || sdp.Time == nil || sdp.Media == nil {
// 		return ErrSDPFormat
// 	}
// 	//
// 	return nil
// }

// // parseA 解析 a 字段
// func (sdp *SDP) parseA(line string) error {
// 	// rtpmap:
// 	value := strings.TrimPrefix(line, "rtpmap:")
// 	if value != line {
// 		p := strings.Fields(value)
// 		if len(p) < 2 {
// 			return ErrSDPRTPMapFormat
// 		}
// 		sdp.RTPMap[p[0]] = p[1]
// 		return nil
// 	}
// 	//
// 	if value == "sendonly" || value == "recvonly" || value == "sendrecv" {
// 		sdp.SendRecv = value
// 		return nil
// 	}
// 	//
// 	p := strings.Split(value, ":")
// 	if len(p) > 1 {
// 		sdp.Attributes[p[0]] = p[1]
// 	} else {
// 		sdp.Attributes[p[0]] = ""
// 	}
// 	//
// 	return nil
// }

// // FormatTo 格式化到 buf
// func (sdp *SDP) FormatTo(buf *bytes.Buffer) {
// 	if sdp.Version != "" {
// 		fmt.Fprintf(buf, "v=%s\r\n", sdp.Version)
// 	} else {
// 		buf.WriteString("v=0\r\n")
// 	}
// 	fmt.Fprintf(buf, "o=%s\r\n", sdp.Origin.String())
// 	fmt.Fprintf(buf, "s=%s\r\n", sdp.SessionName)
// 	fmt.Fprintf(buf, "c=%s\r\n", sdp.Connection.String())
// 	fmt.Fprintf(buf, "t=%s\r\n", sdp.Time.String())
// 	for k, v := range sdp.Attributes {
// 		if k == "" {
// 			continue
// 		}
// 		if v == "" {
// 			fmt.Fprintf(buf, "a=%s\r\n", k)
// 		} else {
// 			fmt.Fprintf(buf, "a=%s:%s\r\n", k, v)
// 		}
// 	}
// 	fmt.Fprintf(buf, "m=%s\r\n", sdp.Media.String())
// 	if sdp.SendRecv != "" {
// 		fmt.Fprintf(buf, "a=%s\r\n", sdp.SendRecv)
// 	}
// 	for k, v := range sdp.RTPMap {
// 		fmt.Fprintf(buf, "a=rtpmap:%s %s\r\n", k, v)
// 	}
// 	if sdp.Codec != "" {
// 		fmt.Fprintf(buf, "f=%s\r\n", sdp.Codec)
// 	}
// 	if sdp.SSRC != "" {
// 		fmt.Fprintf(buf, "y=%s\r\n", sdp.SSRC)
// 	}
// }

// // parseFrom 从 reader 中解析
// func (sdp *SDP) parseFrom(reader io.Reader) error {
// 	scaner := bufio.NewScanner(reader)
// 	for scaner.Scan() {
// 		line := scaner.Text()
// 		if line == "" {
// 			continue
// 		}
// 	}
// 	return scaner.Err()
// }
