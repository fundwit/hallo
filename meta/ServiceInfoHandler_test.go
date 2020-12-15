package meta

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pact-foundation/pact-go/utils"
	"io/ioutil"
	"net/http/httptest"
	"testing"
)

var port, _ = utils.GetFreePort()

func Test_info(t *testing.T) {
	engine := startInstrumentedProvider()

	t.Run("service info", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/", nil)
		w := httptest.NewRecorder()
		engine.ServeHTTP(w, req)

		result := w.Result()
		defer result.Body.Close()

		body, _ := ioutil.ReadAll(result.Body)
		serviceInfo := &ServiceInfo{}

		if json.Unmarshal(body, serviceInfo) != nil {
			t.Errorf("unexpected body %s", string(body))
		}

		expected := GetServiceInfo()
		if serviceInfo.ServiceName != expected.ServiceName || serviceInfo.ServiceInstance != expected.ServiceInstance {
			t.Errorf("expected %v, actual %v", expected, serviceInfo)
		}

		if parseBuildInfo([]byte("xxx")) != nil {
			t.Errorf("unexpected body %s", string(body))
		}
	})
}

func startInstrumentedProvider() *gin.Engine {
	engine := gin.Default()
	Routes(engine.Group("/"))

	go engine.Run(fmt.Sprintf(":%d", port))
	return engine
}
