package main

import (
	"errors"
	"fmt"
	"log"
	"os"

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
	switch record.Type {
	case "instance":
		log.Printf("%+v", newInstance(record))

		return nil
	case "plain":
		log.Printf("%+v", newPlain(record))

		return nil
	case "olh":
		log.Printf("%+v", newOLH(record))

		return nil
	default:
		return fmt.Errorf("%w: %q", errUnsupportedRecordType, record.Type)
	}
}

func newInstance(record *bival.Record) *domain.InstanceEntry {
	return domain.NewInstanceEntry(newDirEntry(record))
}

func newPlain(record *bival.Record) *domain.PlainEntry {
	return domain.NewPlainEntry(newDirEntry(record))
}

func newOLH(record *bival.Record) *domain.OLHEntry {
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
	)
}

func newDirEntry(record *bival.Record) *domain.DirEntry {
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
