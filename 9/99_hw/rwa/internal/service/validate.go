package service

import (
	"github.com/asaskevich/govalidator"
	"rwa/internal/dto"
)

func RegisterValid(userData *dto.UserDataRequest) bool {
	if !govalidator.IsEmail(userData.Email) || !govalidator.IsAlphanumeric(userData.Username) {
		return false
	}

	// minLenPassword (love = 4)
	if len(userData.Password) < 3 {
		return false
	}

	return true
}

func LoginValid(userData *dto.UserDataRequest) bool {
	if !govalidator.IsEmail(userData.Email) {
		return false
	}

	// minLenPassword (love = 4)
	if len(userData.Password) < 3 {
		return false
	}

	return true
}
