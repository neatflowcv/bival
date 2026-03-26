package domain_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

func TestVersionIsMissing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		version *domain.Version
		want    bool
	}{
		{
			name:    "missing version",
			version: domain.NewVersion(-1, 0),
			want:    true,
		},
		{
			name:    "defined version",
			version: domain.NewVersion(1, 1),
			want:    false,
		},
		{
			name:    "negative pool with epoch",
			version: domain.NewVersion(-1, 1),
			want:    false,
		},
		{
			name:    "zero version is not missing",
			version: domain.NewVersion(0, 0),
			want:    false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := testCase.version.IsMissing(); got != testCase.want {
				t.Fatalf("IsMissing() = %t, want %t", got, testCase.want)
			}
		})
	}
}
