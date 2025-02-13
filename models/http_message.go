package models

import "io"

type HTTPMessage struct {
	ObjectName    string
	RequestOrigin string
	Url           string
	Token         string
	Body          io.ReadCloser
	Method        string
	Params        map[string][]string
}

func (m HTTPMessage) String() string {
	return "\n{\n\tObject: " + m.ObjectName + "\n\tRequest: " + m.RequestOrigin + "\n\tUrl: " + m.Url + "\n\tToken : " + m.Token + "\n\tMethod: " + m.Method + "\n}"
}
