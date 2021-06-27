package parser

import (
	"context"
	"strings"
	"testing"

	"github.com/ragpanda/model-ql/util"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func Test(t *testing.T) {
	suite.Run(t, &TestSuite{})
}

func (suite *TestSuite) Test1() {
	ctx := context.Background()
	unit, err := parse(`
model a {
	
	select *, name, age as abc from data join a as b on a.aa = b.bb where a = 1 and b=2

}
	
	`)

	util.Info(ctx, "a %s %v", util.Display(unit), err)
}

func parse(contents string) (*CompileUnit, error) {
	parser := &Parser{}
	thrift, err := parser.Parse(strings.NewReader(contents))
	return thrift, err
}
