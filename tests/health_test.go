package tests

import (
	"gin/app/errcode"
	"gin/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/stretchr/testify/require"
)

func TestHealth(t *testing.T) {
	r := gin.Default()
	router.LoadRouters(r)

	req := httptest.NewRequest("GET", "/ping", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	require.Equal(t, http.StatusOK, w.Code)

	var resp errcode.SuccessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, "pong", resp.Msg)
}
