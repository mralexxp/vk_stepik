package service

import "rwa/internal/dto"

func (s *Service) Add(UserDTO *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	// TODO: Валидируем все поля

	// TODO: в репо отправляем

	// TODO: формируем ответ

	return &dto.UserRegisterResponse{}, nil
}
