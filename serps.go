// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"strings"
)

type (
	// Serps defines the collection to be returned from
	// the client.
	Serps struct {
		Organic        []Organic `json:"serps"`
		Features       []string  `json:"features"`
		mappedFeatures map[string]string
	}
	// Domain are URL specific results returned by
	// Serps.CheckURL
	Domain struct {
		Query   Query     `json:"query"`
		Results []Organic `json:"results"`
	}
	// Query defines the first top level
	Query struct {
		Rank        int    `json:"position"`
		Link        string `json:"url"`
		Description string `json:"text"`
		Features    string `json:"features"`
	}
	// Organic represents a singular organic SERP
	// as defined in Serps.
	Organic struct {
		Rank        int    `json:"position"`
		Description string `json:"text"`
		Link        string `json:"url"`
	}
)

// CheckURL obtains the highest ranking Serp for a given
// URL. Features are also obtained.
func (s *Serps) CheckURL(url string) Domain {
	d := Domain{}

	firstFound := true
	for _, serp := range s.Organic {
		if !strings.Contains(serp.Link, url) {
			continue
		}
		d.Results = append(d.Results, serp)
		if !firstFound {
			continue
		}
		d.Query = Query{
			Rank:        serp.Rank,
			Link:        serp.Link,
			Description: serp.Description,
			Features:    s.getFeatures(url),
		}
		firstFound = false
	}

	return d
}

// getFeatures obtains a comma delimited list of features that
// the domain ranks for.
func (s *Serps) getFeatures(url string) string {
	features := ""
	for key, value := range s.mappedFeatures {
		if strings.Contains(value, url) {
			features += key + ","
		}
	}
	return strings.TrimSuffix(features, ",")
}
