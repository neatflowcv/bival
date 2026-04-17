package domain_test

import (
	"testing"
	"time"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestMetaIsDefault(t *testing.T) {
	t.Parallel()

	require.True(t, newMetaForTest(0, 0, 0, false, time.Time{}, "", "").IsDefault())
	require.False(t, newMetaForTest(1, 0, 0, false, time.Time{}, "", "").IsDefault())
	require.False(t, newMetaForTest(0, 1, 0, false, time.Time{}, "", "").IsDefault())
	require.False(t, newMetaForTest(0, 0, 1, false, time.Time{}, "", "").IsDefault())
	require.False(t, newMetaForTest(0, 0, 0, true, time.Time{}, "", "").IsDefault())
	require.False(t, newMetaForTest(0, 0, 0, false, time.Unix(1, 0), "", "").IsDefault())
	require.False(t, newMetaForTest(0, 0, 0, false, time.Time{}, "STANDARD", "").IsDefault())
	require.False(t, newMetaForTest(0, 0, 0, false, time.Time{}, "", "user").IsDefault())
}

func newMetaForTest(
	category int,
	size int64,
	accountedSize int64,
	appendable bool,
	mTime time.Time,
	storageClass string,
	ownerUserID string,
) *domain.Meta {
	return domain.NewMeta(
		category,
		size,
		accountedSize,
		appendable,
		mTime,
		"",
		storageClass,
		"",
		ownerUserID,
		"",
	)
}
