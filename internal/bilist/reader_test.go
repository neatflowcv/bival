package bilist_test

import (
	"encoding/json"
	"errors"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/stretchr/testify/require"
)

func TestReaderReadSample(t *testing.T) {
	t.Parallel()

	file, err := os.Open("../../sample.json")
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	reader := bilist.NewReader(file)

	var (
		count     int
		totalSize int64
		first     bilist.Record
		last      bilist.Record
	)

	for {
		record, err := reader.Read()
		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err)

		count++
		totalSize += record.Entry.Meta.Size

		if count == 1 {
			first = *record
		}

		last = *record
	}

	require.Equal(t, 4, count)
	require.EqualValues(t, 10, totalSize)
	require.Equal(t, "plain", first.Type)
	require.Equal(t, "test.txt", first.Entry.Name)
	require.False(t, first.Entry.Exists)
	require.Equal(t, "olh", last.Type)
	require.Equal(t, "test.txt", last.Entry.Key.Name)
	require.Equal(t, "PDGqmtJA7imna.RLH.1nsBhSy1ZWf9m", last.Entry.Key.Instance)

	_, err = reader.Read()
	require.ErrorIs(t, err, io.EOF)
}

func TestReaderReadRawSample(t *testing.T) {
	t.Parallel()

	file, err := os.Open("../../sample.json")
	require.NoError(t, err)

	defer func() {
		_ = file.Close()
	}()

	reader := bilist.NewReader(file)

	count := 0

	for {
		raw, err := reader.ReadRaw()
		if errors.Is(err, io.EOF) {
			break
		}

		require.NoError(t, err)
		require.NotEmpty(t, raw)

		var record bilist.Record

		err = json.Unmarshal(raw, &record)
		require.NoError(t, err)

		count++
	}

	require.Equal(t, 4, count)

	raw, err := reader.ReadRaw()
	require.Nil(t, raw)
	require.ErrorIs(t, err, io.EOF)
}

func TestReaderReadEmptyArray(t *testing.T) {
	t.Parallel()

	reader := bilist.NewReader(strings.NewReader("[]"))

	record, err := reader.Read()
	require.Nil(t, record)
	require.ErrorIs(t, err, io.EOF)

	record, err = reader.Read()
	require.Nil(t, record)
	require.ErrorIs(t, err, io.EOF)
}

func TestReaderReadRawEmptyArray(t *testing.T) {
	t.Parallel()

	reader := bilist.NewReader(strings.NewReader("[]"))

	raw, err := reader.ReadRaw()
	require.Nil(t, raw)
	require.ErrorIs(t, err, io.EOF)

	raw, err = reader.ReadRaw()
	require.Nil(t, raw)
	require.ErrorIs(t, err, io.EOF)
}

func TestReaderReadRejectsNonArray(t *testing.T) {
	t.Parallel()

	reader := bilist.NewReader(strings.NewReader(`{"type":"plain"}`))

	record, err := reader.Read()
	require.Nil(t, record)
	require.EqualError(t, err, "expected top-level array: got {")
}

func TestReaderReadRawRejectsNonArray(t *testing.T) {
	t.Parallel()

	reader := bilist.NewReader(strings.NewReader(`{"type":"plain"}`))

	raw, err := reader.ReadRaw()
	require.Nil(t, raw)
	require.EqualError(t, err, "expected top-level array: got {")
}

func TestReaderReadRejectsMissingClosingArray(t *testing.T) {
	t.Parallel()

	reader := bilist.NewReader(strings.NewReader(`[{"type":"plain","idx":"1","entry":{}}`))

	record, err := reader.Read()
	require.NoError(t, err)
	require.NotNil(t, record)

	record, err = reader.Read()
	require.Nil(t, record)
	require.EqualError(t, err, "read closing token: EOF")
}

func TestReaderReadRawRejectsMissingClosingArray(t *testing.T) {
	t.Parallel()

	reader := bilist.NewReader(strings.NewReader(`[{"type":"plain","idx":"1","entry":{}}`))

	raw, err := reader.ReadRaw()
	require.NoError(t, err)
	require.NotNil(t, raw)

	raw, err = reader.ReadRaw()
	require.Nil(t, raw)
	require.EqualError(t, err, "read closing token: EOF")
}
