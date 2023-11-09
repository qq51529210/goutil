package sip

import (
	"bytes"
	"testing"
)

func Test_Msg(t *testing.T) {
	var msg message
	msg.StartLine[0] = MethodMessage
	msg.StartLine[1] = "SIP:34120000002000000001@3412000000"
	msg.StartLine[2] = SIPVersion
	msg.Header.Via = append(msg.Header.Via, &Via{
		Proto:    "TCP",
		Address:  "192.168.31.193:8202",
		Branch:   BranchPrefix + "-ilVbYCX2fCM",
		RProt:    "2312",
		Received: "192.168.31.193:8201",
	})
	msg.Header.From.Name = "a"
	msg.Header.From.URI.Scheme = "sip"
	msg.Header.From.URI.Name = "34120000001310000001"
	msg.Header.From.URI.Domain = "3412000000"
	msg.Header.From.Tag = "ilVbZzetTs4"
	msg.Header.To.Name = "b"
	msg.Header.To.URI.Scheme = "sip"
	msg.Header.To.URI.Name = "34120000002000000001"
	msg.Header.To.URI.Domain = "3412000001"
	msg.Header.To.Tag = "915209491"
	msg.Header.CallID = "cilVbZBr0X16"
	msg.Header.CSeq.SN = "37596"
	msg.Header.CSeq.Method = MethodMessage
	msg.Header.Contact.Scheme = "sip"
	msg.Header.Contact.Name = "34120000001310000001"
	msg.Header.Contact.Domain = "192.168.31.193:8202"
	msg.Header.MaxForwards = "70"
	msg.Header.Expires = "3600"
	msg.Header.ContentType = "Application/MANSCDP+xml"
	msg.Header.UserAgent = "gbs"
	msg.Header.Set("a", "1")
	msg.Header.Set("b", "12")
	msg.Body.WriteString("123123123")
	//
	var buf bytes.Buffer
	msg.Enc(&buf)
	//
	r := newReader(&buf, MaxMessageLen)
	var _msg message
	err := _msg.Dec(r, MaxMessageLen)
	if err != nil {
		t.Fatal(err)
	}
}
