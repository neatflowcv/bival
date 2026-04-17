package domain

func IsVersionPair(plain *PlainEntry, instance *InstanceEntry) bool {
	if plain == nil || instance == nil {
		return false
	}

	return plain.name == instance.name &&
		plain.instance == instance.instance &&
		equalEntryVersion(plain, instance) &&
		equalEntryState(plain, instance) &&
		equalMeta(plain.meta, instance.meta) &&
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

func equalMeta(left *Meta, right *Meta) bool {
	if left == nil || right == nil {
		return left == right
	}

	return equalMetaObject(left, right) &&
		equalMetaAudit(left, right) &&
		equalMetaContent(left, right) &&
		equalMetaOwner(left, right)
}

func equalMetaObject(left *Meta, right *Meta) bool {
	return left.category == right.category &&
		left.size == right.size &&
		left.accountedSize == right.accountedSize &&
		left.appendable == right.appendable
}

func equalMetaAudit(left *Meta, right *Meta) bool {
	return left.mTime.Equal(right.mTime) &&
		left.eTag == right.eTag
}

func equalMetaContent(left *Meta, right *Meta) bool {
	return left.storageClass == right.storageClass &&
		left.contentType == right.contentType
}

func equalMetaOwner(left *Meta, right *Meta) bool {
	return left.ownerUserID == right.ownerUserID &&
		left.ownerDisplayName == right.ownerDisplayName
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
		equalPendingMapVal(left.val, right.val)
}

func equalPendingMapVal(left *PendingMapVal, right *PendingMapVal) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.state == right.state &&
		left.timestamp.Equal(right.timestamp) &&
		left.op == right.op
}
