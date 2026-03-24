package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/neatflowcv/bival"
	"github.com/stretchr/testify/require"
)

func TestSortFileStableByNameAcrossChunks(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	outputPath := filepath.Join(t.TempDir(), "output.json")

	records := []map[string]any{
		recordMap("beta", "plain", "beta-1"),
		recordMap("alpha", "plain", "alpha-1"),
		recordMap("beta", "instance", "beta-2"),
		olhRecordMap("alpha", "alpha-olh"),
		recordMap("alpha", "plain", "alpha-2"),
		recordMap("gamma", "plain", "gamma-1"),
	}

	writeRecords(t, inputPath, records)

	err := sortFile(inputPath, outputPath, 200)
	require.NoError(t, err)

	sorted := readRecords(t, outputPath)
	require.Len(t, sorted, 6)

	names := []string{
		recordName(&sorted[0]),
		recordName(&sorted[1]),
		recordName(&sorted[2]),
		recordName(&sorted[3]),
		recordName(&sorted[4]),
		recordName(&sorted[5]),
	}
	require.Equal(t, []string{"alpha", "alpha", "alpha", "beta", "beta", "gamma"}, names)

	require.Equal(t, "alpha-1", sorted[0].Idx)
	require.Equal(t, "alpha-olh", sorted[1].Idx)
	require.Equal(t, "alpha-2", sorted[2].Idx)
	require.Equal(t, "beta-1", sorted[3].Idx)
	require.Equal(t, "beta-2", sorted[4].Idx)
	require.Equal(t, "gamma-1", sorted[5].Idx)
}

func TestSortFileEmptyInput(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	outputPath := filepath.Join(t.TempDir(), "output.json")

	writeRecords(t, inputPath, []map[string]any{})

	err := sortFile(inputPath, outputPath, 64)
	require.NoError(t, err)

	// #nosec G304 -- test reads a file created in t.TempDir().
	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)
	require.JSONEq(t, "[]", string(data))
	require.Equal(t, "[]\n", string(data))
}

func TestRunRejectsInvalidChunkBytes(t *testing.T) {
	t.Parallel()

	err := run([]string{"--chunk-bytes", "0", "in.json", "out.json"})
	require.ErrorContains(t, err, "chunk-bytes must be greater than zero")
}

func TestSortFileWritesIndentedJSON(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	outputPath := filepath.Join(t.TempDir(), "output.json")

	writeRecords(t, inputPath, []map[string]any{
		recordMap("beta", "plain", "beta-1"),
		recordMap("alpha", "plain", "alpha-1"),
	})

	err := sortFile(inputPath, outputPath, 128)
	require.NoError(t, err)

	// #nosec G304 -- test reads a file created in t.TempDir().
	data, err := os.ReadFile(outputPath)
	require.NoError(t, err)

	require.Contains(t, string(data), "\n    {\n")
	require.Contains(t, string(data), "\"name\": \"alpha\"")
	require.True(t, len(data) > 0 && data[len(data)-1] == '\n')
}

func writeRecords(t *testing.T, path string, records []map[string]any) {
	t.Helper()

	data, err := json.Marshal(records)
	require.NoError(t, err)

	err = os.WriteFile(path, data, 0o600)
	require.NoError(t, err)
}

func readRecords(t *testing.T, path string) []bival.Record {
	t.Helper()

	// #nosec G304 -- test reads a file created in t.TempDir().
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var records []bival.Record

	err = json.Unmarshal(data, &records)
	require.NoError(t, err)

	return records
}

func recordMap(name string, recordType string, idx string) map[string]any {
	return map[string]any{
		"type": recordType,
		"idx":  idx,
		"entry": map[string]any{
			"name":            name,
			"instance":        "",
			"ver":             map[string]any{"pool": 1, "epoch": 1},
			"locator":         "",
			"exists":          true,
			"meta":            map[string]any{"mtime": "2026-03-06T03:34:11.918188Z"},
			"tag":             "",
			"flags":           0,
			"pending_map":     []any{},
			"versioned_epoch": 0,
		},
	}
}

func olhRecordMap(name string, idx string) map[string]any {
	return map[string]any{
		"type": "olh",
		"idx":  idx,
		"entry": map[string]any{
			"key": map[string]any{
				"name":     name,
				"instance": "inst",
			},
			"delete_marker":   false,
			"epoch":           2,
			"pending_log":     []any{},
			"tag":             "",
			"exists":          true,
			"pending_removal": false,
		},
	}
}
