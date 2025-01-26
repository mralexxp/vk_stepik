package users

import "rwa/internal/models"

type Store struct {
	db map[string]*models.User
}

func NewUsersStore() *Store {
	return &Store{db: make(map[string]*models.User)}
}

func (s *Store) AddUser(user *models.User) error {
	s.db[user.Username] = user
	return nil
}

func (s *Store) GetUser(username string) (*models.User, error) {
	return s.db[username], nil
}

func (s *Store) DeleteUser(username string) error {
	delete(s.db, username)

	return nil
}

func (s *Store) UpdateUser(user *models.User) error {
	// todo: обновляются только переданные поля
	const op = "Store.UpdateUser"

	panic(op + ": implement me")
}
