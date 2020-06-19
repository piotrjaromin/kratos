package utils

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/thoas/go-funk"

	"io/ioutil"
)

func TestStringWriteReadShouldPassthourgBytes(t *testing.T) {
	swr := &StringWriterReader{}

	chunk1 := []byte{1, 2, 3, 4}
	chunk2 := []byte{11, 22, 223, 24}
	chunk3 := []byte{1, 2, 3, 213, 123, 123, 123, 123}

	allChunks := [][]byte{chunk1, chunk2, chunk3}

	for index, chunk := range allChunks {
		n, err := swr.Write(chunk)
		require.Nil(t, err)
		require.Equal(t, len(chunk), n, fmt.Sprintf("error for chunk %d", index))
	}

	allRead, err := ioutil.ReadAll(swr)
	require.Nil(t, err)

	expected := funk.FlattenDeep(allChunks)
	require.Equal(t, expected, allRead)
}
