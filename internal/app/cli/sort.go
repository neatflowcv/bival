package cli

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/neatflowcv/bival"
)

const chunkRecordOverhead = 16

var (
	errInvalidChunkBytes = errors.New("chunk-bytes must be greater than zero")
	errExpectedTopArray  = errors.New("expected top-level array")
	errExpectedEndArray  = errors.New("expected closing array")
	errInvalidHeapItem   = errors.New("invalid heap item type")
)

type SortCmd struct {
	Input      string `arg:""             help:"Input BI list JSON file."                                   type:"path"`
	Output     string `arg:""             help:"Output path for the sorted BI list JSON file."              type:"path"`
	ChunkBytes int64  `default:"67108864" help:"Maximum chunk size held in memory before spilling to disk."`
}

func (cmd *SortCmd) Run() error {
	if cmd.ChunkBytes <= 0 {
		return errInvalidChunkBytes
	}

	return sortFile(cmd.Input, cmd.Output, cmd.ChunkBytes)
}

func sortFile(inputPath string, outputPath string, chunkBytes int64) error {
	log.Printf("input=%q output=%q chunk_bytes=%d", inputPath, outputPath, chunkBytes)

	inputFile, err := os.Open(filepath.Clean(inputPath))
	if err != nil {
		return fmt.Errorf("open input: %w", err)
	}

	defer func() {
		_ = inputFile.Close()
	}()

	tempDir, err := os.MkdirTemp("", "bisort-*")
	if err != nil {
		return fmt.Errorf("create temp dir: %w", err)
	}

	defer func() {
		_ = os.RemoveAll(tempDir)
	}()

	chunkPaths, err := writeSortedChunks(inputFile, tempDir, chunkBytes)
	if err != nil {
		return err
	}

	outputFile, err := os.Create(filepath.Clean(outputPath))
	if err != nil {
		return fmt.Errorf("create output: %w", err)
	}

	mergeErr := mergeChunks(outputFile, chunkPaths)
	closeErr := outputFile.Close()

	if mergeErr != nil {
		return mergeErr
	}

	if closeErr != nil {
		return fmt.Errorf("close output: %w", closeErr)
	}

	return nil
}

type chunkRecord struct {
	Seq  int64           `json:"seq"`
	Name string          `json:"name"`
	Raw  json.RawMessage `json:"raw"`
}

func writeSortedChunks(r io.Reader, tempDir string, chunkBytes int64) ([]string, error) {
	dec := json.NewDecoder(r)

	err := expectArrayStart(dec)
	if err != nil {
		return nil, err
	}

	var (
		chunkPaths []string
		records    []chunkRecord
		sizeBytes  int64
		seq        int64
	)

	flush := func() error {
		if len(records) == 0 {
			return nil
		}

		path, err := writeChunkFile(tempDir, records)
		if err != nil {
			return err
		}

		chunkPaths = append(chunkPaths, path)
		records = nil
		sizeBytes = 0

		return nil
	}

	for dec.More() {
		record, recordSizeBytes, readErr := readChunkRecord(dec, seq)
		if readErr != nil {
			return nil, readErr
		}

		records = append(records, *record)
		sizeBytes += recordSizeBytes
		seq++

		err = flushRecordsIfNeeded(sizeBytes, chunkBytes, flush)
		if err != nil {
			return nil, err
		}
	}

	err = expectArrayEnd(dec)
	if err != nil {
		return nil, err
	}

	err = flush()
	if err != nil {
		return nil, err
	}

	return chunkPaths, nil
}

func expectArrayStart(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return fmt.Errorf("read opening token: %w", err)
	}

	return expectDelimiter(tok, '[', errExpectedTopArray)
}

func expectArrayEnd(dec *json.Decoder) error {
	tok, err := dec.Token()
	if err != nil {
		return fmt.Errorf("read closing token: %w", err)
	}

	return expectDelimiter(tok, ']', errExpectedEndArray)
}

