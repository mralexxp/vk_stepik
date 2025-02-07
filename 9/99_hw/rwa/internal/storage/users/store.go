package users

import (
	"fmt"
	"log"
	"reflect"
	"rwa/internal/models"
	"time"
)

type Store struct {
	db map[uint64]*models.User

	// index
	UsernameID map[string]uint64
	EmailID    map[string]uint64

	// increment
	nextID uint64
}

func NewUsersStore() *Store {
	return &Store{
		db:         make(map[uint64]*models.User),
		UsernameID: make(map[string]uint64),
		EmailID:    make(map[string]uint64),
		nextID:     1, // ID начинаются с единицы
	}
}

func (s *Store) Add(user *models.User) (uint64, error) {
	if _, ok := s.UsernameID[user.Username]; ok {
		return 0, fmt.Errorf("username %s already exists", user.Username)
	}

	if _, ok := s.EmailID[user.Email]; ok {
		return 0, fmt.Errorf("username %s already exists", user.Username)
	}

	user.ID = s.nextID
	s.nextID++

	s.UsernameID[user.Username] = user.ID
	s.EmailID[user.Email] = user.ID
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

func (s *Store) GetByEmail(email string) (*models.User, error) {
	if id, ok := s.EmailID[email]; ok {
		if u, ok := s.db[id]; ok {
			return u, nil
		}

		log.Println("user found in relation EmailID, but not found DB: ", email, id)
	}

	return nil, fmt.Errorf("email '%s' not found", email)
}

func (s *Store) GetByID(id uint64) (*models.User, error) {
	if u, ok := s.db[id]; ok {
		return u, nil
	}

	return nil, fmt.Errorf("user %d not found", id)
}

func (s *Store) DeleteByUsername(username string) error {
	if id, ok := s.UsernameID[username]; ok {
		if u, ok := s.db[id]; ok {
			delete(s.EmailID, u.Email)
		}
		delete(s.UsernameID, username)
		delete(s.db, id)
		return nil
	}

	return fmt.Errorf("username '%s' not found", username)
}

func (s *Store) DeleteByEmail(email string) error {
	if id, ok := s.EmailID[email]; ok {
		if u, ok := s.db[id]; ok {
			delete(s.UsernameID, u.Username)
		}
		delete(s.EmailID, email)
		delete(s.db, id)
		return nil
	}

	return fmt.Errorf("email '%s' not found", email)
}

func (s *Store) DeleteByID(id uint64) error {
	if u, ok := s.db[id]; ok {
		delete(s.UsernameID, u.Username)
		delete(s.EmailID, u.Email)
		delete(s.db, id)
		return nil
	}

	return fmt.Errorf("user %d not found", id)
}

func (s *Store) Update(updateUser *models.User) (*models.User, error) {
	currentUser, ok := s.db[updateUser.ID]
	if !ok {
		return nil, fmt.Errorf("user %d not found", updateUser.ID)
	}

	// Сохранение значений для обновления базы по окончанию
	oldUsername := currentUser.Username
	oldEmail := currentUser.Email

	updateV := reflect.ValueOf(updateUser).Elem()
	currentV := reflect.ValueOf(currentUser).Elem()

	for i := 0; i < updateV.NumField(); i++ {
		field := updateV.Field(i)
		fieldName := updateV.Type().Field(i).Name

		if fieldName == "ID" || fieldName == "CreatedAt" || fieldName == "Follows" || fieldName == "UpdatedAt" {
			continue
		}

		if !field.IsZero() {
			currentField := currentV.FieldByName(fieldName)
			if currentField.IsValid() && currentField.CanSet() {
				currentField.Set(field)
			}
		}
	}

	// Обновляем базы:
	if _, ok := s.UsernameID[oldUsername]; ok {
		delete(s.UsernameID, oldUsername)
	}

	s.UsernameID[updateUser.Username] = updateUser.ID

	if _, ok := s.EmailID[oldEmail]; ok {
		delete(s.EmailID, oldEmail)
	}

	s.EmailID[updateUser.Email] = updateUser.ID

	currentUser.UpdatedAt = time.Now()

	return currentUser, nil
}
