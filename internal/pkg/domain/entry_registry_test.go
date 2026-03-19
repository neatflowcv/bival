package domain_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestEntryRegistryAddGroupsByName(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "")
	require.NoError(t, err)

	instance, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1")
	require.NoError(t, err)

	olh, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1")
	require.NoError(t, err)

	instance2, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-2")
	require.NoError(t, err)

	// Act
	registry.Add(plain)
	registry.Add(instance)
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

	entry, err := domain.NewEntry(domain.KindPlain, "test.txt", "")
	require.NoError(t, err)
	registry.Add(entry)

	// Act
	err = registry.Validate()

	// Assert
	require.NoError(t, err)
}

func TestEntryRegistryValidateRejectsInstanceOnly(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "")
	require.NoError(t, err)

	entry, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1")
	require.NoError(t, err)
	registry.Add(plain)
	registry.Add(entry)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "instance entries require an olh entry")
}

func TestEntryRegistryValidateRejectsMissingPlain(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	instance, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1")
	require.NoError(t, err)

	olh, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1")
	require.NoError(t, err)

	registry.Add(instance)
	registry.Add(olh)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "plain entries require at least one plain entry")
}

func TestEntryRegistryValidateRejectsOLHOnly(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "")
	require.NoError(t, err)

	entry, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1")
	require.NoError(t, err)
	registry.Add(plain)
	registry.Add(entry)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "olh entries require an instance entry")
}

func TestEntryRegistryValidateRejectsMoreThanOneOLH(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "")
	require.NoError(t, err)

	instance1, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1")
	require.NoError(t, err)

	instance2, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-2")
	require.NoError(t, err)

	olh1, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1")
	require.NoError(t, err)

	olh2, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-2")
	require.NoError(t, err)

	registry.Add(plain)
	registry.Add(instance1)
	registry.Add(instance2)
	registry.Add(olh1)
	registry.Add(olh2)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "olh entries must be 0 or 1")
}

func TestEntryRegistryValidateRejectsInvalidTotalCount(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	plain, err := domain.NewEntry(domain.KindPlain, "test.txt", "")
	require.NoError(t, err)

	instance, err := domain.NewEntry(domain.KindInstance, "test.txt", "instance-1")
	require.NoError(t, err)

	olh, err := domain.NewEntry(domain.KindOLH, "test.txt", "instance-1")
	require.NoError(t, err)

	registry.Add(plain)
	registry.Add(instance)
	registry.Add(olh)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "total entries must be 1 or 4 plus an even offset")
}

func TestEntryRegistryValidateCollectsAllSetErrors(t *testing.T) {
	t.Parallel()

	// Arrange
	registry := domain.NewEntryRegistry()

	group1Plain, err := domain.NewEntry(domain.KindPlain, "group1.txt", "")
	require.NoError(t, err)

	group1Instance, err := domain.NewEntry(domain.KindInstance, "group1.txt", "instance-1")
	require.NoError(t, err)

	group2Instance, err := domain.NewEntry(domain.KindInstance, "group2.txt", "instance-1")
	require.NoError(t, err)

	group2OLH, err := domain.NewEntry(domain.KindOLH, "group2.txt", "instance-1")
	require.NoError(t, err)

	registry.Add(group1Plain)
	registry.Add(group1Instance)
	registry.Add(group2Instance)
	registry.Add(group2OLH)

	// Act
	err = registry.Validate()

	// Assert
	require.ErrorContains(t, err, "instance entries require an olh entry")
	require.ErrorContains(t, err, "plain entries require at least one plain entry")
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