func readChunkRecord(dec *json.Decoder, seq int64) (*chunkRecord, int64, error) {
	var raw json.RawMessage

	err := dec.Decode(&raw)
	if err != nil {
		return nil, 0, fmt.Errorf("decode raw record: %w", err)
	}

	var record bival.Record

	err = json.Unmarshal(raw, &record)
	if err != nil {
		return nil, 0, fmt.Errorf("decode record: %w", err)
	}

	name := recordName(&record)

	return &chunkRecord{
		Seq:  seq,
		Name: name,
		Raw:  raw,
	}, int64(len(raw) + len(name) + chunkRecordOverhead), nil
}

func flushRecordsIfNeeded(sizeBytes int64, chunkBytes int64, flush func() error) error {
	if sizeBytes < chunkBytes {
		return nil
	}

	return flush()
}

func writeChunkFile(tempDir string, records []chunkRecord) (string, error) {
	slices.SortStableFunc(records, compareChunkRecords)

	file, err := os.CreateTemp(tempDir, "chunk-*.jsonl")
	if err != nil {
		return "", fmt.Errorf("create chunk file: %w", err)
	}

	enc := json.NewEncoder(file)
	for _, record := range records {
		err := enc.Encode(record)
		if err != nil {
			_ = file.Close()

			return "", fmt.Errorf("write chunk file: %w", err)
		}
	}

	err = file.Close()
	if err != nil {
		return "", fmt.Errorf("close chunk file: %w", err)
	}

	return file.Name(), nil
}

func compareChunkRecords(leftRecord chunkRecord, rightRecord chunkRecord) int {
	if leftRecord.Name < rightRecord.Name {
		return -1
	}

	if leftRecord.Name > rightRecord.Name {
		return 1
	}

	if leftRecord.Seq < rightRecord.Seq {
		return -1
	}

	if leftRecord.Seq > rightRecord.Seq {
		return 1
	}

	return 0
}

type chunkReader struct {
	file *os.File
	dec  *json.Decoder
}

func openChunkReader(path string) (*chunkReader, error) {
	// #nosec G304,G703 -- chunk files are created internally under a temporary directory.
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, fmt.Errorf("open chunk file: %w", err)
	}

	return &chunkReader{
		file: file,
		dec:  json.NewDecoder(file),
	}, nil
}

func (r *chunkReader) Next() (*chunkRecord, error) {
	var record chunkRecord

	err := r.dec.Decode(&record)
	if err != nil {
		return nil, fmt.Errorf("decode chunk record: %w", err)
	}

	return &record, nil
}

func (r *chunkReader) Close() error {
	err := r.file.Close()
	if err != nil {
		return fmt.Errorf("close chunk reader: %w", err)
	}

	return nil
}

type heapItem struct {
	record chunkRecord
	reader *chunkReader
}

type mergeHeap []*heapItem

func (h *mergeHeap) Len() int {
	return len(*h)
}

func (h *mergeHeap) Less(i int, j int) bool {
	return compareChunkRecords((*h)[i].record, (*h)[j].record) < 0
}

func (h *mergeHeap) Swap(i int, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]
}

func (h *mergeHeap) Push(x any) {
	item, ok := x.(*heapItem)
	if !ok {
		panic(errInvalidHeapItem)
	}

	*h = append(*h, item)
}

func (h *mergeHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]

	return item
}

func mergeChunks(w io.Writer, chunkPaths []string) error {
	bufferedWriter := bufio.NewWriter(w)

	defer func() {
		_ = bufferedWriter.Flush()
	}()

	readers := make([]*chunkReader, 0, len(chunkPaths))

	defer func() {
		for _, reader := range readers {
			_ = reader.Close()
		}
	}()

	mergeItemHeap := make(mergeHeap, 0, len(chunkPaths))

	readers, err := seedMergeHeap(chunkPaths, &mergeItemHeap)
	if err != nil {
		return err
	}

	return writeMergedRecords(bufferedWriter, &mergeItemHeap)
}

func formatRecord(raw json.RawMessage) ([]byte, error) {
	var formatted bytes.Buffer

	err := json.Indent(&formatted, raw, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("indent JSON: %w", err)
	}

	lines := strings.Split(formatted.String(), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}

		lines[i] = "    " + line
	}

	return []byte(strings.Join(lines, "\n")), nil
}

func expectDelimiter(tok json.Token, expected json.Delim, baseErr error) error {
	delim, ok := tok.(json.Delim)
	if !ok || delim != expected {
		return fmt.Errorf("%w: got %v", baseErr, tok)
	}

	return nil
}

