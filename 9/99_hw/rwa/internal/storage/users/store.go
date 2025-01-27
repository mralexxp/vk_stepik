package users

import (
	"fmt"
	"rwa/internal/models"
)

type Store struct {
	db map[string]*models.User
}

func NewUsersStore() *Store {
	return &Store{db: make(map[string]*models.User)}
}

func (s *Store) AddUser(user *models.User) error {
	if _, ok := s.db[user.Username]; ok {
		return fmt.Errorf("username '%s' already exists", user.Username)
	}

	s.db[user.Username] = user
	return nil
}

func (s *Store) GetUser(username string) (*models.User, error) {
	if user, ok := s.db[username]; ok {
		return user, nil
	}

	return nil, fmt.Errorf("username '%s' not found", username)
}

func (s *Store) DeleteUser(username string) error {
	if user, ok := s.db[username]; ok {
		delete(s.db, user.Username)
		return nil
	}

	return fmt.Errorf("username '%s' not found", username)
}

func (s *Store) UpdateUser(user *models.User) error {
	// todo: обновляются только переданные поля
	const op = "Store.UpdateUser"

	panic(op + ": implement me")
}
