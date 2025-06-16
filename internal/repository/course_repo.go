package repository

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CourseRepository interface {
	GetCourses() ([]dto.CourseDto, error)
	GetCourse(courseId uuid.UUID) (*dto.CourseDto, error)
	AddCourse(course *entity.Course) error
	UpdateCourse(courseId uuid.UUID, updated *dto.UpdateCourseDto) error
	DeleteCourse(courseId uuid.UUID) error
	GetLessonsByCourseId(courseId uuid.UUID) ([]entity.Lesson, error)
	GetLesson(lessonId uuid.UUID) (*entity.Lesson, error)
	AddLesson(lesson *entity.Lesson) error
	UpdateLesson(lesson *entity.Lesson) error
	RemoveLesson(lessonId uuid.UUID) error
	CreateCourseAssignment(assignment *entity.CourseAssignment) error
	GetCourseAssignmentsByCourseId(courseId uuid.UUID) ([]entity.CourseAssignment, error)
	GetCourseAssignmentsByUserId(userId uuid.UUID) ([]entity.CourseAssignment, error)
	GetCourseAssignment(courseId, userId uuid.UUID) (*entity.CourseAssignment, error)
	UpdateCourseAssignment(assignment *entity.CourseAssignment) error
	DeleteCourseAssignment(courseId uuid.UUID, userId uuid.UUID) error
}

type CourseRepo struct {
	config *config.Config
	DB     *gorm.DB
}

func NewCourseRepo(cfg *config.Config) *CourseRepo {
	return &CourseRepo{
		config: cfg,
		DB:     connectDB(cfg),
	}
}

// AddCourse добавляет новый курс в базу
// просто создает новую запись в таблице courses
func (r *CourseRepo) AddCourse(course *entity.Course) error {
	return r.DB.Create(course).Error
}

// DeleteCourse удаляет курс из базы
// просто удаляет запись из таблицы courses по id курса
func (r *CourseRepo) DeleteCourse(courseId uuid.UUID) error {
	return r.DB.Delete(&entity.Course{}, "course_id = ?", courseId).Error
}

