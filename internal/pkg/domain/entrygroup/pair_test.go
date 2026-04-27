package entrygroup_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/neatflowcv/bival/internal/pkg/domain/entrygroup"
	"github.com/stretchr/testify/require"
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
	instanceParams.Instance = "ver-2"
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

func TestPairMTime_UsesPlainFirst(t *testing.T) {
	t.Parallel()

	// Arrange
	plain := domain.NewPlain(versionedPairParams())
	instanceParams := versionedPairParams()
	instanceParams.MTime = "2026-04-28T00:00:00Z"
	instance := domain.NewInstance(instanceParams)
	pair := entrygroup.NewPair(plain, instance)

	// Act
	mtime := pair.MTime()

	// Assert
	require.Equal(t, "2026-04-27T00:00:00Z", mtime)
}

func TestPairMTime_UsesInstanceWhenPlainMissing(t *testing.T) {
	t.Parallel()

	// Arrange
	instance := domain.NewInstance(versionedPairParams())
	pair := entrygroup.NewPair(nil, instance)

	// Act
	mtime := pair.MTime()

	// Assert
	require.Equal(t, "2026-04-27T00:00:00Z", mtime)
}

func TestPairMTime_ReturnsEmptyWhenEntriesMissing(t *testing.T) {
	t.Parallel()

	// Arrange
	pair := entrygroup.NewPair(nil, nil)

	// Act
	mtime := pair.MTime()

	// Assert
	require.Empty(t, mtime)
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
	pairs, err := entrygroup.NewPairsByGroup(group)

	// Assert
	require.NoError(t, err)
	require.Len(t, pairs, 1)
	require.Same(t, plain, pairs[0].Plain())
	require.Same(t, instance, pairs[0].Instance())
}

func TestNewPairsByGroup_AllowsMissingMatchingInstance(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")
	group.AddPlain(domain.NewPlain(versionedHeadParams()))
	group.AddPlain(domain.NewPlain(versionedPairParams()))

	// Act
	pairs, err := entrygroup.NewPairsByGroup(group)

	// Assert
	require.NoError(t, err)
	require.Len(t, pairs, 1)
	require.NotNil(t, pairs[0].Plain())
	require.Nil(t, pairs[0].Instance())
}

func TestNewPairsByGroup_SortsByMTime(t *testing.T) {
	t.Parallel()

	// Arrange
	group := entrygroup.New("alpha")
	group.AddPlain(domain.NewPlain(versionedHeadParams()))

	olderParams := versionedPairParams()
	olderParams.Index = []byte("alpha:ver-1")
	olderParams.Instance = "ver-1"
	olderParams.MTime = "2026-04-27T00:00:00Z"

	newerParams := versionedPairParams()
	newerParams.Index = []byte("alpha:ver-2")
	newerParams.Instance = "ver-2"
	newerParams.MTime = "2026-04-28T00:00:00Z"

	group.AddPlain(domain.NewPlain(newerParams))
	group.AddInstance(domain.NewInstance(newerParams))
	group.AddPlain(domain.NewPlain(olderParams))
	group.AddInstance(domain.NewInstance(olderParams))

	// Act
	pairs, err := entrygroup.NewPairsByGroup(group)

	// Assert
	require.NoError(t, err)
	require.Len(t, pairs, 2)
	require.Equal(t, "2026-04-27T00:00:00Z", pairs[0].MTime())
	require.Equal(t, "2026-04-28T00:00:00Z", pairs[1].MTime())
}

func versionedHeadParams() domain.DirEntryParams {
	return domain.DirEntryParams{
		Kind:             "",
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
		Kind:             "",
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
		MTime:            "2026-04-27T00:00:00Z",
		ETag:             "etag",
		StorageClass:     "",
		ContentType:      "",
		OwnerUserID:      "",
		OwnerDisplayName: "",
		PendingMaps:      nil,
	}
}
