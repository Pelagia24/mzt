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

func (r *CourseRepo) AddCourse(course *entity.Course) error {
	return r.DB.Create(course).Error
}

func (r *CourseRepo) DeleteCourse(courseId uuid.UUID) error {
	return r.DB.Delete(&entity.Course{}, "course_id = ?", courseId).Error
}

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

func (r *CourseRepo) AddLesson(lesson *entity.Lesson) error {
	return r.DB.Create(lesson).Error
}

func (r *CourseRepo) RemoveLesson(lessonId uuid.UUID) error {
	return r.DB.Delete(&entity.Lesson{}, "lesson_id = ?", lessonId).Error
}

func (r *CourseRepo) UpdateLesson(lesson *entity.Lesson) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingLesson entity.Lesson
	if err := tx.Where("lesson_id = ?", lesson.LessonID).First(&existingLesson).Error; err != nil {
		tx.Rollback()
		return err
	}

	existingLesson.Title = lesson.Title
	existingLesson.Summery = lesson.Summery
	existingLesson.VideoURL = lesson.VideoURL
	existingLesson.Text = lesson.Text

	if err := tx.Save(&existingLesson).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *CourseRepo) GetCourse(courseId uuid.UUID) (*dto.CourseDto, error) {
	var course entity.Course
	if err := r.DB.Preload("Price").First(&course, "course_id = ?", courseId).Error; err != nil {
		return nil, err
	}
	result := &dto.CourseDto{
		CourseID:    course.CourseID,
		Name:        course.Title,
		Description: course.Desc,
	}
	if course.Price != nil {
		result.Price.Amount = course.Price.Amount
		result.Price.CurrencyCode = course.Price.CurrencyCode
	}
	return result, nil
}

func (r *CourseRepo) GetLesson(lessonId uuid.UUID) (*entity.Lesson, error) {
	var lesson entity.Lesson
	if err := r.DB.First(&lesson, "lesson_id = ?", lessonId).Error; err != nil {
		return nil, err
	}
	return &lesson, nil
}

func (r *CourseRepo) GetCourses() ([]dto.CourseDto, error) {
	var courses []entity.Course
	if err := r.DB.Preload("Price").Find(&courses).Error; err != nil {
		return nil, err
	}

	result := make([]dto.CourseDto, 0)
	for _, course := range courses {
		courseDto := dto.CourseDto{
			CourseID:    course.CourseID,
			Name:        course.Title,
			Description: course.Desc,
		}
		if course.Price != nil {
			courseDto.Price.Amount = course.Price.Amount
			courseDto.Price.CurrencyCode = course.Price.CurrencyCode
		}
		result = append(result, courseDto)
	}
	return result, nil
}

func (r *CourseRepo) GetLessons() ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	if err := r.DB.Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

func (r *CourseRepo) CreateCourseAssignment(assignment *entity.CourseAssignment) error {
	return r.DB.Create(assignment).Error
}

func (r *CourseRepo) GetCourseAssignment(courseId, userId uuid.UUID) (*entity.CourseAssignment, error) {
	var assignment entity.CourseAssignment
	if err := r.DB.Where("course_id = ? AND user_id = ?", courseId, userId).First(&assignment).Error; err != nil {
		return nil, err
	}
	return &assignment, nil
}

func (r *CourseRepo) UpdateCourseAssignment(assignment *entity.CourseAssignment) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var existingAssignment entity.CourseAssignment
	if err := tx.Where("ca_id = ?", assignment.CaID).First(&existingAssignment).Error; err != nil {
		tx.Rollback()
		return err
	}

	existingAssignment.Progress = assignment.Progress

	if err := tx.Save(&existingAssignment).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *CourseRepo) DeleteCourseAssignment(courseId, userId uuid.UUID) error {
	return r.DB.Where("course_id = ? AND user_id = ?", courseId, userId).Delete(&entity.CourseAssignment{}).Error
}

func (r *CourseRepo) GetCourseAssignmentsByCourseId(courseId uuid.UUID) ([]entity.CourseAssignment, error) {
	var assignments []entity.CourseAssignment
	if err := r.DB.Where("course_id = ?", courseId).Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *CourseRepo) GetCourseAssignmentsByUserId(userId uuid.UUID) ([]entity.CourseAssignment, error) {
	var assignments []entity.CourseAssignment
	if err := r.DB.Preload("Course").Where("user_id = ?", userId).Find(&assignments).Error; err != nil {
		return nil, err
	}
	return assignments, nil
}

func (r *CourseRepo) GetLessonsByCourseId(courseId uuid.UUID) ([]entity.Lesson, error) {
	var lessons []entity.Lesson
	if err := r.DB.Where("course_id = ?", courseId).Find(&lessons).Error; err != nil {
		return nil, err
	}
	return lessons, nil
}

func (r *CourseRepo) GetUserWithDataById(userId uuid.UUID) (*entity.User, error) {
	var user entity.User
	if err := r.DB.Preload("UserData").First(&user, "id = ?", userId).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
