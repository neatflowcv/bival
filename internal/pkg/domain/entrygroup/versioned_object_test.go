package entrygroup_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain/entrygroup"
	"github.com/stretchr/testify/require"
)

func TestVersionedObjectStoresPlaceholderPairsAndOLH(t *testing.T) {
	t.Parallel()

	// Arrange
	placeholder := newVersionedHeadPlainEntry()
	version := defaultVersionedFixture("v1")
	pair := &entrygroup.Pair{
		Plain:    newVersionedPlainEntry(version),
		Instance: newVersionedInstanceEntry(version),
	}
	olh := newVersionedOLHEntry("alpha", "v1", false)

	state := entrygroup.VersionedObject{
		Placeholder: placeholder,
		Pairs:       map[string]*entrygroup.Pair{"v1": pair},
		OLH:         olh,
	}

	// Act
	gotPair := state.Pairs["v1"]

	// Assert
	require.Same(t, placeholder, state.Placeholder)
	require.Len(t, state.Pairs, 1)
	require.Same(t, pair, gotPair)
	require.Same(t, pair.Plain, gotPair.Plain)
	require.Same(t, pair.Instance, gotPair.Instance)
	require.Same(t, olh, state.OLH)
}

func TestNewVersionedObjectBuildsVersionedObjectFromGroup(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")
	versionOne := defaultVersionedFixture("v1")
	versionTwo := defaultVersionedFixture("v2")
	versionTwo.pool = 187
	versionTwo.epoch = 1148
	versionTwo.versionedEpoch = 3
	versionTwo.tag = "tag-v2"

	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionOne)))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(versionTwo)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionOne)))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(versionTwo)))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "v1", false)))

	object, err := entrygroup.NewVersionedObject(group)

	require.NoError(t, err)
	require.NotNil(t, object)
	require.Same(t, group.PlainEntries()[0], object.Placeholder)
	require.Same(t, group.OLHEntries()[0], object.OLH)
	require.Len(t, object.Pairs, 2)
	require.Same(t, group.PlainEntries()[1], object.Pairs["v1"].Plain)
	require.Same(t, group.InstanceEntries()[0], object.Pairs["v1"].Instance)
	require.Same(t, group.PlainEntries()[2], object.Pairs["v2"].Plain)
	require.Same(t, group.InstanceEntries()[1], object.Pairs["v2"].Instance)
}

func TestNewVersionedObjectRejectsUnversionedGroup(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")
	require.NoError(t, group.AddPlain(newUnversionedPlainEntry()))

	object, err := entrygroup.NewVersionedObject(group)

	require.EqualError(t, err, "entry group is not a versioned object")
	require.Nil(t, object)
}

func TestNewVersionedObjectRejectsInvalidVersionedGroup(t *testing.T) {
	t.Parallel()

	group := entrygroup.New("alpha")
	require.NoError(t, group.AddPlain(newVersionedHeadPlainEntry()))
	require.NoError(t, group.AddPlain(newVersionedPlainEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddInstance(newVersionedInstanceEntry(defaultVersionedFixture("v1"))))
	require.NoError(t, group.AddOLH(newVersionedOLHEntry("alpha", "missing", false)))

	object, err := entrygroup.NewVersionedObject(group)

	require.EqualError(t, err, "olh references missing instance")
	require.Nil(t, object)
}
