package xs

import (
	"encoding/xml"
	"strconv"
	"time"
)

const (
	Namespace = "http://www.w3.org/2001/XMLSchema"
)

// NewNamespaceAttr 返回命名空间属性
func NewNamespaceAttr() *xml.Attr {
	return &xml.Attr{
		Name: xml.Name{
			Local: "xmlns:xs",
		},
		Value: Namespace,
	}
}

type Duration time.Duration

func (d Duration) String() string {
	td := time.Duration(d).Nanoseconds()

	seconds := td % int64(time.Minute)
	if seconds > 0 {
		td -= seconds
		seconds = seconds / int64(time.Second)
	}
	minutes := td % int64(time.Hour)
	if minutes > 0 {
		td -= minutes
		minutes = minutes / int64(time.Minute)
	}
	hours := td % int64(24*time.Hour)
	if hours > 0 {
		td -= hours
		hours = hours / int64(time.Hour)
	}

	days := td / int64(24*time.Hour)

	result := "P" // time duration designator
	//years
	// if years > 0 {
	// 	result += strconv.FormatInt(years, 10) + "Y"
	// }
	// if months > 0 {
	// 	result += strconv.FormatInt(months, 10) + "M"
	// }
	if days > 0 {
		result += strconv.FormatInt(days, 10) + "D"
	}

	if hours > 0 || minutes > 0 || seconds > 0 {
		result += "T"
		if hours > 0 {
			result += strconv.FormatInt(hours, 10) + "H"
		}
		if minutes > 0 {
			result += strconv.FormatInt(minutes, 10) + "M"
		}
		if seconds > 0 {
			result += strconv.FormatInt(seconds, 10) + "S"
		}
	}

	if len(result) == 1 {
		result += "T0S"
	}

	return result
}
