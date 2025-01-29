package service

import "rwa/internal/models"

type UserStorer interface {
	AddUser(user *models.User) error
	GetUser(username string) (*models.User, error)
	DeleteUser(username string) error
	UpdateUser(user *models.User) error
}

type SessManager interface {
	Create(string) (string, error)
	Check(string) (string, bool)
	DestroyByToken(string) (string, error)
	DestroyByUsername(string) (int, error)
}

type Service struct {
	Users UserStorer
	SM    SessManager
}

func NewService(users UserStorer, sm SessManager) *Service {
	return &Service{
		Users: users,
		SM:    sm,
	}
}
