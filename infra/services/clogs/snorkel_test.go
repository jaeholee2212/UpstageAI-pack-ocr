package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSnorkel(t *testing.T) {
	s := NewSnorkel("hello", "__snorkel-relay__")
	require.NotNil(t, s, "a snorkel instance should be created")
}

func TestSpec(t *testing.T) {
	t.Run("should add an int type spec", func(t *testing.T) {
		s := NewSnorkel("hello", "__s__")
		s.AddIntField("count")
		ftype, ok := s.spec["count"]
		require.Equal(t, ftype, IntType)
		require.True(t, ok)
	})
	t.Run("should add an string type spec", func(t *testing.T) {
		s := NewSnorkel("hello", "__s__")
		s.AddStrField("id")
		ftype, ok := s.spec["id"]
		require.Equal(t, ftype, StrType)
		require.True(t, ok)
	})
	t.Run("should add an set type spec", func(t *testing.T) {
		s := NewSnorkel("hello", "__s__")
		s.AddSetField("categories")
		ftype, ok := s.spec["categories"]
		require.Equal(t, ftype, SetType)
		require.True(t, ok)
	})
}
