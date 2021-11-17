// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lacuna-seo/luminati/mocks"
	"github.com/lacuna-seo/stash"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

// LuminatiTestSuite defines the helper used for
// Luminati API testing.
type LuminatiTestSuite struct {
	suite.Suite
}

// TestLuminati asserts testing has begun.
func TestLuminati(t *testing.T) {
	suite.Run(t, new(LuminatiTestSuite))
}

// SetupClient creates sa new httptest.Server and cache
// store. It returns the teardown of server.Close.
func (t *LuminatiTestSuite) SetupClient(mock func(m *mocks.Cache), timeout bool) (*Client, func()) {
	// Start a local HTTP Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if timeout {
			time.Sleep(time.Second * 2)
		}
		_, err := w.Write([]byte("test"))
		t.NoError(err)
	}))

	cache := &mocks.Cache{}
	if mock != nil {
		mock(cache)
	}

	server.Client().Timeout = time.Second * 1

	// Set up the client with base URL & client.
	return &Client{
		client:      server.Client(),
		BaseURL:     server.URL,
		cache:       cache,
		bodyReader:  ioutil.ReadAll,
		CacheExpiry: DefaultCacheExpiry,
		HasCache:    true,
	}, server.Close
}

func (t *LuminatiTestSuite) TestNew() {
	tt := map[string]struct {
		proxyURL string
		cache    stash.Store
		want     interface{}
	}{
		"Empty URL": {
			"",
			nil,
			"proxy url cannot be nil, export LUMINATI_URL",
		},
		"Bad URL": {
			"postgres://user:abc{",
			nil,
			"error parsing luminati proxy url",
		},
		"Success": {
			"https://brightdata.com/proxy",
			&mocks.Cache{},
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			err := os.Setenv("LUMINATI_URL", test.proxyURL)
			t.NoError(err)

			got, err := New()
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}

			t.NotNil(got.bodyReader)
			t.Equal("http://www.google.com/search?", got.BaseURL)
			t.Equal(HTTPTimeout, got.client.Timeout)
		})
	}
}

func (t *LuminatiTestSuite) TestNewWithCache() {
	tt := map[string]struct {
		proxyURL string
		cache    stash.Store
		want     interface{}
	}{
		"Bad URL": {
			"postgres://user:abc{",
			nil,
			"error parsing luminati proxy url",
		},
		"Nil Cacher": {
			"https://brightdata.com/proxy",
			nil,
			"cache interface is nil",
		},
		"Success": {
			"https://brightdata.com/proxy",
			&mocks.Cache{},
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			err := os.Setenv("LUMINATI_URL", test.proxyURL)
			t.NoError(err)

			got, err := NewWithCache(test.cache, DefaultCacheExpiry)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}

			t.NotNil(got.cache)
			t.NotNil(got.bodyReader)
			t.True(got.HasCache)
			t.Equal("http://www.google.com/search?", got.BaseURL)
			t.Equal(HTTPTimeout, got.client.Timeout)
		})
	}
}

func (t *LuminatiTestSuite) TestClient_JSON() {
	key := PrefixCacheKey + "-reddico-uk-mobile-json"

	tt := map[string]struct {
		input Options
		mock  func(m *mocks.Cache)
		meta  Meta
		want  interface{}
	}{
		"Validate Error": {
			Options{},
			nil,
			Meta{},
			ErrNoKeywordProvided.Error(),
		},
		"Response Error": {
			Options{Keyword: "reddico"},
			func(m *mocks.Cache) {
				m.On("Get", context.Background(), mock.Anything, mock.Anything).
					Return(fmt.Errorf("error"))
				m.On("Set", context.Background(), key, []byte("test"), stash.Options{Expiration: DefaultCacheExpiry}).
					Return(fmt.Errorf("response error"))
			},
			Meta{CacheKey: "luminati-client-reddico-uk-mobile-json"},
			"response error",
		},
		"Unmarshal Error": {
			Options{Keyword: "reddico"},
			func(m *mocks.Cache) {
				var buf []byte
				m.On("Get", context.Background(), key, &buf).
					Return(nil).
					Run(func(args mock.Arguments) {
						arg := args.Get(2).(*[]byte)
						*arg = []byte("data")
					})
			},
			Meta{CacheKey: "luminati-client-reddico-uk-mobile-json"},
			"error unmarshalling luminati response",
		},
		"Success": {
			Options{Keyword: "reddico"},
			func(m *mocks.Cache) {
				var buf []byte
				m.On("Get", context.Background(), key, &buf).
					Return(nil).
					Run(func(args mock.Arguments) {
						res := map[string]interface{}{"images": 1}
						buf, err := json.Marshal(res)
						t.NoError(err)
						arg := args.Get(2).(*[]byte)
						*arg = buf
					})
			},
			Meta{CacheKey: "luminati-client-reddico-uk-mobile-json"},
			Serps{
				Features: []string{"images"},
				mappedFeatures: map[string]string{
					"images": "1",
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			c, teardown := t.SetupClient(test.mock, false)
			defer teardown()

			got, meta, err := c.JSON(test.input)
			t.Equal(test.meta.CacheKey, meta.CacheKey)
			t.Contains(meta.RequestURL, test.meta.RequestURL)

			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, got)
		})
	}
}

