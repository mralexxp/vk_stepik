package models

type Profile struct {
	ID        uint64
	Follow    map[uint64]struct{}
	Followers map[uint64]struct{}
}

func NewProfile(id uint64) *Profile {
	return &Profile{
		ID:        id,
		Follow:    make(map[uint64]struct{}),
		Followers: make(map[uint64]struct{}),
	}
}
