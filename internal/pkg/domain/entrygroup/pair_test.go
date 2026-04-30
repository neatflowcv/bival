package entrygroup_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/neatflowcv/bival/internal/pkg/domain/entrygroup"
	"github.com/stretchr/testify/require"
)

const (
	versionTwo = "ver-2"
	olderMTime = "2026-04-27T00:00:00Z"
	newerMTime = "2026-04-28T00:00:00Z"
)

func TestNewPair(t *testing.T) {
	t.Parallel()

	// Arrange
	plain := domain.NewPlain(versionedPairParams())

	instanceParams := versionedPairParams()
	instanceParams.Exists = false
	instance := domain.NewInstance(instanceParams)

	// Act
	pair := entrygroup.NewPair(plain, instance)

	// Assert
	require.Same(t, plain, pair.Plain())
	require.Same(t, instance, pair.Instance())
}

func TestNewPair_AllowsDifferentInstance(t *testing.T) {
	t.Parallel()

	// Arrange
	plain := domain.NewPlain(versionedPairParams())

	instanceParams := versionedPairParams()
	instanceParams.Instance = versionTwo
	instance := domain.NewInstance(instanceParams)

	// Act
	pair := entrygroup.NewPair(plain, instance)

	// Assert
	require.Same(t, plain, pair.Plain())
	require.Same(t, instance, pair.Instance())
}

func TestNewPair_AllowsDifferentStateWhenNameAndInstanceMatch(t *testing.T) {
	t.Parallel()

	// Arrange
	plain := domain.NewPlain(versionedPairParams())

	instanceParams := versionedPairParams()
	instanceParams.Epoch++
	instanceParams.Exists = false
	instance := domain.NewInstance(instanceParams)

	// Act
	pair := entrygroup.NewPair(plain, instance)

	// Assert
	require.Same(t, plain, pair.Plain())
	require.Same(t, instance, pair.Instance())
}

func TestNewPair_PanicsWhenEntriesMissing(t *testing.T) {
	t.Parallel()

	require.PanicsWithValue(
		t,
		"entrygroup.NewPair: plain and instance cannot both be nil",
		func() {
			entrygroup.NewPair(nil, nil)
		},
	)
}

func TestPairMTime_UsesPlainFirst(t *testing.T) {
	t.Parallel()

	// Arrange
	plain := domain.NewPlain(versionedPairParams())
	instanceParams := versionedPairParams()
	instanceParams.MTime = newerMTime
	instance := domain.NewInstance(instanceParams)
	pair := entrygroup.NewPair(plain, instance)

	// Act
	mtime := pair.MTime()

	// Assert
	require.Equal(t, olderMTime, mtime)
}

func TestPairMTime_UsesInstanceWhenPlainMissing(t *testing.T) {
	t.Parallel()

	// Arrange
	instance := domain.NewInstance(versionedPairParams())
	pair := entrygroup.NewPair(nil, instance)

	// Act
	mtime := pair.MTime()

	// Assert
	require.Equal(t, olderMTime, mtime)
}

func TestPairVersion_UsesPlainFirst(t *testing.T) {
	t.Parallel()

	// Arrange
	plain := domain.NewPlain(versionedPairParams())
	instanceParams := versionedPairParams()
	instanceParams.Instance = versionTwo
	instance := domain.NewInstance(instanceParams)
	pair := entrygroup.NewPair(plain, instance)

	// Act
	version := pair.Version()

	// Assert
	require.Equal(t, "ver-1", version)
}

func TestPairVersion_UsesInstanceWhenPlainMissing(t *testing.T) {
	t.Parallel()

	// Arrange
	instance := domain.NewInstance(versionedPairParams())
	pair := entrygroup.NewPair(nil, instance)

	// Act
	version := pair.Version()

	// Assert
	require.Equal(t, "ver-1", version)
}

func TestNewPairsByGroup(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")
	group.AddPlain(domain.NewPlain(versionedHeadParams()))

	plain := domain.NewPlain(versionedPairParams())
	instance := domain.NewInstance(versionedPairParams())

	group.AddPlain(plain)
	group.AddInstance(instance)

	// Act
	pairs := entrygroup.NewPairsByGroup(group)
	items := pairs.Items()

	// Assert
	require.Len(t, items, 1)
	require.Same(t, plain, items[0].Plain())
	require.Same(t, instance, items[0].Instance())
}

func TestNewPairsByGroup_AllowsMissingMatchingInstance(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")
	group.AddPlain(domain.NewPlain(versionedHeadParams()))
	group.AddPlain(domain.NewPlain(versionedPairParams()))

	// Act
	pairs := entrygroup.NewPairsByGroup(group)
	items := pairs.Items()

	// Assert
	require.Len(t, items, 1)
	require.NotNil(t, items[0].Plain())
	require.Nil(t, items[0].Instance())
}

