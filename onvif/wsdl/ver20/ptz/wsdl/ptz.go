package ptz

import (
	"encoding/xml"
	"goutil/soap"
)

const (
	Namespace = "http://www.onvif.org/ver20/ptz/wsdl"
)

var (
	envelopeAttr = []*xml.Attr{
		soap.NewNamespaceAttr(),
		NewNamespaceAttr(),
	}
)

// NewNamespaceAttr 返回命名空间属性
func NewNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:tptz",
		},
		Value: Namespace,
	}
}

const (
	PanTiltSpacePositionGenericSpace    = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/PositionGenericSpace"
	PanTiltSpaceTranslationGenericSpace = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/TranslationGenericSpace"
	PanTiltSpaceVelocityGenericSpace    = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/VelocityGenericSpace"
	PanTiltSpaceGenericSpeedSpace       = "http://www.onvif.org/ver10/tptz/PanTiltSpaces/GenericSpeedSpace"
)

const (
	ZoomSpacesPositionGenericSpace    = "http://www.onvif.org/ver10/tptz/ZoomSpaces/PositionGenericSpace"
	ZoomSpacesTranslationGenericSpace = "http://www.onvif.org/ver10/tptz/ZoomSpaces/TranslationGenericSpace"
	ZoomSpacesVelocityGenericSpace    = "http://www.onvif.org/ver10/tptz/ZoomSpaces/VelocityGenericSpace"
	ZoomSpacesZoomGenericSpeedSpace   = "http://www.onvif.org/ver10/tptz/ZoomSpaces/ZoomGenericSpeedSpace"
)
