package repository

import (
	"mzt/config"
	"mzt/internal/dto"
	"mzt/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// интерфейс для работы с событиями
// определяет все методы которые нужны для работы с событиями в базе
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

// репозиторий для работы с событиями
// реализует интерфейс EventRepository
type EventRepo struct {
	config *config.Config
	DB     *gorm.DB
}

// создаем новый репозиторий для работы с событиями
func NewEventRepo(cfg *config.Config) *EventRepo {
	return &EventRepo{
		config: cfg,
		DB:     connectDB(cfg),
	}
}

// получает список всех событий
// берет все события из базы и преобразует их в формат для response
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

// получает информацию о событии
// берет событие из базы по его id и преобразует его в формат для response
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

// получает информацию о событии с секретами
// берет событие из базы по его id и преобразует его в формат для response
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

// создает новое событие
// просто создает новую запись в таблице events
func (r *EventRepo) AddEvent(event *entity.Event) error {
	return r.DB.Create(event).Error
}

// обновляет информацию о событии
// меняет название описание дату и секреты события в базе
func (r *EventRepo) UpdateEvent(eventId uuid.UUID, updated *dto.UpdateEventDto) error {
	// начинаем транзакцию
	tx := r.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// получаем событие
	var event entity.Event
	if err := tx.Where("event_id = ?", eventId).First(&event).Error; err != nil {
		tx.Rollback()
		return err
	}

	// обновляем поля если они не пустые
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

	// сохраняем изменения
	if err := tx.Save(&event).Error; err != nil {
		tx.Rollback()
		return err
	}

	// завершаем транзакцию
	return tx.Commit().Error
}

// удаляет событие
// просто удаляет событие из базы
func (r *EventRepo) DeleteEvent(eventId uuid.UUID) error {
	return r.DB.Delete(&entity.Event{}, "event_id = ?", eventId).Error
}

// получает список событий для курса
// берет все события курса из базы и сортирует их по дате
func (r *EventRepo) GetEventsByCourseId(courseId uuid.UUID) ([]dto.EventDto, error) {
	var events []entity.Event
	if err := r.DB.Where("course_id = ?", courseId).Order("event_date asc").Find(&events).Error; err != nil {
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

// получает список событий с секретами для пользователя
// берет все события для курсов на которые записан пользователь
func (r *EventRepo) GetEventsWithSecretsByUserId(userId uuid.UUID) ([]dto.EventDto, error) {
	// получаем все курсы пользователя
	var courseAssignments []entity.CourseAssignment
	if err := r.DB.Where("user_id = ?", userId).Find(&courseAssignments).Error; err != nil {
		return nil, err
	}

	// если у пользователя нет курсов - возвращаем пустой список
	if len(courseAssignments) == 0 {
		return []dto.EventDto{}, nil
	}

	// собираем id всех курсов пользователя
	courseIDs := make([]uuid.UUID, len(courseAssignments))
	for i, ca := range courseAssignments {
		courseIDs[i] = ca.CourseID
	}

	// получаем все события для этих курсов
	var events []entity.Event
	if err := r.DB.Where("course_id IN ?", courseIDs).Find(&events).Error; err != nil {
		return nil, err
	}

	// преобразуем события в формат для response
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
