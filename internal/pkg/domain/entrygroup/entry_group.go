package entrygroup

import (
	"fmt"
	"slices"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

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

func (g *EntryGroup) AddPlain(entry *domain.Plain) {
	g.ensureNameMatches(entry.Name())
	g.plainEntries = append(g.plainEntries, entry)
}

func (g *EntryGroup) AddInstance(entry *domain.Instance) {
	g.ensureNameMatches(entry.Name())
	g.instanceEntries = append(g.instanceEntries, entry)
}

func (g *EntryGroup) AddOLH(entry *domain.OLH) {
	g.ensureNameMatches(entry.Name())
	g.olhEntries = append(g.olhEntries, entry)
}

func (g *EntryGroup) PlainEntries() []*domain.Plain {
	return slices.Clone(g.plainEntries)
}

func (g *EntryGroup) InstanceEntries() []*domain.Instance {
	return slices.Clone(g.instanceEntries)
}

func (g *EntryGroup) OLHEntries() []*domain.OLH {
	return slices.Clone(g.olhEntries)
}

// 이름 불일치는 복구 대상이 아니라 호출자 버그로 본다.
// 이 제약은 의도된 것이며, 잘못된 엔트리가 추가되면 panic으로 즉시 드러내는 편이 낫다.
func (g *EntryGroup) ensureNameMatches(entryName string) {
	if entryName == g.name {
		return
	}

	panic(fmt.Sprintf("entry name %q does not match group name %q", entryName, g.name))
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
