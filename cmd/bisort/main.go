package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/neatflowcv/bival"
)

const defaultChunkBytes int64 = 64 << 20

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func run(args []string) error {
	fs := flag.NewFlagSet("bisort", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	chunkBytes := fs.Int64("chunk-bytes", defaultChunkBytes, "maximum in-memory chunk size in bytes")

	reorderedArgs := normalizeArgs(args)

	if err := fs.Parse(reorderedArgs); err != nil {
		return fmt.Errorf("parse flags: %w", err)
	}

	if *chunkBytes <= 0 {
		return errors.New("chunk-bytes must be greater than zero")
	}

	if fs.NArg() != 2 {
		return errors.New("usage: bisort <input> <output> --chunk-bytes N")
	}

	return sortFile(fs.Arg(0), fs.Arg(1), *chunkBytes)
}

func normalizeArgs(args []string) []string {
	flags := make([]string, 0, len(args))
	positionals := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "--" {
			positionals = append(positionals, args[i+1:]...)
			break
		}

		if len(arg) > 0 && arg[0] == '-' {
			flags = append(flags, arg)
			if !hasInlineValue(arg) && flagNeedsValue(arg) && i+1 < len(args) {
				i++
				flags = append(flags, args[i])
			}
			continue
		}

		positionals = append(positionals, arg)
	}

	return append(flags, positionals...)
}

func hasInlineValue(arg string) bool {
	for i := 0; i < len(arg); i++ {
		if arg[i] == '=' {
			return true
		}
	}

	return false
}

func flagNeedsValue(arg string) bool {
	return arg == "-chunk-bytes" || arg == "--chunk-bytes"
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

	tok, err := dec.Token()
	if err != nil {
		return nil, fmt.Errorf("read opening token: %w", err)
	}

	delim, ok := tok.(json.Delim)
	if !ok || delim != '[' {
		return nil, fmt.Errorf("expected top-level array: got %v", tok)
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
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			return nil, fmt.Errorf("decode raw record: %w", err)
		}

		var record bival.Record
		if err := json.Unmarshal(raw, &record); err != nil {
			return nil, fmt.Errorf("decode record: %w", err)
		}

		name := recordName(&record)
		records = append(records, chunkRecord{
			Seq:  seq,
			Name: name,
			Raw:  raw,
		})
		sizeBytes += int64(len(raw) + len(name) + 16)
		seq++

		if sizeBytes >= chunkBytes {
			if err := flush(); err != nil {
				return nil, err
			}
		}
	}

	tok, err = dec.Token()
	if err != nil {
		return nil, fmt.Errorf("read closing token: %w", err)
	}

	delim, ok = tok.(json.Delim)
	if !ok || delim != ']' {
		return nil, fmt.Errorf("expected closing array: got %v", tok)
	}

	if err := flush(); err != nil {
		return nil, err
	}

	return chunkPaths, nil
}

func writeChunkFile(tempDir string, records []chunkRecord) (string, error) {
	slices.SortStableFunc(records, compareChunkRecords)

	file, err := os.CreateTemp(tempDir, "chunk-*.jsonl")
	if err != nil {
		return "", fmt.Errorf("create chunk file: %w", err)
	}

	enc := json.NewEncoder(file)
	for _, record := range records {
		if err := enc.Encode(record); err != nil {
			_ = file.Close()
			return "", fmt.Errorf("write chunk file: %w", err)
		}
	}

	if err := file.Close(); err != nil {
		return "", fmt.Errorf("close chunk file: %w", err)
	}

	return file.Name(), nil
}

func compareChunkRecords(a chunkRecord, b chunkRecord) int {
	if a.Name < b.Name {
		return -1
	}
	if a.Name > b.Name {
		return 1
	}
	if a.Seq < b.Seq {
		return -1
	}
	if a.Seq > b.Seq {
		return 1
	}

	return 0
}

func recordName(record *bival.Record) string {
	if record.Type == "olh" {
		return record.Entry.Key.Name
	}

	return record.Entry.Name
}

type chunkReader struct {
	file *os.File
	dec  *json.Decoder
}

func openChunkReader(path string) (*chunkReader, error) {
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
	if err := r.dec.Decode(&record); err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *chunkReader) Close() error {
	return r.file.Close()
}

type heapItem struct {
	record chunkRecord
	reader *chunkReader
}

type mergeHeap []*heapItem

func (h mergeHeap) Len() int {
	return len(h)
}

func (h mergeHeap) Less(i int, j int) bool {
	return compareChunkRecords(h[i].record, h[j].record) < 0
}

func (h mergeHeap) Swap(i int, j int) {
	h[i], h[j] = h[j], h[i]
}

func (h *mergeHeap) Push(x any) {
	*h = append(*h, x.(*heapItem))
}

func (h *mergeHeap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func mergeChunks(w io.Writer, chunkPaths []string) error {
	bw := bufio.NewWriter(w)
	defer func() {
		_ = bw.Flush()
	}()

	readers := make([]*chunkReader, 0, len(chunkPaths))
	defer func() {
		for _, reader := range readers {
			_ = reader.Close()
		}
	}()

	h := make(mergeHeap, 0, len(chunkPaths))
	for _, path := range chunkPaths {
		reader, err := openChunkReader(path)
		if err != nil {
			return err
		}
		readers = append(readers, reader)

		record, err := reader.Next()
		if errors.Is(err, io.EOF) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read chunk file: %w", err)
		}

		heap.Push(&h, &heapItem{
			record: *record,
			reader: reader,
		})
	}

	if h.Len() == 0 {
		if _, err := io.WriteString(bw, "[]\n"); err != nil {
			return fmt.Errorf("write empty array: %w", err)
		}

		if err := bw.Flush(); err != nil {
			return fmt.Errorf("flush output: %w", err)
		}

		return nil
	}

	_, err := io.WriteString(bw, "[\n")
	if err != nil {
		return fmt.Errorf("write opening array: %w", err)
	}

	first := true
	for h.Len() > 0 {
		item := heap.Pop(&h).(*heapItem)

		if !first {
			if _, err := io.WriteString(bw, ",\n"); err != nil {
				return fmt.Errorf("write separator: %w", err)
			}
		}
		first = false

		formatted, err := formatRecord(item.record.Raw)
		if err != nil {
			return fmt.Errorf("format record: %w", err)
		}

		if _, err := bw.Write(formatted); err != nil {
			return fmt.Errorf("write record: %w", err)
		}

		next, err := item.reader.Next()
		if errors.Is(err, io.EOF) {
			continue
		}
		if err != nil {
			return fmt.Errorf("read chunk file: %w", err)
		}

		item.record = *next
		heap.Push(&h, item)
	}

	if !first {
		if _, err := io.WriteString(bw, "\n"); err != nil {
			return fmt.Errorf("write trailing newline: %w", err)
		}
	}

	if _, err := io.WriteString(bw, "]\n"); err != nil {
		return fmt.Errorf("write closing array: %w", err)
	}

	if err := bw.Flush(); err != nil {
		return fmt.Errorf("flush output: %w", err)
	}

	return nil
}

func formatRecord(raw json.RawMessage) ([]byte, error) {
	var formatted bytes.Buffer
	if err := json.Indent(&formatted, raw, "    ", "    "); err != nil {
		return nil, err
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
