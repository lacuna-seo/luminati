// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"encoding/json"
)

func (t *LuminatiTestSuite) TestResponse_ToSerps() {
	tt := map[string]struct {
		input    interface{}
		response response
		want     interface{}
	}{
		"Marshal Error": {
			make(chan byte),
			response{},
			"unmarshalling luminati response",
		},
		"Features": {
			map[string]interface{}{"images": 1},
			response{},
			Serps{
				Features: []string{"images"},
				//mappedFeatures: map[string]string{"images": "1"},
			},
		},
		"Excluded Features": {
			map[string]interface{}{
				"images":  1,
				"general": 1,
			},
			response{},
			Serps{
				Features: []string{"images"},
				//mappedFeatures: map[string]string{"images": "1"},
			},
		},
		"Organic": {
			map[string]interface{}{},
			response{Organic: []responseOrganic{
				{Rank: 1, Link: "https://reddico.co.uk", Description: "SEO"},
			}},
			Serps{
				Organic: []Organic{{Rank: 1, Link: "https://reddico.co.uk", Description: "SEO"}},
				//mappedFeatures: make(map[string]string),
			},
		},
		"Organic Bad URL": {
			map[string]interface{}{},
			response{Organic: []responseOrganic{
				{Rank: 1, Link: "postgres://user:abc{", Description: "SEO"},
			}},
			Serps{
				//mappedFeatures: make(map[string]string),
			},
		},
		"Organic With Features": {
			map[string]interface{}{
				"images": 1,
			},
			response{Organic: []responseOrganic{
				{Rank: 1, Link: "https://reddico.co.uk", Description: "SEO"},
			}},
			Serps{
				Organic:  []Organic{{Rank: 1, Link: "https://reddico.co.uk", Description: "SEO"}},
				Features: []string{"images"},
				//mappedFeatures: map[string]string{"images": "1"},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			buf, _ := json.Marshal(test.input) // Ignore on purpose
			got, err := test.response.ToSerps(buf)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, got)
		})
	}
}
