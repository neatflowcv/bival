package domain

import "reflect"

func IsVersionPair(plain *PlainEntry, instance *InstanceEntry) bool {
	return reflect.DeepEqual(
		payloadWithoutTag(plain.Payload()),
		payloadWithoutTag(instance.Payload()),
	)
}

func payloadWithoutTag(payload *DirPayload) *DirPayload {
	if payload == nil ||
		payload.Key() == nil ||
		payload.VersionInfo() == nil ||
		payload.State() == nil ||
		payload.Meta() == nil {
		return nil
	}

	return NewDirPayload(
		payload.Key(),
		payload.VersionInfo(),
		NewDirState(
			payload.State().Locator(),
			payload.State().Exists(),
			"",
			payload.State().Flags(),
		),
		payload.Meta(),
		payload.PendingMaps(),
	)
}
