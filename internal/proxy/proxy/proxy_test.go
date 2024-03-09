package proxy

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/smakimka/tack/internal/model"
)

type want struct {
	code   int
	body   string
	header string
}

type test struct {
	name          string
	requestMethod string
	requestURL    string
	requestBody   string
	requestHeader string
	want          want
}

func listen(t *testing.T, tests []test) {
	mux := http.NewServeMux()

	count := 0
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		test := tests[count]
		count++

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		header := r.Header.Get("test")

		assert.Equal(t, test.requestBody, string(body))
		assert.Equal(t, test.requestHeader, header)

		w.Header().Add("test", test.want.header)
		w.Write([]byte(test.want.body))
	})

	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		test := tests[count]
		count++

		body, err := io.ReadAll(r.Body)
		require.NoError(t, err)
		header := r.Header.Get("test")

		assert.Equal(t, test.requestBody, string(body))
		assert.Equal(t, test.requestHeader, header)

		w.Header().Add("test", test.want.header)
		w.Write([]byte(test.want.body))
	})

	http.ListenAndServe("localhost:9998", mux)
}

func TestProxy(t *testing.T) {
	proxyAddr := "127.0.0.1:9999"
	tests := []test{
		{
			name:          "get request",
			requestMethod: "GET",
			requestURL:    "http://localhost:9998/",
			requestBody:   "",
			requestHeader: "test",
			want: want{
				code:   200,
				body:   "test",
				header: "test",
			},
		},
		{
			name:          "post request",
			requestMethod: "POST",
			requestURL:    "http://localhost:9998/",
			requestBody:   "testData",
			requestHeader: "testHeader",
			want: want{
				code:   200,
				body:   "testData",
				header: "testHeader",
			},
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	ctx = context.WithValue(ctx, model.EndpointKey, "test")
	defer cancel()
	proxyConfig := model.ProxyEndpoint{
		Endpoint: model.Endpoint{
			Addr:    proxyAddr,
			Workers: 1,
		},
	}
	errChan := make(chan error)

	go Serve(ctx, errChan, nil, &proxyConfig)
	go listen(t, tests)

	time.Sleep(5 * time.Millisecond)

	proxyURL, err := url.Parse(fmt.Sprintf("http://%s", proxyAddr))
	require.NoError(t, err)
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	for _, test := range tests {
		rctx, cancel := context.WithTimeout(ctx, 1*time.Second)
		defer cancel()
		r, err := http.NewRequestWithContext(rctx, test.requestMethod, test.requestURL, bytes.NewReader([]byte(test.requestBody)))
		require.NoError(t, err)

		r.Header.Add("test", test.requestHeader)

		resp, err := client.Do(r)
		require.NoError(t, err)

		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, test.want.code, resp.StatusCode)
		assert.Equal(t, test.want.body, string(body))
		assert.Equal(t, test.want.header, resp.Header.Get("test"))
	}
}
