package addszap

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGinLogger(t *testing.T) {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(GinLogger(NewLogger(true)))

	server := httptest.NewServer(r)
	defer server.Close()

	client := &http.Client{}

	t.Run("simple test", func(t *testing.T) {
		var req *http.Request
		{
			var err error
			req, err = http.NewRequest(
				http.MethodGet,
				server.URL+"/?x=1&y=3",
				strings.NewReader(`{"count": 42, "color": {"green": 255, "red": 0, "blue":0}}`),
			)
			if err != nil {
				t.Fatal(err)
			}
		}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()
	})
}
