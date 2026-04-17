package domain_test

import (
	"testing"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

func TestDirVersionInfoIsMissing(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name        string
		versionInfo *domain.DirVersionInfo
		want        bool
	}{
		{
			name:        "missing version",
			versionInfo: domain.NewDirVersionInfo(-1, 0, 0),
			want:        true,
		},
		{
			name:        "defined version",
			versionInfo: domain.NewDirVersionInfo(1, 1, 0),
			want:        false,
		},
		{
			name:        "negative pool with epoch",
			versionInfo: domain.NewDirVersionInfo(-1, 1, 0),
			want:        false,
		},
		{
			name:        "zero version is not missing",
			versionInfo: domain.NewDirVersionInfo(0, 0, 0),
			want:        false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			if got := testCase.versionInfo.IsMissing(); got != testCase.want {
				t.Fatalf("IsMissing() = %t, want %t", got, testCase.want)
			}
		})
	}
}
