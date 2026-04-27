package entrygroup

import (
	"sort"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

type pairDiagnoser struct{}

func (pairDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	pairedPlainEntries, collectIssue := collectVersionedPlainEntries(group.PlainEntries())
	if collectIssue != nil {
		return nil
	}

	plainByKey, reason := buildPlainEntryMap(pairedPlainEntries)
	if reason != "" {
		return nil
	}

	instanceByKey, reason := buildInstanceEntryMap(group.InstanceEntries())
	if reason != "" {
		return nil
	}

	keys := unionVersionKeys(plainByKey, instanceByKey)
	issues := make([]*Issue, 0)

	for _, key := range keys {
		plainEntry := plainByKey[key]
		instanceEntry := instanceByKey[key]

		switch {
		case plainEntry == nil:
			issues = append(issues, newIssue(
				issueCodeMissingMatchingPlain,
				map[string]string{
					"version": key,
				},
			))
		case instanceEntry == nil:
			issues = append(issues, newIssue(
				issueCodeMissingMatchingInstance,
				map[string]string{
					"version": key,
				},
			))
		case !domain.IsVersionPair(plainEntry, instanceEntry):
			issues = append(issues, newIssue(
				issueCodeMismatchedVersionPair,
				map[string]string{
					"version": key,
				},
			))
		}
	}

	if len(issues) == 0 {
		return nil
	}

	return issues
}

func unionVersionKeys(
	plainByKey map[string]*domain.Plain,
	instanceByKey map[string]*domain.Instance,
) []string {
	keys := make([]string, 0, max(len(plainByKey), len(instanceByKey)))
	seen := make(map[string]struct{}, max(len(plainByKey), len(instanceByKey)))

	for key := range plainByKey {
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		keys = append(keys, key)
	}

	for key := range instanceByKey {
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		keys = append(keys, key)
	}

	sort.Strings(keys)

	return keys
}
