package profile

import (
	"fmt"
	"rwa/internal/models"
)

// TODO: Объединить с Users

type Store struct {
	db map[uint64]*models.Profile
}

func NewStore() *Store {
	return &Store{db: make(map[uint64]*models.Profile)}
}

func (s *Store) AddProfile(profile *models.Profile) error {
	if _, ok := s.db[profile.ID]; ok {
		return fmt.Errorf("profile with id %d already exists", profile.ID)
	}

	s.db[profile.ID] = profile
	return nil
}

func (s *Store) DeleteProfile(id uint64) {
	delete(s.db, id)
}

func (s *Store) GetProfile(id uint64) (*models.Profile, error) {
	if profile, ok := s.db[id]; ok {
		return profile, nil
	}

	return nil, fmt.Errorf("profile with id %d not found", id)
}

func (s *Store) Follow(from uint64, to uint64) error {
	if _, ok := s.db[from]; !ok {
		return fmt.Errorf("profile from %d to %d not found", from, to)
	}

	if _, ok := s.db[to]; !ok {
		return fmt.Errorf("profile from %d to %d not found", from, to)
	}

	s.db[from].Follow[to] = struct{}{}
	s.db[to].Followers[from] = struct{}{}

	return nil
}

func (s *Store) Unfollow(from uint64, to uint64) error {
	if _, ok := s.db[from]; !ok {
		return fmt.Errorf("profile from %d to %d not found", from, to)
	}

	if _, ok := s.db[to]; !ok {
		return fmt.Errorf("profile from %d to %d not found", from, to)
	}

	delete(s.db[from].Follow, to)
	delete(s.db[to].Followers, from)

	return nil
}
