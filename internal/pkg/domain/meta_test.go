package domain_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestObjectSpecIsDefault(t *testing.T) {
	t.Parallel()

	require.True(t, domain.NewObjectSpec(0, 0, 0, false).IsDefault())
	require.False(t, domain.NewObjectSpec(1, 0, 0, false).IsDefault())
	require.False(t, domain.NewObjectSpec(0, 1, 0, false).IsDefault())
	require.False(t, domain.NewObjectSpec(0, 0, 1, false).IsDefault())
	require.False(t, domain.NewObjectSpec(0, 0, 0, true).IsDefault())
}

func TestAuditInfoIsDefault(t *testing.T) {
	t.Parallel()

	require.True(t, domain.NewAuditInfo(time.Time{}, "").IsDefault())
	require.False(t, domain.NewAuditInfo(time.Unix(1, 0), "").IsDefault())
	require.False(t, domain.NewAuditInfo(time.Time{}, "etag").IsDefault())
}

func TestContentInfoIsDefault(t *testing.T) {
	t.Parallel()

	require.True(t, domain.NewContentInfo("", "").IsDefault())
	require.False(t, domain.NewContentInfo("STANDARD", "").IsDefault())
	require.False(t, domain.NewContentInfo("", "text/plain").IsDefault())
}

func TestMetaIsDefault(t *testing.T) {
	t.Parallel()

	defaultMeta := domain.NewMeta(
		domain.NewObjectSpec(0, 0, 0, false),
		domain.NewAuditInfo(time.Time{}, ""),
		domain.NewContentInfo("", ""),
		"",
		"",
	)
	nonDefaultMeta := domain.NewMeta(
		domain.NewObjectSpec(1, 0, 0, false),
		domain.NewAuditInfo(time.Time{}, ""),
		domain.NewContentInfo("", ""),
		"",
		"",
	)
	nonDefaultOwnerMeta := domain.NewMeta(
		domain.NewObjectSpec(0, 0, 0, false),
		domain.NewAuditInfo(time.Time{}, ""),
		domain.NewContentInfo("", ""),
		"user",
		"",
	)
	missingPartMeta := domain.NewMeta(
		nil,
		domain.NewAuditInfo(time.Time{}, ""),
		domain.NewContentInfo("", ""),
		"",
		"",
	)

	require.True(t, defaultMeta.IsDefault())
	require.False(t, nonDefaultMeta.IsDefault())
	require.False(t, nonDefaultOwnerMeta.IsDefault())
	require.False(t, missingPartMeta.IsDefault())
}
