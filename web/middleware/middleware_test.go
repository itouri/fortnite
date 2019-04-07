package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/itouri/fortnite/web/middleware"
	"github.com/stretchr/testify/assert"

	"github.com/labstack/echo"
)

// コピペ

func TestCORS(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(echo.GET, "/", nil)
	//LEARN
	res := httptest.NewRecorder()
	c := e.NewContext(req, res)
	m := middleware.InitMiddleware()

	h := m.CORS(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))
	h(c)

	assert.Equal(t, "*", res.Header().Get("Access-Control-Allow-Origin"))
}
