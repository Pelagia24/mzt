package service

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"mzt/internal/repository"

	"github.com/google/uuid"
)

type CourseServiceInterface interface {
	ListCourses() ([]dto.CourseDto, error)
	GetCourse(courseId uuid.UUID) (*dto.CourseDto, error)
	CreateCourse(course *dto.CreateCourseDto) error
	UpdateCourse(courseId uuid.UUID, updated *dto.UpdateCourseDto) error
	DeleteCourse(courseId uuid.UUID) error

	ListLessons(courseId uuid.UUID) ([]dto.LessonDto, error)
	GetLesson(lessonId uuid.UUID) (*dto.LessonDto, error)
	CreateLesson(courseId uuid.UUID, lesson *dto.CreateLessonDto) error
	UpdateLesson(lessonId uuid.UUID, updated *dto.UpdateLessonDto) error
	DeleteLesson(lessonId uuid.UUID) error

	AssignUserToCourse(courseId uuid.UUID, userId uuid.UUID) error
	ListUsersOnCourse(courseId uuid.UUID) ([]dto.UserInfoAdminDto, error)
	RemoveUserFromCourse(courseId uuid.UUID, userId uuid.UUID) error
	GetProgress(courseId uuid.UUID, userId uuid.UUID) (uint, error)
	UpdateProgress(courseId uuid.UUID, userId uuid.UUID, progress uint) error
}

type CourseService struct {
	config *config.Config
	repo   repository.CourseRepository
}

func NewCourseService(cfg *config.Config, repo repository.CourseRepository) *CourseService {
	return &CourseService{
		config: cfg,
		repo:   repo,
	}
}

// ListCourses получает список всех курсов
// просто берет все курсы из базы и возвращает их
func (s *CourseService) ListCourses() ([]dto.CourseDto, error) {
	return s.repo.GetCourses()
}

// GetCourse получает информацию о курсе
// просто берет курс из базы по его id
func (s *CourseService) GetCourse(courseId uuid.UUID) (*dto.CourseDto, error) {
	return s.repo.GetCourse(courseId)
}

// CreateCourse создает новый курс
// создает новый курс в базе с указанными данными
func (s *CourseService) CreateCourse(course *dto.CreateCourseDto) error {
	// создаем новый курс с уникальным id
	courseEntity := &entity.Course{
		CourseID: uuid.New(),
		Title:    course.Name,
		Desc:     course.Description,
		Price: &entity.CoursePrice{
			Amount:       float64(course.Price),
			CurrencyCode: "RUB",
		},
	}
	return s.repo.AddCourse(courseEntity)
}

// UpdateCourse обновляет информацию о курсе
// меняет название описание и цену курса
func (s *CourseService) UpdateCourse(courseId uuid.UUID, updated *dto.UpdateCourseDto) error {
	return s.repo.UpdateCourse(courseId, updated)
}

// DeleteCourse удаляет курс
// просто удаляет курс из базы по его id
func (s *CourseService) DeleteCourse(courseId uuid.UUID) error {
	return s.repo.DeleteCourse(courseId)
}

// ListLessons получает список всех уроков курса
// просто берет все уроки из базы и преобразует их в формат для response
func (s *CourseService) ListLessons(courseId uuid.UUID) ([]dto.LessonDto, error) {
	lessons, err := s.repo.GetLessonsByCourseId(courseId)
	if err != nil {
		return nil, err
	}

	result := make([]dto.LessonDto, 0)
	for _, lesson := range lessons {
		result = append(result, dto.LessonDto{
			LessonID: lesson.LessonID,
			CourseID: lesson.CourseID,
			Title:    lesson.Title,
			Summery:  lesson.Summery,
			VideoURL: lesson.VideoURL,
			Text:     lesson.Text,
		})
	}
	return result, nil
}

// GetLesson получает информацию об уроке
// просто берет урок из базы и преобразует его в формат для response
func (s *CourseService) GetLesson(lessonId uuid.UUID) (*dto.LessonDto, error) {
	lesson, err := s.repo.GetLesson(lessonId)
	if err != nil {
		return nil, err
	}

	return &dto.LessonDto{
		LessonID: lesson.LessonID,
		CourseID: lesson.CourseID,
		Title:    lesson.Title,
		Summery:  lesson.Summery,
		VideoURL: lesson.VideoURL,
		Text:     lesson.Text,
	}, nil
}

