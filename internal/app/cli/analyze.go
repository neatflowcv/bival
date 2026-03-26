package cli

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/neatflowcv/bival/internal/bilist"
	"github.com/neatflowcv/bival/internal/pkg/domain"
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
	var current *entryGroup

	for {
		record, err := reader.Read()
		if done, readErr := finishAnalyze(err, current, logger); done {
			return readErr
		}

		name := recordName(record)
		if current != nil && current.name != name {
			current.logProblem(logger)
		}

		current = nextGroup(current, name)

		err = addRecordToGroup(current, record)
		if err != nil {
			return err
		}
	}
}

func finishAnalyze(err error, current *entryGroup, logger *log.Logger) (bool, error) {
	if errors.Is(err, io.EOF) {
		if current != nil {
			current.logProblem(logger)
		}

		return true, nil
	}

	if err != nil {
		return true, fmt.Errorf("read record: %w", err)
	}

	return false, nil
}

func nextGroup(current *entryGroup, name string) *entryGroup {
	if current == nil || current.name != name {
		return newEntryGroup(name)
	}

	return current
}

func addRecordToGroup(group *entryGroup, record *bilist.Record) error {
	entry, err := buildEntry(record)
	if err != nil {
		return fmt.Errorf("build entry idx=%q type=%q: %w", record.Idx, record.Type, err)
	}

	switch typed := entry.(type) {
	case *domain.PlainEntry:
		group.addPlain(typed)
	case *domain.InstanceEntry:
		group.addInstance(typed)
	case *domain.OLHEntry:
		group.addOLH(typed)
	default:
		return fmt.Errorf("%w %T", errUnsupportedBuiltEntryType, entry)
	}

	return nil
}

func recordName(record *bilist.Record) string {
	if record.Type == "olh" {
		return record.Entry.Key.Name
	}

	return record.Entry.Name
}

type entryGroup struct {
	name          string
	plainCount    int
	instanceCount int
	olhCount      int
	hasPendingMap bool
	hasPendingLog bool
}

func newEntryGroup(name string) *entryGroup {
	return &entryGroup{
		name:          name,
		plainCount:    0,
		instanceCount: 0,
		olhCount:      0,
		hasPendingMap: false,
		hasPendingLog: false,
	}
}

func (g *entryGroup) addPlain(entry *domain.PlainEntry) {
	g.plainCount++
	if entry.HasPendingMap() {
		g.hasPendingMap = true
	}
}

func (g *entryGroup) addInstance(entry *domain.InstanceEntry) {
	g.instanceCount++
	if entry.HasPendingMap() {
		g.hasPendingMap = true
	}
}

func (g *entryGroup) addOLH(entry *domain.OLHEntry) {
	g.olhCount++
	if entry.HasPendingLog() {
		g.hasPendingLog = true
	}
}

func (g *entryGroup) logProblem(logger *log.Logger) {
	reason := g.problemReason()
	if reason == "" {
		return
	}

	logger.Printf("problem name=%q reason=%q", g.name, reason)
}

func (g *entryGroup) problemReason() string {
	if g.hasPendingMap || g.hasPendingLog {
		return "pending entry exists"
	}

	if g.plainCount == 1 {
		return ""
	}

	if g.olhCount != 1 {
		return "versioning object must have exactly one olh"
	}

	if g.instanceCount+1 != g.plainCount {
		return "versioning object must satisfy instance+1==plain"
	}

	return ""
}
