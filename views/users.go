package views

import "github.com/Bayan2019/hackathon-2025-police-volunteers/repositories/database"

type UpdateProfileRequest struct {
	// Id          int64     `json:"id"`
	Name string `json:"name"`
	// Email       string `json:"email"`
	DateOfBirth     string `json:"date_of_birth"`
	Phone           string `json:"phone"`
	Iin             string `json:"iin"`
	CurrentLocation string `json:"current_location"`
}

type UpdateUserRequest struct {
	// Id          int64     `json:"id"`
	Name string `json:"name"`
	// Email       string  `json:"email"`
	DateOfBirth     string  `json:"date_of_birth"`
	Iin             string  `json:"iin"`
	Phone           string  `json:"phone"`
	RoleIds         []int64 `json:"role_ids"`
	CurrentLocation string  `json:"current_location"`
	// CenterId        int64   `json:"center_id"`
}

type User struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
	// Email       string          `json:"email"`
	DateOfBirth string          `json:"date_of_birth"`
	Phone       string          `json:"phone"`
	Iin         string          `json:"iin"`
	Roles       []database.Role `json:"roles"`
	// Center      database.Center `json:"center"`
}

type SignInRequest struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

type TokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
