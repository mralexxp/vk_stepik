package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"rwa/internal/dto"
	"rwa/internal/models"
	"rwa/internal/password"
	"time"
)

func (s *Service) RegisterUser(UserDTO *dto.UserRequest) (*dto.UserResponse, error) {
	const op = "Service.Add"

	ok, err := govalidator.ValidateStruct(UserDTO)
	if err != nil {
		return nil, err
	}

	if !ok || !RegisterValid(UserDTO.User) {
		return nil, fmt.Errorf(op+": input is invalid: %v", *UserDTO)
	}

	hashedPassword, err := password.Hash(UserDTO.User.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  UserDTO.User.Username,
		Email:     UserDTO.User.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Bio:       "",
	}

	id, err := s.Users.Add(user)
	if err != nil {
		return nil, err
	}

	token, err := s.SM.Create(id)
	if err != nil {
		return nil, err
	}

	responseData := &dto.UserDataResponse{
		Email:     UserDTO.User.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Username:  UserDTO.User.Username,
		Token:     token,
		Bio:       UserDTO.User.Bio,
		Image:     UserDTO.User.Image,
	}

	return &dto.UserResponse{User: responseData}, nil
}

func (s *Service) LoginUser(UserDTO *dto.UserRequest) (*dto.UserResponse, error) {
	const op = "Service.Login"

	ok, err := govalidator.ValidateStruct(UserDTO)
	if err != nil || !ok {
		return nil, fmt.Errorf(op + ": input is invalid")
	}

	// Валидация полей email и password (не чек)
	if !LoginValid(UserDTO.User) {
		return nil, fmt.Errorf(op + ": invalid email or password")
	}

	// PassCheck
	user, err := s.Users.GetByEmail(UserDTO.User.Email)
	if err != nil {
		return nil, err
	}

	if !password.Check(UserDTO.User.Password, user.Password) {
		return nil, fmt.Errorf("incorrect username or password")
	}

	// Успешная аутентификация
	token, err := s.SM.Create(user.ID)
	if err != nil {
		return nil, err
	}

	response := &dto.UserDataResponse{
		Username:  user.Username,
		Email:     user.Email,
		Token:     token,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Bio:       user.Bio,
		Image:     user.Image,
	}

	return &dto.UserResponse{
		User: response,
	}, nil

}

func (s *Service) GetCurrentUser(token string) (*dto.UserResponse, error) {
	id, ok := s.SM.Check(token)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	user, err := s.Users.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.UserDataResponse{
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
	}

	return &dto.UserResponse{User: response}, nil
}
