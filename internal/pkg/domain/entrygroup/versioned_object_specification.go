package entrygroup

import "github.com/neatflowcv/bival/internal/pkg/domain"

type versionedObjectSpecification struct{}

type versionedEntryKey struct {
	name     string
	instance string
	pool     int
	epoch    int
	vEpoch   int
}

func (versionedObjectSpecification) IsSatisfiedBy(group *EntryGroup) bool {
	if !hasVersionedEntryCounts(group) {
		return false
	}

	headCount := 0

	pairedPlainEntries := make([]*domain.Plain, 0, len(group.PlainEntries())-1)
	for _, entry := range group.PlainEntries() {
		if entry.IsPlaceholder() {
			headCount++

			continue
		}

		pairedPlainEntries = append(pairedPlainEntries, entry)
	}

	if headCount != 1 {
		return false
	}

	if !hasValidVersionPairs(pairedPlainEntries, group.InstanceEntries()) {
		return false
	}

	return hasValidOLHReference(group.OLHEntries(), group.InstanceEntries())
}

func hasVersionedEntryCounts(group *EntryGroup) bool {
	return group.PlainCount() >= 2 &&
		group.InstanceCount() >= 1 &&
		group.OLHCount() == 1 &&
		group.PlainCount() == group.InstanceCount()+1
}

func hasValidVersionPairs(plainEntries []*domain.Plain, instanceEntries []*domain.InstanceEntry) bool {
	if len(plainEntries) == 0 || len(plainEntries) != len(instanceEntries) {
		return false
	}

	plainByKey, plainMapOK := buildPlainEntryMap(plainEntries)
	if !plainMapOK {
		return false
	}

	instanceByKey, instanceMapOK := buildInstanceEntryMap(instanceEntries)
	if !instanceMapOK {
		return false
	}

	for key, plainEntry := range plainByKey {
		instanceEntry, exists := instanceByKey[key]
		if !exists {
			return false
		}

		if !domain.IsVersionPair(plainEntry, instanceEntry) {
			return false
		}
	}

	return true
}

func buildPlainEntryMap(entries []*domain.Plain) (map[versionedEntryKey]*domain.Plain, bool) {
	entriesByKey := make(map[versionedEntryKey]*domain.Plain, len(entries))
	for _, entry := range entries {
		if entry.Name() == "" {
			return nil, false
		}

		key := versionedEntryKey{
			name:     entry.Name(),
			instance: entry.Instance(),
			pool:     entry.VersionPool(),
			epoch:    entry.VersionEpoch(),
			vEpoch:   entry.VersionedEpoch(),
		}

		if hasPlainKey(entriesByKey, key) {
			return nil, false
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, true
}

func buildInstanceEntryMap(entries []*domain.InstanceEntry) (map[versionedEntryKey]*domain.InstanceEntry, bool) {
	entriesByKey := make(map[versionedEntryKey]*domain.InstanceEntry, len(entries))
	for _, entry := range entries {
		if entry.Name() == "" {
			return nil, false
		}

		key := versionedEntryKey{
			name:     entry.Name(),
			instance: entry.Instance(),
			pool:     entry.VersionPool(),
			epoch:    entry.VersionEpoch(),
			vEpoch:   entry.VersionedEpoch(),
		}

		if hasInstanceKey(entriesByKey, key) {
			return nil, false
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, true
}

func hasPlainKey(entries map[versionedEntryKey]*domain.Plain, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func hasInstanceKey(entries map[versionedEntryKey]*domain.InstanceEntry, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func hasValidOLHReference(olhEntries []*domain.OLH, instanceEntries []*domain.InstanceEntry) bool {
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

func instanceNameSet(entries []*domain.InstanceEntry) (map[string]struct{}, bool) {
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

func instanceName(entry *domain.InstanceEntry) (string, bool) {
	if entry.Name() == "" {
		return "", false
	}

	return entry.Instance(), true
}