func TestNewPairsByGroup_SortsByMTime(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")
	group.AddPlain(domain.NewPlain(versionedHeadParams()))

	olderParams := versionedPairParams()
	olderParams.Index = []byte("alpha:ver-1")
	olderParams.Instance = "ver-1"
	olderParams.MTime = olderMTime

	newerParams := versionedPairParams()
	newerParams.Index = []byte("alpha:ver-2")
	newerParams.Instance = versionTwo
	newerParams.MTime = newerMTime

	group.AddPlain(domain.NewPlain(newerParams))
	group.AddInstance(domain.NewInstance(newerParams))
	group.AddPlain(domain.NewPlain(olderParams))
	group.AddInstance(domain.NewInstance(olderParams))

	// Act
	pairs := entrygroup.NewPairsByGroup(group)
	items := pairs.Items()

	// Assert
	require.Len(t, items, 2)
	require.Equal(t, olderMTime, items[0].MTime())
	require.Equal(t, newerMTime, items[1].MTime())
}

func TestNewPairs_SortsItemsByMTime(t *testing.T) {
	t.Parallel()

	// Arrange
	olderParams := versionedPairParams()
	olderParams.Index = []byte("alpha:ver-1")
	olderParams.Instance = "ver-1"
	olderParams.MTime = olderMTime

	newerParams := versionedPairParams()
	newerParams.Index = []byte("alpha:ver-2")
	newerParams.Instance = versionTwo
	newerParams.MTime = newerMTime

	newerPair := entrygroup.NewPair(
		domain.NewPlain(newerParams),
		domain.NewInstance(newerParams),
	)
	olderPair := entrygroup.NewPair(
		domain.NewPlain(olderParams),
		domain.NewInstance(olderParams),
	)

	// Act
	pairs := entrygroup.NewPairs([]*entrygroup.Pair{newerPair, olderPair})
	items := pairs.Items()

	// Assert
	require.Len(t, items, 2)
	require.Same(t, olderPair, items[0])
	require.Same(t, newerPair, items[1])
}

func TestPairsPairByVersion(t *testing.T) {
	t.Parallel()

	// Arrange
	matchingParams := versionedPairParams()
	matchingPair := entrygroup.NewPair(
		domain.NewPlain(matchingParams),
		domain.NewInstance(matchingParams),
	)
	missingParams := versionedPairParams()
	missingParams.Index = []byte("alpha:ver-3")
	missingParams.Instance = "ver-3"
	missingInstancePair := entrygroup.NewPair(
		domain.NewPlain(missingParams),
		nil,
	)
	pairs := entrygroup.NewPairs([]*entrygroup.Pair{missingInstancePair, matchingPair})

	// Act
	pair, ok := pairs.PairByVersion("ver-1")

	// Assert
	require.True(t, ok)
	require.Same(t, matchingPair, pair)
}

func TestPairsPairByVersion_ReturnsFalseWhenMissing(t *testing.T) {
	t.Parallel()

	// Arrange
	pairs := entrygroup.NewPairs(nil)

	// Act
	pair, ok := pairs.PairByVersion("ver-404")

	// Assert
	require.False(t, ok)
	require.Nil(t, pair)
}

func versionedHeadParams() domain.DirEntryParams {
	return domain.DirEntryParams{
		Index:            []byte("alpha"),
		Name:             "alpha",
		Instance:         "",
		Pool:             -1,
		Epoch:            0,
		VEpoch:           0,
		Locator:          "",
		Exists:           false,
		Tag:              "",
		Flags:            8,
		Category:         0,
		Size:             0,
		AccountedSize:    0,
		Appendable:       false,
		MTime:            "0.000000",
		ETag:             "",
		StorageClass:     "",
		ContentType:      "",
		OwnerUserID:      "",
		OwnerDisplayName: "",
		PendingMaps:      nil,
	}
}

func versionedPairParams() domain.DirEntryParams {
	return domain.DirEntryParams{
		Index:            []byte("alpha:ver-1"),
		Name:             "alpha",
		Instance:         "ver-1",
		Pool:             1,
		Epoch:            2,
		VEpoch:           3,
		Locator:          "",
		Exists:           true,
		Tag:              "tag",
		Flags:            0,
		Category:         0,
		Size:             0,
		AccountedSize:    0,
		Appendable:       false,
		MTime:            olderMTime,
		ETag:             "etag",
		StorageClass:     "",
		ContentType:      "",
		OwnerUserID:      "",
		OwnerDisplayName: "",
		PendingMaps:      nil,
	}
}
