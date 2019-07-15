package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewNumber(t *testing.T) {
	t.Run("Integer", func(t *testing.T) {
		assert.Equal(t, "123", NewNumber("123", 6).String())
	})

	t.Run("Float", func(t *testing.T) {
		assert.Equal(t, "1.23", NewNumber("1.23", 6).String())
	})

	t.Run("NegativeFloat", func(t *testing.T) {
		assert.Equal(t, "-1.23", NewNumber("-1.23", 6).String())
	})

	t.Run("LargeInteger", func(t *testing.T) {
		assert.Equal(t, "12300", NewNumber("12300", 6).String())
	})

	t.Run("RoundingDown", func(t *testing.T) {
		assert.Equal(t, "123.42", NewNumber("123.421", 2).String())
	})

	t.Run("RoundingUp", func(t *testing.T) {
		assert.Equal(t, "123.43", NewNumber("123.428", 2).String())
	})
}

func TestNumber_Mul(t *testing.T) {
	a := NewNumber("5.5", 1)
	b := NewNumber("6.5", 1)
	c := NewNumber("0", 1)
	c.Mul(a, b) // 35.75 -> 35.8
	assert.Equal(t, "35.8", c.String())

	// To validate that it's not keeping a higher precision internally.
	c.Mul(c, NewNumber("11", 1))

	// This would be 393.25 without correct rounding.
	assert.Equal(t, "393.8", c.String())
}

func TestNumber_Quo(t *testing.T) {
	a := NewNumber("5.5", 1)
	b := NewNumber("6.5", 1)
	c := NewNumber("0", 2)
	c.Quo(a, b) // ~0.8461 = 0.85
	assert.Equal(t, "0.85", c.String())

	// To validate that it's not keeping a higher precision internally.
	c.Mul(c, NewNumber("11", 1))

	// This would be ~9.31 without correct rounding.
	assert.Equal(t, "9.35", c.String())
}
