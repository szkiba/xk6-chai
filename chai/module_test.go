// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: Apache-2.0

package chai

import (
	"testing"

	"github.com/dop251/goja"
	"github.com/stretchr/testify/assert"
	"go.k6.io/k6/js/modulestest"
)

func TestNewModuleInstance(t *testing.T) {
	t.Parallel()

	runtime := modulestest.NewRuntime(t)
	vu := runtime.VU // nolint:varnamelen
	noop := func(call goja.FunctionCall) goja.Value { return goja.Undefined() }

	assert.NoError(t, vu.Runtime().Set("require", noop))
	assert.NoError(t, vu.Runtime().Set("global", vu.Runtime().GlobalObject()))

	root := New()

	var module *Module

	assert.NotPanics(t, func() { module = root.NewModuleInstance(vu).(*Module) }) // nolint:forcetypeassert
	assert.NotNil(t, module)

	exports := module.Exports()

	assert.NotEmpty(t, exports.Default)
	assert.NotEmpty(t, exports.Named)

	assert.Contains(t, exports.Named, "describe")
	assert.Contains(t, exports.Named, "expect")
}

func TestCompile(t *testing.T) {
	t.Parallel()

	_, err := compile("invalid source", "test")

	assert.Error(t, err)

	prog, err := compile("let a=1", "test")

	assert.NoError(t, err)
	assert.NotNil(t, prog)
}

func TestMustCompile(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() { mustCompile("invalid source", "test") })

	prog, err := compile("let a=1", "test")

	assert.NoError(t, err)
	assert.NotNil(t, prog)
}

func TestExecute(t *testing.T) {
	t.Parallel()

	runtime := goja.New()

	prog, err := compile("throw new Error()", "test")

	assert.NoError(t, err)
	assert.NotNil(t, prog)

	_, err = execute(prog, runtime)

	assert.Error(t, err)

	prog, err = goja.Compile("test", "let a=1", true)

	assert.NoError(t, err)
	assert.NotNil(t, prog)

	_, err = execute(prog, runtime)

	assert.Error(t, err)

	prog, err = goja.Compile("test", "(module, exports) => {delete module['exports']}", true)

	assert.NoError(t, err)
	assert.NotNil(t, prog)

	_, err = execute(prog, runtime)

	assert.Error(t, err)
}

func TestRequire(t *testing.T) {
	t.Parallel()

	runtime := goja.New()

	prog, err := compile("throw new Error()", "test")

	assert.NoError(t, err)
	assert.NotNil(t, prog)

	_, err = require(prog, runtime)

	assert.Error(t, err)
}

func TestMustRequire(t *testing.T) {
	t.Parallel()

	runtime := goja.New()

	prog, err := compile("throw new Error()", "test")

	assert.NoError(t, err)
	assert.NotNil(t, prog)

	assert.Panics(t, func() { mustRequire(prog, runtime) })
}
