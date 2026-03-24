package domain

// TODO: 나중에 Entry 로 바꿀거임. 이름이 겹치는게 있어서 일단 MyEntry 로 정의
type MyEntry interface {
	Validate() error
}
