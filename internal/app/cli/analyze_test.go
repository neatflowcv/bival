//nolint:testpackage // Tests exercise internal helpers directly to cover command internals.
package cli

import (
	"bytes"
	"log"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAnalyzeFileSummarizesSortedInput(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha"),
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
		recordMap("alpha", "plain", "alpha"),
		recordMap("beta", "plain", "beta-head"),
		recordMap("beta", "plain", "beta-v1"),
		recordMap("beta", "instance", "beta-i1"),
		olhRecordMap("beta", "beta-olh"),
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
	require.Equal(t, "problem name=\"alpha\" reason=\"object kind is unknown\"\n", buf.String())
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
	require.Equal(t, "problem name=\"alpha\" reason=\"object kind is unknown\"\n", buf.String())
}

func TestAnalyzeFileReportsOnlyProblemGroups(t *testing.T) {
	t.Parallel()

	inputPath := filepath.Join(t.TempDir(), "input.json")
	writeRecords(t, inputPath, []map[string]any{
		recordMap("alpha", "plain", "alpha"),
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
	require.Equal(t, "problem name=\"beta\" reason=\"object kind is unknown\"\n", buf.String())
}
