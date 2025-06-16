package service

import (
	"errors"
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"
	"mzt/internal/repository"

	"github.com/google/uuid"
)

// интерфейс для работы с событиями
// определяет все методы которые нужны для работы с событиями
type EventServiceInterface interface {
	GetEvents() ([]dto.EventDto, error)
	GetEvent(eventId uuid.UUID) (*dto.EventDto, error)
	GetEventWithSecrets(eventId uuid.UUID, userId uuid.UUID) (*dto.EventDto, error)
	CreateEvent(event *dto.CreateEventDto) error
	UpdateEvent(eventId uuid.UUID, updated *dto.UpdateEventDto) error
	DeleteEvent(eventId uuid.UUID) error
	GetEventsByCourseId(courseId uuid.UUID) ([]dto.EventDto, error)
}

// сервис для работы с событиями
// реализует интерфейс EventServiceInterface
type EventService struct {
	config     *config.Config
	eventRepo  repository.EventRepository
	courseRepo repository.CourseRepository
}

// создаем новый сервис для работы с событиями
func NewEventService(cfg *config.Config, eventRepo repository.EventRepository, courseRepo repository.CourseRepository) *EventService {
	return &EventService{
		config:     cfg,
		eventRepo:  eventRepo,
		courseRepo: courseRepo,
	}
}

// получает список всех событий
// просто берет все события из базы
func (s *EventService) GetEvents() ([]dto.EventDto, error) {
	return s.eventRepo.GetEvents()
}

// получает информацию о событии
// просто берет событие из базы по его id
func (s *EventService) GetEvent(eventId uuid.UUID) (*dto.EventDto, error) {
	return s.eventRepo.GetEvent(eventId)
}

// получает информацию о событии вместе с секретами
// проверяет что пользователь имеет доступ к курсу
func (s *EventService) GetEventWithSecrets(eventId uuid.UUID, userId uuid.UUID) (*dto.EventDto, error) {
	// получаем событие
	event, err := s.eventRepo.GetEvent(eventId)
	if err != nil {
		return nil, err
	}

	// проверяем что пользователь записан на курс
	_, err = s.courseRepo.GetCourseAssignment(event.CourseID, userId)
	if err != nil {
		return nil, errors.New("user does not have access to event secrets")
	}

	// получаем событие с секретами
	return s.eventRepo.GetEventWithSecret(eventId)
}

// создает новое событие
// создает новое событие в базе с указанными данными
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

// обновляет информацию о событии
// просто обновляет данные события в базе
func (s *EventService) UpdateEvent(eventId uuid.UUID, updated *dto.UpdateEventDto) error {
	return s.eventRepo.UpdateEvent(eventId, updated)
}

// удаляет событие
// просто удаляет событие из базы
func (s *EventService) DeleteEvent(eventId uuid.UUID) error {
	return s.eventRepo.DeleteEvent(eventId)
}

// получает список событий для курса
// просто берет все события курса из базы
func (s *EventService) GetEventsByCourseId(courseId uuid.UUID) ([]dto.EventDto, error) {
	return s.eventRepo.GetEventsByCourseId(courseId)
}

// получает список событий с секретами для пользователя
// берет все события для курсов на которые записан пользователь
func (s *EventService) GetEventsWithSecretsByUserId(userId uuid.UUID) ([]dto.EventDto, error) {
	return s.eventRepo.GetEventsWithSecretsByUserId(userId)
}
