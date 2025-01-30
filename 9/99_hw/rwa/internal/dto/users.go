package dto

// Register user DTO
type UserRegisterDataRequest struct {
	Username string `json:"username" valid:"required,alphanum"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

type UserRegisterRequest struct {
	User UserRegisterDataRequest `json:"user" valid:"required"`
}

type UserRegisterDataResponse struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UserRegisterResponse struct {
	User UserRegisterDataResponse `json:"user" valid:"required"`
}

// Login user dto
type UserLoginRequest struct {
	Username string `json:"username" valid:"required,alphanum"`
	Password string `json:"password" valid:"required"`
}

type UserLoginResponse struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Token    string `json:"token"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}
