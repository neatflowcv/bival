package domain_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestEntryRegistryAddGroupsByName(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	version1Plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	olh, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-2", time.Time{})
	require.NoError(t, err)

	version2Plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance2, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-2", time.Date(2026, 3, 19, 6, 28, 46, 0, time.UTC))
	require.NoError(t, err)

	// Act
	registry.Add(plain)
	registry.Add(version1Plain)
	registry.Add(instance)
	registry.Add(version2Plain)
	registry.Add(instance2)
	registry.Add(olh)

	err = registry.Validate()

	// Assert
	require.NoError(t, err)
}

func TestEntryRegistryValidateAllowsSinglePlainEntry(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	entry, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)
	registry.Add(entry)

	// Act
	err = registry.Validate()

	// Assert
	require.NoError(t, err)
}

func TestEntryRegistryValidateRejectsVersionedSetWithoutOLH(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	entry, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)
	registry.Add(plain)
	registry.Add(entry)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "versioned set must contain exactly 1 olh entry")
}

func TestEntryRegistryValidateRejectsVersionedSetWithoutHeadPlain(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	instance, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	olh, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1", time.Time{})
	require.NoError(t, err)

	registry.Add(instance)
	registry.Add(olh)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "versioned set must contain exactly 1 head plain entry plus 1 plain entry per instance entry")
}

func TestEntryRegistryValidateRejectsVersionedSetWithoutVersionPlain(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	entry, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1", time.Time{})
	require.NoError(t, err)
	registry.Add(plain)
	registry.Add(instance)
	registry.Add(entry)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "versioned set must contain exactly 1 head plain entry plus 1 plain entry per instance entry")
}

func TestEntryRegistryValidateRejectsMoreThanOneOLH(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	version1Plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance1, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	version2Plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance2, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-2", time.Date(2026, 3, 19, 6, 28, 46, 0, time.UTC))
	require.NoError(t, err)

	olh1, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1", time.Time{})
	require.NoError(t, err)

	olh2, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-2", time.Time{})
	require.NoError(t, err)

	registry.Add(plain)
	registry.Add(version1Plain)
	registry.Add(instance1)
	registry.Add(version2Plain)
	registry.Add(instance2)
	registry.Add(olh1)
	registry.Add(olh2)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "versioned set must contain exactly 1 olh entry")
}

func TestEntryRegistryValidateRejectsNonVersionedSetWithMultiplePlainEntries(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	anotherPlain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	registry.Add(plain)
	registry.Add(anotherPlain)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "non-versioned set must contain exactly 1 plain entry")
}

func TestEntryRegistryValidateCollectsAllSetErrors(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	group1Plain, err := domain.NewEntry(domain.KindPlain, "group1.txt", "", time.Time{})
	require.NoError(t, err)

	group1Instance, err := domain.NewEntry(domain.KindInstance, "group1.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	group2Instance, err := domain.NewEntry(domain.KindInstance, "group2.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	group2OLH, err := domain.NewEntry(domain.KindOLH, "group2.txt", "instance-1", time.Time{})
	require.NoError(t, err)

	registry.Add(group1Plain)
	registry.Add(group1Instance)
	registry.Add(group2Instance)
	registry.Add(group2OLH)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "versioned set must contain exactly 1 olh entry")
	require.ErrorContains(t, err, "versioned set must contain exactly 1 head plain entry plus 1 plain entry per instance entry")
}

func TestEntryRegistryValidateRejectsOLHWithoutMatchingInstance(t *testing.T) {
	t.Parallel()

	registry := domain.NewEntryRegistry()

	headPlain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	versionPlain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	olh, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-2", time.Time{})
	require.NoError(t, err)

	registry.Add(headPlain)
	registry.Add(versionPlain)
	registry.Add(instance)
	registry.Add(olh)

	err = registry.Validate()

	require.ErrorContains(t, err, "versioned set olh must reference an existing instance entry")
}

func TestEntryRegistryValidateRejectsOLHThatDoesNotReferenceLatestInstance(t *testing.T) {
	t.Parallel()

	registry := domain.NewEntryRegistry()

	headPlain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	version1Plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance1, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1", time.Date(2026, 3, 19, 6, 28, 45, 0, time.UTC))
	require.NoError(t, err)

	version2Plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "", time.Time{})
	require.NoError(t, err)

	instance2, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-2", time.Date(2026, 3, 19, 6, 28, 46, 0, time.UTC))
	require.NoError(t, err)

	olh, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1", time.Time{})
	require.NoError(t, err)

	registry.Add(headPlain)
	registry.Add(version1Plain)
	registry.Add(instance1)
	registry.Add(version2Plain)
	registry.Add(instance2)
	registry.Add(olh)

	err = registry.Validate()

	require.ErrorContains(t, err, "versioned set olh must reference the latest instance entry")
}

func TestEntryRegistryAddIgnoresNilEntry(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	// Act
	registry.Add(nil)

	// Assert
	require.NoError(t, registry.Validate())
}
