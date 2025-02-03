package service

import "rwa/internal/models"

type UserStorer interface {
	Add(*models.User) (uint64, error)
	GetByUsername(string) (*models.User, error)
	GetByEmail(string) (*models.User, error)
	GetByID(uint64) (*models.User, error)
	DeleteByUsername(string) error
	DeleteByID(uint64) error
	Update(*models.User) (*models.User, error)
}

type SessManager interface {
	Create(uint64) (string, error)
	Check(string) (uint64, bool)
	DestroyByToken(string) (uint64, error)
	DestroyByID(uint64) (int, error)
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

func (s *Service) GetSessionManager() SessManager {
	return s.SM
}
