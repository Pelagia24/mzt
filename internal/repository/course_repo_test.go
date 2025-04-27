package repository

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCourseRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCourseRepo(&config.Config{DB: config.DB{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "postgres",
		Name:     "mzt_test",
	}})
	repo.DB = db


	t.Run("Create and Get Course", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}

		err := repo.AddCourse(course)
		require.NoError(t, err)

		gotCourse, err := repo.GetCourse(course.CourseID)
		require.NoError(t, err)
		assert.Equal(t, course.CourseID, gotCourse.CourseID)
		assert.Equal(t, course.Title, gotCourse.Name)
		assert.Equal(t, course.Desc, gotCourse.Description)
	})

	t.Run("Update Course", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Original Title",
			Desc:     "Original Description",
		}

		err := repo.AddCourse(course)
		require.NoError(t, err)

		updateDto := &dto.UpdateCourseDto{
			Name:        "Updated Title",
			Description: "Updated Description",
		}

		err = repo.UpdateCourse(course.CourseID, updateDto)
		require.NoError(t, err)

		gotCourse, err := repo.GetCourse(course.CourseID)
		require.NoError(t, err)
		assert.Equal(t, updateDto.Name, gotCourse.Name)
		assert.Equal(t, updateDto.Description, gotCourse.Description)
	})

	t.Run("Delete Course", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}

		err := repo.AddCourse(course)
		require.NoError(t, err)

		err = repo.DeleteCourse(course.CourseID)
		require.NoError(t, err)

		_, err = repo.GetCourse(course.CourseID)
		assert.Error(t, err)
	})

	t.Run("List Courses", func(t *testing.T) {
		courses := []*entity.Course{
			{
				CourseID: uuid.New(),
				Title:    "Course 1",
				Desc:     "Description 1",
			},
			{
				CourseID: uuid.New(),
				Title:    "Course 2",
				Desc:     "Description 2",
			},
		}

		for _, course := range courses {
			err := repo.AddCourse(course)
			require.NoError(t, err)
		}

		gotCourses, err := repo.GetCourses()
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(gotCourses), len(courses))
	})
}

func TestLessonRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCourseRepo(&config.Config{DB: config.DB{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "postgres",
		Name:     "mzt_test",
	}})
	repo.DB = db

	t.Run("Create and Get Lesson", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		lesson := &entity.Lesson{
			LessonID:   uuid.New(),
			CourseID:   course.CourseID,
			Title:      "Test Lesson",
			Desc:       "Test Description",
			VideoURL:   "http://example.com/video",
			SummaryURL: "http://example.com/summary",
		}

		err = repo.AddLesson(lesson)
		require.NoError(t, err)

		gotLesson, err := repo.GetLesson(lesson.LessonID)
		require.NoError(t, err)
		assert.Equal(t, lesson.LessonID, gotLesson.LessonID)
		assert.Equal(t, lesson.CourseID, gotLesson.CourseID)
		assert.Equal(t, lesson.Title, gotLesson.Title)
		assert.Equal(t, lesson.Desc, gotLesson.Desc)
		assert.Equal(t, lesson.VideoURL, gotLesson.VideoURL)
		assert.Equal(t, lesson.SummaryURL, gotLesson.SummaryURL)
	})

	t.Run("Update Lesson", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		lesson := &entity.Lesson{
			LessonID:   uuid.New(),
			CourseID:   course.CourseID,
			Title:      "Original Title",
			Desc:       "Original Description",
			VideoURL:   "http://example.com/video",
			SummaryURL: "http://example.com/summary",
		}

		err = repo.AddLesson(lesson)
		require.NoError(t, err)

		updatedLesson := &entity.Lesson{
			LessonID:   lesson.LessonID,
			CourseID:   lesson.CourseID,
			Title:      "Updated Title",
			Desc:       "Updated Description",
			VideoURL:   "http://example.com/video2",
			SummaryURL: "http://example.com/summary2",
		}

		err = repo.UpdateLesson(updatedLesson)
		require.NoError(t, err)

		gotLesson, err := repo.GetLesson(lesson.LessonID)
		require.NoError(t, err)
		assert.Equal(t, updatedLesson.Title, gotLesson.Title)
		assert.Equal(t, updatedLesson.Desc, gotLesson.Desc)
		assert.Equal(t, updatedLesson.VideoURL, gotLesson.VideoURL)
		assert.Equal(t, updatedLesson.SummaryURL, gotLesson.SummaryURL)
	})

	t.Run("Delete Lesson", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		lesson := &entity.Lesson{
			LessonID:   uuid.New(),
			CourseID:   course.CourseID,
			Title:      "Test Lesson",
			Desc:       "Test Description",
			VideoURL:   "http://example.com/video",
			SummaryURL: "http://example.com/summary",
		}

		err = repo.AddLesson(lesson)
		require.NoError(t, err)

		err = repo.RemoveLesson(lesson.LessonID)
		require.NoError(t, err)

		_, err = repo.GetLesson(lesson.LessonID)
		assert.Error(t, err)
	})

	t.Run("List Lessons by Course", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		lessons := []*entity.Lesson{
			{
				LessonID:   uuid.New(),
				CourseID:   course.CourseID,
				Title:      "Lesson 1",
				Desc:       "Description 1",
				VideoURL:   "http://example.com/video1",
				SummaryURL: "http://example.com/summary1",
			},
			{
				LessonID:   uuid.New(),
				CourseID:   course.CourseID,
				Title:      "Lesson 2",
				Desc:       "Description 2",
				VideoURL:   "http://example.com/video2",
				SummaryURL: "http://example.com/summary2",
			},
		}

		for _, lesson := range lessons {
			err := repo.AddLesson(lesson)
			require.NoError(t, err)
		}

		gotLessons, err := repo.GetLessonsByCourseId(course.CourseID)
		require.NoError(t, err)
		assert.Equal(t, len(lessons), len(gotLessons))
	})
}
