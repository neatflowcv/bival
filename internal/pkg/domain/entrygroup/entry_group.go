package entrygroup

import (
	"errors"
	"fmt"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

var errEntryGroupNameMismatch = errors.New("entry name does not match group name")

const (
	invalidVersionedEntryCountsReason = "invalid versioned entry counts"
	missingVersionedHeadReason        = "missing versioned head"
	duplicateVersionedHeadReason      = "duplicate versioned head"
	invalidVersionedHeadReason        = "invalid versioned head"
	duplicateVersionedEntryKeyReason  = "duplicate versioned entry key"
	missingMatchingInstanceReason     = "plain version has no matching instance"
	mismatchedVersionPairReason       = "plain and instance versions differ"
	invalidOLHReferenceReason         = "olh references missing instance"
	tooManyVersionedEntriesReason     = "too many versioned entries"
	maxVersionedEntryCount            = 8
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

func (g *EntryGroup) ProblemReason() string {
	if g.HasPendingEntries() {
		return "pending entry exists"
	}

	if g.isUnversionedObject() {
		return ""
	}

	if g.ObjectKind() == VersionedObject && g.versionedEntryCount() > maxVersionedEntryCount {
		return tooManyVersionedEntriesReason
	}

	return g.versionedProblemReason()
}

func (g *EntryGroup) ObjectKind() ObjectKind {
	if g.isUnversionedObject() {
		return UnversionedObject
	}

	return VersionedObject
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

func (g *EntryGroup) hasVersionedEntryCounts() bool {
	return g.PlainCount() >= 2 &&
		g.InstanceCount() >= 1 &&
		g.OLHCount() == 1 &&
		g.PlainCount() == g.InstanceCount()+1
}

func (g *EntryGroup) versionedProblemReason() string {
	if !g.hasVersionedEntryCounts() {
		return invalidVersionedEntryCountsReason
	}

	pairedPlainEntries, headReason := g.versionedHeadState()
	if headReason != "" {
		return headReason
	}

	pairReason := g.versionedPairReason(pairedPlainEntries)
	if pairReason != "" {
		return pairReason
	}

	if !hasValidOLHReference(g.OLHEntries(), g.InstanceEntries()) {
		return invalidOLHReferenceReason
	}

	return ""
}

func (g *EntryGroup) versionedHeadState() ([]*domain.Plain, string) {
	headCount := 0
	invalidHeadCount := 0

	pairedPlainEntries := make([]*domain.Plain, 0, len(g.PlainEntries())-1)
	for _, entry := range g.PlainEntries() {
		if entry.IsPlaceholder() {
			headCount++

			continue
		}

		if isVersionedHeadCandidate(entry) {
			invalidHeadCount++

			continue
		}

		pairedPlainEntries = append(pairedPlainEntries, entry)
	}

	if headCount == 0 && invalidHeadCount > 0 {
		return nil, invalidVersionedHeadReason
	}

	if headCount == 0 {
		return nil, missingVersionedHeadReason
	}

	if headCount > 1 {
		return nil, duplicateVersionedHeadReason
	}

	return pairedPlainEntries, ""
}

func (g *EntryGroup) versionedPairReason(pairedPlainEntries []*domain.Plain) string {
	plainByKey, plainMapReason := buildPlainEntryMap(pairedPlainEntries)
	if plainMapReason != "" {
		return plainMapReason
	}

	instanceByKey, instanceMapReason := buildInstanceEntryMap(g.InstanceEntries())
	if instanceMapReason != "" {
		return instanceMapReason
	}

	for key, plainEntry := range plainByKey {
		instanceEntry, exists := instanceByKey[key]
		if !exists {
			return missingMatchingInstanceReason
		}

		if !domain.IsVersionPair(plainEntry, instanceEntry) {
			return mismatchedVersionPairReason
		}
	}

	return ""
}
