//nolint:testpackage // Tests exercise internal helpers directly to cover command internals.
package cli

import (
	"bytes"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeFileSummarizesSortedInput(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha"),
		versionedHeadRecordMap("beta"),
		versionedPlainRecordMap("beta", "v1"),
		versionedInstanceRecordMap("beta", "v1"),
		versionedOLHRecordMap("beta", "v1"),
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
		recordMap("alpha", "plain", "alpha"),
		versionedHeadRecordMap("beta"),
		versionedPlainRecordMap("beta", "v1"),
		versionedInstanceRecordMap("beta", "v1"),
		versionedOLHRecordMap("beta", "v1"),
		recordMap("alpha", "plain", "alpha"),
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
		versionedHeadRecordMap("alpha"),
		versionedPlainRecordMapWithPendingMap("alpha", "v1"),
		versionedInstanceRecordMap("alpha", "v1"),
		versionedOLHRecordMap("alpha", "v1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(
		t,
		"problem name=\"alpha\" code=\"entry.pending.exists\" code=\"pair.version.mismatched\"\n",
		buf.String(),
	)
}

func TestAnalyzeFileReportsPendingLogGroupOnce(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		versionedHeadRecordMap("alpha"),
		versionedPlainRecordMap("alpha", "v1"),
		versionedInstanceRecordMap("alpha", "v1"),
		versionedOLHRecordMapWithPendingLog("alpha", "v1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" code=\"entry.pending.exists\" code=\"olh.invalid\"\n", buf.String())
}

func TestAnalyzeFileReportsInvalidOLHCount(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		versionedHeadRecordMap("alpha"),
		versionedPlainRecordMap("alpha", "v1"),
		versionedInstanceRecordMap("alpha", "v1"),
		versionedInstanceRecordMap("alpha", "v2"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(
		t,
		"problem name=\"alpha\" code=\"pair.plain.missing\" code=\"olh.missing\"\n",
		buf.String(),
	)
}

func TestAnalyzeFileReportsInvalidInstanceCount(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		versionedHeadRecordMap("alpha"),
		versionedPlainRecordMap("alpha", "v1"),
		versionedPlainRecordMap("alpha", "v2"),
		versionedInstanceRecordMap("alpha", "v1"),
		versionedOLHRecordMap("alpha", "v1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" code=\"pair.instance.missing\"\n", buf.String())
}

func TestAnalyzeFileReportsTooManyVersionedEntries(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		versionedHeadRecordMap("alpha"),
		versionedPlainRecordMap("alpha", "v1"),
		versionedPlainRecordMap("alpha", "v2"),
		versionedPlainRecordMap("alpha", "v3"),
		versionedPlainRecordMap("alpha", "v4"),
		versionedInstanceRecordMap("alpha", "v1"),
		versionedInstanceRecordMap("alpha", "v2"),
		versionedInstanceRecordMap("alpha", "v3"),
		versionedInstanceRecordMap("alpha", "v4"),
		versionedOLHRecordMap("alpha", "v1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" code=\"entry.versioned.count.exceeded\"\n", buf.String())
}

func TestAnalyzeFileReportsPendingEntryAndTooManyVersionedEntries(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		versionedHeadRecordMap("alpha"),
		versionedPlainRecordMapWithPendingMap("alpha", "v1"),
		versionedPlainRecordMap("alpha", "v2"),
		versionedPlainRecordMap("alpha", "v3"),
		versionedInstanceRecordMap("alpha", "v1"),
		versionedInstanceRecordMap("alpha", "v2"),
		versionedInstanceRecordMap("alpha", "v3"),
		versionedOLHRecordMap("alpha", "v1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(
		t,
		"problem name=\"alpha\" code=\"entry.pending.exists\" code=\"pair.version.mismatched\"\n",
		buf.String(),
	)
}

func TestAnalyzeFileReportsStaleDeleteMarkerOLH(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		versionedHeadRecordMap("alpha"),
		versionedDeleteMarkerRecordMap("alpha", "delete-v1", sampleStaleVersionedRecordMTime()),
		versionedDeleteMarkerInstanceRecordMap("alpha", "delete-v1", sampleStaleVersionedRecordMTime()),
		versionedOLHDeleteMarkerRecordMap("alpha", "delete-v1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"alpha\" code=\"version.stale\"\n", buf.String())
}

func TestAnalyzeFileReportsOnlyProblemGroups(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha"),
		versionedHeadRecordMap("beta"),
		versionedPlainRecordMap("beta", "v1"),
		versionedPlainRecordMap("beta", "v2"),
		versionedInstanceRecordMap("beta", "v1"),
		versionedOLHRecordMap("beta", "v1"),
		versionedHeadRecordMap("gamma"),
		versionedPlainRecordMap("gamma", "v1"),
		versionedInstanceRecordMap("gamma", "v1"),
		versionedOLHRecordMap("gamma", "v1"),
	})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Equal(t, "problem name=\"beta\" code=\"pair.instance.missing\"\n", buf.String())
}

