package repository

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventRepository interface {
	GetEvents() ([]dto.EventDto, error)
	GetEvent(eventId uuid.UUID) (*dto.EventDto, error)
	GetEventWithSecret(eventId uuid.UUID) (*dto.EventDto, error)
	AddEvent(event *entity.Event) error
	UpdateEvent(eventId uuid.UUID, updated *dto.UpdateEventDto) error
	DeleteEvent(eventId uuid.UUID) error
	GetEventsByCourseId(courseId uuid.UUID) ([]dto.EventDto, error)
	GetEventsWithSecretsByUserId(userId uuid.UUID) ([]dto.EventDto, error)
}

type EventRepo struct {
	config *config.Config
	DB     *gorm.DB
}

func NewEventRepo(cfg *config.Config) *EventRepo {
	return &EventRepo{
		config: cfg,
		DB:     connectDB(cfg),
	}
}

func (r *EventRepo) GetEvents() ([]dto.EventDto, error) {
	var events []entity.Event
	if err := r.DB.Find(&events).Error; err != nil {
		return nil, err
	}

	result := make([]dto.EventDto, 0, len(events))
	for _, event := range events {
		result = append(result, dto.EventDto{
			EventID:     event.EventID,
			CourseID:    event.CourseID,
			Title:       event.Title,
			Description: event.Description,
			EventDate:   event.EventDate,
			SecretInfo:  event.SecretInfo,
		})
	}
	return result, nil
}

func (r *EventRepo) GetEvent(eventId uuid.UUID) (*dto.EventDto, error) {
	var event entity.Event
	if err := r.DB.First(&event, "event_id = ?", eventId).Error; err != nil {
		return nil, err
	}
	return &dto.EventDto{
		EventID:     event.EventID,
		CourseID:    event.CourseID,
		Title:       event.Title,
		Description: event.Description,
		EventDate:   event.EventDate,
		SecretInfo:  event.SecretInfo,
	}, nil
}

func (r *EventRepo) GetEventWithSecret(eventId uuid.UUID) (*dto.EventDto, error) {
	var event entity.Event
	if err := r.DB.First(&event, "event_id = ?", eventId).Error; err != nil {
		return nil, err
	}
	return &dto.EventDto{
		EventID:     event.EventID,
		CourseID:    event.CourseID,
		Title:       event.Title,
		Description: event.Description,
		EventDate:   event.EventDate,
		SecretInfo:  event.SecretInfo,
	}, nil
}

func (r *EventRepo) AddEvent(event *entity.Event) error {
	return r.DB.Create(event).Error
}

func (r *EventRepo) UpdateEvent(eventId uuid.UUID, updated *dto.UpdateEventDto) error {
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	var event entity.Event
	if err := tx.Where("event_id = ?", eventId).First(&event).Error; err != nil {
		tx.Rollback()
		return err
	}

	if updated.Title != "" {
		event.Title = updated.Title
	}
	if updated.Description != "" {
		event.Description = updated.Description
	}
	if !updated.EventDate.IsZero() {
		event.EventDate = updated.EventDate
	}
	if updated.SecretInfo != "" {
		event.SecretInfo = updated.SecretInfo
	}

	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *EventRepo) DeleteEvent(eventId uuid.UUID) error {
	return r.DB.Delete(&entity.Event{}, "event_id = ?", eventId).Error
}

func (r *EventRepo) GetEventsByCourseId(courseId uuid.UUID) ([]dto.EventDto, error) {
	var events []entity.Event
	if err := r.DB.Where("course_id = ?", courseId).Find(&events).Error; err != nil {
		return nil, err
	}

	result := make([]dto.EventDto, 0, len(events))
	for _, event := range events {
		result = append(result, dto.EventDto{
			EventID:     event.EventID,
			CourseID:    event.CourseID,
			Title:       event.Title,
			Description: event.Description,
			EventDate:   event.EventDate,
			SecretInfo:  event.SecretInfo,
		})
	}
	return result, nil
}

func (r *EventRepo) GetEventsWithSecretsByUserId(userId uuid.UUID) ([]dto.EventDto, error) {
	var courseAssignments []entity.CourseAssignment
	if err := r.DB.Where("user_id = ?", userId).Find(&courseAssignments).Error; err != nil {
		return nil, err
	}

	if len(courseAssignments) == 0 {
		return []dto.EventDto{}, nil
	}

	courseIDs := make([]uuid.UUID, len(courseAssignments))
	for i, ca := range courseAssignments {
		courseIDs[i] = ca.CourseID
	}

	var events []entity.Event
	if err := r.DB.Where("course_id IN ?", courseIDs).Find(&events).Error; err != nil {
		return nil, err
	}

	result := make([]dto.EventDto, len(events))
	for i, event := range events {
		result[i] = dto.EventDto{
			EventID:     event.EventID,
			CourseID:    event.CourseID,
			Title:       event.Title,
			Description: event.Description,
			EventDate:   event.EventDate,
			SecretInfo:  event.SecretInfo,
		}
	}

	return result, nil
}
