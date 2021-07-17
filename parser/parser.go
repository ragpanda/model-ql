// Copyright 2012-2015 Samuel Stauffer. All rights reserved.
// Use of this source code is governed by a 3-clause BSD
// license that can be found in the LICENSE file.

package parser

//go:generate pigeon -o grammar.peg.go ./grammar.peg
//go:generate goimports -w ./grammar.peg.go

import (
	"context"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ragpanda/model-ql/util"
)

type Filesystem interface {
	Open(filename string) (io.ReadCloser, error)
	Abs(path string) (string, error)
}

type Parser struct {
	ctx        context.Context
	filesystem Filesystem
	// For handling includes. Can be set to nil to fall back to os package.
}

func NewParser(ctx context.Context) *Parser {
	return &Parser{
		ctx:        ctx,
		filesystem: nil,
	}
}

func (p *Parser) Parse(r io.Reader, opts ...Option) (*CompileUnit, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	name := "<reader>"
	if named, ok := r.(namedReader); ok {
		name = named.Name()
	}
	t, err := Parse(name, b, opts...)
	if err != nil {
		return nil, err
	}

	compileUnit := t.(*CompileUnit)
	p.process(compileUnit)

	return compileUnit, nil
}

func (p *Parser) ParseFile(filename string) (map[string]*CompileUnit, string, error) {
	files := make(map[string]*CompileUnit)

	absPath, err := p.abs(filename)
	if err != nil {
		return nil, "", err
	}

	path := absPath
	for path != "" {
		rd, err := p.open(path)
		if err != nil {
			return nil, "", err
		}
		unit, err := p.Parse(rd)
		if err != nil {
			return nil, "", err
		}
		files[path] = unit

	}

	return files, absPath, nil
}

func (p *Parser) open(path string) (io.ReadCloser, error) {
	if p.filesystem == nil {
		return os.Open(path)
	}
	return p.filesystem.Open(path)
}

func (p *Parser) abs(path string) (string, error) {
	if p.filesystem == nil {
		absPath, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		return filepath.Clean(absPath), nil
	}
	return p.filesystem.Abs(path)
}

type namedReader interface {
	Name() string
}

func (p *Parser) process(compileUnit *CompileUnit) error {
	ctx := p.ctx
	for _, view := range compileUnit.ViewList {
		err := view.process(ctx)
		if err != nil {
			util.Info(ctx, "Err: %s", err.Error())
			return err
		}
	}

	return nil
}
