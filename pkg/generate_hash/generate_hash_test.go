package generate_hash

import (
	"testing"
	"github.com/stretchr/testify/suite"
)

type GenerateHashTestSuite struct {
	suite.Suite
}

func (suite *GenerateHashTestSuite) TestGenerateHashShouldNotFail() {
	suite.Equal(
		GenerateHash("salt", "secret"),
		Hashed("2gDsLm/57U00KyShbiYsgvPIsQs="),
	)
}

func (suite *GenerateHashTestSuite) TestGenerateHashShouldProduceDifferentOutputForDifferentSalt() {
	suite.NotEqual(
		GenerateHash("salt1", "secret"),
		GenerateHash("salt2", "secret"),
	)
}

func (suite *GenerateHashTestSuite) TestGenerateHashShouldProduceDifferentOutputForDifferentSecret() {
	suite.NotEqual(
		GenerateHash("salt", "secret"),
		GenerateHash("salt", "secreT"),
	)
}

func TestGenerateHash(t *testing.T) {
	suite.Run(t, new(GenerateHashTestSuite))
}
