package domain

import (
	"errors"
	"fmt"
)

var errEntryGroupNameMismatch = errors.New("entry name does not match group name")

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

func (g *EntryGroup) ProblemReason() string {
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

func (g *EntryGroup) validateName(entryName string) error {
	if entryName == g.name {
		return nil
	}

	return fmt.Errorf("%w: entry name %q does not match group name %q", errEntryGroupNameMismatch, entryName, g.name)
}
