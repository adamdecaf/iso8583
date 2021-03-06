package encoding

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHex(t *testing.T) {
	enc := &hexEncoder{}

	got, read, err := enc.Decode([]byte("aabbccdd"), 3)
	require.NoError(t, err)
	require.Equal(t, 6, read)
	require.Equal(t, []byte{0xAA, 0xBB, 0xCC}, got)

	got, err = enc.Encode([]byte{0xAA, 0xBB, 0xCC})
	require.NoError(t, err)
	require.Equal(t, []byte("AABBCC"), got)
}
