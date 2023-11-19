package soap

import (
	"encoding/xml"
	"fmt"
	"testing"
)

func Test_XML(t *testing.T) {
	var m Envelope[Header[any], Body[any]]
	m.Attr = append(m.Attr, NewNamespaceAttr())
	m.Header = new(Header[any])
	m.Body = new(Body[any])
	d, _ := xml.MarshalIndent(m, "", " ")
	fmt.Println(string(d))
}
