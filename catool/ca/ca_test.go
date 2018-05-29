package ca

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionCmd(t *testing.T) {
	cmd := versionCmd()
	assert.NoError(t, cmd.Execute(), "expected version command to succeed")
}
