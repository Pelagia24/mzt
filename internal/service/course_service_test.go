package service

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/mocks"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCourseService_CreateCourse(t *testing.T) {
	mockRepo := mocks.NewMockCourseRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:       "test-access-key",
			RefreshKey:      "test-refresh-key",
			AccessExpiresIn: time.Minute * 30,
		},
	}
	service := NewCourseService(cfg, mockRepo)
	courseDto := &dto.CreateCourseDto{
		Name:        "Test Course",
		Description: "Test Description",
		Price:       100,
	}

	err := service.CreateCourse(courseDto)

	assert.NoError(t, err)
	courses, err := service.ListCourses()
	assert.NoError(t, err)
	assert.Len(t, courses, 1)
	assert.Equal(t, courseDto.Name, courses[0].Name)
	assert.Equal(t, courseDto.Description, courses[0].Description)
}

func TestCourseService_GetCourse(t *testing.T) {
	mockRepo := mocks.NewMockCourseRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:       "test-access-key",
			RefreshKey:      "test-refresh-key",
			AccessExpiresIn: time.Minute * 30,
		},
	}
	service := NewCourseService(cfg, mockRepo)
	courseDto := &dto.CreateCourseDto{
		Name:        "Test Course",
		Description: "Test Description",
		Price:       100,
	}
	err := service.CreateCourse(courseDto)
	assert.NoError(t, err)

	courses, err := service.ListCourses()
	assert.NoError(t, err)
	courseId := courses[0].CourseID

	course, err := service.GetCourse(courseId)

	assert.NoError(t, err)
	assert.NotNil(t, course)
	assert.Equal(t, courseDto.Name, course.Name)
	assert.Equal(t, courseDto.Description, course.Description)
}

func TestCourseService_UpdateCourse(t *testing.T) {
	mockRepo := mocks.NewMockCourseRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:       "test-access-key",
			RefreshKey:      "test-refresh-key",
			AccessExpiresIn: time.Minute * 30,
		},
	}
	service := NewCourseService(cfg, mockRepo)
	courseDto := &dto.CreateCourseDto{
		Name:        "Test Course",
		Description: "Test Description",
		Price:       100,
	}
	err := service.CreateCourse(courseDto)
	assert.NoError(t, err)

	courses, err := service.ListCourses()
	assert.NoError(t, err)
	courseId := courses[0].CourseID

	updatedDto := &dto.UpdateCourseDto{
		Name:        "Updated Course",
		Description: "Updated Description",
		Price:       200,
	}

	err = service.UpdateCourse(courseId, updatedDto)

	assert.NoError(t, err)
	course, err := service.GetCourse(courseId)
	assert.NoError(t, err)
	assert.Equal(t, updatedDto.Name, course.Name)
	assert.Equal(t, updatedDto.Description, course.Description)
}

func TestCourseService_DeleteCourse(t *testing.T) {
	mockRepo := mocks.NewMockCourseRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:       "test-access-key",
			RefreshKey:      "test-refresh-key",
			AccessExpiresIn: time.Minute * 30,
		},
	}
	service := NewCourseService(cfg, mockRepo)
	courseDto := &dto.CreateCourseDto{
		Name:        "Test Course",
		Description: "Test Description",
		Price:       100,
	}
	err := service.CreateCourse(courseDto)
	assert.NoError(t, err)

	courses, err := service.ListCourses()
	assert.NoError(t, err)
	courseId := courses[0].CourseID

	err = service.DeleteCourse(courseId)

	assert.NoError(t, err)
	courses, err = service.ListCourses()
	assert.NoError(t, err)
	assert.Len(t, courses, 0)
}

func TestCourseService_EnrollUserToCourse(t *testing.T) {
	mockRepo := mocks.NewMockCourseRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:       "test-access-key",
			RefreshKey:      "test-refresh-key",
			AccessExpiresIn: time.Minute * 30,
		},
	}
	service := NewCourseService(cfg, mockRepo)
	courseDto := &dto.CreateCourseDto{
		Name:        "Test Course",
		Description: "Test Description",
		Price:       100,
	}
	err := service.CreateCourse(courseDto)
	assert.NoError(t, err)

	courses, err := service.ListCourses()
	assert.NoError(t, err)
	courseId := courses[0].CourseID
	userId := uuid.New()

	err = service.AssignUserToCourse(courseId, userId)

	assert.NoError(t, err)
	users, err := service.ListUsersOnCourse(courseId)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.Equal(t, userId, users[0].ID)
}

func TestCourseService_UpdateProgress(t *testing.T) {
	mockRepo := mocks.NewMockCourseRepository()
	cfg := &config.Config{
		Jwt: config.Jwt{
			AccessKey:       "test-access-key",
			RefreshKey:      "test-refresh-key",
			AccessExpiresIn: time.Minute * 30,
		},
	}
	service := NewCourseService(cfg, mockRepo)
	courseDto := &dto.CreateCourseDto{
		Name:        "Test Course",
		Description: "Test Description",
		Price:       100,
	}
	err := service.CreateCourse(courseDto)
	assert.NoError(t, err)

	courses, err := service.ListCourses()
	assert.NoError(t, err)
	courseId := courses[0].CourseID
	userId := uuid.New()

	err = service.AssignUserToCourse(courseId, userId)
	assert.NoError(t, err)

	err = service.UpdateProgress(courseId, userId, 50)

	assert.NoError(t, err)
	progress, err := service.GetProgress(courseId, userId)
	assert.NoError(t, err)
	assert.Equal(t, uint(50), progress)
}
