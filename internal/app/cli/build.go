package cli

import (
	"errors"
	"fmt"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errUnsupportedRecordType = errors.New("unsupported record type")

const recordTypeInstance = "instance"
const recordTypePlain = "plain"
const recordTypeOLH = "olh"

func buildEntry(record *bilist.Record) (any, error) {
	switch record.Type {
	case recordTypeInstance:
		params := newDirEntryParams(record)

		return domain.NewInstance(params), nil
	case recordTypePlain:
		params := newDirEntryParams(record)

		return domain.NewPlain(params), nil
	case recordTypeOLH:
		return domain.NewOLH(domain.OLHParams{
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
	var maps []*domain.PendingMap
	for _, pendingMap := range record.Entry.PendingMap {
		maps = append(
			maps,
			domain.NewPendingMap(
				pendingMap.Key,
				pendingMap.Val.State,
				pendingMap.Val.Timestamp,
				pendingMap.Val.Op,
			),
		)
	}

	return maps
}

func newPendingLogs(record *bilist.Record) []*domain.PendingLog {
	var logs []*domain.PendingLog

	for _, pendingLog := range record.Entry.PendingLog {
		var vals []*domain.PendingLogVal
		for _, pendingVal := range pendingLog.Val {
			vals = append(
				vals,
				domain.NewPendingLogVal(
					pendingVal.Epoch,
					pendingVal.Op,
					pendingVal.OpTag,
					pendingVal.Key.Name,
					pendingVal.Key.Instance,
					pendingVal.DeleteMarker,
				),
			)
		}

		logs = append(logs, domain.NewPendingLog(pendingLog.Key, vals))
	}

	return logs
}
