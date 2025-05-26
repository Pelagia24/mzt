package service

import (
	"errors"
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"mzt/internal/repository"

	"github.com/google/uuid"
)

type EventServiceInterface interface {
	GetEvents() ([]dto.EventDto, error)
	GetEvent(eventId uuid.UUID) (*dto.EventDto, error)
	GetEventWithSecrets(eventId uuid.UUID, userId uuid.UUID) (*dto.EventDto, error)
	CreateEvent(event *dto.CreateEventDto) error
	UpdateEvent(eventId uuid.UUID, updated *dto.UpdateEventDto) error
	DeleteEvent(eventId uuid.UUID) error
	GetEventsByCourseId(courseId uuid.UUID) ([]dto.EventDto, error)
}

type EventService struct {
	config     *config.Config
	eventRepo  repository.EventRepository
	courseRepo repository.CourseRepository
}

func NewEventService(cfg *config.Config, eventRepo repository.EventRepository, courseRepo repository.CourseRepository) *EventService {
	return &EventService{
		config:     cfg,
		eventRepo:  eventRepo,
		courseRepo: courseRepo,
	}
}

func (s *EventService) GetEvents() ([]dto.EventDto, error) {
	return s.eventRepo.GetEvents()
}

func (s *EventService) GetEvent(eventId uuid.UUID) (*dto.EventDto, error) {
	return s.eventRepo.GetEvent(eventId)
}

func (s *EventService) GetEventWithSecrets(eventId uuid.UUID, userId uuid.UUID) (*dto.EventDto, error) {
	event, err := s.eventRepo.GetEvent(eventId)
	if err != nil {
		return nil, err
	}

	_, err = s.courseRepo.GetCourseAssignment(event.CourseID, userId)
	if err != nil {
		return nil, errors.New("user does not have access to event secrets")
	}

	return s.eventRepo.GetEventWithSecret(eventId)
}

func (s *EventService) CreateEvent(event *dto.CreateEventDto) error {
	eventEntity := &entity.Event{
		EventID:     uuid.New(),
		CourseID:    event.CourseID,
		Title:       event.Title,
		Description: event.Description,
		EventDate:   event.EventDate,
		SecretInfo:  event.SecretInfo,
	}
	return s.eventRepo.AddEvent(eventEntity)
}

func (s *EventService) UpdateEvent(eventId uuid.UUID, updated *dto.UpdateEventDto) error {
	return s.eventRepo.UpdateEvent(eventId, updated)
}

func (s *EventService) DeleteEvent(eventId uuid.UUID) error {
	return s.eventRepo.DeleteEvent(eventId)
}

func (s *EventService) GetEventsByCourseId(courseId uuid.UUID) ([]dto.EventDto, error) {
	return s.eventRepo.GetEventsByCourseId(courseId)
}

func (s *EventService) GetEventsWithSecretsByUserId(userId uuid.UUID) ([]dto.EventDto, error) {
	return s.eventRepo.GetEventsWithSecretsByUserId(userId)
}
