package entrygroup

import "maps"

type Issue struct {
	code string
	meta map[string]string
}

func newIssue(code string, meta map[string]string) *Issue {
	return &Issue{
		code: code,
		meta: cloneIssueMeta(meta),
	}
}

func (i *Issue) Code() string {
	return i.code
}

func (i *Issue) Meta() map[string]string {
	return cloneIssueMeta(i.meta)
}

func cloneIssueMeta(meta map[string]string) map[string]string {
	if len(meta) == 0 {
		return nil
	}

	cloned := make(map[string]string, len(meta))
	maps.Copy(cloned, meta)

	return cloned
}
