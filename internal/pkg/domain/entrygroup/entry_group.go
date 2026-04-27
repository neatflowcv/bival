package entrygroup

import (
	"errors"
	"fmt"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errEntryGroupNameMismatch = errors.New("entry name does not match group name")

const (
	missingVersionedHeadReason       = "missing versioned head"
	duplicateVersionedHeadReason     = "duplicate versioned head"
	invalidVersionedHeadReason       = "invalid versioned head"
	duplicateVersionedEntryKeyReason = "duplicate versioned entry key"
	missingMatchingPlainReason       = "instance version has no matching plain"
	missingMatchingInstanceReason    = "plain version has no matching instance"
	mismatchedVersionPairReason      = "plain and instance versions differ"
	missingOLHReason                 = "missing olh"
	invalidOLHReason                 = "invalid olh"
	invalidOLHReferenceReason        = "olh references missing instance"
	staleOLHReferenceReason          = "stale olh reference allows only one version"
	staleDeleteMarkerOLHReason       = "stale delete-marker olh allows no versions"
	tooManyVersionedEntriesReason    = "too many versioned entries"
	maxVersionedEntryCount           = 8
)

type EntryGroup struct {
	name            string
	plainEntries    []*domain.Plain
	instanceEntries []*domain.Instance
	olhEntries      []*domain.OLH
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

func (g *EntryGroup) ProblemReason() []*Issue {
	return g.Issues()
}

func (g *EntryGroup) Issues() []*Issue {
	return diagnoseObject(g)
}

func (g *EntryGroup) AddPlain(entry *domain.Plain) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.plainEntries = append(g.plainEntries, entry)

	return nil
}

func (g *EntryGroup) AddInstance(entry *domain.Instance) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.instanceEntries = append(g.instanceEntries, entry)

	return nil
}

func (g *EntryGroup) AddOLH(entry *domain.OLH) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.olhEntries = append(g.olhEntries, entry)

	return nil
}

func (g *EntryGroup) PlainEntries() []*domain.Plain {
	return g.plainEntries
}

func (g *EntryGroup) InstanceEntries() []*domain.Instance {
	return g.instanceEntries
}

func (g *EntryGroup) OLHEntries() []*domain.OLH {
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

func (g *EntryGroup) isUnversionedObject() bool {
	if !g.hasUnversionedEntryCounts() {
		return false
	}

	plainEntries := g.PlainEntries()
	if len(plainEntries) != 1 {
		return false
	}

	return plainEntries[0].IsUnversioned()
}

func (g *EntryGroup) hasUnversionedEntryCounts() bool {
	return g.PlainCount() == 1 &&
		g.InstanceCount() == 0 &&
		g.OLHCount() == 0
}
