package entrygroup

import (
	"strconv"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

type olhDiagnoser struct{}

func (olhDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	olhEntries := group.OLHEntries()

	countIssue := diagnoseOLHCount(olhEntries)
	if countIssue != nil {
		return []*Issue{countIssue}
	}

	olh := olhEntries[0]

	olhIssue := diagnoseInvalidOLH(olh)
	if olhIssue != nil {
		return []*Issue{olhIssue}
	}

	referencedVersion := olh.Instance()

	instanceSet, ok := instanceNameSet(group.InstanceEntries())
	if !ok {
		return []*Issue{newInvalidOLHReferenceIssue(referencedVersion)}
	}

	if _, exists := instanceSet[referencedVersion]; exists {
		return nil
	}

	return []*Issue{
		newInvalidOLHReferenceIssue(referencedVersion),
	}
}

func diagnoseOLHCount(olhEntries []*domain.OLH) *Issue {
	switch len(olhEntries) {
	case 0:
		return newIssue(
			issueCodeMissingOLH,
			nil,
		)
	case 1:
		return nil
	default:
		return newIssue(
			issueCodeInvalidOLH,
			map[string]string{
				"olh_count": strconv.Itoa(len(olhEntries)),
			},
		)
	}
}

func diagnoseInvalidOLH(olh *domain.OLH) *Issue {
	if olh == nil {
		return newIssue(
			issueCodeInvalidOLH,
			nil,
		)
	}

	if olh.Name() != "" && !olh.HasPendingLog() {
		return nil
	}

	return newIssue(
		issueCodeInvalidOLH,
		map[string]string{
			"referenced_version": olh.Instance(),
		},
	)
}

func newInvalidOLHReferenceIssue(referencedVersion string) *Issue {
	return newIssue(
		issueCodeInvalidOLHReference,
		map[string]string{
			"referenced_version": referencedVersion,
		},
	)
}
