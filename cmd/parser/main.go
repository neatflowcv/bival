package main

import (
	"encoding/json"
	"errors"
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
		err := validateTypedEntry(record)
		if err != nil {
			return fmt.Errorf("validate record type %q: %w", record.Type, err)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}
}

var errUnsupportedRecordType = errors.New("unsupported record type")

func validateTypedEntry(record *bival.Record) error {
	switch domain.Kind(record.Type) {
	case domain.KindInstance:
		entry, err := newInstance(record)
		if err != nil {
			return err
		}

		log.Printf("%+v", entry)

		err = entry.Validate()
		if err != nil {
			return fmt.Errorf("validate instance: %w", err)
		}

		return nil
	case domain.KindPlain:
		entry, err := newPlain(record)
		if err != nil {
			return err
		}

		log.Printf("%+v", entry)

		err = entry.Validate()
		if err != nil {
			return fmt.Errorf("validate plain: %w", err)
		}

		return nil
	case domain.KindOLH:
		entry := newOLH(record)
		log.Printf("%+v", entry)

		err := entry.Validate()
		if err != nil {
			return fmt.Errorf("validate olh: %w", err)
		}

		return nil
	default:
		return fmt.Errorf("%w: %q", errUnsupportedRecordType, record.Type)
	}
}

//nolint:dupl // Plain and Instance share field layout but produce distinct domain types.
func newInstance(record *bival.Record) (*domain.Instance, error) {
	common, err := newEntryCommon(record, "instance")
	if err != nil {
		return nil, err
	}

	return &domain.Instance{
		Idx:              common.idx,
		Name:             common.name,
		Instance:         common.instance,
		Pool:             common.pool,
		Epoch:            common.epoch,
		Locator:          common.locator,
		Exists:           common.exists,
		Category:         common.category,
		Size:             common.size,
		Mtime:            common.mtime,
		Etag:             common.etag,
		StorageClass:     common.storageClass,
		Owner:            common.owner,
		OwnerDisplayName: common.ownerDisplayName,
		ContentType:      common.contentType,
		AccountedSize:    common.accountedSize,
		UserData:         common.userData,
		Appendable:       common.appendable,
		Tag:              common.tag,
		Flags:            common.flags,
		PendingMap:       common.pendingMap,
		VersionedEpoch:   common.versionedEpoch,
	}, nil
}

//nolint:dupl // Plain and Instance share field layout but produce distinct domain types.
func newPlain(record *bival.Record) (*domain.Plain, error) {
	common, err := newEntryCommon(record, "plain")
	if err != nil {
		return nil, err
	}

	return &domain.Plain{
		Idx:              common.idx,
		Name:             common.name,
		Instance:         common.instance,
		Pool:             common.pool,
		Epoch:            common.epoch,
		Locator:          common.locator,
		Exists:           common.exists,
		Category:         common.category,
		Size:             common.size,
		Mtime:            common.mtime,
		Etag:             common.etag,
		StorageClass:     common.storageClass,
		Owner:            common.owner,
		OwnerDisplayName: common.ownerDisplayName,
		ContentType:      common.contentType,
		AccountedSize:    common.accountedSize,
		UserData:         common.userData,
		Appendable:       common.appendable,
		Tag:              common.tag,
		Flags:            common.flags,
		PendingMap:       common.pendingMap,
		VersionedEpoch:   common.versionedEpoch,
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

type entryCommon struct {
	idx              []byte
	name             string
	instance         string
	pool             int
	epoch            int
	locator          string
	exists           bool
	category         int
	size             int
	mtime            time.Time
	etag             string
	storageClass     string
	owner            string
	ownerDisplayName string
	contentType      string
	accountedSize    int
	userData         string
	appendable       bool
	tag              string
	flags            int
	pendingMap       []any
	versionedEpoch   int
}

func newEntryCommon(record *bival.Record, kind string) (*entryCommon, error) {
	mtime, err := parseMTime(record.Entry.Meta.MTime)
	if err != nil {
		return nil, fmt.Errorf("parse %s mtime: %w", kind, err)
	}

	return &entryCommon{
		idx:              []byte(record.Idx),
		name:             record.Entry.Name,
		instance:         record.Entry.Instance,
		pool:             record.Entry.Ver.Pool,
		epoch:            record.Entry.Ver.Epoch,
		locator:          record.Entry.Locator,
		exists:           record.Entry.Exists,
		category:         record.Entry.Meta.Category,
		size:             int(record.Entry.Meta.Size),
		mtime:            mtime,
		etag:             record.Entry.Meta.ETag,
		storageClass:     record.Entry.Meta.StorageClass,
		owner:            record.Entry.Meta.Owner,
		ownerDisplayName: record.Entry.Meta.OwnerDisplayName,
		contentType:      record.Entry.Meta.ContentType,
		accountedSize:    int(record.Entry.Meta.AccountedSize),
		userData:         record.Entry.Meta.UserData,
		appendable:       record.Entry.Meta.Appendable,
		tag:              record.Entry.Tag,
		flags:            record.Entry.Flags,
		pendingMap:       decodeRawMessages(record.Entry.PendingMap),
		versionedEpoch:   record.Entry.VersionedEpoch,
	}, nil
}

func parseMTime(value string) (time.Time, error) {
	if value == "" || value == "0.000000" {
		return time.Time{}, nil
	}

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse RFC3339Nano time: %w", err)
	}

	return parsed, nil
}
