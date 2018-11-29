package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInmemStorage_Put(t *testing.T) {
	storage := NewInmemStorage()
	require.NoError(t, storage.Put([]byte("hello,world")))
}

func TestInmemStorage_Get(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		storage := NewInmemStorage()
		value, err := storage.Get()
		require.Error(t, err)
		require.Nil(t, value)
		require.Equal(t, ErrStorageIsEmpty, err)
	})

	t.Run("ok", func(t *testing.T) {
		storage := NewInmemStorage()
		// fifo
		storage.Put([]byte("hello,world1"))
		storage.Put([]byte("hello,world2"))
		value, err := storage.Get()
		require.NoError(t, err)
		require.Equal(t, []byte("hello,world1"), value)
	})
}
