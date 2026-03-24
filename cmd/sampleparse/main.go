package main

import (
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

	counts := map[string]int{}

	err := bival.ParseFile(path, func(record *bival.Record) error {
		if describeBuiltEntry(record) {
			counts[record.Type]++
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	for recordType, count := range counts {
		log.Printf("type=%s count=%d", recordType, count)
	}
}

func describeBuiltEntry(record *bival.Record) bool {
	switch record.Type {
	case "instance":
		entry := domain.NewInstanceEntry(newDirEntry(record))
		log.Printf("type=%s entry=%+v", record.Type, entry)
		return true
	case "plain":
		entry := domain.NewPlainEntry(newDirEntry(record))
		log.Printf("type=%s entry=%+v", record.Type, entry)
		return true
	case "olh":
		entry := domain.NewOLHEntry(
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
		log.Printf("type=%s entry=%+v", record.Type, entry)
		return true
	default:
		log.Printf("skip unsupported type=%s idx=%s", record.Type, record.Idx)
		return false
	}
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
