package domain

// MyEntry is the shared validation contract for parsed entry models.
type MyEntry interface {
	Validate() error
}
