package checksum

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyLuhn(t *testing.T) {
	assert.True(t, VerifyLuhn("1234567897"))
}
