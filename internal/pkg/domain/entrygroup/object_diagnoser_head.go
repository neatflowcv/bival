package entrygroup

import "strconv"

type headDiagnoser struct{}

func (headDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	headCount := 0
	invalidHeadCount := 0

	for _, entry := range group.PlainEntries() {
		if entry.IsPlaceholder() {
			headCount++

			continue
		}

		if isVersionedHeadCandidate(entry) {
			invalidHeadCount++
		}
	}

	switch {
	case headCount == 0 && invalidHeadCount > 0:
		return []*Issue{
			newIssue(
				issueCodeInvalidVersionedHead,
				map[string]string{
					"invalid_head_count": strconv.Itoa(invalidHeadCount),
				},
			),
		}
	case headCount == 0:
		return []*Issue{
			newIssue(
				issueCodeMissingVersionedHead,
				nil,
			),
		}
	case headCount > 1:
		return []*Issue{
			newIssue(
				issueCodeDuplicateVersionedHead,
				map[string]string{
					"head_count": strconv.Itoa(headCount),
				},
			),
		}
	default:
		return nil
	}
}
