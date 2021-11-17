// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

var (
	TestURL         = "https://www.apple.com"
	OrganicTestData = []Organic{
		{Rank: 1, Description: "MacBook Pro. Our most powerful notebooks. Fast M1 processors, incredible graphics, and spectacular Retina displays. Now available in a 14-inch model.", Link: "https://www.apple.com/macbook-pro/"},
		{Rank: 2, Description: "Explore the world of Mac. Check out the all-new MacBook Pro, MacBook Air, iMac, Mac mini, and more.", Link: "https://www.apple.com/mac/"},
		{Rank: 3, Description: "MacBook Air is completely transformed by the power of Apple-designed M1 chip. Up to 3.5x faster CPU, 5x faster graphics, and 18 hours of battery life.", Link: "https://www.apple.com/macbook-air/"},
		{Rank: 4, Description: "The M1 chip and macOS Monterey work together to make the entire system snappier. MacBook Pro wakes instantly from sleep. Everyday tasks from flipping through ...", Link: "https://www.apple.com/macbook-pro-13/"},
		{Rank: 5, Description: "Pay over time, interest-free for your MacBook Air with Apple Card Monthly Installments. Free delivery. Select a model or customize your own.", Link: "https://www.apple.com/shop/buy-mac/macbook-air"},
		{Rank: 6, Description: "Shop Best Buy for a new Apple MacBook. Choose a MacBook Air, Pro, or Retina Display model along with all the accessories you need.", Link: "https://www.bestbuy.com/site/all-laptops/macbooks/pcmcat247400050001.c"},
		{Rank: 7, Description: "Apple has increased the weight and thickness of these MacBooks. And the overall shape and ergonomics are reminiscent of MacBook Pro models from ...", Link: "https://9to5mac.com/2021/11/10/macbook-pro-14-1999-m1-pro-review-video/"},
		{Rank: 8, Description: "The next-generation MacBook Air refresh coming in 2022 will see Apple introduce the biggest design update to the MacBook Air since 2010, ...", Link: "https://www.macrumors.com/guide/2022-macbook-air/"},
		{Rank: 9, Description: "Deals on the MacBook Air, Pro, and Mac Mini ... The 2020 MacBook Air, one of the best laptops you can get, is now available at Amazon ...", Link: "https://www.theverge.com/22399419/apple-macbook-air-pro-mac-mini-imac-deals"},
		{Rank: 10, Description: "The MacBook is a brand of Macintosh notebook computers designed and marketed by Apple Inc. that use Apple's macOS operating system since 2006.", Link: "https://en.wikipedia.org/wiki/MacBook"},
	}
)

func (t *LuminatiTestSuite) TestSerps_CheckURL() {
	tt := map[string]struct {
		serps Serps
		want  Domain
	}{
		"Organic": {
			Serps{
				Organic:  OrganicTestData,
				Features: nil,
			},
			Domain{
				Query: Query{
					Rank:        1,
					Description: "MacBook Pro. Our most powerful notebooks. Fast M1 processors, incredible graphics, and spectacular Retina displays. Now available in a 14-inch model.",
					Link:        "https://www.apple.com/macbook-pro/",
				},
				Results: []Organic{
					OrganicTestData[0], OrganicTestData[1], OrganicTestData[2], OrganicTestData[3], OrganicTestData[4],
				},
			},
		},
		"Organic with Features": {
			Serps{
				Organic: OrganicTestData,
				mappedFeatures: map[string]string{
					"images": TestURL,
				},
			},
			Domain{
				Query: Query{
					Rank:        1,
					Description: "MacBook Pro. Our most powerful notebooks. Fast M1 processors, incredible graphics, and spectacular Retina displays. Now available in a 14-inch model.",
					Link:        "https://www.apple.com/macbook-pro/",
					Features:    "images",
				},
				Results: []Organic{
					OrganicTestData[0], OrganicTestData[1], OrganicTestData[2], OrganicTestData[3], OrganicTestData[4],
				},
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.serps.CheckURL(TestURL)
			t.Equal(test.want.Results, got.Results)
			t.Equal(test.want.Query, got.Query)
		})
	}
}

func (t *LuminatiTestSuite) TestSerps_GetFeatures() {
	tt := map[string]struct {
		serps Serps
		want  []string
	}{
		"One": {
			Serps{
				mappedFeatures: map[string]string{
					"images": TestURL,
				},
			},
			[]string{"images"},
		},
		"Two": {
			Serps{
				mappedFeatures: map[string]string{
					"images":            TestURL,
					"people_also_asked": TestURL,
				},
			},
			[]string{"images", "people_also_asked"},
		},
		"Excluded": {
			Serps{
				mappedFeatures: map[string]string{
					"images":            "wrong",
					"people_also_asked": "wrong",
				},
			},
			nil,
		},
	}

	for name, test := range tt {
		t.Run(name, func() {
			got := test.serps.getFeatures(TestURL)
			for _, feature := range test.want {
				t.Contains(got, feature)
			}
		})
	}
}
