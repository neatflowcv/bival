package domain_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestEntryGroupAddPlainTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddPlain(newPlainEntry("alpha", "", true))
	require.NoError(t, err)
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupAddInstanceTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddInstance(newInstanceEntry("alpha", true))
	require.NoError(t, err)
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupAddOLHTracksPendingLog(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddOLH(newOLHEntry("alpha", true))
	require.NoError(t, err)
	require.True(t, group.HasPendingEntries())
}

func TestEntryGroupClassifierReturnsUnversionedObject(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	err := group.AddPlain(newUnversionedPlainEntry())
	require.NoError(t, err)
	require.Equal(t, domain.UnversionedObject, classifier.Classify(group))
}

func TestEntryGroupClassifierReturnsUnknownObjectWhenNoRuleMatches(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newUnversionedPlainEntry()))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierReturnsUnknownObjectWhenVersionedRuleDoesNotMatch(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newUnversionedPlainEntry()))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v2", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.NoError(t, group.AddOLH(newOLHEntry("alpha", false)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierReturnsVersionedObject(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newUnversionedPlainEntry()))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.NoError(t, group.AddOLH(newOLHEntry("alpha", false)))
	require.Equal(t, domain.VersionedObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenIdxDiffersFromName(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha-idx",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenInstanceIsNotEmpty(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"v1",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenVersionPoolIsBelowMinimum(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		0,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenVersionEpochIsBelowMinimum(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		0,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenExistsIsFalse(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		false,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenMTimeIsZero(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		time.Time{},
		"etag",
		"tag",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenETagIsEmpty(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"",
		"tag",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenTagIsEmpty(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"",
		0,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierRejectsUnversionedWhenFlagsAreNotZero(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		1,
	)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupRejectsMismatchedPlainName(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddPlain(newPlainEntry("beta", "", false))
	require.EqualError(
		t,
		err,
		"entry name does not match group name: entry name \"beta\" does not match group name \"alpha\"",
	)
}

func TestEntryGroupRejectsMismatchedInstanceName(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddInstance(newInstanceEntry("beta", false))
	require.EqualError(
		t,
		err,
		"entry name does not match group name: entry name \"beta\" does not match group name \"alpha\"",
	)
}

func TestEntryGroupRejectsMismatchedOLHName(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddOLH(newOLHEntry("beta", false))
	require.EqualError(
		t,
		err,
		"entry name does not match group name: entry name \"beta\" does not match group name \"alpha\"",
	)
}

func newPlainEntry(name string, instance string, pending bool) *domain.PlainEntry {
	return newCustomPlainEntry(
		name,
		name,
		instance,
		pending,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)
}

func newUnversionedPlainEntry() *domain.PlainEntry {
	return newCustomPlainEntry(
		"alpha",
		"alpha",
		"",
		false,
		1,
		1,
		true,
		sampleMTime(),
		"etag",
		"tag",
		0,
	)
}

func newInstanceEntry(name string, pending bool) *domain.InstanceEntry {
	return domain.NewInstanceEntry(newDirEntry(name, "v1", pending))
}

func newDirEntry(name string, instance string, pending bool) *domain.DirEntry {
	var pendingMaps []*domain.PendingMap
	if pending {
		pendingMaps = []*domain.PendingMap{nil}
	}

	return domain.NewDirEntry(
		"plain",
		[]byte(name),
		domain.NewDirPayload(
			domain.NewKey(name, instance),
			nil,
			nil,
			nil,
			pendingMaps,
		),
	)
}

func newCustomPlainEntry(
	idx string,
	name string,
	instance string,
	pending bool,
	pool int,
	epoch int,
	exists bool,
	mtime time.Time,
	etag string,
	tag string,
	flags int,
) *domain.PlainEntry {
	var pendingMaps []*domain.PendingMap
	if pending {
		pendingMaps = []*domain.PendingMap{nil}
	}

	return domain.NewPlainEntry(
		domain.NewDirEntry(
			"plain",
			[]byte(idx),
			domain.NewDirPayload(
				domain.NewKey(name, instance),
				domain.NewDirVersionInfo(
					domain.NewVersion(pool, epoch),
					0,
				),
				domain.NewDirState(
					"",
					exists,
					tag,
					flags,
				),
				domain.NewMeta(
					domain.NewObjectSpec(1, 4, 4, false),
					domain.NewAuditInfo(mtime, etag),
					domain.NewContentInfo("", ""),
					domain.NewOwner("", ""),
				),
				pendingMaps,
			),
		),
	)
}

func sampleMTime() time.Time {
	return time.Date(2026, time.March, 6, 3, 34, 11, 918188000, time.UTC)
}

func newOLHEntry(name string, pending bool) *domain.OLHEntry {
	var pendingLogs []*domain.PendingLog
	if pending {
		pendingLogs = []*domain.PendingLog{nil}
	}

	return domain.NewOLHEntry(
		"olh",
		[]byte(name),
		domain.NewOLHPayload(
			domain.NewKey(name, "v1"),
			nil,
			0,
			pendingLogs,
			"",
		),
	)
}
