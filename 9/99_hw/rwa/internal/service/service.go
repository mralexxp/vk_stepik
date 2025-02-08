package service

const (
	DefaultArticleLimit = 20
)

type SessManager interface {
	Create(uint64) (string, error)
	Check(string) (uint64, bool)
	DestroyByToken(string) (uint64, error)
	DestroyByID(uint64) (int, error)
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
