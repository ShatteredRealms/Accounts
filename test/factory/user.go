package factory

import (
	"github.com/ShatteredRealms/Accounts/pkg/model"
	"github.com/ShatteredRealms/GoUtils/pkg/helpers"
	"time"
)

type UserFactory interface {
	User() model.User
}

type userFactory struct{}

func NewUserFactory() userFactory {
	return userFactory{}
}

func (u userFactory) BaseUser() model.User {
	return model.User{
		FirstName: helpers.RandString(10),
		LastName:  helpers.RandString(10),
		Email:     helpers.RandString(10) + "@test.com",
		Password:  helpers.RandString(10),
		Username:  helpers.RandString(10),
	}
}

func (u userFactory) User() model.User {
	user := u.BaseUser()
	user.ID = helpers.RandInt(10000)
	user.CreatedAt = time.Now().Add(-time.Second * time.Duration(helpers.RandInt(10000)))
	user.UpdatedAt = user.CreatedAt

	return user
}
