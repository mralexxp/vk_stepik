package service

import "rwa/internal/models"

const (
	DefaultArticleLimit = 20
)

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

type ArticlesStore interface {
	Add(article *models.Article) uint64
	Delete(id uint64) error
	Get(id uint64) (*models.Article, error)
	GetSlugID(slug string) (uint64, error)
	GetByFilter(filter *models.ArticleFilter) ([]*models.Article, error)
}

type Service struct {
	Users    UserStore
	Articles ArticlesStore
	SM       SessManager
}

func NewService(
	a ArticlesStore,
	u UserStore,
	sm SessManager) *Service {
	return &Service{
		Users:    u,
		Articles: a,
		SM:       sm,
	}
}

func (s *Service) GetSessionManager() SessManager {
	return s.SM
}
