package dto

// Register response
func NewRegisterResponse(email, token, username, bio, image string) *UserRegisterResponse {
	return &UserRegisterResponse{
		User: UserRegisterDataResponse{
			Email:    email,
			Token:    token,
			Username: username,
			Bio:      bio,
			Image:    image,
		},
	}
}
