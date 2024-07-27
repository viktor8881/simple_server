package user

import (
	"context"
	"errors"
	"github.com/viktor8881/service-utilities/http/server"
	"net/http"
	"simpleserver/domain/user"
	generated "simpleserver/generated/http/server"
)

type UserService struct {
	repository *user.Repository
}

func NewService(repository *user.Repository) *UserService {
	return &UserService{repository: repository}
}

func (s *UserService) ListUser(ctx context.Context, in *generated.EmptyRequest) (*generated.ListUserResponse, error) {
	users, err := s.repository.FetchAll(ctx)
	if err != nil {
		return nil, err
	}

	return listUserToResponse(users), nil
}

func (s *UserService) GetUser(ctx context.Context, in *generated.GetUserRequest) (*generated.GetUserResponse, error) {
	userModel, err := s.repository.Get(ctx, in.ID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, &server.CustomError{
				HttpCode:    http.StatusBadRequest,
				HttpMessage: "user not found",
				Err:         err,
			}
		}
		return nil, err
	}

	return userToResponse(userModel), nil
}

func (s *UserService) CreateUser(ctx context.Context, in *generated.CreateUserRequest) (*generated.GetUserResponse, error) {
	user := user.Model{
		Name:  in.Name,
		Email: in.Email,
	}

	newID, err := s.repository.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	user.ID = newID

	return userToResponse(&user), nil
}

func (s *UserService) UpdateUser(ctx context.Context, in *generated.UpdateUserRequest) (*generated.GetUserResponse, error) {
	userModel := user.Model{
		ID:    in.ID,
		Name:  in.Name,
		Email: in.Email,
	}

	rowAffect, err := s.repository.Update(ctx, userModel)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, &server.CustomError{
				HttpCode:    http.StatusBadRequest,
				HttpMessage: "user not found",
				Err:         err,
			}
		}
		return nil, err
	}

	if rowAffect == 0 {
		return nil, &server.CustomError{
			HttpCode: http.StatusNoContent,
		}
	}

	return userToResponse(&userModel), nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *generated.DeleteUserRequest) (*generated.EmptyResponse, error) {
	_, err := s.repository.Delete(ctx, in.ID)
	if err != nil {
		if errors.Is(err, user.ErrNotFound) {
			return nil, &server.CustomError{
				HttpCode:    http.StatusBadRequest,
				HttpMessage: "user not found",
				Err:         err,
			}
		}
		return nil, err
	}

	return &generated.EmptyResponse{}, nil
}

func (s *UserService) ListUserByEmail(ctx context.Context, in *generated.ListUserByEmailRequest) (*generated.ListUserResponse, error) {
	users, err := s.repository.FetchAllByEmail(ctx, in.Email)
	if err != nil {
		return nil, err
	}

	return listUserToResponse(users), nil
}
