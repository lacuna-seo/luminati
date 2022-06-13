// Copyright 2020 The Reddico Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package stasher

import (
	"github.com/ainsleyclark/redigo"
)

// Cache interface used for mocking.
type Cache interface {
	redigo.Store
}
