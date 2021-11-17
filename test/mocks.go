// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stasher

import "github.com/lacuna-seo/stash"

// Cache interface used for mocking.
type Cache interface {
	stash.Store
}