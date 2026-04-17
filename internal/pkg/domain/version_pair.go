package domain

func IsVersionPair(plain *PlainEntry, instance *InstanceEntry) bool {
	if plain == nil || instance == nil {
		return false
	}

	return plain.name == instance.name &&
		plain.instance == instance.instance &&
		equalEntryVersion(plain, instance) &&
		equalEntryState(plain, instance) &&
		equalEntryMeta(plain, instance) &&
		equalPendingMaps(plain.pendingMaps, instance.pendingMaps)
}

func equalEntryVersion(plain *PlainEntry, instance *InstanceEntry) bool {
	return plain.pool == instance.pool &&
		plain.epoch == instance.epoch &&
		plain.vEpoch == instance.vEpoch
}

func equalEntryState(plain *PlainEntry, instance *InstanceEntry) bool {
	return plain.locator == instance.locator &&
		plain.exists == instance.exists &&
		plain.flags == instance.flags
}

func equalEntryMeta(plain *PlainEntry, instance *InstanceEntry) bool {
	return plain.category == instance.category &&
		plain.size == instance.size &&
		plain.accountedSize == instance.accountedSize &&
		plain.appendable == instance.appendable &&
		plain.mTime.Equal(instance.mTime) &&
		plain.eTag == instance.eTag &&
		plain.storageClass == instance.storageClass &&
		plain.contentType == instance.contentType &&
		plain.ownerUserID == instance.ownerUserID &&
		plain.ownerDisplayName == instance.ownerDisplayName
}

func equalPendingMaps(left []*PendingMap, right []*PendingMap) bool {
	if len(left) != len(right) {
		return false
	}

	for i := range left {
		if !equalPendingMap(left[i], right[i]) {
			return false
		}
	}

	return true
}

func equalPendingMap(left *PendingMap, right *PendingMap) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.key == right.key &&
		left.state == right.state &&
		left.timestamp.Equal(right.timestamp) &&
		left.op == right.op
}
