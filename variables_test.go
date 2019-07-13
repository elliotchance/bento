package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNumber(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		assert.Equal(t, "123.000000", NewNumber("123").FloatString(6))
	})

	t.Run("Float", func(t *testing.T) {
		assert.Equal(t, "1.230000", NewNumber("1.23").FloatString(6))
	})

	t.Run("NegativeFloat", func(t *testing.T) {
		assert.Equal(t, "-1.230000", NewNumber("-1.23").FloatString(6))
	})
}
