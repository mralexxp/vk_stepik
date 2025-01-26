package service

import "rwa/internal/models"

type UserStorer interface {
	AddUser(user *models.User) error
	GetUser(username string) (*models.User, error)
	DeleteUser(username string) error
	UpdateUser(user *models.User) error
}

type Service struct {
	Users UserStorer
}

func NewService(users UserStorer) *Service {
	return &Service{Users: users}
}
