package models

type Profile struct {
	ID        uint64
	Follow    map[uint64]struct{}
	Followers map[uint64]struct{}
}
