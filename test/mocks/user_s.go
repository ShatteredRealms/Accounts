package mocks

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/Accounts/pkg/service"
	"gorm.io/gorm"
)

type UserService struct {
	CreateReturn      error
	SaveReturn        error
	FindByIdReturn    model.User
	FindByEmailReturn model.User
}

func (t UserService) Create(u model.User) (model.User, error) {
	return u, t.CreateReturn
}
func (t UserService) Save(u model.User) (model.User, error) {
	return u, t.SaveReturn
}
func (t UserService) WithTrx(*gorm.DB) service.UserService {
	return t
}

func (t UserService) FindById(id uint) model.User {
	return t.FindByIdReturn
}
func (t UserService) FindByEmail(email string) model.User {
	return t.FindByEmailReturn
}
