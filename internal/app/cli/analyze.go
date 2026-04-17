package cli

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/neatflowcv/bival/internal/pkg/domain/entrygroup"
)

var errUnsupportedBuiltEntryType = errors.New("unsupported built entry type")

type AnalyzeCmd struct {
	Input string `arg:"" help:"Path to a BI list JSON file to analyze by name group." name:"input" type:"path"`
}

func (cmd *AnalyzeCmd) Run() error {
	return analyzeFile(cmd.Input, log.Default())
}

func analyzeFile(path string, logger *log.Logger) error {
	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}

	defer func() {
		_ = file.Close()
	}()

	err = analyzeRecords(bilist.NewReader(file), logger)
	if err != nil {
		return err
	}

	return nil
}

func analyzeRecords(reader *bilist.Reader, logger *log.Logger) error {
	var (
		current     *entrygroup.EntryGroup
		currentName string
	)

	for {
		record, err := reader.Read()
		if done, readErr := finishAnalyze(err, current, logger); done {
			return readErr
		}

		name := recordName(record)
		if current != nil && currentName != name {
			logProblem(current, logger)
		}

		current = nextGroup(current, currentName, name)
		currentName = name

		err = addRecordToGroup(current, record)
		if err != nil {
			return err
		}
	}
}

func finishAnalyze(err error, current *entrygroup.EntryGroup, logger *log.Logger) (bool, error) {
	if errors.Is(err, io.EOF) {
		if current != nil {
			logProblem(current, logger)
		}

		return true, nil
	}

	if err != nil {
		return true, fmt.Errorf("read record: %w", err)
	}

	return false, nil
}

func nextGroup(current *entrygroup.EntryGroup, currentName string, nextName string) *entrygroup.EntryGroup {
	if current == nil || currentName != nextName {
		return entrygroup.New(nextName)
	}

	return current
}

func addRecordToGroup(group *entrygroup.EntryGroup, record *bilist.Record) error {
	entry, err := buildEntry(record)
	if err != nil {
		return fmt.Errorf("build entry idx=%q type=%q: %w", record.Idx, record.Type, err)
	}

	switch typed := entry.(type) {
	case *domain.Plain:
		err = group.AddPlain(typed)
	case *domain.Instance:
		err = group.AddInstance(typed)
	case *domain.OLH:
		err = group.AddOLH(typed)
	default:
		return fmt.Errorf("%w %T", errUnsupportedBuiltEntryType, entry)
	}

	if err != nil {
		return fmt.Errorf("group entry idx=%q type=%q: %w", record.Idx, record.Type, err)
	}

	return nil
}

func recordName(record *bilist.Record) string {
	if record.Type == "olh" {
		return record.Entry.Key.Name
	}

	return record.Entry.Name
}

func logProblem(group *entrygroup.EntryGroup, logger *log.Logger) {
	reasons := group.ProblemReason()
	if len(reasons) == 0 {
		return
	}

	format := "problem name=%q"
	args := make([]any, 0, len(reasons)+1)
	args = append(args, group.Name())

	var formatSb137 strings.Builder
	for _, reason := range reasons {
		formatSb137.WriteString(" reason=%q")

		args = append(args, reason)
	}
	format += formatSb137.String()

	logger.Printf(format, args...)
}
