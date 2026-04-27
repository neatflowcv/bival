package entrygroup_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain/entrygroup"
	"github.com/stretchr/testify/require"
)

func TestEntryGroupIssuesDescribeMissingMatchingInstance(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1")))
	group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v2")))
	group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1")))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))

	issues := group.Issues()
	require.Len(t, issues, 1)
	require.Equal(t, "pair.instance.missing", issues[0].Code())
	require.Equal(t, map[string]string{"version": "v2"}, issues[0].Meta())
}

func TestEntryGroupIssuesDescribeInvalidOLHReference(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1")))
	group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1")))
	group.AddOLH(newVersionedOLHEntry("alpha", "missing", false))

	issues := group.Issues()
	require.Len(t, issues, 1)
	require.Equal(t, "olh.reference.invalid", issues[0].Code())
	require.Equal(t, map[string]string{"referenced_version": "missing"}, issues[0].Meta())
}

func TestEntryGroupIssuesDescribeOutdatedOLHReference(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	versionOne := defaultVersionedFixture("v1")
	versionOne.mtime = olderMTime

	versionTwo := defaultVersionedFixture("v2")
	versionTwo.mtime = newerMTime

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(versionOne))
	group.AddPlain(newVersionedPlainEntry(versionTwo))
	group.AddInstance(newVersionedInstanceEntry(versionOne))
	group.AddInstance(newVersionedInstanceEntry(versionTwo))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))

	issues := group.Issues()
	require.Len(t, issues, 1)
	require.Equal(t, "olh.reference.outdated", issues[0].Code())
	require.Equal(
		t,
		map[string]string{
			"referenced_version": "v1",
			"version":            "v2",
		},
		issues[0].Meta(),
	)
}

func TestEntryGroupIssuesMetaIsDefensivelyCopied(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1")))
	group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v2")))
	group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1")))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))

	issues := group.Issues()
	require.Len(t, issues, 1)

	meta := issues[0].Meta()
	meta["version"] = "mutated"

	require.Equal(t, map[string]string{"version": "v2"}, issues[0].Meta())
}

func TestEntryGroupIssuesDescribePendingEntryMeta(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")

	group.AddPlain(newVersionedHeadPlainEntry())
	group.AddPlain(newVersionedPlainEntry(fixtureWithPendingMap(defaultVersionedFixture("v1"), true)))
	group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1")))
	group.AddOLH(newVersionedOLHEntry("alpha", "v1", false))

	issues := group.Issues()
	require.Len(t, issues, 2)
	require.Equal(t, "plain.pending.exists", issues[0].Code())
	require.Equal(
		t,
		map[string]string{
			"instance": "v1",
			"index":    "alpha\x00v913\x00iv1",
		},
		issues[0].Meta(),
	)
}
