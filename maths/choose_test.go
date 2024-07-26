package maths

import (
	"errors"
	"testing"

	assert "github.com/stretchr/testify/require"
)

func TestChoose(t *testing.T) {
	assert.Equal(t, "yes", Choose(true, "yes", "no"))
	assert.Equal(t, "no", Choose(false, "yes", "no"))
}

func TestChooseFunc(t *testing.T) {
	yes := func() string { return "yes" }
	no := func() string { return "no" }
	assert.Equal(t, "yes", ChooseFunc(true, yes, no))
	assert.Equal(t, "no", ChooseFunc(false, yes, no))
}

func TestChooseFuncE(t *testing.T) {
	yes := func() (string, error) { return "yes", nil }
	no := func() (string, error) { return "", errors.New("no") }

	res, err := ChooseFuncE(true, yes, no)
	assert.NoError(t, err)
	assert.Equal(t, "yes", res)

	res, err = ChooseFuncE(false, yes, no)
	assert.Error(t, err, "no")
	assert.Equal(t, "", res)
}
