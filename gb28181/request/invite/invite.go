package invite

import (
	"fmt"
	"goutil/gb28181/request"
	"goutil/sdp"
	"goutil/sip"
	"net"
)

// Invite 类型
const (
	InvitePlay     string = "Play"
	InvitePlayback string = "Playback"
	InviteDownload string = "Download"
	InviteVideo    string = "video"
	InviteAudio    string = "audio"
)

// sdp m 字段
const (
	SDPMediaRTPMap96      = "rtpmap:96 PS/90000"
	SDPMediaRTPMap97      = "rtpmap:97 MPEG4/90000"
	SDPMediaRTPMap98      = "rtpmap:98 H264/90000"
	SDPMediaRTPMap99      = "rtpmap:99 H265/90000"
	SDPMediaConnectionNew = "connection:new"
	SDPMediaSetupActive   = "setup:active"
	SDPMediaSetupPassive  = "setup:passive"
	SDPMediaFMT           = "96 97 98 99"
	SDPMediaFMT96         = "96"
	SDPMediaFMT97         = "97"
	SDPMediaFMT98         = "98"
	SDPMediaFMT99         = "99"
)

type StreamMode string

// 流传输模式
const (
	StreamModePassive StreamMode = "passive"
	StreamModeActive  StreamMode = "active"
	StreamModeUDP     StreamMode = "udp"
)

type InviteData interface {
	GetSSRC() string
	GetStreamMode() StreamMode
	GetLocalIP() string
	GetLocalPort() string
}

// Invite 是 SendXX 的参数
type Invite struct {
	Ser       *sip.Server
	Device    request.Request
	ChannelID string
	BitStream string
	Invite    InviteData
	// sdp: u=
	SDPU  string
	AckID string
}

func (m *Invite) Message(sdp *sdp.Session, action, downloadSpeed string) (*sip.Message, net.Addr, error) {
	// 网络地址
	addr, err := m.Device.GetNetAddr()
	if err != nil {
		return nil, nil, err
	}
	// 消息
	msg := request.NewInvite(m.Device, addr.Network(), m.ChannelID)
	// subject
	msg.Header.Set("Subject", fmt.Sprintf("%s:%s,%s:0", m.ChannelID, m.Invite.GetSSRC(), m.Device.GetFromID()))
	// sdp
	m.initSDP(sdp, action, downloadSpeed)
	//
	return msg, addr, nil
}

func (m *Invite) initSDP(s *sdp.Session, action string, downloadSpeed string) {
	s.Init()
	s.S = action
	s.O.Username = m.ChannelID
	s.O.Address = m.Invite.GetLocalIP()
	s.C.Address = m.Invite.GetLocalIP()
	s.O.AddrType = sdp.AddrTypeIP4
	s.C.AddrType = sdp.AddrTypeIP4
	s.U = m.SDPU
	mm := new(sdp.Media)
	mm.Type = InviteVideo
	mm.Port = m.Invite.GetLocalPort()
	mm.FMT = SDPMediaFMT
	mm.A = append(mm.A, sdp.RecvOnly)
	mm.A = append(mm.A, SDPMediaRTPMap96)
	mm.A = append(mm.A, SDPMediaRTPMap97)
	mm.A = append(mm.A, SDPMediaRTPMap98)
	mm.A = append(mm.A, SDPMediaRTPMap99)
	// 码流
	if m.BitStream != "" {
		mm.A = append(mm.A, m.BitStream)
	}
	// 流传输模式
	mm.Proto = sdp.ProtoTCP
	streamMode := m.Invite.GetStreamMode()
	if streamMode == StreamModePassive {
		mm.A = append(mm.A, SDPMediaSetupPassive)
		mm.A = append(mm.A, SDPMediaConnectionNew)
	} else if streamMode == StreamModeActive {
		mm.A = append(mm.A, SDPMediaSetupActive)
		mm.A = append(mm.A, SDPMediaConnectionNew)
	} else {
		mm.Proto = sdp.ProtoUDP
	}
	// 下载速度
	if downloadSpeed != "" {
		mm.A = append(mm.A, "downloadspeed:"+downloadSpeed)
	}
	// ssrc
	mm.AddOther("y", m.Invite.GetSSRC())
	s.M = append(s.M, mm)

}
