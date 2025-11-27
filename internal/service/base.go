package service

type Repository interface{}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

type service struct {
	repository Repository
}
