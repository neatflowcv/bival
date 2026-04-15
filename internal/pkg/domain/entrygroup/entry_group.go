package entrygroup

import (
	"errors"
	"fmt"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errEntryGroupNameMismatch = errors.New("entry name does not match group name")

const (
	unknownObjectReason           = "object kind is unknown"
	tooManyVersionedEntriesReason = "too many versioned entries"
	maxVersionedEntryCount        = 6
)

type EntryGroup struct {
	name            string
	plainEntries    []*domain.PlainEntry
	instanceEntries []*domain.InstanceEntry
	olhEntries      []*domain.OLHEntry
}

func New(name string) *EntryGroup {
	return &EntryGroup{
		name:            name,
		plainEntries:    nil,
		instanceEntries: nil,
		olhEntries:      nil,
	}
}

func (g *EntryGroup) Name() string {
	return g.name
}

func (g *EntryGroup) PlainCount() int {
	return len(g.plainEntries)
}

func (g *EntryGroup) InstanceCount() int {
	return len(g.instanceEntries)
}

func (g *EntryGroup) OLHCount() int {
	return len(g.olhEntries)
}

func (g *EntryGroup) HasPendingMap() bool {
	for _, entry := range g.plainEntries {
		if entry.HasPendingMap() {
			return true
		}
	}

	for _, entry := range g.instanceEntries {
		if entry.HasPendingMap() {
			return true
		}
	}

	return false
}

func (g *EntryGroup) HasPendingLog() bool {
	for _, entry := range g.olhEntries {
		if entry.HasPendingLog() {
			return true
		}
	}

	return false
}

func (g *EntryGroup) HasPendingEntries() bool {
	return g.HasPendingMap() || g.HasPendingLog()
}

func (g *EntryGroup) ProblemReason() string {
	if g.HasPendingEntries() {
		return "pending entry exists"
	}

	objectKind := NewClassifier().Classify(g)

	if objectKind == VersionedObject && g.versionedEntryCount() > maxVersionedEntryCount {
		return tooManyVersionedEntriesReason
	}

	if objectKind == UnknownObject {
		return unknownObjectReason
	}

	return ""
}

func (g *EntryGroup) AddPlain(entry *domain.PlainEntry) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.plainEntries = append(g.plainEntries, entry)

	return nil
}

func (g *EntryGroup) AddInstance(entry *domain.InstanceEntry) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.instanceEntries = append(g.instanceEntries, entry)

	return nil
}

func (g *EntryGroup) AddOLH(entry *domain.OLHEntry) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.olhEntries = append(g.olhEntries, entry)

	return nil
}

func (g *EntryGroup) PlainEntries() []*domain.PlainEntry {
	return g.plainEntries
}

func (g *EntryGroup) InstanceEntries() []*domain.InstanceEntry {
	return g.instanceEntries
}

func (g *EntryGroup) OLHEntries() []*domain.OLHEntry {
	return g.olhEntries
}

func (g *EntryGroup) validateName(entryName string) error {
	if entryName == g.name {
		return nil
	}

	return fmt.Errorf("%w: entry name %q does not match group name %q", errEntryGroupNameMismatch, entryName, g.name)
}

func (g *EntryGroup) versionedEntryCount() int {
	return g.PlainCount() + g.InstanceCount() + g.OLHCount()
}
