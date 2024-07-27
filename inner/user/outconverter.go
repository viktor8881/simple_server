package user

import (
	"simpleserver/domain/user"
	generated "simpleserver/generated/http/server"
)

func userToResponse(u *user.Model) *generated.GetUserResponse {
	return &generated.GetUserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}

func listUserToResponse(users []user.Model) *generated.ListUserResponse {
	dest := generated.ListUserResponse{
		Users: make([]generated.UserResponse, 0, len(users)),
	}
	for _, u := range users {
		dest.Users = append(dest.Users, generated.UserResponse{
			ID:    u.ID,
			Name:  u.Name,
			Email: u.Email,
		})
	}
	return &dest
}
