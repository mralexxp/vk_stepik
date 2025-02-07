package service

import (
	"fmt"
	"github.com/asaskevich/govalidator"
	"rwa/internal/dto"
	"rwa/internal/models"
	"rwa/internal/password"
	"time"
)

func (s *Service) RegisterUser(userDTO *dto.UserRequest) (*dto.UserResponse, error) {
	ok, err := govalidator.ValidateStruct(userDTO)
	if err != nil {
		return nil, err
	}

	if !ok || !RegisterValid(userDTO.User) {
		return nil, fmt.Errorf(": input is invalid: %v", *userDTO)
	}

	hashedPassword, err := password.Hash(userDTO.User.Password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:  userDTO.User.Username,
		Email:     userDTO.User.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Bio:       "",
	}

	id, err := s.Users.Add(user)
	if err != nil {
		return nil, err
	}

	// profile register
	err = s.Profile.AddProfile(models.NewProfile(id))
	if err != nil {
		return nil, err
	}

	token, err := s.SM.Create(id)
	if err != nil {
		return nil, err
	}

	responseData := &dto.UserDataResponse{
		Email:     userDTO.User.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
		UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		Username:  userDTO.User.Username,
		Token:     token,
		Bio:       userDTO.User.Bio,
		Image:     userDTO.User.Image,
	}

	return &dto.UserResponse{User: responseData}, nil
}

func (s *Service) LoginUser(userDTO *dto.UserRequest) (*dto.UserResponse, error) {
	ok, err := govalidator.ValidateStruct(userDTO)
	if err != nil || !ok {
		return nil, fmt.Errorf(": input is invalid")
	}

	// Валидация полей email и password (не чек)
	if !LoginValid(userDTO.User) {
		return nil, fmt.Errorf(": invalid email or password")
	}

	// PassCheck
	user, err := s.Users.GetByEmail(userDTO.User.Email)
	if err != nil {
		return nil, err
	}

	if !password.Check(userDTO.User.Password, user.Password) {
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

func (s *Service) UpdateUser(userDTO *dto.UserRequest) (*dto.UserResponse, error) {
	id, ok := s.SM.Check(userDTO.User.Token)
	if !ok {
		return nil, fmt.Errorf("invalid token")
	}

	user := &models.User{
		ID:       id,
		Username: userDTO.User.Username,
		Email:    userDTO.User.Email,
		Bio:      userDTO.User.Bio,
		Image:    userDTO.User.Image,
	}

	newUser, err := s.Users.Update(user)
	if err != nil {
		return nil, err
	}

	response := &dto.UserDataResponse{
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
		UpdatedAt: newUser.UpdatedAt.Format(time.RFC3339),
		Username:  newUser.Username,
		Bio:       newUser.Bio,
		Image:     newUser.Image,
		Token:     userDTO.User.Token,
	}

	return &dto.UserResponse{User: response}, nil
}

func (s *Service) LogoutUser(token string) (*dto.UserResponse, error) {
	_, err := s.SM.DestroyByToken(token)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
