//nolint:testpackage // Tests exercise internal helpers directly to cover command internals.
package cli

import (
	"bytes"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/neatflowcv/bival/internal/pkg/domain"
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

func TestSortCmdRejectsInvalidChunkBytes(t *testing.T) {
	t.Parallel()

	cmd := SortCmd{
		Input:      "in.json",
		Output:     "out.json",
		ChunkBytes: 0,
	}

	err := cmd.Run()
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

func TestRunSortCommand(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	outputPath := filepath.Join(t.TempDir(), "output.json")

	writeRecords(t, inputPath, []map[string]any{
		recordMap("beta", "plain", "beta-1"),
		recordMap("alpha", "plain", "alpha-1"),
	})

	err := Run([]string{"sort", inputPath, outputPath})
	require.NoError(t, err)

	sorted := readRecords(t, outputPath)
	require.Equal(t, "alpha-1", sorted[0].Idx)
	require.Equal(t, "beta-1", sorted[1].Idx)
}

func TestAnalyzeFileSummarizesSortedInput(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-1"),
		recordMap("beta", "plain", "beta-head"),
		recordMap("beta", "plain", "beta-v1"),
		recordMap("beta", "instance", "beta-i1"),
		olhRecordMap("beta", "beta-olh"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Empty(t, buf.String())
}

func TestAnalyzeFileAcceptsUnsortedInput(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-1"),
		recordMap("beta", "plain", "beta-head"),
		recordMap("beta", "plain", "beta-v1"),
		recordMap("beta", "instance", "beta-i1"),
		olhRecordMap("beta", "beta-olh"),
		recordMap("alpha", "plain", "alpha-2"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Empty(t, buf.String())
}

func TestAnalyzeFileReportsPendingMapGroupOnce(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-head"),
		recordMapWithPendingMap("alpha", "plain", "alpha-v1"),
		recordMap("alpha", "instance", "alpha-i1"),
		olhRecordMap("alpha", "alpha-olh"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" reason=\"pending entry exists\"\n", buf.String())
}

func TestAnalyzeFileReportsPendingLogGroupOnce(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-head"),
		recordMap("alpha", "plain", "alpha-v1"),
		recordMap("alpha", "instance", "alpha-i1"),
		olhRecordMapWithPendingLog("alpha", "alpha-olh"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" reason=\"pending entry exists\"\n", buf.String())
}

func TestAnalyzeFileReportsInvalidOLHCount(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-head"),
		recordMap("alpha", "plain", "alpha-v1"),
		recordMap("alpha", "instance", "alpha-i1"),
		recordMap("alpha", "instance", "alpha-i2"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" reason=\"versioning object must have exactly one olh\"\n", buf.String())
}

func TestAnalyzeFileReportsInvalidInstanceCount(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-head"),
		recordMap("alpha", "plain", "alpha-v1"),
		recordMap("alpha", "plain", "alpha-v2"),
		recordMap("alpha", "instance", "alpha-i1"),
		olhRecordMap("alpha", "alpha-olh"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" reason=\"versioning object must satisfy instance+1==plain\"\n", buf.String())
}

func TestAnalyzeFileReportsOnlyProblemGroups(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-1"),
		recordMap("beta", "plain", "beta-head"),
		recordMap("beta", "plain", "beta-v1"),
		recordMap("beta", "plain", "beta-v2"),
		recordMap("beta", "instance", "beta-i1"),
		olhRecordMap("beta", "beta-olh"),
		recordMap("gamma", "plain", "gamma-head"),
		recordMap("gamma", "plain", "gamma-v1"),
		recordMap("gamma", "instance", "gamma-i1"),
		olhRecordMap("gamma", "gamma-olh"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"beta\" reason=\"versioning object must satisfy instance+1==plain\"\n", buf.String())
}

func TestEntryGroupAddPlainTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := newEntryGroup("alpha")

	group.addPlain(domain.NewPlainEntry(domain.NewDirEntry(
		"plain",
		[]byte("alpha"),
		domain.NewDirPayload(
			domain.NewKey("alpha", ""),
			nil,
			nil,
			nil,
			[]*domain.PendingMap{nil},
		),
	)))

	require.Equal(t, 1, group.plainCount)
	require.True(t, group.hasPendingMap)
}

func TestEntryGroupAddInstanceTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := newEntryGroup("alpha")

	group.addInstance(domain.NewInstanceEntry(domain.NewDirEntry(
		"instance",
		[]byte("alpha"),
		domain.NewDirPayload(
			domain.NewKey("alpha", "v1"),
			nil,
			nil,
			nil,
			[]*domain.PendingMap{nil},
		),
	)))

	require.Equal(t, 1, group.instanceCount)
	require.True(t, group.hasPendingMap)
}

func TestEntryGroupAddOLHTracksPendingLog(t *testing.T) {
	t.Parallel()

	group := newEntryGroup("alpha")

	group.addOLH(domain.NewOLHEntry(
		"olh",
		[]byte("alpha"),
		domain.NewOLHPayload(
			domain.NewKey("alpha", "v1"),
			nil,
			0,
			[]*domain.PendingLog{nil},
			"",
		),
	))

	require.Equal(t, 1, group.olhCount)
	require.True(t, group.hasPendingLog)
}

func TestBuildEntryPreservesPendingPresence(t *testing.T) {
	t.Parallel()

	plainEntry, err := buildEntry(&bilist.Record{
		Type:  "plain",
		Idx:   "alpha",
		Entry: pendingMapEntry("alpha", ""),
	})
	require.NoError(t, err)

	instanceEntry, err := buildEntry(&bilist.Record{
		Type:  "instance",
		Idx:   "alpha-instance",
		Entry: pendingMapEntry("alpha", "v1"),
	})
	require.NoError(t, err)

	olhEntry, err := buildEntry(&bilist.Record{
		Type:  "olh",
		Idx:   "alpha-olh",
		Entry: pendingLogEntry("alpha", "v1"),
	})
	require.NoError(t, err)

	plainTyped, isPlain := plainEntry.(*domain.PlainEntry)
	require.True(t, isPlain)
	require.True(t, plainTyped.HasPendingMap())

	instanceTyped, isInstance := instanceEntry.(*domain.InstanceEntry)
	require.True(t, isInstance)
	require.True(t, instanceTyped.HasPendingMap())

	olhTyped, isOLH := olhEntry.(*domain.OLHEntry)
	require.True(t, isOLH)
	require.True(t, olhTyped.HasPendingLog())
}

func TestAnalyzeFileRejectsUnsupportedType(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "directory", "alpha-1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.ErrorContains(t, err, "unsupported record type")
}

func TestRunAnalyzeCommand(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha-1"),
		recordMap("beta", "instance", "beta-1"),
	})

	err := Run([]string{"analyze", inputPath})
	require.NoError(t, err)
}

