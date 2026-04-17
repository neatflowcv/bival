package entrygroup

import "github.com/neatflowcv/bival/internal/pkg/domain"

func isVersionedHeadCandidate(entry *domain.Plain) bool {
	return entry.Index() == entry.Name() &&
		entry.Instance() == ""
}

func buildPlainEntryMap(entries []*domain.Plain) (map[string]*domain.Plain, string) {
	entriesByKey := make(map[string]*domain.Plain, len(entries))
	for _, entry := range entries {
		if entry.Name() == "" {
			return nil, duplicateVersionedEntryKeyReason
		}

		key := entry.Instance()

		if hasPlainKey(entriesByKey, key) {
			return nil, duplicateVersionedEntryKeyReason
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, ""
}

func buildInstanceEntryMap(entries []*domain.Instance) (map[string]*domain.Instance, string) {
	entriesByKey := make(map[string]*domain.Instance, len(entries))
	for _, entry := range entries {
		if entry.Name() == "" {
			return nil, duplicateVersionedEntryKeyReason
		}

		key := entry.Instance()

		if hasInstanceKey(entriesByKey, key) {
			return nil, duplicateVersionedEntryKeyReason
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, ""
}

func hasPlainKey(entries map[string]*domain.Plain, key string) bool {
	_, exists := entries[key]

	return exists
}

func hasInstanceKey(entries map[string]*domain.Instance, key string) bool {
	_, exists := entries[key]

	return exists
}

func hasValidOLHReference(olhEntry *domain.OLH, instanceEntries []*domain.Instance) bool {
	instanceSet, instanceSetOK := instanceNameSet(instanceEntries)
	if !instanceSetOK {
		return false
	}

	referencedInstance := olhEntry.Instance()
	_, exists := instanceSet[referencedInstance]

	return exists
}

func singleValidOLHEntry(entries []*domain.OLH) (*domain.OLH, string) {
	if len(entries) == 0 {
		return nil, missingOLHReason
	}

	if len(entries) != 1 {
		return nil, invalidOLHReason
	}

	entry := entries[0]
	if entry == nil || entry.Name() == "" {
		return nil, invalidOLHReason
	}

	if entry.HasPendingLog() {
		return nil, invalidOLHReason
	}

	return entry, ""
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
