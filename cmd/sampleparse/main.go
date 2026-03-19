package main

import (
	"log"
	"os"
	"strconv"

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

	const previewLimit = 3

	err := bival.ParseFile(path, func(record bival.Record) error {
		count++
		totalSize += record.Entry.Meta.Size

		if count <= previewLimit {
			log.Printf(
				"%s: type=%s idx=%q size=%d exists=%t",
				strconv.Itoa(count),
				record.Type,
				record.Idx,
				record.Entry.Meta.Size,
				record.Entry.Exists,
			)
		}

		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("records=%d total_size=%d", count, totalSize)
}
