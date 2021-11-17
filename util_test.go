// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

func (t *LuminatiTestSuite) TestCleanURI() {
	tt := map[string]struct {
		input string
		want  interface{}
	}{
		"Error": {
			"postgres://user:abc{",
			"error parsing luminati response url",
		},
		"With Params": {
			"https://google.com?query=test",
			"https://google.com",
		},
		"With She Bang": {
			"https://google.com?query=test#test",
			"https://google.com",
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got, err := cleanURL(test.input)
			if err != nil {
				t.Contains(err.Error(), test.want)
				return
			}
			t.Equal(test.want, got)
		})
	}
}

func (t *LuminatiTestSuite) TestStringInSlice() {
	tt := map[string]struct {
		input string
		list  []string
		want  bool
	}{
		"Truthy": {
			"test",
			[]string{"test"},
			true,
		},
		"Falsey": {
			"test",
			[]string{""},
			false,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := stringInSlice(test.input, test.list)
			t.Equal(test.want, got)
		})
	}
}
