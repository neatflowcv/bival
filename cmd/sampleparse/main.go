package main

import (
	"log"
	"os"

	"github.com/neatflowcv/bival"
)

func main() {
	path := "sample.json"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	var (
		count     int
		totalSize int64
	)

	err := bival.ParseFile(path, func(record bival.Record) error {
		count++
		totalSize += record.Entry.Meta.Size

		describeRecord(record)

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("records=%d total_size=%d", count, totalSize)
}

func describeRecord(record bival.Record) {
	name := record.Entry.Name
	instance := record.Entry.Instance

	if record.Type == "olh" {
		name = record.Entry.Key.Name
		instance = record.Entry.Key.Instance
	}

	log.Printf("name=%q type=%s instance=%q", name, record.Type, instance)
}
