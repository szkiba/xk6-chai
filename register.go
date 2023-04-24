// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: Apache-2.0

package chai

import (
	"github.com/szkiba/xk6-chai/chai"
	"go.k6.io/k6/js/modules"
)

func register() {
	modules.Register("k6/x/chai", chai.New())
}

func init() { //nolint:gochecknoinits
	register()
}