func writeRecords(t *testing.T, path string, records []map[string]any) {
	t.Helper()

	data, err := json.Marshal(records)
	require.NoError(t, err)

	err = os.WriteFile(path, data, 0o600)
	require.NoError(t, err)
}

func readRecords(t *testing.T, path string) []bilist.Record {
	t.Helper()

	// #nosec G304 -- test reads a file created in t.TempDir().
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var records []bilist.Record

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

func recordMapWithPendingMap(name string, recordType string, idx string) map[string]any {
	record := recordMap(name, recordType, idx)

	entry, ok := record["entry"].(map[string]any)
	if !ok {
		panic("record entry must be a map")
	}

	entry["pending_map"] = []any{
		map[string]any{"op": "test"},
	}

	return record
}

func olhRecordMapWithPendingLog(name string, idx string) map[string]any {
	record := olhRecordMap(name, idx)

	entry, ok := record["entry"].(map[string]any)
	if !ok {
		panic("record entry must be a map")
	}

	entry["pending_log"] = []any{
		map[string]any{"op": "test"},
	}

	return record
}

func pendingMapEntry(name string, instance string) bilist.Entry {
	entry := bilist.Entry{
		Name:     name,
		Instance: instance,
		Ver:      bilist.Version{Pool: 0, Epoch: 0},
		Locator:  "",
		Exists:   false,
		Meta: bilist.Meta{
			Category:         0,
			Size:             0,
			MTime:            "",
			ETag:             "",
			StorageClass:     "",
			Owner:            "",
			OwnerDisplayName: "",
			ContentType:      "",
			AccountedSize:    0,
			UserData:         "",
			Appendable:       false,
		},
		Tag:            "",
		Flags:          0,
		PendingMap:     []json.RawMessage{json.RawMessage(`{"op":"x"}`)},
		VersionedEpoch: 0,
		Key:            bilist.Key{Name: "", Instance: ""},
		DeleteMarker:   false,
		Epoch:          0,
		PendingLog:     nil,
		PendingRemoval: false,
	}

	return entry
}

func pendingLogEntry(name string, instance string) bilist.Entry {
	entry := bilist.Entry{
		Name:     "",
		Instance: "",
		Ver:      bilist.Version{Pool: 0, Epoch: 0},
		Locator:  "",
		Exists:   false,
		Meta: bilist.Meta{
			Category:         0,
			Size:             0,
			MTime:            "",
			ETag:             "",
			StorageClass:     "",
			Owner:            "",
			OwnerDisplayName: "",
			ContentType:      "",
			AccountedSize:    0,
			UserData:         "",
			Appendable:       false,
		},
		Tag:            "",
		Flags:          0,
		PendingMap:     nil,
		VersionedEpoch: 0,
		Key: bilist.Key{
			Name:     name,
			Instance: instance,
		},
		DeleteMarker:   false,
		Epoch:          0,
		PendingLog:     []json.RawMessage{json.RawMessage(`{"op":"x"}`)},
		PendingRemoval: false,
	}

	return entry
}
