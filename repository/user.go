package repository

import (
	"github.com/rnikrozoft/promptpay-core-service/model"
)

type UserRepository interface {
	GetUser() (*model.Users, error)
}

func (r *repository) GetUser() (*model.Users, error) {
	users := &model.Users{}
	if tx := r.db.First(users, "username = ?", "ber").Error; tx != nil {
		return nil, tx
	}
	return users, nil
}
