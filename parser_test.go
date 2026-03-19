package bival_test

import (
	"testing"

	"github.com/neatflowcv/bival"
	"github.com/stretchr/testify/require"
)

func TestParseFileSample(t *testing.T) {
	t.Parallel()

	// Arrange
	var (
		count     int
		totalSize int64
		first     bival.Record
		last      bival.Record
	)

	// Act
	err := bival.ParseFile("sample.json", func(record *bival.Record) error {
		count++
		totalSize += record.Entry.Meta.Size

		if count == 1 {
			first = *record
		}

		last = *record

		return nil
	})

	// Assert
	require.NoError(t, err)
	require.Equal(t, 4, count)
	require.EqualValues(t, 10, totalSize)
	require.Equal(t, "plain", first.Type)
	require.Equal(t, "test.txt", first.Entry.Name)
	require.False(t, first.Entry.Exists)
	require.Equal(t, "olh", last.Type)
	require.Equal(t, "test.txt", last.Entry.Key.Name)
	require.Equal(t, "PDGqmtJA7imna.RLH.1nsBhSy1ZWf9m", last.Entry.Key.Instance)
}
