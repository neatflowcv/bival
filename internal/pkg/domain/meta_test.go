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

func TestMetaIsDefault(t *testing.T) {
	t.Parallel()

	require.True(t, newMetaForTest(domain.NewObjectSpec(0, 0, 0, false), time.Time{}, "", "").IsDefault())
	require.False(t, newMetaForTest(domain.NewObjectSpec(1, 0, 0, false), time.Time{}, "", "").IsDefault())
	require.False(t, newMetaForTest(domain.NewObjectSpec(0, 0, 0, false), time.Unix(1, 0), "", "").IsDefault())
	require.False(t, newMetaForTest(domain.NewObjectSpec(0, 0, 0, false), time.Time{}, "STANDARD", "").IsDefault())
	require.False(t, newMetaForTest(domain.NewObjectSpec(0, 0, 0, false), time.Time{}, "", "user").IsDefault())
	require.False(t, newMetaForTest(nil, time.Time{}, "", "").IsDefault())
}

func newMetaForTest(
	objectSpec *domain.ObjectSpec,
	mTime time.Time,
	storageClass string,
	ownerUserID string,
) *domain.Meta {
	return domain.NewMeta(
		objectSpec,
		mTime,
		"",
		storageClass,
		"",
		ownerUserID,
		"",
	)
}
