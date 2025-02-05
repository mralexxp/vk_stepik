package service

import "rwa/internal/models"

type UserStore interface {
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

type ProfileStore interface {
	AddProfile(*models.Profile) error
	DeleteProfile(uint64)
	GetProfile(uint64) (*models.Profile, error)
	Follow(uint64, uint64) error
	Unfollow(uint64, uint64) error
}

type Service struct {
	Users   UserStore
	Profile ProfileStore
	SM      SessManager
}

func NewService(users UserStore, sm SessManager, profile ProfileStore) *Service {
	return &Service{
		Users:   users,
		Profile: profile,
		SM:      sm,
	}
}

func (s *Service) GetSessionManager() SessManager {
	return s.SM
}
