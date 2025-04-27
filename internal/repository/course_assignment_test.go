package repository

import (
	"mzt/config"
	"mzt/internal/entity"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCourseAssignmentRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewCourseRepo(&config.Config{DB: config.DB{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "postgres",
		Name:     "mzt_test",
	}})
	repo.DB = db

	t.Run("Create and Get Course Assignment", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		user := &entity.User{
			ID:         uuid.New(),
			PasswdHash: "test_hash",
			Role:       1,
		}
		err = db.Create(user).Error
		require.NoError(t, err)

		assignment := &entity.CourseAssignment{
			CaID:     uuid.New(),
			UserID:   user.ID,
			CourseID: course.CourseID,
			Progress: 0,
		}

		err = repo.CreateCourseAssignment(assignment)
		require.NoError(t, err)

		gotAssignment, err := repo.GetCourseAssignment(course.CourseID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, assignment.CaID, gotAssignment.CaID)
		assert.Equal(t, assignment.UserID, gotAssignment.UserID)
		assert.Equal(t, assignment.CourseID, gotAssignment.CourseID)
		assert.Equal(t, assignment.Progress, gotAssignment.Progress)
	})

	t.Run("Update Course Assignment Progress", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		user := &entity.User{
			ID:         uuid.New(),
			PasswdHash: "test_hash",
			Role:       1,
		}
		err = db.Create(user).Error
		require.NoError(t, err)

		assignment := &entity.CourseAssignment{
			CaID:     uuid.New(),
			UserID:   user.ID,
			CourseID: course.CourseID,
			Progress: 0,
		}

		err = repo.CreateCourseAssignment(assignment)
		require.NoError(t, err)

		assignment.Progress = 50
		err = repo.UpdateCourseAssignment(assignment)
		require.NoError(t, err)

		gotAssignment, err := repo.GetCourseAssignment(course.CourseID, user.ID)
		require.NoError(t, err)
		assert.Equal(t, uint(50), gotAssignment.Progress)
	})

	t.Run("Delete Course Assignment", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		user := &entity.User{
			ID:         uuid.New(),
			PasswdHash: "test_hash",
			Role:       1,
		}
		err = db.Create(user).Error
		require.NoError(t, err)

		assignment := &entity.CourseAssignment{
			CaID:     uuid.New(),
			UserID:   user.ID,
			CourseID: course.CourseID,
			Progress: 0,
		}

		err = repo.CreateCourseAssignment(assignment)
		require.NoError(t, err)

		err = repo.DeleteCourseAssignment(course.CourseID, user.ID)
		require.NoError(t, err)

		_, err = repo.GetCourseAssignment(course.CourseID, user.ID)
		assert.Error(t, err)
	})

	t.Run("List Course Assignments", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := repo.AddCourse(course)
		require.NoError(t, err)

		users := []*entity.User{
			{
				ID:         uuid.New(),
				PasswdHash: "test_hash1",
				Role:       1,
			},
			{
				ID:         uuid.New(),
				PasswdHash: "test_hash2",
				Role:       1,
			},
		}

		for _, user := range users {
			err = db.Create(user).Error
			require.NoError(t, err)
		}

		for _, user := range users {
			assignment := &entity.CourseAssignment{
				CaID:     uuid.New(),
				UserID:   user.ID,
				CourseID: course.CourseID,
				Progress: 0,
			}
			err = repo.CreateCourseAssignment(assignment)
			require.NoError(t, err)
		}

		assignments, err := repo.GetCourseAssignmentsByCourseId(course.CourseID)
		require.NoError(t, err)
		assert.Equal(t, len(users), len(assignments))
	})
}
