package main

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestSnorkelWrite(t *testing.T) {
	t.Run("should write integer fields", func(t *testing.T) {
		s := NewSnorkel("hello", "__s__")
		s.AddIntField("bytes")

		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		s.Write(w, map[string]interface{}{
			"bytes": 1,
		})
		w.Flush()
		assert.Equal(t, `{"table":"hello","token":"__s__","data":{"bytes":1}}`, b.String())
	})
	t.Run("should write string fields", func(t *testing.T) {
		s := NewSnorkel("hello", "__s__")
		s.AddStrField("id")

		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		s.Write(w, map[string]interface{}{
			"id": "s00",
		})
		w.Flush()
		assert.Equal(t, `{"table":"hello","token":"__s__","data":{"id":"s00"}}`, b.String())
	})
	t.Run("should write set fields", func(t *testing.T) {
		s := NewSnorkel("hello", "__s__")
		s.AddSetField("categories")

		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		s.Write(w, map[string]interface{}{
			"categories": []string{"fruits", "books"},
		})
		w.Flush()
		assert.Equal(t, `{"table":"hello","token":"__s__","data":{"categories":["fruits","books"]}}`, b.String())
	})
	t.Run("type mismatches", func(t *testing.T) {
		t.Run("when int fields get other than integers", func(t *testing.T) {
			s := NewSnorkel("hello", "__s__")
			s.AddIntField("bytes")

			badValues := []interface{}{nil, "300", 45.0, true, []interface{}{}}
			for _, bad := range badValues {
				t.Run(fmt.Sprintf("bad case: %v", bad), func(t *testing.T) {
					var b bytes.Buffer
					w := bufio.NewWriter(&b)
					err := s.Write(w, map[string]interface{}{"bytes": bad})
					assert.Error(t, err)
				})
			}
		})
		t.Run("when string fields get other than string values", func(t *testing.T) {
			s := NewSnorkel("hello", "__s__")
			s.AddStrField("id")

			badValues := []interface{}{nil, 300, 45.0, true, []interface{}{}}
			for _, bad := range badValues {
				t.Run(fmt.Sprintf("bad case: %v", bad), func(t *testing.T) {
					var b bytes.Buffer
					w := bufio.NewWriter(&b)
					err := s.Write(w, map[string]interface{}{"id": bad})
					assert.Error(t, err)
				})
			}
		})
		t.Run("when set fields get other than an array of string", func(t *testing.T) {
			s := NewSnorkel("hello", "__s__")
			s.AddSetField("categories")

			badValues := []interface{}{nil, 300, 45.0, true, []int{}, "cat"}
			for _, bad := range badValues {
				t.Run(fmt.Sprintf("bad case: %v", bad), func(t *testing.T) {
					var b bytes.Buffer
					w := bufio.NewWriter(&b)
					err := s.Write(w, map[string]interface{}{"categories": bad})
					assert.Error(t, err)
				})
			}
		})
	})
	t.Run("should filter out unregistered fields", func(t *testing.T) {
		s := NewSnorkel("hello", "__s__")
		s.
			AddIntField("bytes").
			AddStrField("id").
			AddSetField("categories")

		var b bytes.Buffer
		w := bufio.NewWriter(&b)
		s.Write(w, map[string]interface{}{
			"bytes":      10,
			"id":         "abc",
			"categories": []string{"fruits", "books"},
			"source":     "http://example.com",
		})
		w.Flush()
		assert.Equal(
			t,
			`{"table":"hello","token":"__s__","data":{"bytes":10,"categories":["fruits","books"],"id":"abc"}}`,
			b.String())
	})
}
