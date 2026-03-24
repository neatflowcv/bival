package bilist_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/stretchr/testify/require"
)

func TestReadFileSample(t *testing.T) {
	t.Parallel()

	var (
		count     int
		totalSize int64
		first     bilist.Record
		last      bilist.Record
	)

	err := bilist.ReadFile("../../sample.json", func(record *bilist.Record) error {
		count++
		totalSize += record.Entry.Meta.Size

		if count == 1 {
			first = *record
		}

		last = *record

		return nil
	})

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
