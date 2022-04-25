package factory

type Factory interface {
	UserFactory() UserFactory
}

type factory struct{}

func NewFactory() factory {
	return factory{}
}

func (f factory) UserFactory() userFactory {
	return NewUserFactory()
}
