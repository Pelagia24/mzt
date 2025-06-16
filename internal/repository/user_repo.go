package repository

import (
	"fmt"
	"mzt/config"
	"mzt/internal/entity"
	"time"

	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// интерфейс для работы с пользователями
// определяет все методы которые нужны для работы с пользователями в базе
type UserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserWithDataById(id uuid.UUID) (*entity.User, error)
	GetUserWithRefreshById(id uuid.UUID) (*entity.User, error)
	CreateUser(user *entity.User, userData *entity.UserData, auth *entity.Auth) error
	UpdateToken(userId uuid.UUID, token string) error
	DeleteUser(userId uuid.UUID) error
	UpdateUser(userId uuid.UUID, updated *entity.UserData) error
	GetUsers() ([]entity.User, error)
	GetUserById(userId uuid.UUID) (*entity.User, error)
}

// репозиторий для работы с пользователями
// реализует интерфейс UserRepository
type UserRepo struct {
	config *config.Config
	DB     *gorm.DB
}

func NewUserRepo(cfg *config.Config) *UserRepo {
	db := connectDB(cfg)
	return &UserRepo{
		config: cfg,
		DB:     db,
	}
}

// получает список всех пользователей
// загружает данные пользователей и их курсы
func (r *UserRepo) GetUsers() ([]entity.User, error) {
	var users []entity.User
	err := r.DB.Preload("UserData").Preload("CourseAssignments").Preload("CourseAssignments.Course").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

// получает пользователя по почте
// сначала ищет данные пользователя по почте, потом получает самого пользователя
func (r *UserRepo) GetUserByEmail(email string) (*entity.User, error) {
	var userdata entity.UserData
	err := r.DB.Where("email = ?", email).First(&userdata).Error
	if err != nil {
		return nil, err
	}

	user, err := r.GetUserById(userdata.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// получает пользователя по id
// просто берет пользователя из базы по его id
func (r *UserRepo) GetUserById(userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.DB.Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// создает нового пользователя
// создает записи в таблицах users, user_data и auth
func (r *UserRepo) CreateUser(user *entity.User, userData *entity.UserData, auth *entity.Auth) error {
	// начинаем транзакцию
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// создаем пользователя
	err := tx.Create(&user).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// создаем данные пользователя
	err = tx.Create(&userData).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// создаем запись для авторизации
	err = tx.Create(&auth).Error
	if err != nil {
		tx.Rollback()
		return err
	}

	// завершаем транзакцию
	return tx.Commit().Error
}

// удаляет пользователя
// просто удаляет пользователя из базы
func (r *UserRepo) DeleteUser(userId uuid.UUID) error {
	result := r.DB.Delete(&entity.User{}, userId)
	return result.Error
}

// обновляет данные пользователя
// обновляет запись в таблице user_data
func (r *UserRepo) UpdateUser(userId uuid.UUID, updated *entity.UserData) error {
	// начинаем транзакцию
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// получаем текущие данные пользователя
	var data entity.UserData
	if err := tx.Where("user_id = ?", userId).First(&data).Error; err != nil {
		tx.Rollback()
		return err
	}

	// обновляем id
	updated.ID = data.ID
	updated.UserID = data.UserID

	data = *updated

	// сохраняем изменения
	if err := tx.Save(&data).Error; err != nil {
		tx.Rollback()
		return err
	}

	// завершаем транзакцию
	return tx.Commit().Error
}

// получает пользователя с данными
// загружает данные пользователя вместе с его профилем
func (r *UserRepo) GetUserWithDataById(userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("UserData").First(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// получает пользователя с токеном
// загружает данные пользователя вместе с его токеном
func (r *UserRepo) GetUserWithRefreshById(userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	err := r.DB.Preload("Auth").First(&user, "id = ?", userId).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// обновляет токен пользователя
// обновляет запись в таблице auth
func (r *UserRepo) UpdateToken(userId uuid.UUID, token string) error {
	// начинаем транзакцию
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// получаем текущий токен
	var auth entity.Auth
	if err := tx.Where("user_id = ?", userId).First(&auth).Error; err != nil {
		tx.Rollback()
		return err
	}

	// обновляем токен
	auth.Key = token

	// сохраняем изменения
	if err := tx.Save(&auth).Error; err != nil {
		tx.Rollback()
		return err
	}

	// завершаем транзакцию
	return tx.Commit().Error
}

// подключение к базе данных
// пытается подключиться несколько раз с задержкой
func connectDB(config *config.Config) *gorm.DB {
	// формируем строку подключения
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		config.DB.Host,
		config.DB.Port,
		config.DB.User,
		config.DB.Name,
		config.DB.Password)

	var postrgresDB *gorm.DB
	var err error
	maxRetries := 5
	backoff := time.Second

	// пытаемся подключиться несколько раз(вдруг не получится сразу)
	for i := 0; i < maxRetries; i++ {
		postrgresDB, err = gorm.Open(postgres.Open(connectionString), &gorm.Config{})
		if err == nil {
			return postrgresDB
		}

		// ждем перед следующей попыткой
		if i < maxRetries-1 {
			time.Sleep(backoff)
		}
	}

	// если не удалось подключиться - паникуем
	panic(fmt.Sprintf("failed to connect database after %d retries: %v", maxRetries, err))
}
