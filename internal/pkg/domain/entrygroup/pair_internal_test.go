package entrygroup

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain"
	"github.com/stretchr/testify/require"
)

func TestPairIsFull(t *testing.T) {
	t.Parallel()

	params := domain.DirEntryParams{
		Kind:             "",
		Index:            nil,
		Name:             "",
		Instance:         "",
		Pool:             0,
		Epoch:            0,
		VEpoch:           0,
		Locator:          "",
		Exists:           false,
		Tag:              "",
		Flags:            0,
		Category:         0,
		Size:             0,
		AccountedSize:    0,
		Appendable:       false,
		MTime:            "",
		ETag:             "",
		StorageClass:     "",
		ContentType:      "",
		OwnerUserID:      "",
		OwnerDisplayName: "",
		PendingMaps:      nil,
	}

	fullPair := &Pair{
		Plain:    domain.NewPlain(params),
		Instance: domain.NewInstance(params),
	}
	plainOnlyPair := &Pair{
		Plain:    domain.NewPlain(params),
		Instance: nil,
	}
	instanceOnlyPair := &Pair{
		Plain:    nil,
		Instance: domain.NewInstance(params),
	}
	emptyPair := &Pair{
		Plain:    nil,
		Instance: nil,
	}

	require.True(t, fullPair.isFull())
	require.False(t, plainOnlyPair.isFull())
	require.False(t, instanceOnlyPair.isFull())
	require.False(t, emptyPair.isFull())
}
