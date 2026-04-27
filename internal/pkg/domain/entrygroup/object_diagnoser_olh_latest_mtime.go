package entrygroup

import "time"

type olhLatestMTimeDiagnoser struct{}

func (olhLatestMTimeDiagnoser) Diagnose(group *EntryGroup) []*Issue {
	olh, err := extractOLH(group)
	if err != nil {
		return nil
	}

	pairs := NewPairsByGroup(group)

	olhPair, ok := pairs.PairByVersion(olh.Instance())
	if !ok {
		return nil
	}

	olhMTime, err := time.Parse(time.RFC3339Nano, olhPair.MTime())
	if err != nil {
		return nil
	}

	latestPair, ok := latestPairByMTime(pairs.Items())
	if !ok {
		return nil
	}

	latestMTime, err := time.Parse(time.RFC3339Nano, latestPair.MTime())
	if err != nil || olhMTime.Equal(latestMTime) {
		return nil
	}

	return []*Issue{
		newIssue(issueCodeOutdatedOLHReference, map[string]string{
			"referenced_version": olh.Instance(),
			"version":            latestPair.Version(),
		}),
	}
}

func latestPairByMTime(items []*Pair) (*Pair, bool) {
	var (
		latestPair *Pair
		latestTime time.Time
	)

	for _, pair := range items {
		mtime, err := time.Parse(time.RFC3339Nano, pair.MTime())
		if err != nil {
			return nil, false
		}

		if latestPair == nil || mtime.After(latestTime) {
			latestPair = pair
			latestTime = mtime
		}
	}

	if latestPair == nil {
		return nil, false
	}

	return latestPair, true
}
