// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"github.com/stretchr/testify/assert"
	"net/url"
	"testing"
)

func (t *LuminatiTestSuite) TestOptions_Validate() {
	tt := map[string]struct {
		input Options
		want  interface{}
	}{
		"No Keyword": {
			Options{},
			ErrNoKeywordProvided.Error(),
		},
		"No Country": {
			Options{
				Keyword: "reddico",
				Params:  url.Values{},
			},
			Options{
				Keyword: "reddico",
				Country: DefaultCountry,
				Params:  url.Values{"gl": []string{DefaultCountry}, "lum_json": []string{"1"}, "lum_mobile": []string{"1"}, "num": []string{"100"}, "pws": []string{"0"}, "q": []string{"reddico"}},
			},
		},
		"No Params": {
			Options{
				Keyword: "reddico",
			},
			Options{
				Keyword: "reddico",
				Country: DefaultCountry,
				Params:  url.Values{"gl": []string{DefaultCountry}, "lum_json": []string{"1"}, "lum_mobile": []string{"1"}, "num": []string{"100"}, "pws": []string{"0"}, "q": []string{"reddico"}},
			},
		},
		"Desktop": {
			Options{
				Keyword: "reddico",
				Desktop: true,
			},
			Options{
				Keyword: "reddico",
				Country: DefaultCountry,
				Params:  url.Values{"gl": []string{DefaultCountry}, "lum_json": []string{"1"}, "lum_mobile": []string{"0"}, "num": []string{"100"}, "pws": []string{"0"}, "q": []string{"reddico"}},
				Desktop: true,
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			err := test.input.Validate()
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, test.input)
		})
	}
}

func (t *LuminatiTestSuite) TestOptions_CacheKey() {
	tt := map[string]struct {
		html     bool
		hasCache bool
		input    Options
		want     string
	}{
		"Mobile": {
			false,
			true,
			Options{Keyword: "reddico", Country: "uk"},
			PrefixCacheKey + "-reddico-uk-mobile-json",
		},
		"Desktop": {
			false,
			true,
			Options{Keyword: "reddico", Country: "uk", Desktop: true},
			PrefixCacheKey + "-reddico-uk-desktop-json",
		},
		"HTML": {
			true,
			true,
			Options{Keyword: "reddico", Country: "uk"},
			PrefixCacheKey + "-reddico-uk-mobile-html",
		},
		"No Cache": {
			false,
			false,
			Options{Keyword: "reddico", Country: "uk"},
			"",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.input.cacheKey(test.html, test.hasCache)
			t.Equal(test.want, got)
		})
	}
}

func (t *LuminatiTestSuite) TestOptions_SetDefaultParam() {
	tt := map[string]struct {
		input Options
		key   string
		value string
		want  Options
	}{
		"Exists": {
			Options{Params: url.Values{"key": []string{"value"}}},
			"key",
			"value",
			Options{Params: url.Values{"key": []string{"value"}}},
		},
		"Set": {
			Options{Params: url.Values{}},
			"key",
			"value",
			Options{Params: url.Values{"key": []string{"value"}}},
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			test.input.setDefaultParam(test.key, test.value)
			t.Equal(test.want, test.input)
		})
	}
}

func TestAlphaNum(t *testing.T) {
	got := alphaNum(`"Disabled driving" site:uk intitle:"external sites" -wordpress.org -blogspot -pinterest -pdf -docx -doc -.info -.biz`)
	want := "Disabled-driving-site-uk-intitle-external-sites-wordpress-org-blogspot-pinterest-pdf-docx-doc-info-biz"
	assert.Equal(t, want, got)
}