// UpdateCourse обновляет информацию о курсе
// меняет название описание и цену курса в базе
func (r *CourseRepo) UpdateCourse(courseId uuid.UUID, updated *dto.UpdateCourseDto) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var course entity.Course
	if err := tx.Where("course_id = ?", courseId).First(&course).Error; err != nil {
		tx.Rollback()
		return err
	}

	course.Title = updated.Name
	course.Desc = updated.Description

	if err := tx.Save(&course).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update or create price
	var price entity.CoursePrice
	if err := tx.Where("course_id = ?", courseId).First(&price).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			price = entity.CoursePrice{
				CourseID:     courseId,
				Amount:       float64(updated.Price),
				CurrencyCode: "RUB",
			}
			if err := tx.Create(&price).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	} else {
		price.Amount = float64(updated.Price)
		if err := tx.Save(&price).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

// AddLesson добавляет новый урок в базу
// просто создает новую запись в таблице lessons
func (r *CourseRepo) AddLesson(lesson *entity.Lesson) error {
	return r.DB.Create(lesson).Error
}

// RemoveLesson удаляет урок из базы
// просто удаляет запись из таблицы lessons по id урока
func (r *CourseRepo) RemoveLesson(lessonId uuid.UUID) error {
	return r.DB.Delete(&entity.Lesson{}, "lesson_id = ?", lessonId).Error
}

// UpdateLesson обновляет информацию об уроке
// меняет название описание видео и текст урока в базе
func (r *CourseRepo) UpdateLesson(lesson *entity.Lesson) error {
	// начинаем транзакцию чтобы все изменения сохранились вместе
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// ищем урок который хотим обновить
	var existingLesson entity.Lesson
	if err := tx.Where("lesson_id = ?", lesson.LessonID).First(&existingLesson).Error; err != nil {
		tx.Rollback()
		return err
	}

	// обновляем все поля урока
	existingLesson.Title = lesson.Title
	existingLesson.Summery = lesson.Summery
	existingLesson.VideoURL = lesson.VideoURL
	existingLesson.Text = lesson.Text

	// сохраняем изменения
	if err := tx.Save(&existingLesson).Error; err != nil {
		tx.Rollback()
		return err
	}

	// подтверждаем все изменения
	return tx.Commit().Error
}

// GetCourse получает информацию о курсе
// ищет курс в базе по id и возвращает его данные
func (r *CourseRepo) GetCourse(courseId uuid.UUID) (*dto.CourseDto, error) {
	var course entity.Course
	// загружаем курс вместе с его ценой
	if err := r.DB.Preload("Price").First(&course, "course_id = ?", courseId).Error; err != nil {
		return nil, err
	}
	// преобразуем данные курса в формат для response
	result := &dto.CourseDto{
		CourseID:    course.CourseID,
		Name:        course.Title,
		Description: course.Desc,
	}
	// если у курса есть цена добавляем ее в response
	if course.Price != nil {
		result.Price.Amount = course.Price.Amount
		result.Price.CurrencyCode = course.Price.CurrencyCode
	}
	return result, nil
}

// GetLesson получает информацию об уроке
// ищет урок в базе по id и возвращает его данные
func (r *CourseRepo) GetLesson(lessonId uuid.UUID) (*entity.Lesson, error) {
	var lesson entity.Lesson
	if err := r.DB.First(&lesson, "lesson_id = ?", lessonId).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

// GetCourses получает список всех курсов
// берет все курсы из базы и возвращает их данные
func (r *CourseRepo) GetCourses() ([]dto.CourseDto, error) {
	var courses []entity.Course
	// загружаем курсы вместе с их ценами
	if err := r.DB.Preload("Price").Find(&courses).Error; err != nil {
		return nil, err
	}

	// преобразуем каждый курс в формат для response
	result := make([]dto.CourseDto, 0)
	for _, course := range courses {
		courseDto := dto.CourseDto{
			CourseID:    course.CourseID,
			Name:        course.Title,
			Description: course.Desc,
		}
		// если у курса есть цена добавляем ее в response
		if course.Price != nil {
			courseDto.Price.Amount = course.Price.Amount
			courseDto.Price.CurrencyCode = course.Price.CurrencyCode
		}
		result = append(result, courseDto)
	}
	return result, nil
}

// GetLessons получает список всех уроков
// просто берет все уроки из базы
func (r *CourseRepo) GetLessons() ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	if err := r.DB.Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

// CreateCourseAssignment создает запись о том что пользователь записан на курс
// просто создает новую запись в таблице course_assignments
func (r *CourseRepo) CreateCourseAssignment(assignment *entity.CourseAssignment) error {
	return r.DB.Create(assignment).Error
}

// GetCourseAssignment получает информацию о записи пользователя на курс
// ищет запись в базе по id курса и id пользователя
func (r *CourseRepo) GetCourseAssignment(courseId, userId uuid.UUID) (*entity.CourseAssignment, error) {
	var assignment entity.CourseAssignment
	if err := r.DB.Where("course_id = ? AND user_id = ?", courseId, userId).First(&assignment).Error; err != nil {
		return nil, err
	}
	return &assignment, nil
}

// UpdateCourseAssignment обновляет прогресс пользователя по курсу
// меняет значение прогресса в базе
func (r *CourseRepo) UpdateCourseAssignment(assignment *entity.CourseAssignment) error {
	// начинаем транзакцию чтобы все изменения сохранились вместе
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// ищем запись которую хотим обновить
	var existingAssignment entity.CourseAssignment
	if err := tx.Where("ca_id = ?", assignment.CaID).First(&existingAssignment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// обновляем прогресс
	existingAssignment.Progress = assignment.Progress

	// сохраняем изменения
	if err := tx.Save(&existingAssignment).Error; err != nil {
		tx.Rollback()
		return err
	}

	// подтверждаем все изменения
	return tx.Commit().Error
}

// DeleteCourseAssignment удаляет запись о том что пользователь записан на курс
// просто удаляет запись из таблицы course_assignments
func (r *CourseRepo) DeleteCourseAssignment(courseId, userId uuid.UUID) error {
	return r.DB.Where("course_id = ? AND user_id = ?", courseId, userId).Delete(&entity.CourseAssignment{}).Error
}

// GetCourseAssignmentsByCourseId получает список всех записей на курс
// ищет все записи в базе по id курса
func (r *CourseRepo) GetCourseAssignmentsByCourseId(courseId uuid.UUID) ([]entity.CourseAssignment, error) {
	var assignments []entity.CourseAssignment
	if err := r.DB.Where("course_id = ?", courseId).Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

// GetCourseAssignmentsByUserId получает список всех курсов пользователя
// ищет все записи в базе по id пользователя и загружает информацию о курсах
func (r *CourseRepo) GetCourseAssignmentsByUserId(userId uuid.UUID) ([]entity.CourseAssignment, error) {
	var assignments []entity.CourseAssignment
	if err := r.DB.Preload("Course").Where("user_id = ?", userId).Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

// GetLessonsByCourseId получает список всех уроков курса
// ищет все уроки в базе по id курса
func (r *CourseRepo) GetLessonsByCourseId(courseId uuid.UUID) ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	if err := r.DB.Where("course_id = ?", courseId).Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

// GetUserWithDataById получает информацию о пользователе
// ищет пользователя в базе по id и загружает его данные
func (r *CourseRepo) GetUserWithDataById(userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Preload("UserData").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
