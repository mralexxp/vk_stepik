package users

import (
	"fmt"
	"log"
	"rwa/internal/models"
)

type Store struct {
	db         map[uint64]*models.User
	UsernameID map[string]uint64
	nextID     uint64
}

func NewUsersStore() *Store {
	return &Store{
		db:         make(map[uint64]*models.User),
		UsernameID: make(map[string]uint64),
		nextID:     1, // ID начинаются с единицы
	}
}

func (s *Store) Add(user *models.User) (uint64, error) {
	if _, ok := s.UsernameID[user.Username]; ok {
		return 0, fmt.Errorf("username %s already exists", user.Username)
	}

	user.ID = s.nextID
	s.nextID++

	s.UsernameID[user.Username] = user.ID
	s.db[user.ID] = user

	return user.ID, nil
}

func (s *Store) GetByUsername(username string) (*models.User, error) {
	if id, ok := s.UsernameID[username]; ok {
		if u, ok := s.db[id]; ok {
			return u, nil
		}

		log.Println("user found in relation UsernameID, but not found DB: ", username, id)
	}

	return nil, fmt.Errorf("username '%s' not found", username)
}

func (s *Store) GetByID(id uint64) (*models.User, error) {
	if u, ok := s.db[id]; ok {
		return u, nil
	}

	return nil, fmt.Errorf("user %d not found", id)
}

func (s *Store) DeleteByUsername(username string) error {
	if id, ok := s.UsernameID[username]; ok {
		delete(s.UsernameID, username)
		delete(s.db, id)
		return nil
	}

	return fmt.Errorf("username '%s' not found", username)
}

func (s *Store) DeleteByID(id uint64) error {
	if u, ok := s.db[id]; ok {
		delete(s.UsernameID, u.Username)
		delete(s.db, id)
		return nil
	}

	return fmt.Errorf("user %d not found", id)
}

func (s *Store) Update(user *models.User) error {
	// todo: обновляются только переданные поля
	const op = "Store.UpdateUser"

	panic(op + ": implement me")
}
