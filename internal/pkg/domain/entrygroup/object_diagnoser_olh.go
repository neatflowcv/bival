package entrygroup

import "strconv"

type olhDiagnoser struct{}

func (olhDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	olhEntries := group.OLHEntries()

	switch len(olhEntries) {
	case 0:
		return []*Issue{
			newIssue(
				issueCodeMissingOLH,
				nil,
			),
		}
	case 1:
	default:
		return []*Issue{
			newIssue(
				issueCodeInvalidOLH,
				map[string]string{
					"olh_count": strconv.Itoa(len(olhEntries)),
				},
			),
		}
	}

	olh := olhEntries[0]
	if olh == nil {
		return []*Issue{
			newIssue(
				issueCodeInvalidOLH,
				nil,
			),
		}
	}

	if olh.Name() == "" {
		return []*Issue{
			newIssue(
				issueCodeInvalidOLH,
				map[string]string{
					"referenced_version": olh.Instance(),
				},
			),
		}
	}

	if olh.HasPendingLog() {
		return []*Issue{
			newIssue(
				issueCodeInvalidOLH,
				map[string]string{
					"referenced_version": olh.Instance(),
				},
			),
		}
	}

	instanceSet, ok := instanceNameSet(group.InstanceEntries())
	if !ok {
		return []*Issue{
			newIssue(
				issueCodeInvalidOLHReference,
				map[string]string{
					"referenced_version": olh.Instance(),
				},
			),
		}
	}

	referencedVersion := olh.Instance()
	if _, exists := instanceSet[referencedVersion]; exists {
		return nil
	}

	return []*Issue{
		newIssue(
			issueCodeInvalidOLHReference,
			map[string]string{
				"referenced_version": referencedVersion,
			},
		),
	}
}
