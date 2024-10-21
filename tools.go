// MIT No Attribution
//
// Copyright (c) 2024 Eric Cornelissen
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

//go:build tools

package main

import (
	_ "4d63.com/gochecknoinits"
	_ "github.com/alexkohler/dogsled/cmd/dogsled"
	_ "github.com/alexkohler/nakedret/v2/cmd/nakedret"
	_ "github.com/alexkohler/prealloc"
	_ "github.com/alexkohler/unimport"
	_ "github.com/butuzov/ireturn/cmd/ireturn"
	_ "github.com/catenacyber/perfsprint"
	_ "github.com/dkorunic/betteralign/cmd/betteralign"
	_ "github.com/go-critic/go-critic/cmd/gocritic"
	_ "github.com/google/capslock/cmd/capslock"
	_ "github.com/gordonklaus/ineffassign"
	_ "github.com/jgautheron/goconst/cmd/goconst"
	_ "github.com/kisielk/errcheck"
	_ "github.com/kunwardeep/paralleltest"
	_ "github.com/mdempsky/unconvert"
	_ "github.com/nishanths/exhaustive/cmd/exhaustive"
	_ "github.com/polyfloyd/go-errorlint"
	_ "github.com/remyoudompheng/go-misc/deadcode"
	_ "github.com/tetafro/godot/cmd/godot"
	_ "github.com/tomarrell/wrapcheck/v2/cmd/wrapcheck"
	_ "github.com/ultraware/whitespace/cmd/whitespace"
	_ "gitlab.com/bosi/decorder/cmd/decorder"
	_ "go.uber.org/nilaway/cmd/nilaway"
	_ "golang.org/x/tools/cmd/goimports"
	_ "golang.org/x/tools/go/analysis/passes/shadow/cmd/shadow"
	_ "golang.org/x/vuln/cmd/govulncheck"
	_ "honnef.co/go/tools/cmd/staticcheck"
	_ "mvdan.cc/unparam"
)
