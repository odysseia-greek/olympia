package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateNewConfig(t *testing.T) {
	env := "LOCAL"

	result, err := CreateNewConfig(env)

	// Assert the expected results
	assert.NotNil(t, result)
	assert.Nil(t, err)
	assert.Equal(t, "dictionary", result.Index)
	assert.NotNil(t, result.Elastic)
}
