package service

import (
	"fmt"
	"practice8/repository"
	"testing"

	"github.com/go-openapi/testify/v2/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan agai legenda"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	res, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, res)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan agai legenda"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestRegisterUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("User already exists", func(t *testing.T) {
		user := &repository.User{ID: 2, Name: "Existing User"}
		existing := &repository.User{ID: 1, Name: "Existing"}
		mockRepo.EXPECT().GetByEmail("existing@example.com").Return(existing, nil)

		err := userService.RegisterUser(user, "existing@example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user with this email already exists")
	})

	t.Run("New User -> Success", func(t *testing.T) {
		user := &repository.User{ID: 2, Name: "New User"}
		mockRepo.EXPECT().GetByEmail("new@example.com").Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(nil)

		err := userService.RegisterUser(user, "new@example.com")
		assert.NoError(t, err)
	})

	t.Run("Repository error on CreateUser", func(t *testing.T) {
		user := &repository.User{ID: 2, Name: "New User"}
		mockRepo.EXPECT().GetByEmail("new@example.com").Return(nil, nil)
		mockRepo.EXPECT().CreateUser(user).Return(fmt.Errorf("db error"))

		err := userService.RegisterUser(user, "new@example.com")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "db error")
	})
}

func TestUpdateUserName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("Empty name", func(t *testing.T) {
		err := userService.UpdateUserName(1, "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "name cannot be empty")
	})

	t.Run("User not found/repo error", func(t *testing.T) {
		mockRepo.EXPECT().GetUserByID(1).Return(nil, fmt.Errorf("user not found"))

		err := userService.UpdateUserName(1, "New Name")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "user not found")
	})

	t.Run("Successful update", func(t *testing.T) {
		oldUser := &repository.User{ID: 1, Name: "Old Name"}
		mockRepo.EXPECT().GetUserByID(1).Return(oldUser, nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(user *repository.User) error {
			// Verify that name was changed
			assert.Equal(t, 1, user.ID)
			assert.Equal(t, "New Name", user.Name)
			return nil
		})

		err := userService.UpdateUserName(1, "New Name")
		assert.NoError(t, err)
	})

	t.Run("UpdateUser Fails", func(t *testing.T) {
		user := &repository.User{ID: 1, Name: "Old Name"}
		mockRepo.EXPECT().GetUserByID(1).Return(user, nil)
		mockRepo.EXPECT().UpdateUser(gomock.Any()).Return(fmt.Errorf("update failed"))

		err := userService.UpdateUserName(1, "New Name")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "update failed")
	})
}

func TestDeleteUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	t.Run("Attempt to delete admin", func(t *testing.T) {
		err := userService.DeleteUser(1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "it is not allowed to delete admin user")
	})

	t.Run("Successful delete", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(2).Return(nil)

		err := userService.DeleteUser(2)
		assert.NoError(t, err)
	})

	t.Run("Repository Error", func(t *testing.T) {
		mockRepo.EXPECT().DeleteUser(2).Return(fmt.Errorf("delete failed"))

		err := userService.DeleteUser(2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delete failed")
	})
}
