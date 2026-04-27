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
