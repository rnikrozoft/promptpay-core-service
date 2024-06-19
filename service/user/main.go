package user

import (
	"github.com/rnikrozoft/promptpay-core-service/model"
	"github.com/rnikrozoft/promptpay-core-service/repository"
)

type UserService interface {
	GetUser() (*model.Users, error)
}

type user struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &user{
		userRepo: userRepo,
	}
}
