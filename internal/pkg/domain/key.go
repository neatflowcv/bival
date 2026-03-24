package domain

type Key struct {
	name     string
	instance string
}

func NewKey(name string, instance string) *Key {
	return &Key{
		name:     name,
		instance: instance,
	}
}
