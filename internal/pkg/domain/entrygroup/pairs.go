package entrygroup

import (
	"cmp"
	"slices"

	"github.com/neatflowcv/bival/internal/pkg/domain"
)

type Pairs struct {
	items []*Pair
}

func NewPairs(items []*Pair) *Pairs {
	sortedItems := slices.Clone(items)
	slices.SortFunc(sortedItems, func(left *Pair, right *Pair) int {
		return cmp.Compare(left.MTime(), right.MTime())
	})

	return &Pairs{
		items: sortedItems,
	}
}

func NewPairsByGroup(group *EntryGroup) (*Pairs, error) {
	versionMap := map[string]struct{}{}
	plains := slices.DeleteFunc(group.PlainEntries(), func(entry *domain.Plain) bool {
		return entry.IsPlaceholder()
	})
	plainMap := map[string]*domain.Plain{}

	for _, plain := range plains {
		versionMap[plain.Instance()] = struct{}{}
		plainMap[plain.Instance()] = plain
	}

	instances := group.InstanceEntries()
	instanceMap := map[string]*domain.Instance{}

	for _, instance := range instances {
		versionMap[instance.Instance()] = struct{}{}
		instanceMap[instance.Instance()] = instance
	}

	var items []*Pair

	for version := range versionMap {
		pair := NewPair(plainMap[version], instanceMap[version])
		items = append(items, pair)
	}

	return NewPairs(items), nil
}

func (p *Pairs) Items() []*Pair {
	return slices.Clone(p.items)
}

func (p *Pairs) PairByVersion(instance string) (*Pair, bool) {
	for _, pair := range p.items {
		if pair.Version() == instance {
			return pair, true
		}
	}

	return nil, false
}
