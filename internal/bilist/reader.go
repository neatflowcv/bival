package bilist

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
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
	dec              *json.Decoder
	started          bool
	finished         bool
	closingTokenRead bool
}

// NewReader returns a Reader that decodes bilist records from r.
func NewReader(r io.Reader) *Reader {
	return &Reader{dec: json.NewDecoder(r)}
}

// Read decodes one top-level array item at a time.
func (r *Reader) Read() (*Record, error) {
	if r.finished {
		return nil, io.EOF
	}

	if !r.started {
		tok, err := r.dec.Token()
		if err != nil {
			return nil, fmt.Errorf("read opening token: %w", err)
		}

		delim, isDelim := tok.(json.Delim)
		if !isDelim || delim != '[' {
			return nil, fmt.Errorf("%w: got %v", errExpectedTopLevelArray, tok)
		}

		r.started = true
	}

	if !r.dec.More() {
		err := r.readClosingToken()
		if err != nil {
			return nil, err
		}

		r.finished = true

		return nil, io.EOF
	}

	var record Record

	err := r.dec.Decode(&record)
	if err != nil {
		return nil, fmt.Errorf("decode record: %w", err)
	}

	return &record, nil
}

func (r *Reader) readClosingToken() error {
	if r.closingTokenRead {
		return nil
	}

	tok, err := r.dec.Token()
	if err != nil {
		return fmt.Errorf("read closing token: %w", err)
	}

	delim, isDelim := tok.(json.Delim)
	if !isDelim || delim != ']' {
		return fmt.Errorf("%w: got %v", errExpectedClosingArray, tok)
	}

	r.closingTokenRead = true

	return nil
}
