package feedripper

import "testing"

func TestSize(t *testing.T) {
	t.Run("small", func(t *testing.T) {
		s := sizeOf(512)
		if s != "512B" {
			t.Error(s)
		}
	})

	t.Run("1k", func(t *testing.T) {
		s := sizeOf(1000)
		if s != "1.0K" {
			t.Error(s)
		}
	})

	t.Run("med", func(t *testing.T) {
		s := sizeOf(51200)
		if s != "51.2K" {
			t.Error(s)
		}
	})

	t.Run("large", func(t *testing.T) {
		s := sizeOf(5120000)
		if s != "5.1M" {
			t.Error(s)
		}
	})
}
