// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package luminati

import "time"

func (t *LuminatiTestSuite) TestMeta_Process() {
	meta := Meta{
		CacheKey:    "cache-key",
		RequestURL:  "https://google.com",
		RequestTime: time.Now(),
	}
	time.Sleep(1 * time.Second)
	got := meta.process()
	t.Equal("cache-key", got.CacheKey)
	t.Equal("https://google.com", got.RequestURL)
	t.WithinDuration(time.Now(), got.ResponseTime, 100*time.Millisecond)
}
