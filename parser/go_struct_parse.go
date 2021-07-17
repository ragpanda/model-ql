package parser

import (
	"context"
	"errors"
	"fmt"
	"go/ast"
	goparser "go/parser"
	"go/token"

	"github.com/ragpanda/model-ql/util"
)

type ParseGolangStruct struct {
	ctx context.Context
}

func NewParseGolangStruct(ctx context.Context) *ParseGolangStruct {
	p := &ParseGolangStruct{
		ctx: ctx,
	}
	return p
}

func (p *ParseGolangStruct) ParseType(pkgPath, ident string) (*Type, error) {
	fset := token.NewFileSet()
	pkgMap, err := goparser.ParseDir(fset, pkgPath, nil, goparser.ParseComments)
	if err != nil {
		util.Error(p.ctx, "parse dir err %s", err.Error())
		return nil, err
	}

	pkgScope := pkgMap[pkgPath]
	obj := pkgScope.Scope.Lookup(ident)
	if obj == nil {
		return nil, errors.New(fmt.Sprintf("Not find %s %s", pkgPath, ident))
	}

	dSpec := obj.Decl.(*ast.TypeSpec)
	tSpec := dSpec.Type.(*ast.StructType)

	return p.parseObject(tSpec)
}

func (p *ParseGolangStruct) parseObject(sType *ast.StructType) (*Type, error) {

	for _, field := range sType.Fields.List {
		if ident, ok := field.Type.(*ast.StarExpr); ok {
			x, ok := ident.X.(*ast.Ident)
			if !ok {
				continue
			}

			util.Info(p.ctx, "%s", x)

		}
	}

	return nil, nil
}

/*

input file or package

locate struct iter





*/
