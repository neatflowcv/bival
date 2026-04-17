package entrygroup

import "errors"

type Rule interface {
	Check(group *EntryGroup) error
}

func newVersionedObjectRules() []Rule {
	return []Rule{
		versionedHeadRule{},
		versionedEntryKeyRule{},
		versionedPairRule{},
		versionedOLHRule{},
	}
}

func checkRules(group *EntryGroup, rules []Rule) error {
	errs := make([]error, 0)

	for _, rule := range rules {
		err := rule.Check(group)
		if err != nil {
			errs = append(errs, splitJoinedErrors(err)...)
		}
	}

	if len(errs) == 0 {
		return nil
	}

	return errors.Join(errs...)
}
