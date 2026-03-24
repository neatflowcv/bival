package bilist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var (
	errExpectedTopLevelArray = errors.New("expected top-level array")
	errExpectedClosingArray  = errors.New("expected closing array")
)

// Record is one item from the top-level JSON array.
type Record struct {
	Type  string `json:"type"`
	Idx   string `json:"idx"`
	Entry Entry  `json:"entry"`
}

// Entry is a superset of the shapes observed in sample data.
// Fields that do not exist for a given record type remain zero-valued.
type Entry struct {
	Name           string            `json:"name"`
	Instance       string            `json:"instance"`
	Ver            Version           `json:"ver"`
	Locator        string            `json:"locator"`
	Exists         bool              `json:"exists"`
	Meta           Meta              `json:"meta"`
	Tag            string            `json:"tag"`
	Flags          int               `json:"flags"`
	PendingMap     []json.RawMessage `json:"pending_map"`
	VersionedEpoch int               `json:"versioned_epoch"`

	Key            Key               `json:"key"`
	DeleteMarker   bool              `json:"delete_marker"`
	Epoch          int               `json:"epoch"`
	PendingLog     []json.RawMessage `json:"pending_log"`
	PendingRemoval bool              `json:"pending_removal"`
}

type Version struct {
	Pool  int `json:"pool"`
	Epoch int `json:"epoch"`
}

type Meta struct {
	Category         int    `json:"category"`
	Size             int64  `json:"size"`
	MTime            string `json:"mtime"`
	ETag             string `json:"etag"`
	StorageClass     string `json:"storage_class"`
	Owner            string `json:"owner"`
	OwnerDisplayName string `json:"owner_display_name"`
	ContentType      string `json:"content_type"`
	AccountedSize    int64  `json:"accounted_size"`
	UserData         string `json:"user_data"`
	Appendable       bool   `json:"appendable"`
}

type Key struct {
	Name     string `json:"name"`
	Instance string `json:"instance"`
}

// Reader streams bilist records from a JSON array source.
type Reader struct {
	dec *json.Decoder
}

// NewReader returns a Reader that decodes bilist records from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{dec: json.NewDecoder(r)}
}

// ReadFile opens a JSON file whose top-level value is an array and decodes
// each element one at a time. It keeps memory usage bounded because it never
// builds the full slice in memory.
func ReadFile(path string, visit func(*Record) error) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}

	defer func() {
		_ = file.Close()
	}()

	return Read(file, visit)
}

// Read streams records from r and calls visit for each decoded item.
func Read(r io.Reader, visit func(*Record) error) error {
	return NewReader(r).Read(visit)
}

// Read decodes each top-level array item and passes it to visit.
func (r *Reader) Read(visit func(*Record) error) error {
	tok, err := r.dec.Token()
	if err != nil {
		return fmt.Errorf("read opening token: %w", err)
	}

	delim, isDelim := tok.(json.Delim)
	if !isDelim || delim != '[' {
		return fmt.Errorf("%w: got %v", errExpectedTopLevelArray, tok)
	}

	for r.dec.More() {
		var record Record

		err := r.dec.Decode(&record)
		if err != nil {
			return fmt.Errorf("decode record: %w", err)
		}

		err = visit(&record)
		if err != nil {
			return err
		}
	}

	tok, err = r.dec.Token()
	if err != nil {
		return fmt.Errorf("read closing token: %w", err)
	}

	delim, isDelim = tok.(json.Delim)
	if !isDelim || delim != ']' {
		return fmt.Errorf("%w: got %v", errExpectedClosingArray, tok)
	}

	return nil
}
