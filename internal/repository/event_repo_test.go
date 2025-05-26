package repository

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEventRepository(t *testing.T) {
	db := setupTestDB(t)
	repo := NewEventRepo(&config.Config{DB: config.DB{
		Host:     "localhost",
		Port:     "5433",
		User:     "postgres",
		Password: "postgres",
		Name:     "mzt_test",
	}})
	repo.DB = db

	t.Run("Create and Get Event", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := db.Create(course).Error
		require.NoError(t, err)

		event := &entity.Event{
			EventID:     uuid.New(),
			CourseID:    course.CourseID,
			Title:       "Test Event",
			Description: "Test Description",
			EventDate:   time.Now().Add(24 * time.Hour),
			SecretInfo:  "Secret Information",
		}

		err = repo.AddEvent(event)
		require.NoError(t, err)

		gotEvent, err := repo.GetEvent(event.EventID)
		require.NoError(t, err)
		assert.Equal(t, event.EventID, gotEvent.EventID)
		assert.Equal(t, event.CourseID, gotEvent.CourseID)
		assert.Equal(t, event.Title, gotEvent.Title)
		assert.Equal(t, event.Description, gotEvent.Description)
		assert.Empty(t, gotEvent.SecretInfo)

		eventWithSecret, err := repo.GetEventWithSecret(event.EventID)
		require.NoError(t, err)
		assert.Equal(t, event.EventID, eventWithSecret.EventID)
		assert.Equal(t, event.SecretInfo, eventWithSecret.SecretInfo)
	})

	t.Run("Update Event", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := db.Create(course).Error
		require.NoError(t, err)

		event := &entity.Event{
			EventID:     uuid.New(),
			CourseID:    course.CourseID,
			Title:       "Original Title",
			Description: "Original Description",
			EventDate:   time.Now().Add(24 * time.Hour),
			SecretInfo:  "Original Secret",
		}
		err = db.Create(event).Error
		require.NoError(t, err)

		updateDto := &dto.UpdateEventDto{
			Title:       "Updated Title",
			Description: "Updated Description",
			SecretInfo:  "Updated Secret",
		}

		err = repo.UpdateEvent(event.EventID, updateDto)
		require.NoError(t, err)

		updatedEvent, err := repo.GetEventWithSecret(event.EventID)
		require.NoError(t, err)
		assert.Equal(t, updateDto.Title, updatedEvent.Title)
		assert.Equal(t, updateDto.Description, updatedEvent.Description)
		assert.Equal(t, updateDto.SecretInfo, updatedEvent.SecretInfo)
	})

	t.Run("Delete Event", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := db.Create(course).Error
		require.NoError(t, err)

		event := &entity.Event{
			EventID:     uuid.New(),
			CourseID:    course.CourseID,
			Title:       "Test Event",
			Description: "Test Description",
			EventDate:   time.Now().Add(24 * time.Hour),
		}
		err = db.Create(event).Error
		require.NoError(t, err)

		err = repo.DeleteEvent(event.EventID)
		require.NoError(t, err)

		var count int64
		db.Model(&entity.Event{}).Where("event_id = ?", event.EventID).Count(&count)
		assert.Equal(t, int64(0), count)
	})

	t.Run("Get Events By Course ID", func(t *testing.T) {
		course := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Test Course",
			Desc:     "Test Description",
		}
		err := db.Create(course).Error
		require.NoError(t, err)

		events := []*entity.Event{
			{
				EventID:     uuid.New(),
				CourseID:    course.CourseID,
				Title:       "Event 1",
				Description: "Description 1",
				EventDate:   time.Now().Add(24 * time.Hour),
			},
			{
				EventID:     uuid.New(),
				CourseID:    course.CourseID,
				Title:       "Event 2",
				Description: "Description 2",
				EventDate:   time.Now().Add(48 * time.Hour),
			},
		}

		for _, event := range events {
			err := db.Create(event).Error
			require.NoError(t, err)
		}

		otherCourse := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Other Course",
			Desc:     "Other Description",
		}
		err = db.Create(otherCourse).Error
		require.NoError(t, err)

		otherEvent := &entity.Event{
			EventID:     uuid.New(),
			CourseID:    otherCourse.CourseID,
			Title:       "Other Event",
			Description: "Other Description",
			EventDate:   time.Now().Add(24 * time.Hour),
		}
		err = db.Create(otherEvent).Error
		require.NoError(t, err)

		courseEvents, err := repo.GetEventsByCourseId(course.CourseID)
		require.NoError(t, err)
		assert.Len(t, courseEvents, 2)

		eventTitles := []string{courseEvents[0].Title, courseEvents[1].Title}
		assert.Contains(t, eventTitles, "Event 1")
		assert.Contains(t, eventTitles, "Event 2")
	})

	t.Run("Get Events With Secrets By User ID", func(t *testing.T) {
		userID := uuid.New()
		user := &entity.User{
			ID:         userID,
			PasswdHash: "test_hash",
			Role:       0,
		}
		err := db.Create(user).Error
		require.NoError(t, err)

		course1 := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Course 1",
			Desc:     "Description 1",
		}
		err = db.Create(course1).Error
		require.NoError(t, err)

		course2 := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Course 2",
			Desc:     "Description 2",
		}
		err = db.Create(course2).Error
		require.NoError(t, err)

		ca1 := &entity.CourseAssignment{
			CaID:     uuid.New(),
			UserID:   userID,
			CourseID: course1.CourseID,
			Progress: 0,
		}
		err = db.Create(ca1).Error
		require.NoError(t, err)

		ca2 := &entity.CourseAssignment{
			CaID:     uuid.New(),
			UserID:   userID,
			CourseID: course2.CourseID,
			Progress: 0,
		}
		err = db.Create(ca2).Error
		require.NoError(t, err)

		event1 := &entity.Event{
			EventID:     uuid.New(),
			CourseID:    course1.CourseID,
			Title:       "Event 1",
			Description: "Description 1",
			EventDate:   time.Now().Add(24 * time.Hour),
			SecretInfo:  "Secret 1",
		}
		err = db.Create(event1).Error
		require.NoError(t, err)

		event2 := &entity.Event{
			EventID:     uuid.New(),
			CourseID:    course2.CourseID,
			Title:       "Event 2",
			Description: "Description 2",
			EventDate:   time.Now().Add(48 * time.Hour),
			SecretInfo:  "Secret 2",
		}
		err = db.Create(event2).Error
		require.NoError(t, err)

		otherCourse := &entity.Course{
			CourseID: uuid.New(),
			Title:    "Other Course",
			Desc:     "Other Description",
		}
		err = db.Create(otherCourse).Error
		require.NoError(t, err)

		otherEvent := &entity.Event{
			EventID:     uuid.New(),
			CourseID:    otherCourse.CourseID,
			Title:       "Other Event",
			Description: "Other Description",
			EventDate:   time.Now().Add(24 * time.Hour),
			SecretInfo:  "Other Secret",
		}
		err = db.Create(otherEvent).Error
		require.NoError(t, err)

		userEvents, err := repo.GetEventsWithSecretsByUserId(userID)
		require.NoError(t, err)
		assert.Len(t, userEvents, 2)

		secretMap := make(map[string]string)
		for _, event := range userEvents {
			secretMap[event.Title] = event.SecretInfo
		}

		assert.Equal(t, "Secret 1", secretMap["Event 1"])
		assert.Equal(t, "Secret 2", secretMap["Event 2"])

		titles := []string{userEvents[0].Title, userEvents[1].Title}
		assert.NotContains(t, titles, "Other Event")
	})
}
