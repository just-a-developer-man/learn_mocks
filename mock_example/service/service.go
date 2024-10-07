package service

import (
	"context"
	"errors"
	"fmt"
	"mock_example/models"
)

//go:generate go run github.com/vektra/mockery/v2@v2.46.2 --name=UserCreator
type UserCreator interface {
	Create(ctx context.Context, u models.User) (int, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.2 --name=UserProvider
type UserProvider interface {
	User(ctx context.Context, email string) (*models.User, error)
}

//go:generate go run github.com/vektra/mockery/v2@v2.46.2 --name=EventNotifier
type EventNotifier interface {
	NotifyUserCreated(ctx context.Context, user models.User) error
}

type Service struct {
	userCreator   UserCreator
	userProvider  UserProvider
	eventNotifier EventNotifier
}

func (s *Service) CreateUser(ctx context.Context, user models.User) (int, error) {
	// check if user exists
	foundUser, err := s.userProvider.User(ctx, user.Email)
	if err != nil {
		return 0, fmt.Errorf("can't get user: %w", err)
	}

	if foundUser != nil {
		return 0, errors.New("user already exists")
	}

	// create user
	uid, err := s.userCreator.Create(ctx, user)
	if err != nil {
		return 0, fmt.Errorf("can't create user: %w", err)
	}

	// notify user created
	if err := s.eventNotifier.NotifyUserCreated(ctx, user); err != nil {
		return 0, fmt.Errorf("can't notify user created: %w", err)
	}

	return uid, nil
}
