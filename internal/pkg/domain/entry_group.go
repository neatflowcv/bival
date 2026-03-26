package domain

import (
	"errors"
	"fmt"
)

var errEntryGroupNameMismatch = errors.New("entry name does not match group name")

const unknownObjectReason = "object kind is unknown"

type EntryGroup struct {
	name          string
	plainCount    int
	instanceCount int
	olhCount      int
	hasPendingMap bool
	hasPendingLog bool
}

func NewEntryGroup(name string) *EntryGroup {
	return &EntryGroup{
		name:          name,
		plainCount:    0,
		instanceCount: 0,
		olhCount:      0,
		hasPendingMap: false,
		hasPendingLog: false,
	}
}

func (g *EntryGroup) Name() string {
	return g.name
}

func (g *EntryGroup) PlainCount() int {
	return g.plainCount
}

func (g *EntryGroup) InstanceCount() int {
	return g.instanceCount
}

func (g *EntryGroup) OLHCount() int {
	return g.olhCount
}

func (g *EntryGroup) HasPendingMap() bool {
	return g.hasPendingMap
}

func (g *EntryGroup) HasPendingLog() bool {
	return g.hasPendingLog
}

func (g *EntryGroup) HasPendingEntries() bool {
	return g.hasPendingMap || g.hasPendingLog
}

func (g *EntryGroup) ProblemReason() string {
	if g.HasPendingEntries() {
		return "pending entry exists"
	}

	if NewEntryGroupClassifier().Classify(g) == UnknownObject {
		return unknownObjectReason
	}

	return ""
}

func (g *EntryGroup) AddPlain(entry *PlainEntry) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.plainCount++
	if entry.HasPendingMap() {
		g.hasPendingMap = true
	}

	return nil
}

func (g *EntryGroup) AddInstance(entry *InstanceEntry) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.instanceCount++
	if entry.HasPendingMap() {
		g.hasPendingMap = true
	}

	return nil
}

func (g *EntryGroup) AddOLH(entry *OLHEntry) error {
	err := g.validateName(entry.Name())
	if err != nil {
		return err
	}

	g.olhCount++
	if entry.HasPendingLog() {
		g.hasPendingLog = true
	}

	return nil
}

func (g *EntryGroup) validateName(entryName string) error {
	if entryName == g.name {
		return nil
	}

	return fmt.Errorf("%w: entry name %q does not match group name %q", errEntryGroupNameMismatch, entryName, g.name)
}
