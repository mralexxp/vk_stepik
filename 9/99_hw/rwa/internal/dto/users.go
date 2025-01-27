package dto

// Register user DTO
type UserRegisterRequest struct {
	Username string `json:"username" valid:"required,alphanum"`
	Email    string `json:"email" valid:"required,email"`
	Password string `json:"password" valid:"required"`
}

type UserRegisterResponse struct {
	Email    string `json:"email"`
	Token    string `json:"token"`
	Username string `json:"username"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}
