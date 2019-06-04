package wputil

import "testing"

func TestTrim(t *testing.T) {
	s := "1234567890123456789012345678901234567890"

	t.Run("left", func(t *testing.T) {
		got := TrimLeft(30, s)
		want := "...456789012345678901234567890"

		if got != want {
			t.Error(got)
		}
	})

	t.Run("right", func(t *testing.T) {
		got := TrimRight(30, s)
		want := "123456789012345678901234567..."

		if got != want {
			t.Error(got)
		}
	})
}

func TestSize(t *testing.T) {
	t.Run("small", func(t *testing.T) {
		s := FileSize(512)
		if s != "512B" {
			t.Error(s)
		}
	})

	t.Run("1k", func(t *testing.T) {
		s := FileSize(1000)
		if s != "1.0K" {
			t.Error(s)
		}
	})

	t.Run("med", func(t *testing.T) {
		s := FileSize(51200)
		if s != "51.2K" {
			t.Error(s)
		}
	})

	t.Run("large", func(t *testing.T) {
		s := FileSize(5120000)
		if s != "5.1M" {
			t.Error(s)
		}
	})
}
