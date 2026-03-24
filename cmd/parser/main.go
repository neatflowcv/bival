package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/neatflowcv/bival"
	"github.com/neatflowcv/bival/internal/pkg/domain"
)

func main() {
	path := "sample.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	err := bival.ParseFile(path, func(record *bival.Record) error {
		entry, err := newTypedEntry(record)
		if err != nil {
			return err
		}

		err = entry.Validate()
		if err != nil {
			return fmt.Errorf("validate record type %q: %w", record.Type, err)
		}

		fmt.Printf("%+v\n", entry)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

func newTypedEntry(record *bival.Record) (domain.MyEntry, error) {
	switch domain.Kind(record.Type) {
	case domain.KindInstance:
		return newInstance(record)
	case domain.KindPlain:
		return newPlain(record)
	case domain.KindOLH:
		return newOLH(record), nil
	default:
		return nil, fmt.Errorf("unsupported record type: %q", record.Type)
	}
}

func newInstance(record *bival.Record) (*domain.Instance, error) {
	mtime, err := parseMTime(record.Entry.Meta.MTime)
	if err != nil {
		return nil, fmt.Errorf("parse instance mtime: %w", err)
	}

	return &domain.Instance{
		Idx:              []byte(record.Idx),
		Name:             record.Entry.Name,
		Instance:         record.Entry.Instance,
		Pool:             record.Entry.Ver.Pool,
		Epoch:            record.Entry.Ver.Epoch,
		Locator:          record.Entry.Locator,
		Exists:           record.Entry.Exists,
		Category:         record.Entry.Meta.Category,
		Size:             int(record.Entry.Meta.Size),
		Mtime:            mtime,
		Etag:             record.Entry.Meta.ETag,
		StorageClass:     record.Entry.Meta.StorageClass,
		Owner:            record.Entry.Meta.Owner,
		OwnerDisplayName: record.Entry.Meta.OwnerDisplayName,
		ContentType:      record.Entry.Meta.ContentType,
		AccountedSize:    int(record.Entry.Meta.AccountedSize),
		UserData:         record.Entry.Meta.UserData,
		Appendable:       record.Entry.Meta.Appendable,
		Tag:              record.Entry.Tag,
		Flags:            record.Entry.Flags,
		PendingMap:       decodeRawMessages(record.Entry.PendingMap),
		VersionedEpoch:   record.Entry.VersionedEpoch,
	}, nil
}

func newPlain(record *bival.Record) (*domain.Plain, error) {
	mtime, err := parseMTime(record.Entry.Meta.MTime)
	if err != nil {
		return nil, fmt.Errorf("parse plain mtime: %w", err)
	}

	return &domain.Plain{
		Idx:              []byte(record.Idx),
		Name:             record.Entry.Name,
		Instance:         record.Entry.Instance,
		Pool:             record.Entry.Ver.Pool,
		Epoch:            record.Entry.Ver.Epoch,
		Locator:          record.Entry.Locator,
		Exists:           record.Entry.Exists,
		Category:         record.Entry.Meta.Category,
		Size:             int(record.Entry.Meta.Size),
		Mtime:            mtime,
		Etag:             record.Entry.Meta.ETag,
		StorageClass:     record.Entry.Meta.StorageClass,
		Owner:            record.Entry.Meta.Owner,
		OwnerDisplayName: record.Entry.Meta.OwnerDisplayName,
		ContentType:      record.Entry.Meta.ContentType,
		AccountedSize:    int(record.Entry.Meta.AccountedSize),
		UserData:         record.Entry.Meta.UserData,
		Appendable:       record.Entry.Meta.Appendable,
		Tag:              record.Entry.Tag,
		Flags:            record.Entry.Flags,
		PendingMap:       decodeRawMessages(record.Entry.PendingMap),
		VersionedEpoch:   record.Entry.VersionedEpoch,
	}, nil
}

func newOLH(record *bival.Record) *domain.OLH {
	return &domain.OLH{
		Idx:            []byte(record.Idx),
		Name:           record.Entry.Key.Name,
		Instance:       record.Entry.Key.Instance,
		DeleteMarker:   record.Entry.DeleteMarker,
		Epoch:          record.Entry.Epoch,
		PendingLog:     decodeRawMessages(record.Entry.PendingLog),
		Tag:            record.Entry.Tag,
		Exists:         record.Entry.Exists,
		PendingRemoval: record.Entry.PendingRemoval,
	}
}

func decodeRawMessages(values []json.RawMessage) []any {
	if len(values) == 0 {
		return nil
	}

	decoded := make([]any, 0, len(values))
	for _, value := range values {
		var item any

		err := json.Unmarshal(value, &item)
		if err != nil {
			decoded = append(decoded, string(value))
			continue
		}

		decoded = append(decoded, item)
	}

	return decoded
}

func parseMTime(value string) (time.Time, error) {
	if value == "" || value == "0.000000" {
		return time.Time{}, nil
	}

	return time.Parse(time.RFC3339Nano, value)
}
