package entrygroup

import (
	"fmt"
	"sort"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

type entryKeyDiagnoser struct{}

func (entryKeyDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	pairedPlainEntries, collectIssue := collectVersionedPlainEntries(group.PlainEntries())
	if collectIssue != nil {
		return nil
	}

	issues := make([]*Issue, 0)
	issues = append(issues, diagnoseDuplicateEntryKeysInPlain(pairedPlainEntries)...)
	issues = append(issues, diagnoseDuplicateEntryKeysInInstance(group.InstanceEntries())...)

	if len(issues) == 0 {
		return nil
	}

	return issues
}

func diagnoseDuplicateEntryKeysInPlain(entries []*domain.Plain) []*Issue {
	return diagnoseDuplicateEntryKeys(
		plainEntryKeys(entries),
		"plain",
	)
}

func diagnoseDuplicateEntryKeysInInstance(entries []*domain.Instance) []*Issue {
	return diagnoseDuplicateEntryKeys(
		instanceEntryKeys(entries),
		"instance",
	)
}

func diagnoseDuplicateEntryKeys(keys []entryKey, entryType string) []*Issue {
	counts := make(map[string]int, len(keys))
	issues := make([]*Issue, 0)

	for _, key := range keys {
		if !key.valid {
			issues = append(issues, newIssue(
				issueCodeInvalidEntryKey,
				map[string]string{
					"entry_type": entryType,
					"version":    key.version,
				},
			))

			continue
		}

		counts[key.version]++
	}

	duplicateVersions := make([]string, 0)
	for version, count := range counts {
		if count > 1 {
			duplicateVersions = append(duplicateVersions, version)
		}
	}
	sort.Strings(duplicateVersions)

	for _, version := range duplicateVersions {
		issues = append(issues, newIssue(
			issueCodeDuplicateEntryKey,
			map[string]string{
				"entry_type": entryType,
				"version":    version,
				"count":      fmt.Sprintf("%d", counts[version]),
			},
		))
	}

	if len(issues) == 0 {
		return nil
	}

	return issues
}

type entryKey struct {
	version string
	valid   bool
}

func plainEntryKeys(entries []*domain.Plain) []entryKey {
	keys := make([]entryKey, 0, len(entries))
	for _, entry := range entries {
		keys = append(keys, entryKey{
			version: entry.Instance(),
			valid:   entry.Name() != "",
		})
	}

	return keys
}

func instanceEntryKeys(entries []*domain.Instance) []entryKey {
	keys := make([]entryKey, 0, len(entries))
	for _, entry := range entries {
		keys = append(keys, entryKey{
			version: entry.Instance(),
			valid:   entry.Name() != "",
		})
	}

	return keys
}
