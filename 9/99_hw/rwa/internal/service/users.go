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
	hashedPassword, err := password.Hash(UserDTO.User.Password)
	if err != nil {
		return nil, err
	}

	id, err := s.Users.Add(&models.User{
		Username: UserDTO.User.Username,
		Email:    UserDTO.User.Email,
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

	return dto.NewRegisterResponse(
		UserDTO.User.Email,
		token,
		UserDTO.User.Username,
		"",
		"",
	), nil
}

func (s *Service) Login(UserDTO *dto.UserLoginRequest) (*dto.UserLoginResponse, error) {
	const op = "Service.Login"

	ok, err := govalidator.ValidateStruct(UserDTO)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, fmt.Errorf(op+": input is invalid: %v", *UserDTO)
	}

	user, err := s.Users.GetByUsername(UserDTO.Username)
	if err != nil {
		return nil, err
	}

	if !password.Check(UserDTO.Password, user.Password) {
		return nil, fmt.Errorf("incorrect username or password")
	}

	token, err := s.SM.Create(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.UserLoginResponse{
		Username: user.Username,
		Email:    user.Email,
		Token:    token,
		Bio:      user.Bio,
		Image:    user.Image,
	}, nil

}
