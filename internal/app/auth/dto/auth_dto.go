package dto

type RegisterUserRequest struct {
	FirstName string `json:"first_name" validate:"required,min=3"`
	LastName  string `json:"last_name" validate:"required,min=3"`
	Email     string `json:"email" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
