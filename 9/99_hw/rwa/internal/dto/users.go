package dto

type UserRequest struct {
	User *UserDataRequest `json:"user" valid:"required"`
}

type UserResponse struct {
	User *UserDataResponse `json:"user" valid:"required"`
}

type UserDataRequest struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
	Bio      string `json:"bio,omitempty"`
	Image    string `json:"image,omitempty"`
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
	Follow    bool   `json:"follow,omitempty"`
}
