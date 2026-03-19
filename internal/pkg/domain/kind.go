package domain

type Kind string

const (
	KindPlain    Kind = "plain"
	KindInstance Kind = "instance"
	KindOLH      Kind = "olh"
)

func (k Kind) IsValid() bool {
	switch k {
	case KindPlain, KindInstance, KindOLH:
		return true
	default:
		return false
	}
}
