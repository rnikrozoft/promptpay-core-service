package user

import (
	"github.com/rnikrozoft/promptpay-core-service/model"
	"github.com/rnikrozoft/promptpay-core-service/pkg/errs"
	"github.com/rnikrozoft/promptpay-core-service/pkg/logs"
	"gorm.io/gorm"
)

func (s *user) GetUser() (*model.Users, error) {
	users, err := s.userRepo.GetUser()
	if err != nil {
		logs.Error(err)
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewNotFoundError(err.Error())
		}
		return nil, errs.NewUnexpectedError()
	}
	return users, nil
}