func (t *LuminatiTestSuite) TestClient_HTML() {
	key := PrefixCacheKey + "-reddico-uk-mobile-html"

	tt := map[string]struct {
		input Options
		mock  func(m *mocks.Cache)
		meta  Meta
		want  interface{}
	}{
		"Validate Error": {
			Options{},
			nil,
			Meta{},
			ErrNoKeywordProvided.Error(),
		},
		"Response Error": {
			Options{Keyword: "reddico"},
			func(m *mocks.Cache) {
				m.On("Get", context.Background(), mock.Anything, mock.Anything).
					Return(fmt.Errorf("error"))
				m.On("Set", context.Background(), key, []byte("test"), stash.Options{Expiration: DefaultCacheExpiry}).
					Return(fmt.Errorf("response error"))
			},
			Meta{CacheKey: key, RequestURL: "gl=uk&lum_json=0&lum_mobile=1&num=100&pws=0&q=reddico"},
			"response error",
		},
		"Success": {
			Options{Keyword: "reddico", Country: "uk"},
			func(m *mocks.Cache) {
				var buf []byte
				m.On("Get", context.Background(), key, &buf).
					Return(nil).
					Run(func(args mock.Arguments) {
						arg := args.Get(2).(*[]byte)
						*arg = []byte("data")
					})
			},
			Meta{CacheKey: key, RequestURL: "gl=uk&lum_json=0&lum_mobile=1&num=100&pws=0&q=reddico"},
			"data",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			c, teardown := t.SetupClient(test.mock, false)
			defer teardown()

			got, meta, err := c.HTML(test.input)
			t.Equal(test.meta.CacheKey, meta.CacheKey)
			t.Contains(meta.RequestURL, test.meta.RequestURL)

			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}

			t.Equal(meta.CacheKey, key)
			t.Equal(test.want, got)
		})
	}
}

func (t *LuminatiTestSuite) TestClient_GetResponse() {
	var (
		cacheFail = func(m *mocks.Cache) {
			m.On("Get", context.Background(), mock.Anything, mock.Anything).
				Return(fmt.Errorf("error"))
		}
		key = PrefixCacheKey + "-reddico-uk-mobile-json"
	)

	tt := map[string]struct {
		url        string
		mock       func(m *mocks.Cache)
		bodyReader func(io.Reader) ([]byte, error)
		timeout    bool
		withCache  bool
		want       interface{}
	}{
		"From Cache": {
			"",
			func(m *mocks.Cache) {
				var buf []byte
				m.On("Get", context.Background(), key, &buf).
					Return(nil).
					Run(func(args mock.Arguments) {
						arg := args.Get(2).(*[]byte)
						*arg = []byte("data")
					})
			},
			ioutil.ReadAll,
			false,
			true,
			"data",
		},
		"Bad Request": {
			"@#@#$$%$",
			cacheFail,
			ioutil.ReadAll,
			false,
			true,
			"error creating request",
		},
		"Do Error": {
			"doerror",
			cacheFail,
			ioutil.ReadAll,
			false,
			true,
			"luminati client request failed",
		},
		"Read Error": {
			"https://google.com",
			cacheFail,
			func(reader io.Reader) ([]byte, error) {
				return nil, fmt.Errorf("error")
			},
			false,
			true,
			"luminati body read failed",
		},
		"Cache Error": {
			"",
			func(m *mocks.Cache) {
				m.On("Get", context.Background(), mock.Anything, mock.Anything).
					Return(fmt.Errorf("error"))
				m.On("Set", context.Background(), key, []byte("test"), stash.Options{Expiration: DefaultCacheExpiry}).
					Return(fmt.Errorf("cache error"))
			},
			ioutil.ReadAll,
			false,
			true,
			"cache error",
		},
		"Timeout": {
			"",
			cacheFail,
			ioutil.ReadAll,
			true,
			true,
			ErrClientTimeout.Error(),
		},
		"Prevent Cache": {
			"",
			func(m *mocks.Cache) {
				m.On("Get", context.Background(), mock.Anything, mock.Anything).
					Return(fmt.Errorf("error"))
			},
			ioutil.ReadAll,
			false,
			false,
			"test",
		},
		"Success": {
			"",
			func(m *mocks.Cache) {
				m.On("Get", context.Background(), mock.Anything, mock.Anything).
					Return(fmt.Errorf("error"))
				m.On("Set", context.Background(), key, []byte("test"), stash.Options{Expiration: DefaultCacheExpiry}).
					Return(nil)
			},
			ioutil.ReadAll,
			false,
			true,
			"test",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			c, teardown := t.SetupClient(test.mock, test.timeout)
			defer teardown()

			if test.url == "" {
				test.url = c.BaseURL
			}

			c.bodyReader = test.bodyReader

			if !test.withCache {
				c.cache = nil
				c.HasCache = false
			}

			got, err := c.getResponse(key, test.url)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}

			t.Equal(test.want, string(got))
		})
	}
}
