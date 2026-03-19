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

	registry := domain.NewEntryRegistry()

	err := bival.ParseFile(path, func(record *bival.Record) error {
		name := record.Entry.Name

		instance := record.Entry.Instance
		if record.Type == "olh" {
			name = record.Entry.Key.Name
			instance = record.Entry.Key.Instance
		}

		if len(record.Entry.PendingMap) > 0 {
			describeRecord(record)

			return nil
		}

		if len(record.Entry.PendingLog) > 0 {
			describeRecord(record)

			return nil
		}

		entry := newDomainEntry(record, name, instance)
		if entry == nil {
			describeRecord(record)

			return nil
		}

		registry.Add(entry)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	err = registry.Validate()
	if err != nil {
		log.Fatal(err)
	}
}

func describeRecord(record *bival.Record) {
	name := record.Entry.Name
	instance := record.Entry.Instance

	if record.Type == "olh" {
		name = record.Entry.Key.Name
		instance = record.Entry.Key.Instance
	}

	log.Printf("name=%q type=%s instance=%q", name, record.Type, instance)
}

func newDomainEntry(record *bival.Record, name string, instance string) *domain.Entry {
	entry, validationErr := domain.NewEntry(domain.Kind(record.Type), name, instance)
	if validationErr != nil {
		return nil
	}

	return entry
}
