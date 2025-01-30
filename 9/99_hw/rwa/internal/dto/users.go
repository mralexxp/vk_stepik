package dto

// Register user DTO
type UserDataRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Token    string `json:"token"`
	Bio      string `json:"bio"`
	Image    string `json:"image"`
}

type UserRegisterRequest struct {
	User *UserDataRequest `json:"user" valid:"required"`
}

type UserDataResponse struct {
	ID        uint64 `json:"id,omitempty"`
	Email     string `json:"email,omitempty"`
	Token     string `json:"token,omitempty"`
	CreatedAt string `json:"createdAt,omitempty"`
	UpdatedAt string `json:"updatedAt,omitempty"`
	Username  string `json:"username,omitempty"`
	Bio       string `json:"bio,omitempty"`
	Image     string `json:"image,omitempty"`
	Following bool   `json:"following,omitempty"`
}

type UserResponse struct {
	User *UserDataResponse `json:"user" valid:"required"`
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
