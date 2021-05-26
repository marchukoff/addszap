package addszap

import (
	"fmt"
	"strings"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

type message struct {
	ReturnCode     int           `json:"return-code"`
	HttpMethod     string        `json:"http-method"`
	RequestHeaders header        `json:"request-headers"`
	RemoteAddr     string        `json:"remote-addr"`
	Proto          string        `json:"proto"`
	Path           string        `json:"path"`
	Latency        time.Duration `json:"latency"`
	RequestParams  params        `json:"request-params"`
	Body           body          `json:"body"`
}

func (m message) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	var err error
	enc.AddInt("return-code", m.ReturnCode)
	enc.AddString("http-method", m.HttpMethod)
	err = multierr.Append(err, enc.AddObject("request-headers", m.RequestHeaders))
	enc.AddString("remote-addr", m.RemoteAddr)
	enc.AddString("proto", m.Proto)
	enc.AddString("path", m.Path)
	enc.AddDuration("latency", m.Latency)
	err = multierr.Append(err, enc.AddObject("request-params", m.RequestParams))
	err = multierr.Append(err, enc.AddObject("body", m.Body))
	return err
}

type header struct {
	ContentType   string `json:"Content-Type"`
	ContentLength string `json:"Content-Length"`
	UserAgent     string `json:"User-Agent"`
	Server        string `json:"Server"`
	Via           string `json:"Via"`
	Accept        string `json:"Accept"`
	XForwardedFor string `json:"X-Forwarded-For"`
}

func (h header) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("content-type", h.ContentType)
	enc.AddString("content-length", h.ContentLength)
	enc.AddString("user-agent", h.UserAgent)
	enc.AddString("server", h.Server)
	enc.AddString("via", h.Via)
	enc.AddString("accept", h.Accept)
	enc.AddString("x-forwarded-for", h.XForwardedFor)
	return nil
}

type params map[string][]string

func (p params) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range p {
		enc.AddString(k, strings.Join(v, ", "))
	}
	return nil
}

type body map[string]interface{}

func (b body) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	for k, v := range b {
		switch t := v.(type) {
		case int:
			enc.AddInt(k, t)
		case string:
			enc.AddString(k, t)
		default:
			enc.AddString(k, fmt.Sprint(v))
		}
	}
	return nil
}
