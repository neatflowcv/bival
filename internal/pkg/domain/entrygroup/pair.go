package entrygroup

import "github.com/neatflowcv/bival/internal/pkg/domain"

type Pair struct {
	plain    *domain.Plain
	instance *domain.Instance
}

func NewPair(plain *domain.Plain, instance *domain.Instance) *Pair {
	// Both values being nil indicates a programmer error, so panic intentionally.
	if plain == nil && instance == nil {
		panic("entrygroup.NewPair: plain and instance cannot both be nil")
	}

	return &Pair{
		plain:    plain,
		instance: instance,
	}
}

func (p *Pair) Plain() *domain.Plain {
	return p.plain
}

func (p *Pair) Instance() *domain.Instance {
	return p.instance
}

func (p *Pair) Version() string {
	if p.plain != nil {
		return p.plain.Instance()
	}

	return p.instance.Instance()
}

func (p *Pair) IsSoftDeleted() bool {
	return p.instance.IsSoftDeleted()
}

func (p *Pair) MTime() string {
	if p.plain != nil {
		return p.plain.MTime()
	}

	return p.instance.MTime()
}