func TestAnalyzeFileHandlesZeroFloatMTime(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	record := recordMap("alpha", "plain", "alpha")
	entry, ok := record["entry"].(map[string]any)
	require.True(t, ok)

	meta, ok := entry["meta"].(map[string]any)
	require.True(t, ok)

	meta["mtime"] = "0.000000"

	writeRecords(t, inputPath, []map[string]any{record})

	var buf bytes.Buffer

	logger := log.New(&buf, "", 0)

	err := analyzeFile(inputPath, logger)
	require.NoError(t, err)
	require.Empty(t, buf.String())
}

func versionedHeadRecordMap(name string) map[string]any {
	meta := map[string]any{
		"category":           0,
		"size":               0,
		"mtime":              "0.000000",
		"etag":               "",
		"storage_class":      "",
		"owner":              "",
		"owner_display_name": "",
		"content_type":       "",
		"accounted_size":     0,
		"user_data":          "",
		"appendable":         false,
	}

	return map[string]any{
		"type": "plain",
		"idx":  name,
		"entry": map[string]any{
			"name":            name,
			"instance":        "",
			"ver":             map[string]any{"pool": -1, "epoch": 0},
			"locator":         "",
			"exists":          false,
			"meta":            meta,
			"tag":             "",
			"flags":           8,
			"pending_map":     []any{},
			"versioned_epoch": 0,
		},
	}
}

func versionedPlainRecordMap(name string, instance string) map[string]any {
	meta := map[string]any{
		"category":           1,
		"size":               4,
		"mtime":              sampleVersionedRecordMTime(),
		"etag":               "etag",
		"storage_class":      "",
		"owner":              "test",
		"owner_display_name": "test",
		"content_type":       "text/plain",
		"accounted_size":     4,
		"user_data":          "",
		"appendable":         false,
	}

	return map[string]any{
		"type": "plain",
		"idx":  name + "\u0000v913\u0000i" + instance,
		"entry": map[string]any{
			"name":            name,
			"instance":        instance,
			"ver":             map[string]any{"pool": 186, "epoch": 1147},
			"locator":         "",
			"exists":          true,
			"meta":            meta,
			"tag":             "tag",
			"flags":           1,
			"pending_map":     []any{},
			"versioned_epoch": 2,
		},
	}
}

func sampleVersionedRecordMTime() string {
	now := time.Now().UTC().Truncate(time.Second)

	return now.Add(-24 * time.Hour).Format(time.RFC3339Nano)
}

func sampleStaleVersionedRecordMTime() string {
	now := time.Now().UTC().Truncate(time.Second)

	return now.Add(-8 * 24 * time.Hour).Format(time.RFC3339Nano)
}

func versionedInstanceRecordMap(name string, instance string) map[string]any {
	record := versionedPlainRecordMap(name, instance)
	record["type"] = "instance"
	record["idx"] = "\u00801000_" + name + "\u0000i" + instance

	return record
}

func versionedDeleteMarkerRecordMap(name string, instance string, mtime string) map[string]any {
	record := versionedPlainRecordMap(name, instance)
	entry := mustEntryMap(record)
	meta := mustMetaMap(entry)

	entry["exists"] = false
	entry["tag"] = "delete-marker"
	entry["flags"] = 7
	entry["versioned_epoch"] = 3
	entry["ver"] = map[string]any{"pool": -1, "epoch": 0}
	meta["category"] = 0
	meta["size"] = 0
	meta["mtime"] = mtime
	meta["etag"] = ""
	meta["content_type"] = ""
	meta["accounted_size"] = 0

	return record
}

func versionedDeleteMarkerInstanceRecordMap(name string, instance string, mtime string) map[string]any {
	record := versionedDeleteMarkerRecordMap(name, instance, mtime)
	record["type"] = "instance"
	record["idx"] = "\u00801000_" + name + "\u0000i" + instance

	return record
}

func versionedOLHRecordMap(name string, instance string) map[string]any {
	return map[string]any{
		"type": "olh",
		"idx":  "\u00801001_" + name,
		"entry": map[string]any{
			"key": map[string]any{
				"name":     name,
				"instance": instance,
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

func versionedOLHDeleteMarkerRecordMap(name string, instance string) map[string]any {
	record := versionedOLHRecordMap(name, instance)
	entry := mustEntryMap(record)
	entry["delete_marker"] = true

	return record
}

func versionedPlainRecordMapWithPendingMap(name string, instance string) map[string]any {
	record := versionedPlainRecordMap(name, instance)

	entry, ok := record["entry"].(map[string]any)
	if !ok {
		panic("record entry must be a map")
	}

	entry["pending_map"] = []any{map[string]any{"op": "test"}}

	return record
}

func versionedOLHRecordMapWithPendingLog(name string, instance string) map[string]any {
	record := versionedOLHRecordMap(name, instance)

	entry := mustEntryMap(record)

	entry["pending_log"] = []any{map[string]any{"op": "test"}}

	return record
}

func mustEntryMap(record map[string]any) map[string]any {
	entry, ok := record["entry"].(map[string]any)
	if !ok {
		panic("record entry must be a map")
	}

	return entry
}

func mustMetaMap(entry map[string]any) map[string]any {
	meta, ok := entry["meta"].(map[string]any)
	if !ok {
		panic("record meta must be a map")
	}

	return meta
}
