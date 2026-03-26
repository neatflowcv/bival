package domain

import "reflect"

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

	pairedPlainEntries := make([]*PlainEntry, 0, len(group.PlainEntries())-1)
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

func isValidVersionedHeadPlainEntry(entry *PlainEntry) bool {
	payload, ok := headPlainPayload(entry)
	if !ok {
		return false
	}

	return hasHeadIdentity(entry) &&
		hasHeadVersion(payload) &&
		hasHeadState(payload) &&
		hasZeroValueHeadMeta(payload.meta)
}

func headPlainPayload(entry *PlainEntry) (*DirPayload, bool) {
	if entry == nil || entry.entry == nil || entry.entry.payload == nil {
		return nil, false
	}

	payload := entry.entry.payload
	if payload.key == nil || payload.versionInfo == nil || payload.versionInfo.version == nil {
		return nil, false
	}

	if payload.state == nil || payload.meta == nil {
		return nil, false
	}

	return payload, true
}

func hasHeadIdentity(entry *PlainEntry) bool {
	return entry.Index() == entry.Name() && entry.Instance() == ""
}

func hasHeadVersion(payload *DirPayload) bool {
	return payload.versionInfo.Version().IsMissing() &&
		payload.versionInfo.VersionedEpoch() == 0
}

func hasHeadState(payload *DirPayload) bool {
	return !payload.state.exists &&
		payload.state.locator == "" &&
		payload.state.tag == "" &&
		payload.state.flags == 8 &&
		len(payload.pendingMaps) == 0
}

func hasZeroValueHeadMeta(meta *Meta) bool {
	if !hasHeadMetaParts(meta) {
		return false
	}

	return hasZeroObjectSpec(meta.objectSpec) &&
		hasZeroAuditInfo(meta.auditInfo) &&
		hasZeroContentInfo(meta.contentInfo) &&
		hasZeroOwner(meta.owner)
}

func hasHeadMetaParts(meta *Meta) bool {
	if meta == nil {
		return false
	}

	return meta.objectSpec != nil &&
		meta.auditInfo != nil &&
		meta.contentInfo != nil &&
		meta.owner != nil
}

func hasZeroObjectSpec(spec *ObjectSpec) bool {
	return spec.category == 0 &&
		spec.size == 0 &&
		spec.accountedSize == 0 &&
		!spec.appendable
}

func hasZeroAuditInfo(info *AuditInfo) bool {
	return info.mTime.IsZero() && info.eTag == ""
}

func hasZeroContentInfo(info *ContentInfo) bool {
	return info.storageClass == "" && info.contentType == ""
}

func hasZeroOwner(owner *Owner) bool {
	return owner.userID == "" && owner.displayName == ""
}

func hasValidVersionPairs(plainEntries []*PlainEntry, instanceEntries []*InstanceEntry) bool {
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

func buildPlainEntryMap(entries []*PlainEntry) (map[versionedEntryKey]*PlainEntry, bool) {
	entriesByKey := make(map[versionedEntryKey]*PlainEntry, len(entries))
	for _, entry := range entries {
		key, ok := versionedPlainEntryKey(entry)
		if !ok || hasPlainKey(entriesByKey, key) {
			return nil, false
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, true
}

func buildInstanceEntryMap(entries []*InstanceEntry) (map[versionedEntryKey]*InstanceEntry, bool) {
	entriesByKey := make(map[versionedEntryKey]*InstanceEntry, len(entries))
	for _, entry := range entries {
		key, ok := versionedInstanceEntryKey(entry)
		if !ok || hasInstanceKey(entriesByKey, key) {
			return nil, false
		}

		entriesByKey[key] = entry
	}

	return entriesByKey, true
}

func hasPlainKey(entries map[versionedEntryKey]*PlainEntry, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func hasInstanceKey(entries map[versionedEntryKey]*InstanceEntry, key versionedEntryKey) bool {
	_, exists := entries[key]

	return exists
}

func versionedPlainEntryKey(entry *PlainEntry) (versionedEntryKey, bool) {
	return entryKeyFromPayload(dirPayloadFromPlainEntry(entry))
}

func versionedInstanceEntryKey(entry *InstanceEntry) (versionedEntryKey, bool) {
	return entryKeyFromPayload(dirPayloadFromInstanceEntry(entry))
}

func dirPayloadFromPlainEntry(entry *PlainEntry) (*DirPayload, bool) {
	if entry == nil || entry.entry == nil || entry.entry.payload == nil {
		return nil, false
	}

	return entry.entry.payload, true
}

func dirPayloadFromInstanceEntry(entry *InstanceEntry) (*DirPayload, bool) {
	if entry == nil || entry.entry == nil || entry.entry.payload == nil {
		return nil, false
	}

	return entry.entry.payload, true
}

func entryKeyFromPayload(payload *DirPayload, ok bool) (versionedEntryKey, bool) {
	if !ok || payload == nil || payload.key == nil {
		return invalidVersionedEntryKey(), false
	}

	if payload.key.name == "" || payload.versionInfo == nil || payload.versionInfo.version == nil {
		return invalidVersionedEntryKey(), false
	}

	return versionedEntryKey{
		name:     payload.key.name,
		instance: payload.key.instance,
		pool:     payload.versionInfo.Version().Pool(),
		epoch:    payload.versionInfo.Version().Epoch(),
		vEpoch:   payload.versionInfo.VersionedEpoch(),
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

func hasEquivalentVersionPayload(plainEntry *PlainEntry, instanceEntry *InstanceEntry) bool {
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

func payloadWithoutTag(payload *DirPayload) *DirPayload {
	if payload == nil || payload.key == nil || payload.versionInfo == nil || payload.state == nil || payload.meta == nil {
		return nil
	}

	return &DirPayload{
		key:         payload.key,
		versionInfo: payload.versionInfo,
		state: &DirState{
			locator: payload.state.locator,
			exists:  payload.state.exists,
			tag:     "",
			flags:   payload.state.flags,
		},
		meta:        payload.meta,
		pendingMaps: payload.pendingMaps,
	}
}

func hasValidOLHReference(olhEntries []*OLHEntry, instanceEntries []*InstanceEntry) bool {
	olhEntry, olhEntryOK := singleValidOLHEntry(olhEntries)
	if !olhEntryOK {
		return false
	}

	instanceSet, instanceSetOK := instanceNameSet(instanceEntries)
	if !instanceSetOK {
		return false
	}

	referencedInstance := olhEntry.payload.key.instance
	_, exists := instanceSet[referencedInstance]

	return exists
}

func singleValidOLHEntry(entries []*OLHEntry) (*OLHEntry, bool) {
	if len(entries) != 1 {
		return nil, false
	}

	entry := entries[0]
	if entry == nil || entry.payload == nil || entry.payload.key == nil {
		return nil, false
	}

	if entry.HasPendingLog() {
		return nil, false
	}

	return entry, true
}

func instanceNameSet(entries []*InstanceEntry) (map[string]struct{}, bool) {
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

func instanceName(entry *InstanceEntry) (string, bool) {
	payload, ok := dirPayloadFromInstanceEntry(entry)
	if !ok || payload.key == nil {
		return "", false
	}

	return payload.key.instance, true
}
