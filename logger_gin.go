package addszap

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func GinLogger(log *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var buf bytes.Buffer
		startTime := time.Now()
		tee := io.TeeReader(c.Request.Body, &buf)
		reqBody := make(map[string]interface{})
		_ = json.NewDecoder(tee).Decode(&reqBody)
		c.Request.Body = io.NopCloser(&buf)

		defer func() {
			headers := header{}
			t := reflect.ValueOf(&headers).Elem()
			for i := 0; i < t.NumField(); i++ {
				tag := t.Type().Field(i).Tag.Get("json")
				name := strings.Split(tag, ",")[0]
				if t.Field(i).Kind() == reflect.String && tag != "-" {
					t.Field(i).SetString(strings.Join(c.Request.Header.Values(name), ","))
				}
			}

			msg := message{
				ReturnCode:     c.Writer.Status(),
				HttpMethod:     c.Request.Method,
				RequestHeaders: headers,
				RemoteAddr:     c.Request.RemoteAddr,
				Proto:          c.Request.Proto,
				Path:           c.Request.URL.Path,
				Latency:        time.Since(startTime),
				RequestParams:  params(c.Request.URL.Query()),
				Body:           reqBody,
			}

			log.Info("", zap.Object("message", msg))
		}()
		c.Next()
	}
}
