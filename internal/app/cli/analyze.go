package cli

import (
	"errors"
	"fmt"
	"log"
	"slices"

	"github.com/neatflowcv/bival/internal/bilist"
)

var errInputNotSorted = errors.New("input must be sorted by name")

type AnalyzeCmd struct {
	Input string `arg:"" help:"Path to a BI list JSON file sorted by name." name:"input" type:"path"`
}

func (cmd *AnalyzeCmd) Run() error {
	return analyzeFile(cmd.Input, log.Default())
}

func analyzeFile(path string, logger *log.Logger) error {
	counts := map[string]int{}
	total := 0

	var prevName string

	err := bilist.ParseFile(path, func(record *bilist.Record) error {
		name := recordName(record)
		if prevName != "" && name < prevName {
			return fmt.Errorf(
				"%w: previous=%q current=%q idx=%q",
				errInputNotSorted,
				prevName,
				name,
				record.Idx,
			)
		}

		prevName = name
		total++

		entry, err := buildEntry(record)
		if err != nil {
			return fmt.Errorf("build entry idx=%q type=%q: %w", record.Idx, record.Type, err)
		}

		counts[record.Type]++
		logger.Printf("ok type=%s name=%q idx=%q entry=%T", record.Type, name, record.Idx, entry)

		return nil
	})
	if err != nil {
		return fmt.Errorf("analyze file: %w", err)
	}

	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}

	slices.Sort(keys)
	logger.Printf("summary total=%d", total)

	for _, key := range keys {
		logger.Printf("summary type=%s count=%d", key, counts[key])
	}

	return nil
}

func recordName(record *bilist.Record) string {
	if record.Type == "olh" {
		return record.Entry.Key.Name
	}

	return record.Entry.Name
}
