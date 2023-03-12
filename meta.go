// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"time"
)

// Meta defines the information sent back from the client.
// It contains a cache key (if the client is using the cache).
// The request URL used to perform the request and response,
// request and latency times.
type Meta struct {
	// CacheKey defines the cache key that was created
	// when requesting. It's empty if no cache was used.
	CacheKey string
	// RequestURL is request URI that was sent to Luminati.
	RequestURL string
	// RequestTime is the time in which the request was
	// started.
	RequestTime time.Time
	// ResponseTime is the time in which all processing
	// was finished.
	ResponseTime time.Time
	// LatencyTime is the duration in which the client took
	// to perform the request.
	LatencyTime time.Duration
	// WasCached determines if the request was cached.
	WasCached bool
	// Body is the request body sent back from Luminati.
	Body string
}

// process adds the ResponseTime & LatencyTime to the
// Meta struct.
func (m *Meta) process() Meta {
	return Meta{
		CacheKey:     m.CacheKey,
		RequestURL:   m.RequestURL,
		RequestTime:  m.RequestTime.UTC(),
		ResponseTime: time.Now().UTC(),
		LatencyTime:  time.Since(m.RequestTime).Round(time.Microsecond),
	}
}
