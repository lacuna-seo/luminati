// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"errors"
	"fmt"
	"net/url"
)

// Options contains the data used for obtaining serp
// results from the Luminati API.
type Options struct {
	// Keyword is the search term used for lookup.
	// NOTE: This is a required field.
	Keyword string
	// Country is the country to obtain data from for SERP's.
	// If nothing is passed, DefaultCountry will be used.
	Country string
	// Params is the url.Values to be sent to google. Default
	// parameters will be added if none are set such as
	// lum_mobile.
	Params url.Values
	// Desktop is the bool defining if desktop results should
	// be obtained as opposed to mobile.
	Desktop bool
}

var (
	// ErrNoKeywordProvided is returned by validate when no keyword
	// was provided to the Options struct.
	ErrNoKeywordProvided = errors.New("error: no keyword provided to options")
)

// Validate checks to see if the options passed are valid.
// And assigns default values if some arguments are
// missing.
func (o *Options) Validate() error {
	if o.Keyword == "" {
		return ErrNoKeywordProvided
	}

	if o.Country == "" {
		o.Country = DefaultCountry
	}

	if len(o.Params) == 0 {
		o.Params = url.Values{}
	}

	o.setDefaultParam("q", o.Keyword)
	o.setDefaultParam("gl", o.Country)
	o.setDefaultParam("num", "100")
	o.setDefaultParam("pws", "0")
	o.setDefaultParam("lum_json", "1")
	o.setDefaultParam("lum_mobile", "1")

	if o.Desktop {
		o.Params.Set("lum_mobile", "0")
	}

	return nil
}

// cacheKey obtains the key for storing response data in the cache.
// It will be unique per keyword and country.
func (o *Options) cacheKey(html, hasCache bool) string {
	if !hasCache {
		return ""
	}
	device := "mobile"
	if o.Desktop {
		device = "desktop"
	}
	format := "json"
	if html {
		format = "html"
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s", PrefixCacheKey, o.Keyword, o.Country, device, format)
}

// getRequestURL returns the URL for the request to Luminati.
func (o *Options) getRequestURL(baseURL string) string {
	return baseURL + "?" + o.Params.Encode()
}

// setDefaultParam checks to see if a query parameter exists by key
// and appends if it doesn't have it.
func (o *Options) setDefaultParam(key, value string) {
	_, ok := o.Params[key]
	if ok {
		return
	}
	o.Params.Set(key, value)
}
