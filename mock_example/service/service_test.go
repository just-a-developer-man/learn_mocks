package service

import (
	"context"
	"mock_example/models"
	"mock_example/service/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestCreateUser(t *testing.T) {
	type args struct {
		user models.User
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "base test",
			args: args{
				models.User{
					Name:  "user",
					Email: "user@mail.com",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// для каждого из интерфейсов создаем мок-объект
			userCreator := mocks.NewUserCreator(t)
			userProvider := mocks.NewUserProvider(t)
			eventNotifier := mocks.NewEventNotifier(t)

			// задаем поведение мок-объектов
			userProvider.On("User", mock.AnythingOfType("context.Context"), tt.args.user.Email).Once().Return(nil, nil)
			userCreator.On("Create", mock.AnythingOfType("context.Context"), tt.args.user).Once().Return(0, nil)
			eventNotifier.On("NotifyUserCreated", mock.AnythingOfType("context.Context"), tt.args.user).Once().Return(nil)

			s := &Service{
				userCreator:   userCreator,
				userProvider:  userProvider,
				eventNotifier: eventNotifier,
			}

			// тестируем требуемый функционал
			_, err := s.CreateUser(context.Background(), tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateUser() err = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
