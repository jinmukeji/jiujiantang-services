package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type GenerateTokenTestSuite struct {
	suite.Suite
}

func (suite *GenerateTokenTestSuite) TestToken() {
	t := suite.T()
	ctx := context.Background()
	ctx = AddContextToken(ctx, "token")
	tk, ok := TokenFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, "token", tk)
}

func TestGenerateToken(t *testing.T) {
	suite.Run(t, new(GenerateTokenTestSuite))
}
