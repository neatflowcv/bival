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

func (k *Key) Name() string {
	if k == nil {
		return ""
	}

	return k.name
}

func (k *Key) Instance() string {
	if k == nil {
		return ""
	}

	return k.instance
}
