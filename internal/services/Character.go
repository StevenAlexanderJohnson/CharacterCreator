package services

import (
	"dndcc/internal/models"
	"dndcc/internal/repositories"
)

type CharacterService struct {
	repo *repositories.CharacterRepository
}

func NewCharacterService(repo *repositories.CharacterRepository) *CharacterService {
	return &CharacterService{repo: repo}
}

func (s *CharacterService) Create(data *models.Character) (*models.Character, error) {
	return s.repo.Create(data)
}

func (s *CharacterService) Get(id, userId int) (*models.Character, error) {
	return s.repo.Get(id, userId)
}

func (s *CharacterService) List(userId int) ([]models.Character, error) {
	return s.repo.GetAll(userId)
}

func (s *CharacterService) Update(data *models.Character, id, userId int) (*models.Character, error) {
	return s.repo.Update(data, id, userId)
}

func (s *CharacterService) Delete(id, userId int) error {
	return s.repo.Delete(id, userId)
}
