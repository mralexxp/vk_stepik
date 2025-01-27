package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"rwa/internal/dto"
	"rwa/internal/models"
	"time"
)

func (s *Service) Add(UserDTO *dto.UserRegisterRequest) (*dto.UserRegisterResponse, error) {
	const op = "Service.Add"

	ok, err := govalidator.ValidateStruct(UserDTO)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf(op+": input is invalid: %v", *UserDTO)
	}

	err = s.Users.AddUser(&models.User{
		Username:  UserDTO.Username,
		Email:     UserDTO.Email,
		Password:  UserDTO.Password,
		Created:   time.Now().Unix(),
		Updated:   0,
		ProfileID: 0, // TODO: Сразу генерация профиля
	})
	if err != nil {
		return nil, err
	}

	return &dto.UserRegisterResponse{
		Email:    UserDTO.Email,
		Token:    "GENERATE TOKEN!!!", // TODO: подключить auth
		Username: UserDTO.Username,
		Bio:      "",
		Image:    "",
	}, nil
}
