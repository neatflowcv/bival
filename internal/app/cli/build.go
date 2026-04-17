package cli

import (
	"errors"
	"fmt"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errUnsupportedRecordType = errors.New("unsupported record type")

func buildEntry(record *bilist.Record) (any, error) {
	switch record.Type {
	case "instance":
		params := newDirEntryParams(record)

		return domain.NewInstance(params), nil
	case "plain":
		params := newDirEntryParams(record)

		return domain.NewPlain(params), nil
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

func newDirEntryParams(record *bilist.Record) domain.DirEntryParams {
	return domain.DirEntryParams{
		Kind:             record.Type,
		Index:            []byte(record.Idx),
		Name:             record.Entry.Name,
		Instance:         record.Entry.Instance,
		Pool:             record.Entry.Ver.Pool,
		Epoch:            record.Entry.Ver.Epoch,
		VEpoch:           record.Entry.VersionedEpoch,
		Locator:          record.Entry.Locator,
		Exists:           record.Entry.Exists,
		Tag:              record.Entry.Tag,
		Flags:            record.Entry.Flags,
		Category:         record.Entry.Meta.Category,
		Size:             record.Entry.Meta.Size,
		AccountedSize:    record.Entry.Meta.AccountedSize,
		Appendable:       record.Entry.Meta.Appendable,
		MTime:            record.Entry.Meta.MTime,
		ETag:             record.Entry.Meta.ETag,
		StorageClass:     record.Entry.Meta.StorageClass,
		ContentType:      record.Entry.Meta.ContentType,
		OwnerUserID:      record.Entry.Meta.Owner,
		OwnerDisplayName: record.Entry.Meta.OwnerDisplayName,
		PendingMaps:      newPendingMaps(record),
	}
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
