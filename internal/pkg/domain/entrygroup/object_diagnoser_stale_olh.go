package entrygroup

import (
	"strconv"
	"time"
)

type staleOLHDiagnoser struct {
	now time.Time
}

func (d staleOLHDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	olh, pairs, instanceByKey := versionedOLHContext(group)
	if olh == nil || pairs == nil || instanceByKey == nil {
		return nil
	}

	referencedInstance, ok := instanceByKey[olh.Instance()]
	if !ok || !isStaleInstance(referencedInstance, d.now) {
		return nil
	}

	referencedVersion := olh.Instance()

	if olh.DeleteMarker() {
		if len(pairs) == 0 {
			return nil
		}

		return []*Issue{
			newIssue(
				issueCodeStaleDeleteMarkerOLH,
				map[string]string{
					"referenced_version": referencedVersion,
					"pair_count":         strconv.Itoa(len(pairs)),
				},
			),
		}
	}

	if len(pairs) <= 1 {
		return nil
	}

	return []*Issue{
		newIssue(
			issueCodeStaleOLHReference,
			map[string]string{
				"referenced_version": referencedVersion,
				"pair_count":         strconv.Itoa(len(pairs)),
			},
		),
	}
}
