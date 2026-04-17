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
		params, err := newDirEntryParams(record)
		if err != nil {
			return nil, err
		}

		return domain.NewInstanceEntry(params), nil
	case "plain":
		params, err := newDirEntryParams(record)
		if err != nil {
			return nil, err
		}

		return domain.NewPlainEntry(domain.NewDirEntry(params)), nil
	case "olh":
		return domain.NewOLH(domain.OLHParams{
			Kind:           record.Type,
			Index:          []byte(record.Idx),
			Name:           record.Entry.Key.Name,
			Instance:       record.Entry.Key.Instance,
			DeleteMarker:   record.Entry.DeleteMarker,
			PendingRemoval: record.Entry.PendingRemoval,
			Exists:         record.Entry.Exists,
			Epoch:          record.Entry.Epoch,
			PendingLogs:    newPendingLogs(record),
			Tag:            record.Entry.Tag,
		}), nil
	default:
		return nil, fmt.Errorf("%w %q", errUnsupportedRecordType, record.Type)
	}
}

func newDirEntryParams(record *bilist.Record) (domain.DirEntryParams, error) {
	mTime, err := parseMTime(record.Entry.Meta.MTime)
	if err != nil {
		return domain.DirEntryParams{}, fmt.Errorf("parse mtime %q: %w", record.Entry.Meta.MTime, err)
	}

	return domain.DirEntryParams{
		Kind:  record.Type,
		Index: []byte(record.Idx),
		Key:   domain.NewKey(record.Entry.Name, record.Entry.Instance),
		VersionInfo: domain.NewDirVersionInfo(
			domain.NewVersion(record.Entry.Ver.Pool, record.Entry.Ver.Epoch),
			record.Entry.VersionedEpoch,
		),
		State: domain.NewDirState(
			record.Entry.Locator,
			record.Entry.Exists,
			record.Entry.Tag,
			record.Entry.Flags,
		),
		Meta: domain.NewMeta(
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
		PendingMaps: newPendingMaps(record),
	}, nil
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
