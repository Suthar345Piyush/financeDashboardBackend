package models

import "time"

// this file contains all the structs related to user

type Role string

// providing role to user's

const (
	RoleViewer  Role = "viewer"
	RoleAnalyst Role = "analyst"
	RoleAdmin   Role = "admin"
)

// function to check that user choosing valid role on not

func (r Role) isValid() bool {

	switch r {
	case RoleViewer, RoleAnalyst, RoleAdmin:
		return true
	}

	return false
}

// has permission function to check the user doing what is allowed to do

func (r Role) HasPermission(required Role) bool {

	order := map[Role]int{
		RoleViewer:  1,
		RoleAnalyst: 2,
		RoleAdmin:   3,
	}

	return order[r] >= order[required]
}

// struct for user

type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // password hash
	Role      Role      `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// struct for response comes from user

type UserResponse struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Role      Role      `json:"role"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

// function for , which response  sent to someone

func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Email:     u.Email,
		Role:      u.Role,
		IsActive:  u.IsActive,
		CreatedAt: u.CreatedAt,
	}
}
