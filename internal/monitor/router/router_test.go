package router

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/tack/internal/model"
	"github.com/smakimka/tack/internal/monitor/storage"
)

func TestNewRouter(t *testing.T) {
	type want struct {
		code        int
		contentType string
		body        string
	}
	tests := []struct {
		method string
		path   string
		body   *model.SendData
		want   want
	}{
		{"POST", "/api/update/", &model.SendData{Name: "test", Endpoints: map[string]model.EndpointData{"test": {}}}, want{200, "", ""}},
		{"GET", "/", nil, want{200, "text/html; charset=utf-8", "\n<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n  <meta charset=\"utf-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>Tacks</title>\n</head>\n<body>\n  <h1><a href=\"tacks/test\">test</a></h1>\n</body>\n</html>"}},
		{"GET", "/tacks/test/proxies/test/", nil, want{200, "text/html; charset=utf-8", "\n<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n  <meta charset=\"utf-8\">\n  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n  <title>test</title>\n</head>\n<body>\n  <h1>test<h1>\n<h2>TotalRecieved: 0 bytes<h2>\n<h2>TotalSent: 0 bytes<h2>\n<h1>Clients:<h1>\n</body>\n</html>"}},
	}

	s := storage.NewMemStorage()
	router := New(s)

	for _, test := range tests {
		var body []byte
		if test.body != nil {
			byteBody, err := json.Marshal(test.body)
			require.NoError(t, err)
			body = byteBody
		}
		req := httptest.NewRequest(test.method, test.path, bytes.NewReader(body))
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		assert.Equal(t, test.want.code, w.Code)
		assert.Equal(t, test.want.contentType, w.Header().Get("Content-type"))
		assert.Equal(t, test.want.body, w.Body.String())
	}
}
