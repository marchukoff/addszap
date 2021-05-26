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

func HttpLogger(next http.Handler, log *zap.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var buf bytes.Buffer
			startTime := time.Now()
			tee := io.TeeReader(r.Body, &buf)
			reqBody := make(map[string]interface{})
			_ = json.NewDecoder(tee).Decode(&reqBody)
			r.Body = io.NopCloser(&buf)

			headers := header{}
			t := reflect.ValueOf(&headers).Elem()
			for i := 0; i < t.NumField(); i++ {
				tag := t.Type().Field(i).Tag.Get("json")
				name := strings.Split(tag, ",")[0]
				if t.Field(i).Kind() == reflect.String && tag != "-" {
					t.Field(i).SetString(strings.Join(r.Header.Values(name), ","))
				}
			}

			wrapped := wrapResponseWriter(w)
			next.ServeHTTP(wrapped, r)
			msg := message{
				ReturnCode:     wrapped.Status(),
				HttpMethod:     r.Method,
				RequestHeaders: headers,
				RemoteAddr:     r.RemoteAddr,
				Proto:          r.Proto,
				Path:           r.RequestURI,
				Latency:        time.Since(startTime),
				RequestParams:  params(r.URL.Query()),
				Body:           reqBody,
			}

			log.Info("", zap.Object("message", msg))
		},
	)
}
