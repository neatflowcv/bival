package entrygroup

import "github.com/neatflowcv/bival/internal/pkg/domain"

func isVersionedHeadCandidate(entry *domain.Plain) bool {
	return entry.Index() == entry.Name() &&
		entry.Instance() == ""
}

type versionedEntryKey struct {
	name     string
	instance string
	pool     int
	epoch    int
	vEpoch   int
}

func buildPlainEntryMap(entries []*domain.Plain) (map[versionedEntryKey]*domain.Plain, string) {
	entriesByKey := make(map[versionedEntryKey]*domain.Plain, len(entries))
	for _, entry := range entries {
		if entry.Name() == "" {
			return nil, duplicateVersionedEntryKeyReason
		}

		key := versionedEntryKey{
			name:     entry.Name(),
			instance: entry.Instance(),
			pool:     entry.VersionPool(),
			epoch:    entry.VersionEpoch(),
			vEpoch:   entry.VersionedEpoch(),
		}

		if hasPlainKey(entriesByKey, key) {
			return nil, duplicateVersionedEntryKeyReason
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, ""
}

func buildInstanceEntryMap(entries []*domain.Instance) (map[versionedEntryKey]*domain.Instance, string) {
	entriesByKey := make(map[versionedEntryKey]*domain.Instance, len(entries))
	for _, entry := range entries {
		if entry.Name() == "" {
			return nil, duplicateVersionedEntryKeyReason
		}

		key := versionedEntryKey{
			name:     entry.Name(),
			instance: entry.Instance(),
			pool:     entry.VersionPool(),
			epoch:    entry.VersionEpoch(),
			vEpoch:   entry.VersionedEpoch(),
		}

		if hasInstanceKey(entriesByKey, key) {
			return nil, duplicateVersionedEntryKeyReason
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, ""
}

func hasPlainKey(entries map[versionedEntryKey]*domain.Plain, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func hasInstanceKey(entries map[versionedEntryKey]*domain.Instance, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func hasValidOLHReference(olhEntries []*domain.OLH, instanceEntries []*domain.Instance) bool {
	olhEntry, olhEntryOK := singleValidOLHEntry(olhEntries)
	if !olhEntryOK {
		return false
	}

	instanceSet, instanceSetOK := instanceNameSet(instanceEntries)
	if !instanceSetOK {
		return false
	}

	referencedInstance := olhEntry.Instance()
	_, exists := instanceSet[referencedInstance]

	return exists
}

func singleValidOLHEntry(entries []*domain.OLH) (*domain.OLH, bool) {
	if len(entries) != 1 {
		return nil, false
	}

	entry := entries[0]
	if entry == nil || entry.Name() == "" {
		return nil, false
	}

	if entry.HasPendingLog() {
		return nil, false
	}

	return entry, true
}

func instanceNameSet(entries []*domain.Instance) (map[string]struct{}, bool) {
	set := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		instance, ok := instanceName(entry)
		if !ok {
			return nil, false
		}

		set[instance] = struct{}{}
	}

	return set, true
}

func instanceName(entry *domain.Instance) (string, bool) {
	if entry.Name() == "" {
		return "", false
	}

	return entry.Instance(), true
}
