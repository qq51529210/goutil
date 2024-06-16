package sip

import (
	"fmt"
	"sync/atomic"
)

// 方法
const (
	MethodRegister  string = "REGISTER"
	MethodInvite    string = "INVITE"
	MethodACK       string = "ACK"
	MethodBye       string = "BYE"
	MethodMessage   string = "MESSAGE"
	MethodNotify    string = "NOTIFY"
	MethodSubscribe string = "SUBSCRIBE"
	MethodInfo      string = "INFO"
)

// 一些常量
const (
	BranchPrefix  = "z9hG4bK"
	SIPVersion    = "SIP/2.0"
	MaxMessageLen = 1024 * 10
	TCP           = "SIP/2.0/TCP"
	TCPS          = "SIPS/2.0/TCP"
	UDP           = "SIP/2.0/UDP"
	UDPS          = "SIPS/2.0/UDP"
	SIP           = "sip"
	SIPS          = "sips"
)

var (
	// CSeq 的递增 SN
	cseqSN = int64(0)
)

// GetSN 返回全局递增的 sn
func GetSN() int64 {
	return atomic.AddInt64(&cseqSN, 1)
}

// GetSNString 返回字符串形式的全局递增的 sn
func GetSNString() string {
	return fmt.Sprintf("%d", atomic.AddInt64(&cseqSN, 1))
}
