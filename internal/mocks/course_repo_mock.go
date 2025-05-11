package mocks

import (
	"mzt/internal/dto"
	"mzt/internal/entity"
	"mzt/internal/repository"

	"github.com/google/uuid"
)

type MockCourseRepository struct {
	Courses     map[uuid.UUID]*entity.Course
	Lessons     map[uuid.UUID]*entity.Lesson
	Assignments map[uuid.UUID]map[uuid.UUID]*entity.CourseAssignment
}

func NewMockCourseRepository() repository.CourseRepository {
	return &MockCourseRepository{
		Courses:     make(map[uuid.UUID]*entity.Course),
		Lessons:     make(map[uuid.UUID]*entity.Lesson),
		Assignments: make(map[uuid.UUID]map[uuid.UUID]*entity.CourseAssignment),
	}
}

func (m *MockCourseRepository) GetCourses() ([]dto.CourseDto, error) {
	courses := make([]dto.CourseDto, 0, len(m.Courses))
	for _, course := range m.Courses {
		courses = append(courses, dto.CourseDto{
			CourseID:    course.CourseID,
			Name:        course.Title,
			Description: course.Desc,
		})
	}
	return courses, nil
}

func (m *MockCourseRepository) GetCourse(courseId uuid.UUID) (*dto.CourseDto, error) {
	if course, exists := m.Courses[courseId]; exists {
		return &dto.CourseDto{
			CourseID:    course.CourseID,
			Name:        course.Title,
			Description: course.Desc,
		}, nil
	}
	return nil, nil
}

func (m *MockCourseRepository) AddCourse(course *entity.Course) error {
	m.Courses[course.CourseID] = course
	return nil
}

func (m *MockCourseRepository) UpdateCourse(courseId uuid.UUID, updated *dto.UpdateCourseDto) error {
	if course, exists := m.Courses[courseId]; exists {
		course.Title = updated.Name
		course.Desc = updated.Description
		return nil
	}
	return nil
}

func (m *MockCourseRepository) DeleteCourse(courseId uuid.UUID) error {
	delete(m.Courses, courseId)
	return nil
}

func (m *MockCourseRepository) GetLessonsByCourseId(courseId uuid.UUID) ([]entity.Lesson, error) {
	lessons := make([]entity.Lesson, 0)
	for _, lesson := range m.Lessons {
		if lesson.CourseID == courseId {
			lessons = append(lessons, *lesson)
		}
	}
	return lessons, nil
}

func (m *MockCourseRepository) GetLesson(lessonId uuid.UUID) (*entity.Lesson, error) {
	if lesson, exists := m.Lessons[lessonId]; exists {
		return lesson, nil
	}
	return nil, nil
}

func (m *MockCourseRepository) AddLesson(lesson *entity.Lesson) error {
	m.Lessons[lesson.LessonID] = lesson
	return nil
}

func (m *MockCourseRepository) UpdateLesson(lesson *entity.Lesson) error {
	if _, exists := m.Lessons[lesson.LessonID]; exists {
		m.Lessons[lesson.LessonID] = lesson
		return nil
	}
	return nil
}

func (m *MockCourseRepository) RemoveLesson(lessonId uuid.UUID) error {
	delete(m.Lessons, lessonId)
	return nil
}

func (m *MockCourseRepository) CreateCourseAssignment(assignment *entity.CourseAssignment) error {
	if _, exists := m.Assignments[assignment.CourseID]; !exists {
		m.Assignments[assignment.CourseID] = make(map[uuid.UUID]*entity.CourseAssignment)
	}
	m.Assignments[assignment.CourseID][assignment.UserID] = assignment
	return nil
}

func (m *MockCourseRepository) GetCourseAssignments(courseId uuid.UUID) ([]entity.CourseAssignment, error) {
	assignments := make([]entity.CourseAssignment, 0)
	if courseAssignments, exists := m.Assignments[courseId]; exists {
		for _, assignment := range courseAssignments {
			assignments = append(assignments, *assignment)
		}
	}
	return assignments, nil
}

func (m *MockCourseRepository) GetCourseAssignmentsByCourseId(courseId uuid.UUID) ([]entity.CourseAssignment, error) {
	assignments := make([]entity.CourseAssignment, 0)
	if courseAssignments, exists := m.Assignments[courseId]; exists {
		for _, assignment := range courseAssignments {
			assignments = append(assignments, *assignment)
		}
	}
	return assignments, nil
}

func (m *MockCourseRepository) GetCourseAssignmentsByUserId(userId uuid.UUID) ([]entity.CourseAssignment, error) {
	assignments := make([]entity.CourseAssignment, 0)
	if courseAssignments, exists := m.Assignments[userId]; exists {
		for _, assignment := range courseAssignments {
			assignments = append(assignments, *assignment)
		}
	}
	return assignments, nil
}

func (m *MockCourseRepository) GetCourseAssignment(courseId uuid.UUID, userId uuid.UUID) (*entity.CourseAssignment, error) {
	if courseAssignments, exists := m.Assignments[courseId]; exists {
		if assignment, exists := courseAssignments[userId]; exists {
			return assignment, nil
		}
	}
	return nil, nil
}
func (m *MockCourseRepository) UpdateCourseAssignment(assignment *entity.CourseAssignment) error {
	if courseAssignments, exists := m.Assignments[assignment.CourseID]; exists {
		courseAssignments[assignment.UserID] = assignment
		return nil
	}
	return nil
}

func (m *MockCourseRepository) DeleteCourseAssignment(courseId uuid.UUID, userId uuid.UUID) error {
	if courseAssignments, exists := m.Assignments[courseId]; exists {
		delete(courseAssignments, userId)
	}
	return nil
}
