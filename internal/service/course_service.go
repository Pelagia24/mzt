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

func (s *CourseService) ListCourses() ([]dto.CourseDto, error) {
	return s.repo.GetCourses()
}

func (s *CourseService) GetCourse(courseId uuid.UUID) (*dto.CourseDto, error) {
	return s.repo.GetCourse(courseId)
}

func (s *CourseService) CreateCourse(course *dto.CreateCourseDto) error {
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

func (s *CourseService) UpdateCourse(courseId uuid.UUID, updated *dto.UpdateCourseDto) error {
	return s.repo.UpdateCourse(courseId, updated)
}

func (s *CourseService) DeleteCourse(courseId uuid.UUID) error {
	return s.repo.DeleteCourse(courseId)
}

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

func (s *CourseService) DeleteLesson(lessonId uuid.UUID) error {
	return s.repo.RemoveLesson(lessonId)
}

func (s *CourseService) AssignUserToCourse(courseId uuid.UUID, userId uuid.UUID) error {
	assignment := &entity.CourseAssignment{
		CaID:     uuid.New(),
		UserID:   userId,
		CourseID: courseId,
		Progress: 0,
	}
	return s.repo.CreateCourseAssignment(assignment)
}

func (s *CourseService) ListUserCourses(userId uuid.UUID) ([]dto.CourseDto, error) {
	assignments, err := s.repo.GetCourseAssignmentsByUserId(userId)
	if err != nil {
		return nil, err
	}
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

func (s *CourseService) ListUsersOnCourse(courseId uuid.UUID) ([]dto.UserInfoAdminDto, error) {
	assignments, err := s.repo.GetCourseAssignmentsByCourseId(courseId)
	if err != nil {
		return nil, err
	}

	users := make([]dto.UserInfoAdminDto, 0)
	for _, assignment := range assignments {
		users = append(users, dto.UserInfoAdminDto{
			ID: assignment.UserID,
		})
	}
	return users, nil
}

func (s *CourseService) RemoveUserFromCourse(courseId uuid.UUID, userId uuid.UUID) error {
	return s.repo.DeleteCourseAssignment(courseId, userId)
}

func (s *CourseService) GetProgress(courseId uuid.UUID, userId uuid.UUID) (uint, error) {
	assignment, err := s.repo.GetCourseAssignment(courseId, userId)
	if err != nil {
		return 0, err
	}
	return assignment.Progress, nil
}

func (s *CourseService) UpdateProgress(courseId uuid.UUID, userId uuid.UUID, progress uint) error {
	assignment, err := s.repo.GetCourseAssignment(courseId, userId)
	if err != nil {
		return err
	}
	assignment.Progress = progress
	return s.repo.UpdateCourseAssignment(assignment)
}
