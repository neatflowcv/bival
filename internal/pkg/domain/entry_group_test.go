package domain_test

import (
	"testing"

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

	err := group.AddPlain(newPlainEntry("alpha", "", false))
	require.NoError(t, err)
	require.Equal(t, domain.UnversionedObject, classifier.Classify(group))
}

func TestEntryGroupClassifierReturnsUnknownObjectWhenNoRuleMatches(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "", false)))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.Equal(t, domain.UnknownObject, classifier.Classify(group))
}

func TestEntryGroupClassifierReturnsUnknownObjectWhenVersionedRuleDoesNotMatch(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")
	classifier := domain.NewEntryGroupClassifier()

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "", false)))
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

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "", false)))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.NoError(t, group.AddOLH(newOLHEntry("alpha", false)))
	require.Equal(t, domain.VersionedObject, classifier.Classify(group))
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
	return domain.NewPlainEntry(newDirEntry(name, instance, pending))
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
