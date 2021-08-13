package parser

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type GoSrtructParseTestSuite struct {
	suite.Suite
}

func TestGoSrtructParseSuite(t *testing.T) {
	suite.Run(t, &GoSrtructParseTestSuite{})
}
func (suite *GoSrtructParseTestSuite) TestSuccess() {
	ctx := context.Background()

	NewParseGolangStruct(ctx).ParseType("parse", "a")
}
func (suite *GoSrtructParseTestSuite) TestFail() {
	// ctx := context.Background()
}
