package mocks

import (
	"mzt/internal/entity"
	"mzt/internal/repository"

	"github.com/google/uuid"
)

type MockUserRepository struct {
	Users     map[uuid.UUID]*entity.User
	UserData  map[uuid.UUID]*entity.UserData
	UserAuth  map[uuid.UUID]*entity.Auth
	UserEmail map[string]uuid.UUID
}

func NewMockUserRepository() repository.UserRepository {
	return &MockUserRepository{
		Users:     make(map[uuid.UUID]*entity.User),
		UserData:  make(map[uuid.UUID]*entity.UserData),
		UserAuth:  make(map[uuid.UUID]*entity.Auth),
		UserEmail: make(map[string]uuid.UUID),
	}
}

func (m *MockUserRepository) GetUserByEmail(email string) (*entity.User, error) {
	if userId, exists := m.UserEmail[email]; exists {
		return m.Users[userId], nil
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserWithDataById(id uuid.UUID) (*entity.User, error) {
	if user, exists := m.Users[id]; exists {
		if userData, exists := m.UserData[id]; exists {
			user.UserData = userData
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) GetUserWithRefreshById(id uuid.UUID) (*entity.User, error) {
	if user, exists := m.Users[id]; exists {
		if auth, exists := m.UserAuth[id]; exists {
			user.Auth = auth
			return user, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) CreateUser(user *entity.User, userData *entity.UserData, auth *entity.Auth) error {
	m.Users[user.ID] = user
	m.UserData[user.ID] = userData
	m.UserAuth[user.ID] = auth
	m.UserEmail[userData.Email] = user.ID
	return nil
}

func (m *MockUserRepository) UpdateToken(userId uuid.UUID, token string) error {
	if auth, exists := m.UserAuth[userId]; exists {
		auth.Key = token
		return nil
	}
	return nil
}

func (m *MockUserRepository) GetUsers() ([]entity.User, error) {
	users := make([]entity.User, 0)
	for id, user := range m.Users {
		if userData, exists := m.UserData[id]; exists {
			user.UserData = userData
			users = append(users, *user)
		}
	}
	return users, nil
}

func (m *MockUserRepository) UpdateUser(userId uuid.UUID, updated *entity.UserData) error {
	if _, exists := m.Users[userId]; exists {
		m.UserData[userId] = updated
		return nil
	}
	return nil
}

func (m *MockUserRepository) DeleteUser(userId uuid.UUID) error {
	if _, exists := m.Users[userId]; exists {
		delete(m.Users, userId)
		delete(m.UserData, userId)
		delete(m.UserAuth, userId)
		for email, id := range m.UserEmail {
			if id == userId {
				delete(m.UserEmail, email)
				break
			}
		}
		return nil
	}
	return nil
}

func (m *MockUserRepository) GetUserById(userId uuid.UUID) (*entity.User, error) {
	if user, exists := m.Users[userId]; exists {
		return user, nil
	}
	return nil, nil
}
