package addszap

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"time"

	"go.uber.org/zap"
)

type LoggerRoundTripper struct {
	next http.RoundTripper
	log  *zap.Logger
}

func NewLoggerRoundTripper(next http.RoundTripper, log *zap.Logger) LoggerRoundTripper {
	return LoggerRoundTripper{next: next, log: log}
}

func (rt LoggerRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	var buf bytes.Buffer
	startTime := time.Now()
	reqBody := make(map[string]interface{})
	if r.Body != nil {
		tee := io.TeeReader(r.Body, &buf)
		_ = json.NewDecoder(tee).Decode(&reqBody)
		r.Body = io.NopCloser(&buf)
	}

	headers := header{}
	t := reflect.ValueOf(&headers).Elem()
	for i := 0; i < t.NumField(); i++ {
		tag := t.Type().Field(i).Tag.Get("json")
		name := strings.Split(tag, ",")[0]
		if t.Field(i).Kind() == reflect.String && tag != "-" {
			t.Field(i).SetString(strings.Join(r.Header.Values(name), ","))
		}
	}

	res, err := rt.next.RoundTrip(r)
	var status int
	if res != nil {
		status = res.StatusCode
	}

	msg := message{
		ReturnCode:     status,
		HttpMethod:     r.Method,
		RequestHeaders: headers,
		RemoteAddr:     r.RemoteAddr,
		Proto:          r.Proto,
		Path:           r.RequestURI,
		Latency:        time.Since(startTime),
		RequestParams:  params(r.URL.Query()),
		Body:           reqBody,
	}

	rt.log.Info("", zap.Object("message", msg))
	return res, err
}