// CreateLesson создает новый урок
// просто создает новый урок в базе с указанным id курса
func (s *CourseService) CreateLesson(courseId uuid.UUID, lesson *dto.CreateLessonDto) error {
	lessonEntity := &entity.Lesson{
		LessonID: uuid.New(),
		CourseID: courseId,
		Title:    lesson.Title,
		Summery:  lesson.Description,
		VideoURL: lesson.VideoURL,
		Text:     lesson.SummaryURL,
	}
	return s.repo.AddLesson(lessonEntity)
}

// UpdateLesson обновляет информацию об уроке
// просто обновляет данные урока в базе
func (s *CourseService) UpdateLesson(lessonId uuid.UUID, updated *dto.UpdateLessonDto) error {
	lessonEntity := &entity.Lesson{
		LessonID: lessonId,
		Title:    updated.Title,
		Summery:  updated.Description,
		VideoURL: updated.VideoURL,
		Text:     updated.SummaryURL,
	}
	return s.repo.UpdateLesson(lessonEntity)
}

// DeleteLesson удаляет урок
// просто удаляет урок из базы
func (s *CourseService) DeleteLesson(lessonId uuid.UUID) error {
	return s.repo.RemoveLesson(lessonId)
}

// AssignUserToCourse записывает пользователя на курс
// создает новую запись о том что пользователь записан на курс
func (s *CourseService) AssignUserToCourse(courseId uuid.UUID, userId uuid.UUID) error {
	assignment := &entity.CourseAssignment{
		CaID:     uuid.New(),
		UserID:   userId,
		CourseID: courseId,
		Progress: 0,
	}
	return s.repo.CreateCourseAssignment(assignment)
}

// ListUserCourses получает список курсов пользователя
// берет все курсы на которые записан пользователь и преобразует их в формат для response
func (s *CourseService) ListUserCourses(userId uuid.UUID) ([]dto.CourseDto, error) {
	// получаем все записи о курсах пользователя
	assignments, err := s.repo.GetCourseAssignmentsByUserId(userId)
	if err != nil {
		return nil, err
	}

	// преобразуем каждую запись в формат для response
	courses := make([]dto.CourseDto, 0)
	for _, assignment := range assignments {
		courses = append(courses, dto.CourseDto{
			CourseID:    assignment.CourseID,
			Name:        assignment.Course.Title,
			Description: assignment.Course.Desc,
		})
	}

	return courses, nil
}

// ListUsersOnCourse получает список пользователей на курсе
// берет всех пользователей записанных на курс и преобразует их в формат для response
func (s *CourseService) ListUsersOnCourse(courseId uuid.UUID) ([]dto.UserInfoAdminDto, error) {
	// получаем все записи о пользователях на курсе
	assignments, err := s.repo.GetCourseAssignmentsByCourseId(courseId)
	if err != nil {
		return nil, err
	}

	// преобразуем каждую запись в формат для response
	users := make([]dto.UserInfoAdminDto, 0)
	for _, assignment := range assignments {
		users = append(users, dto.UserInfoAdminDto{
			ID: assignment.UserID,
		})
	}
	return users, nil
}

// RemoveUserFromCourse отписывает пользователя от курса(пока не используется на frontend)
// просто удаляет запись о том что пользователь записан на курс
func (s *CourseService) RemoveUserFromCourse(courseId uuid.UUID, userId uuid.UUID) error {
	return s.repo.DeleteCourseAssignment(courseId, userId)
}

// GetProgress получает прогресс пользователя по курсу
// просто берет значение прогресса из базы
func (s *CourseService) GetProgress(courseId uuid.UUID, userId uuid.UUID) (uint, error) {
	// получаем запись о прогрессе пользователя
	assignment, err := s.repo.GetCourseAssignment(courseId, userId)
	if err != nil {
		return 0, err
	}
	return assignment.Progress, nil
}

// UpdateProgress обновляет прогресс пользователя по курсу
// просто меняет значение прогресса в базе
func (s *CourseService) UpdateProgress(courseId uuid.UUID, userId uuid.UUID, progress uint) error {
	// получаем запись о прогрессе пользователя
	assignment, err := s.repo.GetCourseAssignment(courseId, userId)
	if err != nil {
		return err
	}
	// обновляем значение прогресса
	assignment.Progress = progress
	return s.repo.UpdateCourseAssignment(assignment)
}
