package model

type UserResponse struct {
	Username  string         `json:"username,omitempty"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Avatar    AvatarResponse `json:"avatar"`
	CreatedAt string         `json:"created_at,omitempty"`
	UpdatedAt string         `json:"updated_at,omitempty"`
}

type RegisterUserRequest struct {
	Username string `json:"username" validate:"required,max=100"`
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" validate:"required,max=100"`
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"omitempty,min=3,max=100"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,max=100"`
	Password string `json:"password" validate:"required,max=100"`
}

type LoginUserResponse struct {
	Username    string      `json:"username"`
	Name        string      `json:"name"`
	Email       string      `json:"email"`
	AccessToken DetailToken `json:"access_token"`
	CreatedAt   string      `json:"created_at,omitempty"`
	UpdatedAt   string      `json:"updated_at,omitempty"`
}
