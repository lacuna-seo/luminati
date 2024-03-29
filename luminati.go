// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"context"
	"encoding/json"
	"github.com/ainsleyclark/redigo"
	"github.com/pkg/errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client is an HTTP Client for returning and obtaining data
// from the Luminati API.
type Client struct {
	client      *http.Client
	cache       redigo.Store
	bodyReader  func(io.Reader) ([]byte, error)
	BaseURL     string
	CacheExpiry time.Duration
	HasCache    bool
}

// KeywordFinder defines the methods used for finding Serp data
// through the Luminati API.
type KeywordFinder interface {
	// JSON Retrieves json from the search and returns a return struct
	// after processing.
	//
	// Returns an if the options failed validation, the request failed
	// or if there was a problem unmarshalling the response.
	JSON(ctx context.Context, o Options) (Serps, Meta, error)

	// HTML Retrieves raw HTML from the search and returns a string
	// of the result.
	//
	// Returns an if the options failed validation or the request
	// failed.
	HTML(ctx context.Context, o Options) (string, Meta, error)
}

const (
	// HTTPTimeout is the time limit for requests made by this
	// Client.
	HTTPTimeout = time.Minute * 5
	// IdleConnections controls the maximum number of idle (keep-alive)
	// connections across all hosts. Zero means no limit.
	IdleConnections = 50
	// DefaultCountry is the country used to obtain serps when
	// no country is passed via the Options.
	DefaultCountry = "uk"
	// DefaultCacheExpiry is the amount of time the response data
	// will live in the cache.
	DefaultCacheExpiry = 8 * time.Hour
	// PrefixCacheKey is the string prepended before the cache key.
	PrefixCacheKey = "luminati-client"
)

var (
	// ErrClientTimeout is returned the by the client when the
	// context deadline has exceeded.
	ErrClientTimeout = errors.New(context.DeadlineExceeded.Error() + " (Luminati.Client.Timeout exceeded while awaiting headers)")
)

// New creates a new Luminati client, an error will be returned if
// there was an issue parsing the proxy URL.
func New(uri string) (*Client, error) {
	if uri == "" {
		return nil, errors.New("proxy url cannot be nil, export LUMINATI_URL")
	}

	proxy, err := url.Parse(uri)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing luminati proxy url")
	}

	client := &Client{
		bodyReader: io.ReadAll,
		BaseURL:    "http://www.google.com/search",
		client: &http.Client{
			Timeout: HTTPTimeout,
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
				DialContext: (&net.Dialer{
					Timeout: HTTPTimeout,
				}).DialContext,
				MaxIdleConns:    IdleConnections,
				IdleConnTimeout: HTTPTimeout,
			},
		},
	}

	return client, nil
}

// NewWithCache creates a new Luminati client, with a cache store and
// default CacheExpiry, If the redigo.Store (Cache) interface
// passed is nil and error will be returned.
func NewWithCache(uri string, cache redigo.Store, cacheExpiry time.Duration) (*Client, error) {
	c, err := New(uri)
	if err != nil {
		return nil, err
	}
	if cache == nil {
		return nil, errors.New("cache interface is nil")
	}
	c.cache = cache
	c.CacheExpiry = cacheExpiry
	c.HasCache = true
	return c, nil
}

// JSON Retrieves json from the search and returns a return struct
// after processing.
//
// Returns an if the options failed validation, the request failed
// or if there was a problem unmarshalling the response.
func (c *Client) JSON(ctx context.Context, o Options) (Serps, Meta, error) {
	// For request/response times.
	now := time.Now()

	// Check the options are valid and assign defaults.
	err := o.Validate()
	if err != nil {
		return Serps{}, Meta{RequestTime: now}, err
	}

	// Setup the return meta.
	meta := Meta{
		CacheKey:    o.cacheKey(false, c.HasCache),
		RequestURL:  o.getRequestURL(c.BaseURL),
		RequestTime: now,
	}

	defer func() { meta = meta.process() }()

	// Try and retrieve in cache.
	if c.HasCache {
		var s Serps
		err = c.cache.Get(ctx, meta.CacheKey, &s)
		if err == nil {
			return s, meta, err
		}
	}

	// Obtain the response from either cache or the API.
	buf, err := c.fromLuminati(ctx, meta.RequestURL)
	if err != nil {
		return Serps{}, meta, err
	}
	meta.Body = string(buf)

	// Unmarshal into a response struct.
	res := response{}
	err = json.Unmarshal(buf, &res)
	if err != nil {
		return Serps{}, meta, errors.Wrap(err, "error unmarshalling luminati response")
	}

	// Get Serp data from the response.
	serps, err := res.ToSerps(buf)

	// Store in cache
	if c.HasCache && len(serps.Organic) > 0 {
		_ = c.cache.Set(context.Background(), meta.CacheKey, serps, redigo.Options{
			Expiration: c.CacheExpiry,
		})
	}

	return serps, meta, err
}

// HTML Retrieves raw HTML from the search and returns a string
// of the result.
//
// Returns an if the options failed validation or the request
// failed.
func (c *Client) HTML(ctx context.Context, o Options) (string, Meta, error) {
	// For request/response times.
	now := time.Now()

	// Check the options are valid and assign defaults.
	err := o.Validate()
	if err != nil {
		return "", Meta{RequestTime: now}, err
	}

	// Set html to true for options and set the
	// query for json to be false.
	o.Params.Set("lum_json", "0")

	// Setup the return meta.
	meta := Meta{
		CacheKey:    o.cacheKey(true, c.HasCache),
		RequestURL:  o.getRequestURL(c.BaseURL),
		RequestTime: now,
	}

	wasCached := false
	defer func() {
		meta = meta.process()
		meta.WasCached = wasCached
	}()

	// Try and retrieve in cache.
	if c.HasCache {
		var html string
		err = c.cache.Get(context.Background(), meta.CacheKey, &html)
		if err == nil {
			wasCached = true
			return html, meta, err
		}
	}

	// Obtain the response from either cache or the API.
	html, err := c.fromLuminati(ctx, meta.RequestURL)
	if err != nil {
		return "", meta.process(), err
	}
	meta.Body = string(html)

	// Store in cache
	if c.HasCache {
		_ = c.cache.Set(context.Background(), meta.CacheKey, html, redigo.Options{
			Expiration: c.CacheExpiry,
		})
	}

	return string(html), meta, nil
}

// fromLuminati obtains the response data from the luminati API
// if there is nothing stored in the cache.
//
// Returns errors.INTERNAL if the request could not be created, the
// request failed, the body could not be read or the cache
// could not be set.
func (c *Client) fromLuminati(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating request")
	}
	req = req.WithContext(ctx)

	resp, err := c.client.Do(req)
	if err != nil && err == context.Canceled {
		return nil, context.Canceled
	} else if err != nil && strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		return nil, ErrClientTimeout
	} else if err != nil {
		return nil, errors.Wrap(err, "luminati client request failed")
	}

	defer resp.Body.Close()

	buf, err := c.bodyReader(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "luminati body read failed")
	}

	return buf, nil
}
