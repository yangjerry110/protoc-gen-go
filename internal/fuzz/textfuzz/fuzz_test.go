// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package textfuzz

import (
	"testing"

	"github.com/yangjerry110/protoc-gen-go/internal/fuzztest"
)

func Test(t *testing.T) {
	fuzztest.Test(t, Fuzz)
}
