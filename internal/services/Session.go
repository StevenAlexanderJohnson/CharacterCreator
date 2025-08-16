package services

import (
	"dndcc/internal/models"
	"dndcc/internal/repositories"
)

type SessionService struct {
	repo *repositories.SessionRepository
}

func NewSessionService(repo *repositories.SessionRepository) *SessionService {
	return &SessionService{repo: repo}
}

func (s *SessionService) Create(session *models.Session) (*models.Session, error) {
	return s.repo.Create(session)
}

func (s *SessionService) Get(token string) (*models.Session, error) {
	return s.repo.GetByToken(token)
}

func (s *SessionService) List(userId int) ([]models.Session, error) {
	return s.repo.GetAllUserSessions(userId)
}

func (s *SessionService) Update(data *models.Session, userId int) (*models.Session, error) {
	return s.repo.Update(data, userId)
}

func (s *SessionService) Delete(id int) error {
	return s.repo.Delete(id)
}
