package service

import (
	"context"
	"redis-rate-limiter/dto"
	userModel "redis-rate-limiter/model"
	"redis-rate-limiter/repository"
)

type UserService struct {
	repository *repository.UserRepository
}
type Service interface {
	Profile(ctx context.Context, id string) *userModel.User
}

func NewUserService(repository *repository.UserRepository) *UserService {
	return &UserService{
		repository: repository,
	}
}

func (s *UserService) Profile(ctx context.Context, id string) *userModel.User {
	var RequestUser *dto.UserRequest

	//get the user from the repository using the id
	for _, user := range *s.repository.DB { //ek ek object ko iterate kar raha hai s->repository->db slice of user model
		if user.Id == id {
			RequestUser = &dto.UserRequest{
				ID:   user.Id,
				Name: user.Name,
			}
		}
	}
	return &userModel.User{
		Id:   RequestUser.ID,
		Name: RequestUser.Name,
	}

}
