package domain

type Owner struct {
	userID      string
	displayName string
}

func NewOwner(userID string, displayName string) *Owner {
	return &Owner{
		userID:      userID,
		displayName: displayName,
	}
}

func (o *Owner) IsDefault() bool {
	return o.userID == "" && o.displayName == ""
}
