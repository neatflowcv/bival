package domain

func IsVersionPair(plain *PlainEntry, instance *InstanceEntry) bool {
	if plain == nil || instance == nil {
		return false
	}

	return plain.name == instance.name &&
		plain.instance == instance.instance &&
		equalVersionInfo(plain.versionInfo, instance.versionInfo) &&
		equalStateWithoutTag(plain.state, instance.state) &&
		equalMeta(plain.meta, instance.meta) &&
		equalPendingMaps(plain.pendingMaps, instance.pendingMaps)
}

func equalVersionInfo(left *DirVersionInfo, right *DirVersionInfo) bool {
	if left == nil || right == nil {
		return left == right
	}

	return equalVersion(left.version, right.version) &&
		left.versionedEpoch == right.versionedEpoch
}

func equalVersion(left *Version, right *Version) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.pool == right.pool &&
		left.epoch == right.epoch
}

func equalStateWithoutTag(left *DirState, right *DirState) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.locator == right.locator &&
		left.exists == right.exists &&
		left.flags == right.flags
}

func equalMeta(left *Meta, right *Meta) bool {
	if left == nil || right == nil {
		return left == right
	}

	return equalObjectSpec(left.objectSpec, right.objectSpec) &&
		equalAuditInfo(left.auditInfo, right.auditInfo) &&
		equalContentInfo(left.contentInfo, right.contentInfo) &&
		equalOwner(left.owner, right.owner)
}

func equalObjectSpec(left *ObjectSpec, right *ObjectSpec) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.category == right.category &&
		left.size == right.size &&
		left.accountedSize == right.accountedSize &&
		left.appendable == right.appendable
}

func equalAuditInfo(left *AuditInfo, right *AuditInfo) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.mTime.Equal(right.mTime) &&
		left.eTag == right.eTag
}

func equalContentInfo(left *ContentInfo, right *ContentInfo) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.storageClass == right.storageClass &&
		left.contentType == right.contentType
}

func equalOwner(left *Owner, right *Owner) bool {
	if left == nil || right == nil {
		return left == right
	}

	return left.userID == right.userID &&
		left.displayName == right.displayName
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