func seedMergeHeap(
	chunkPaths []string,
	mergeItemHeap *mergeHeap,
) ([]*chunkReader, error) {
	readers := make([]*chunkReader, 0, len(chunkPaths))

	for _, path := range chunkPaths {
		reader, err := openChunkReader(path)
		if err != nil {
			return readers, err
		}

		readers = append(readers, reader)

		record, readErr := reader.Next()
		if errors.Is(readErr, io.EOF) {
			continue
		}

		if readErr != nil {
			return readers, fmt.Errorf("read chunk file: %w", readErr)
		}

		heap.Push(mergeItemHeap, &heapItem{
			record: *record,
			reader: reader,
		})
	}

	return readers, nil
}

func writeMergedRecords(bufferedWriter *bufio.Writer, mergeItemHeap *mergeHeap) error {
	if mergeItemHeap.Len() == 0 {
		return writeEmptyArray(bufferedWriter)
	}

	err := writeOpeningArray(bufferedWriter)
	if err != nil {
		return err
	}

	firstRecord := true

	for mergeItemHeap.Len() > 0 {
		firstRecord, err = writeNextMergedRecord(bufferedWriter, mergeItemHeap, firstRecord)
		if err != nil {
			return err
		}
	}

	return finishMergedArray(bufferedWriter, firstRecord)
}

func writeEmptyArray(bufferedWriter *bufio.Writer) error {
	_, err := io.WriteString(bufferedWriter, "[]\n")
	if err != nil {
		return fmt.Errorf("write empty array: %w", err)
	}

	err = bufferedWriter.Flush()
	if err != nil {
		return fmt.Errorf("flush output: %w", err)
	}

	return nil
}

func writeOpeningArray(bufferedWriter *bufio.Writer) error {
	_, err := io.WriteString(bufferedWriter, "[\n")
	if err != nil {
		return fmt.Errorf("write opening array: %w", err)
	}

	return nil
}

func writeNextMergedRecord(
	bufferedWriter *bufio.Writer,
	mergeItemHeap *mergeHeap,
	firstRecord bool,
) (bool, error) {
	item, err := popHeapItem(mergeItemHeap)
	if err != nil {
		return firstRecord, err
	}

	if !firstRecord {
		_, err = io.WriteString(bufferedWriter, ",\n")
		if err != nil {
			return false, fmt.Errorf("write separator: %w", err)
		}
	}

	err = writeFormattedRecord(bufferedWriter, item.record.Raw)
	if err != nil {
		return false, err
	}

	err = pushNextRecord(mergeItemHeap, item)
	if err != nil {
		return false, err
	}

	return false, nil
}

func writeFormattedRecord(bufferedWriter *bufio.Writer, raw json.RawMessage) error {
	formatted, err := formatRecord(raw)
	if err != nil {
		return fmt.Errorf("format record: %w", err)
	}

	_, err = bufferedWriter.Write(formatted)
	if err != nil {
		return fmt.Errorf("write record: %w", err)
	}

	return nil
}

func pushNextRecord(mergeItemHeap *mergeHeap, item *heapItem) error {
	nextRecord, err := item.reader.Next()
	if errors.Is(err, io.EOF) {
		return nil
	}

	if err != nil {
		return fmt.Errorf("read chunk file: %w", err)
	}

	item.record = *nextRecord
	heap.Push(mergeItemHeap, item)

	return nil
}

func finishMergedArray(bufferedWriter *bufio.Writer, firstRecord bool) error {
	if !firstRecord {
		_, err := io.WriteString(bufferedWriter, "\n")
		if err != nil {
			return fmt.Errorf("write trailing newline: %w", err)
		}
	}

	_, err := io.WriteString(bufferedWriter, "]\n")
	if err != nil {
		return fmt.Errorf("write closing array: %w", err)
	}

	err = bufferedWriter.Flush()
	if err != nil {
		return fmt.Errorf("flush output: %w", err)
	}

	return nil
}

func popHeapItem(mergeItemHeap *mergeHeap) (*heapItem, error) {
	popped := heap.Pop(mergeItemHeap)

	item, ok := popped.(*heapItem)
	if !ok {
		return nil, errInvalidHeapItem
	}

	return item, nil
}
