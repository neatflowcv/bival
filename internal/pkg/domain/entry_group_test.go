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
	require.Equal(t, "pending entry exists", group.ProblemReason())
}

func TestEntryGroupAddInstanceTracksPendingMap(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddInstance(newInstanceEntry("alpha", true))
	require.NoError(t, err)
	require.Equal(t, "pending entry exists", group.ProblemReason())
}

func TestEntryGroupAddOLHTracksPendingLog(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddOLH(newOLHEntry("alpha", true))
	require.NoError(t, err)
	require.Equal(t, "pending entry exists", group.ProblemReason())
}

func TestEntryGroupSinglePlainHasNoProblem(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	err := group.AddPlain(newPlainEntry("alpha", "", false))
	require.NoError(t, err)
	require.Empty(t, group.ProblemReason())
}

func TestEntryGroupReportsInvalidOLHCount(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "", false)))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.Equal(t, "versioning object must have exactly one olh", group.ProblemReason())
}

func TestEntryGroupReportsInvalidInstanceCount(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "", false)))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v2", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.NoError(t, group.AddOLH(newOLHEntry("alpha", false)))
	require.Equal(t, "versioning object must satisfy instance+1==plain", group.ProblemReason())
}

func TestEntryGroupAcceptsValidVersionedObject(t *testing.T) {
	t.Parallel()

	group := domain.NewEntryGroup("alpha")

	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "", false)))
	require.NoError(t, group.AddPlain(newPlainEntry("alpha", "v1", false)))
	require.NoError(t, group.AddInstance(newInstanceEntry("alpha", false)))
	require.NoError(t, group.AddOLH(newOLHEntry("alpha", false)))
	require.Empty(t, group.ProblemReason())
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
