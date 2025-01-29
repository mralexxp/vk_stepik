package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"rwa/internal/dto"
	"rwa/internal/models"
	"rwa/internal/password"
	"time"
)

func (s *Service) Register(UserDTO *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	const op = "Service.Add"

	// Валидируем полученные данные
	ok, err := govalidator.ValidateStruct(UserDTO)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf(op+": input is invalid: %v", *UserDTO)
	}

	// Хешируем пароль
	hashedPassword, err := password.Hash(UserDTO.Password)
	if err != nil {
		return nil, err
	}

	id, err := s.Users.Add(&models.User{
		Username: UserDTO.Username,
		Email:    UserDTO.Email,
		Password: hashedPassword,
		Created:  time.Now().Unix(),
		Updated:  0,
		Bio:      "",
	})
	if err != nil {
		return nil, err
	}

	token, err := s.SM.Create(id)
	if err != nil {
		return nil, err
	}

	return &dto.UserRegisterResponse{
		Email:    UserDTO.Email,
		Token:    token,
		Username: UserDTO.Username,
		Bio:      "",
		Image:    "",
	}, nil
}
