package repository

import (
	"mzt/config"
	"mzt/internal/entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreloadRelationships(t *testing.T) {
	db := setupTestDB(t)
	userRepo := NewUserRepo(&config.Config{DB: config.DB{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "postgres",
		Name:     "mzt_test",
	}})
	userRepo.DB = db

	courseRepo := NewCourseRepo(&config.Config{DB: config.DB{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "postgres",
		Name:     "mzt_test",
	}})
	courseRepo.DB = db

	t.Run("Test User with UserData preload", func(t *testing.T) {
		userId := uuid.New()
		user := &entity.User{
			ID:         userId,
			PasswdHash: "test_hash",
			Role:       1,
		}
		userData := &entity.UserData{
			UserID:      userId,
			Email:       "test@example.com",
			Name:        "Test User",
			Birthdate:   time.Now(),
			PhoneNumber: "+1234567890",
		}
		auth := &entity.Auth{
			UserID: userId,
			Key:    "test_key",
		}

		err := userRepo.CreateUser(user, userData, auth)
		require.NoError(t, err)

		gotUser, err := userRepo.GetUserWithDataById(userId)
		require.NoError(t, err)
		assert.NotNil(t, gotUser.UserData)
		assert.Equal(t, userData.Email, gotUser.UserData.Email)
		assert.Equal(t, userData.Name, gotUser.UserData.Name)
	})

	t.Run("Test User with Auth preload", func(t *testing.T) {
		userId := uuid.New()
		user := &entity.User{
			ID:         userId,
			PasswdHash: "test_hash",
			Role:       1,
		}
		userData := &entity.UserData{
			UserID:      userId,
			Email:       "test2@example.com",
			Name:        "Test User 2",
			Birthdate:   time.Now(),
			PhoneNumber: "+1234567891",
		}
		auth := &entity.Auth{
			UserID: userId,
			Key:    "test_key_2",
		}

		err := userRepo.CreateUser(user, userData, auth)
		require.NoError(t, err)

		gotUser, err := userRepo.GetUserWithRefreshById(userId)
		require.NoError(t, err)
		assert.NotNil(t, gotUser.Auth)
		assert.Equal(t, auth.Key, gotUser.Auth.Key)
	})

	t.Run("Test User with CourseAssignments preload", func(t *testing.T) {
		userId := uuid.New()
		user := &entity.User{
			ID:         userId,
			PasswdHash: "test_hash",
			Role:       1,
		}
		userData := &entity.UserData{
			UserID:      userId,
			Email:       "test3@example.com",
			Name:        "Test User 3",
			Birthdate:   time.Now(),
			PhoneNumber: "+1234567892",
		}
		auth := &entity.Auth{
			UserID: userId,
			Key:    "test_key_3",
		}

		err := userRepo.CreateUser(user, userData, auth)
		require.NoError(t, err)

		courseId := uuid.New()
		course := &entity.Course{
			CourseID: courseId,
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err = courseRepo.AddCourse(course)
		require.NoError(t, err)

		assignment := &entity.CourseAssignment{
			CaID:     uuid.New(),
			UserID:   userId,
			CourseID: courseId,
			Progress: 50,
		}
		err = courseRepo.CreateCourseAssignment(assignment)
		require.NoError(t, err)

		users, err := userRepo.GetUsers()
		require.NoError(t, err)

		var foundUser *entity.User
		for _, u := range users {
			if u.ID == userId {
				foundUser = &u
				break
			}
		}
		require.NotNil(t, foundUser)
		assert.NotEmpty(t, foundUser.CourseAssignments)
		assert.Equal(t, assignment.CaID, foundUser.CourseAssignments[0].CaID)
		assert.Equal(t, assignment.Progress, foundUser.CourseAssignments[0].Progress)
	})

	t.Run("Test Course with Lessons preload", func(t *testing.T) {
		courseId := uuid.New()
		course := &entity.Course{
			CourseID: courseId,
			Title:    "Test Course with Lessons",
			Desc:     "Test Description",
		}
		err := courseRepo.AddCourse(course)
		require.NoError(t, err)

		lessons := []*entity.Lesson{
			{
				LessonID:   uuid.New(),
				CourseID:   courseId,
				Title:      "Lesson 1",
				Desc:       "Description 1",
				VideoURL:   "http://example.com/video1",
				SummaryURL: "http://example.com/summary1",
			},
			{
				LessonID:   uuid.New(),
				CourseID:   courseId,
				Title:      "Lesson 2",
				Desc:       "Description 2",
				VideoURL:   "http://example.com/video2",
				SummaryURL: "http://example.com/summary2",
			},
		}

		for _, lesson := range lessons {
			err = courseRepo.AddLesson(lesson)
			require.NoError(t, err)
		}

		gotLessons, err := courseRepo.GetLessonsByCourseId(courseId)
		require.NoError(t, err)
		assert.Equal(t, len(lessons), len(gotLessons))

		for i, lesson := range lessons {
			assert.Equal(t, lesson.LessonID, gotLessons[i].LessonID)
			assert.Equal(t, lesson.Title, gotLessons[i].Title)
			assert.Equal(t, lesson.Desc, gotLessons[i].Desc)
			assert.Equal(t, lesson.VideoURL, gotLessons[i].VideoURL)
			assert.Equal(t, lesson.SummaryURL, gotLessons[i].SummaryURL)
		}
	})

	t.Run("Test Course with CourseAssignments preload", func(t *testing.T) {
		courseId := uuid.New()
		course := &entity.Course{
			CourseID: courseId,
			Title:    "Test Course with Assignments",
			Desc:     "Test Description",
		}
		err := courseRepo.AddCourse(course)
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
			userData := &entity.UserData{
				UserID:      user.ID,
				Email:       "test" + user.ID.String() + "@example.com",
				Name:        "Test User " + user.ID.String(),
				Birthdate:   time.Now(),
				PhoneNumber: "+1234567890",
			}
			auth := &entity.Auth{
				UserID: user.ID,
				Key:    "test_key_" + user.ID.String(),
			}
			err = userRepo.CreateUser(user, userData, auth)
			require.NoError(t, err)
		}

		for _, user := range users {
			assignment := &entity.CourseAssignment{
				CaID:     uuid.New(),
				UserID:   user.ID,
				CourseID: courseId,
				Progress: 75,
			}
			err = courseRepo.CreateCourseAssignment(assignment)
			require.NoError(t, err)
		}

		assignments, err := courseRepo.GetCourseAssignmentsByCourseId(courseId)
		require.NoError(t, err)
		assert.Equal(t, len(users), len(assignments))

		for i, assignment := range assignments {
			assert.Equal(t, courseId, assignment.CourseID)
			assert.Equal(t, users[i].ID, assignment.UserID)
			assert.Equal(t, uint(75), assignment.Progress)
		}
	})

	t.Run("Test GetUsers preloads all relationships", func(t *testing.T) {
		db = setupTestDB(t)
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

		for _, user := range users {
			userData := &entity.UserData{
				UserID:      user.ID,
				Email:       "test" + user.ID.String() + "@example.com",
				Name:        "Test User " + user.ID.String(),
				Birthdate:   time.Now(),
				PhoneNumber: "+1234567890",
			}
			auth := &entity.Auth{
				UserID: user.ID,
				Key:    "test_key_" + user.ID.String(),
			}
			err := userRepo.CreateUser(user, userData, auth)
			require.NoError(t, err)
		}

		for _, course := range courses {
			err := courseRepo.AddCourse(course)
			require.NoError(t, err)
		}

		for _, user := range users {
			for _, course := range courses {
				assignment := &entity.CourseAssignment{
					CaID:     uuid.New(),
					UserID:   user.ID,
					CourseID: course.CourseID,
					Progress: 50,
				}
				err := courseRepo.CreateCourseAssignment(assignment)
				require.NoError(t, err)
			}
		}

		gotUsers, err := userRepo.GetUsers()
		require.NoError(t, err)
		assert.Equal(t, len(users), len(gotUsers))

		for _, gotUser := range gotUsers {
			assert.NotNil(t, gotUser.UserData)
			assert.NotEmpty(t, gotUser.UserData.Email)
			assert.NotEmpty(t, gotUser.UserData.Name)

			assert.NotEmpty(t, gotUser.CourseAssignments)
			assert.Equal(t, len(courses), len(gotUser.CourseAssignments))

			for _, assignment := range gotUser.CourseAssignments {
				assert.NotNil(t, assignment.Course)
				assert.NotEmpty(t, assignment.Course.Title)
				assert.NotEmpty(t, assignment.Course.Desc)
				assert.Equal(t, uint(50), assignment.Progress)
			}
		}
	})
}
