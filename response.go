// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import (
	"encoding/json"
	"github.com/pkg/errors"
)

type (
	// response defines the data received back from the
	// Luminati API.
	response struct {
		Organic []responseOrganic `json:"organic"`
	}
	// responseOrganic is the collection of organic items.
	responseOrganic struct {
		Rank        int    `json:"rank"`
		Link        string `json:"link"`
		DisplayLink string `json:"display_link"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Extensions  []struct {
			Text   string `json:"text"`
			Inline bool   `json:"inline"`
		} `json:"extensions,omitempty"`
	}
)

// ToSerps transforms a buffer with options to a collection
// of serps. Top level features will be found and a query will
// be built up dependent on the URL passed in options.
func (r *response) ToSerps(buf []byte) (Serps, error) {
	var excluded = []string{"general", "organic", "pagination", "related"}

	m := map[string]json.RawMessage{}
	err := json.Unmarshal(buf, &m)
	if err != nil {
		return Serps{}, errors.Wrap(err, "error unmarshalling luminati response")
	}

	serps := r.GetSerps()

	// Find features before continuing on to get organic
	// results.
	for key, _ := range m { //nolint
		if stringInSlice(key, excluded) {
			continue
		}
		//serps.mappedFeatures[key] = string(value)
		serps.Features = append(serps.Features, key)
	}

	return serps, nil
}

// GetSerps returns Organic results and Features that
// appear for the keyword. URLs are cleaned and
// Organic results are appended.
func (r *response) GetSerps() Serps {
	s := Serps{
		//mappedFeatures: make(map[string]string),
	}
	for _, v := range r.Organic {
		link, err := cleanURL(v.Link)
		if err != nil {
			continue
		}
		serp := Organic{
			Rank:        v.Rank,
			Description: v.Description,
			Link:        link,
		}
		s.Organic = append(s.Organic, serp)
	}
	return s
}
