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
		return domain.NewInstanceEntry(newDirEntry(record)), nil
	case "plain":
		return domain.NewPlainEntry(newDirEntry(record)), nil
	case "olh":
		return domain.NewOLHEntry(
			record.Type,
			[]byte(record.Idx),
			domain.NewOLHPayload(
				domain.NewKey(record.Entry.Key.Name, record.Entry.Key.Instance),
				domain.NewOLHState(record.Entry.DeleteMarker, record.Entry.PendingRemoval, record.Entry.Exists),
				record.Entry.Epoch,
				nil,
				record.Entry.Tag,
			),
		), nil
	default:
		return nil, fmt.Errorf("%w %q", errUnsupportedRecordType, record.Type)
	}
}

func newDirEntry(record *bilist.Record) *domain.DirEntry {
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
				domain.NewAuditInfo(record.Entry.Meta.MTime, record.Entry.Meta.ETag),
				domain.NewContentInfo(record.Entry.Meta.StorageClass, record.Entry.Meta.ContentType),
				domain.NewOwner(record.Entry.Meta.Owner, record.Entry.Meta.OwnerDisplayName),
			),
			nil,
		),
	)
}
