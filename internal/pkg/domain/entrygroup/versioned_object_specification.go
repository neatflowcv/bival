package entrygroup

import (
	"reflect"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

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

	pairedPlainEntries := make([]*domain.PlainEntry, 0, len(group.PlainEntries())-1)
	for _, entry := range group.PlainEntries() {
		if isValidVersionedHeadPlainEntry(entry) {
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

func isValidVersionedHeadPlainEntry(entry *domain.PlainEntry) bool {
	payload, ok := headPlainPayload(entry)
	if !ok {
		return false
	}

	return hasHeadIdentity(entry) &&
		hasHeadVersion(payload) &&
		hasHeadState(payload) &&
		hasHeadMetaParts(payload.Meta()) &&
		payload.Meta().IsDefault()
}

func headPlainPayload(entry *domain.PlainEntry) (*domain.DirPayload, bool) {
	payload := entry.Payload()
	if payload == nil || payload.Key() == nil || payload.VersionInfo() == nil || payload.VersionInfo().Version() == nil {
		return nil, false
	}

	if payload.State() == nil || payload.Meta() == nil {
		return nil, false
	}

	return payload, true
}

func hasHeadIdentity(entry *domain.PlainEntry) bool {
	return entry.Index() == entry.Name() && entry.Instance() == ""
}

func hasHeadVersion(payload *domain.DirPayload) bool {
	return payload.VersionInfo().Version().IsMissing() &&
		payload.VersionInfo().VersionedEpoch() == 0
}

func hasHeadState(payload *domain.DirPayload) bool {
	return !payload.State().Exists() &&
		payload.State().Locator() == "" &&
		payload.State().Tag() == "" &&
		payload.State().Flags() == 8 &&
		len(payload.PendingMaps()) == 0
}

func hasHeadMetaParts(meta *domain.Meta) bool {
	return meta.HasParts()
}

func hasValidVersionPairs(plainEntries []*domain.PlainEntry, instanceEntries []*domain.InstanceEntry) bool {
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

		if !hasEquivalentVersionPayload(plainEntry, instanceEntry) {
			return false
		}
	}

	return true
}

func buildPlainEntryMap(entries []*domain.PlainEntry) (map[versionedEntryKey]*domain.PlainEntry, bool) {
	entriesByKey := make(map[versionedEntryKey]*domain.PlainEntry, len(entries))
	for _, entry := range entries {
		key, ok := versionedPlainEntryKey(entry)
		if !ok || hasPlainKey(entriesByKey, key) {
			return nil, false
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, true
}

func buildInstanceEntryMap(entries []*domain.InstanceEntry) (map[versionedEntryKey]*domain.InstanceEntry, bool) {
	entriesByKey := make(map[versionedEntryKey]*domain.InstanceEntry, len(entries))
	for _, entry := range entries {
		key, ok := versionedInstanceEntryKey(entry)
		if !ok || hasInstanceKey(entriesByKey, key) {
			return nil, false
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, true
}

func hasPlainKey(entries map[versionedEntryKey]*domain.PlainEntry, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func hasInstanceKey(entries map[versionedEntryKey]*domain.InstanceEntry, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func versionedPlainEntryKey(entry *domain.PlainEntry) (versionedEntryKey, bool) {
	return entryKeyFromPayload(dirPayloadFromPlainEntry(entry))
}

func versionedInstanceEntryKey(entry *domain.InstanceEntry) (versionedEntryKey, bool) {
	return entryKeyFromPayload(dirPayloadFromInstanceEntry(entry))
}

func dirPayloadFromPlainEntry(entry *domain.PlainEntry) (*domain.DirPayload, bool) {
	payload := entry.Payload()
	if payload == nil {
		return nil, false
	}

	return payload, true
}

func dirPayloadFromInstanceEntry(entry *domain.InstanceEntry) (*domain.DirPayload, bool) {
	payload := entry.Payload()
	if payload == nil {
		return nil, false
	}

	return payload, true
}

func entryKeyFromPayload(payload *domain.DirPayload, ok bool) (versionedEntryKey, bool) {
	if !ok || payload == nil || payload.Key() == nil {
		return invalidVersionedEntryKey(), false
	}

	if payload.Key().Name() == "" || payload.VersionInfo() == nil || payload.VersionInfo().Version() == nil {
		return invalidVersionedEntryKey(), false
	}

	return versionedEntryKey{
		name:     payload.Key().Name(),
		instance: payload.Key().Instance(),
		pool:     payload.VersionInfo().Version().Pool(),
		epoch:    payload.VersionInfo().Version().Epoch(),
		vEpoch:   payload.VersionInfo().VersionedEpoch(),
	}, true
}

func invalidVersionedEntryKey() versionedEntryKey {
	return versionedEntryKey{
		name:     "",
		instance: "",
		pool:     0,
		epoch:    0,
		vEpoch:   0,
	}
}

func hasEquivalentVersionPayload(plainEntry *domain.PlainEntry, instanceEntry *domain.InstanceEntry) bool {
	plainPayload, plainOK := dirPayloadFromPlainEntry(plainEntry)

	instancePayload, instanceOK := dirPayloadFromInstanceEntry(instanceEntry)
	if !plainOK || !instanceOK {
		return false
	}

	return reflect.DeepEqual(
		payloadWithoutTag(plainPayload),
		payloadWithoutTag(instancePayload),
	)
}

func payloadWithoutTag(payload *domain.DirPayload) *domain.DirPayload {
	if payload == nil ||
		payload.Key() == nil ||
		payload.VersionInfo() == nil ||
		payload.State() == nil ||
		payload.Meta() == nil {
		return nil
	}

	return domain.NewDirPayload(
		payload.Key(),
		payload.VersionInfo(),
		domain.NewDirState(
			payload.State().Locator(),
			payload.State().Exists(),
			"",
			payload.State().Flags(),
		),
		payload.Meta(),
		payload.PendingMaps(),
	)
}

func hasValidOLHReference(olhEntries []*domain.OLHEntry, instanceEntries []*domain.InstanceEntry) bool {
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

func singleValidOLHEntry(entries []*domain.OLHEntry) (*domain.OLHEntry, bool) {
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
	payload, ok := dirPayloadFromInstanceEntry(entry)
	if !ok || payload.Key() == nil {
		return "", false
	}

	return payload.Key().Instance(), true
}
