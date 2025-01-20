package services

import (
	"database/sql"
	"errors"
	"websocket-server/internal/models"
	"websocket-server/internal/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) RegisterUser(username, password string) (string, error) {
	if username == "" || password == "" {
		return "", errors.New("username and password cannot be empty")
	}

	user := models.User{
		Username: username,
		Password: password, // Password should be hashed before storing
	}

	return s.repo.CreateUser(user)
}

func (s *UserService) AuthenticateUser(username, password string) (string, error) {
	return s.repo.Authenticate(username, password)
}

func (s *UserService) GetUserProfile(userID string) (*models.User, error) {
	return s.repo.GetUserByID(userID)
}

func (s *UserService) UpdateUserProfile(userID string, updatedUser models.User) error {
	return s.repo.UpdateUser(userID, updatedUser)
}