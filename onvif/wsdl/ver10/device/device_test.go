package device

import (
	"encoding/xml"
	"os"
	"testing"
)

func print(v any) {
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	enc.Encode(v)
	os.Stdout.WriteString("\n")
}

func Test_GetSystemDateAndTime(t *testing.T) {
	var _req getSystemDateAndTimeReq
	_req.Envelope.Attr = envelopeAttr
	print(&_req)
}

func Test_GetDeviceInformation(t *testing.T) {
	var _req getDeviceInformationReq
	_req.Envelope.Attr = envelopeAttr
	print(&_req)
}

func Test_GetCapabilities(t *testing.T) {
	var _req getCapabilitiesReq
	_req.Envelope.Attr = envelopeAttr
	print(&_req)
}
