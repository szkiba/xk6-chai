// SPDX-FileCopyrightText: 2023 Iván Szkiba
//
// SPDX-License-Identifier: Apache-2.0

package chai

import (
	_ "embed"
	"errors"

	"github.com/dop251/goja"
	"github.com/sirupsen/logrus"
	"go.k6.io/k6/js/compiler"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/lib"
)

const filename = "k6chaijs.min.js"

//go:embed k6chaijs.min.js
var content string

type RootModule struct {
	program *goja.Program
}

func New() modules.Module {
	return &RootModule{program: mustCompile(content, filename)}
}

func (root *RootModule) NewModuleInstance(vu modules.VU) modules.Instance { // nolint:varnamelen
	exports := mustRequire(root.program, vu.Runtime())

	return &Module{exports: *exports}
}

type Module struct {
	exports modules.Exports
}

func (mod *Module) Exports() modules.Exports {
	return mod.exports
}

var (
	_ modules.Module   = (*RootModule)(nil)
	_ modules.Instance = (*Module)(nil)
)

func mustRequire(prog *goja.Program, runtime *goja.Runtime) *modules.Exports {
	exports, err := require(prog, runtime)
	if err != nil {
		panic(err)
	}

	return exports
}

func require(prog *goja.Program, runtime *goja.Runtime) (*modules.Exports, error) {
	exports, err := execute(prog, runtime)
	if err != nil {
		return nil, err
	}

	named, assertOK := exports.Export().(map[string]interface{})
	if !assertOK {
		return nil, errInvalidModule
	}

	return &modules.Exports{
		Default: exports.Get("default"),
		Named:   named,
	}, nil
}

func mustCompile(source string, name string) *goja.Program {
	prog, err := compile(source, name)
	if err != nil {
		panic(err)
	}

	return prog
}

func compile(source string, name string) (*goja.Program, error) {
	comp := compiler.New(logrus.StandardLogger())

	comp.Options.CompatibilityMode = lib.CompatibilityModeBase
	comp.Options.Strict = true

	prog, _, err := comp.Compile(source, name, false)
	if err != nil {
		return nil, err
	}

	return prog, nil
}

func execute(prog *goja.Program, runtime *goja.Runtime) (*goja.Object, error) {
	module := runtime.NewObject()
	exports := runtime.NewObject()

	if err := module.Set("exports", exports); err != nil {
		return nil, err
	}

	value, err := runtime.RunProgram(prog)
	if err != nil {
		return nil, err
	}

	callable, assertOK := goja.AssertFunction(value)
	if !assertOK {
		return nil, errInvalidModule
	}

	_, err = callable(exports, module, exports)
	if err != nil {
		return nil, err
	}

	exports, assertOK = module.Get("exports").(*goja.Object)
	if !assertOK {
		return nil, errInvalidModule
	}

	return exports, nil
}

var errInvalidModule = errors.New("invalid module")
