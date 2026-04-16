package cli

import (
	"errors"
	"fmt"
	"time"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errUnsupportedRecordType = errors.New("unsupported record type")

func buildEntry(record *bilist.Record) (any, error) {
	switch record.Type {
	case "instance":
		entry, err := newDirEntry(record)
		if err != nil {
			return nil, err
		}

		return domain.NewInstanceEntry(entry), nil
	case "plain":
		entry, err := newDirEntry(record)
		if err != nil {
			return nil, err
		}

		return domain.NewPlainEntry(entry), nil
	case "olh":
		return domain.NewOLHEntry(domain.OLHEntryParams{
			Kind:  record.Type,
			Index: []byte(record.Idx),
			Payload: domain.NewOLHPayload(
				domain.NewKey(record.Entry.Key.Name, record.Entry.Key.Instance),
				domain.NewOLHState(record.Entry.DeleteMarker, record.Entry.PendingRemoval, record.Entry.Exists),
				record.Entry.Epoch,
				newPendingLogs(record),
				record.Entry.Tag,
			),
		}), nil
	default:
		return nil, fmt.Errorf("%w %q", errUnsupportedRecordType, record.Type)
	}
}

func newDirEntry(record *bilist.Record) (*domain.DirEntry, error) {
	mTime, err := parseMTime(record.Entry.Meta.MTime)
	if err != nil {
		return nil, fmt.Errorf("parse mtime %q: %w", record.Entry.Meta.MTime, err)
	}

	return domain.NewDirEntry(
		record.Type,
		[]byte(record.Idx),
		domain.NewDirPayload(
			domain.NewKey(record.Entry.Name, record.Entry.Instance),
			domain.NewDirVersionInfo(
				domain.NewVersion(record.Entry.Ver.Pool, record.Entry.Ver.Epoch),
				record.Entry.VersionedEpoch,
			),
			domain.NewDirState(
				record.Entry.Locator,
				record.Entry.Exists,
				record.Entry.Tag,
				record.Entry.Flags,
			),
			domain.NewMeta(
				domain.NewObjectSpec(
					record.Entry.Meta.Category,
					record.Entry.Meta.Size,
					record.Entry.Meta.AccountedSize,
					record.Entry.Meta.Appendable,
				),
				domain.NewAuditInfo(mTime, record.Entry.Meta.ETag),
				domain.NewContentInfo(record.Entry.Meta.StorageClass, record.Entry.Meta.ContentType),
				domain.NewOwner(record.Entry.Meta.Owner, record.Entry.Meta.OwnerDisplayName),
			),
			newPendingMaps(record),
		),
	), nil
}

func parseMTime(value string) (time.Time, error) {
	if value == "" || value == "0.000000" {
		return time.Time{}, nil
	}

	parsed, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse RFC3339Nano: %w", err)
	}

	return parsed, nil
}

func newPendingMaps(record *bilist.Record) []*domain.PendingMap {
	if len(record.Entry.PendingMap) == 0 {
		return nil
	}

	return make([]*domain.PendingMap, len(record.Entry.PendingMap))
}

func newPendingLogs(record *bilist.Record) []*domain.PendingLog {
	if len(record.Entry.PendingLog) == 0 {
		return nil
	}

	return make([]*domain.PendingLog, len(record.Entry.PendingLog))
}
