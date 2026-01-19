package services

import (
	"github.com/google/uuid"
	"github.com/hoyci/gopher-chaos/example/server/internal/repositories"
)

type UserUseCase struct {
	Repo repositories.UserRepository
}

func (s *UserUseCase) CreateUser(name, email string, age int32) (*repositories.User, error) {
	user := &repositories.User{
		ID:    uuid.New().String(),
		Name:  name,
		Email: email,
		Age:   age,
	}
	err := s.Repo.Save(user)
	return user, err
}

func (s *UserUseCase) GetByID(id string) (*repositories.User, error) {
	user, err := s.Repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserUseCase) UpdateByID(id, name string) (*repositories.User, error) {
	user, err := s.Repo.UpdateByID(id, name)
	if err != nil {
		return nil, err
	}

	return user, err
}

func (s *UserUseCase) DeleteByID(id string) error {
	err := s.Repo.DeleteByID(id)
	if err != nil {
		return err
	}

	return nil
}
