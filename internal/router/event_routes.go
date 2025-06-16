package router

import (
	"net/http"

	"mzt/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// получает список всех событий
// доступно только админам
func (r *Router) ListEvents(c *gin.Context) {
	// получаем список событий из сервиса
	events, err := r.eventService.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

// получает информацию о событии
// берет событие из базы по его id
func (r *Router) GetEvent(c *gin.Context) {
	// достаем id события из параметров запроса
	eventId := c.Param("event_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	// получаем информацию о событии из сервиса
	event, err := r.eventService.GetEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем информацию о событии клиенту
	c.JSON(http.StatusOK, gin.H{"event": event})
}

// получает информацию о событии со всеми секретами
// доступно только пользователям записанным на курс
func (r *Router) GetEventWithSecrets(c *gin.Context) {
	// достаем id события из параметров запроса
	eventId := c.Param("event_id")
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	// достаем id пользователя из контекста
	self, ok := c.Get("self")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	selfId, ok := self.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown sender"})
		return
	}

	// получаем информацию о событии со всеми секретами из сервиса
	event, err := r.eventService.GetEventWithSecrets(id, selfId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

// получает список всех событий пользователя со всеми секретами
// доступно только авторизованным пользователям
func (r *Router) GetMyEventsWithSecrets(c *gin.Context) {
	// достаем id пользователя из контекста
	self, ok := c.Get("self")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	selfId, ok := self.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown sender"})
		return
	}

	// получаем список событий со всеми секретами из сервиса
	events, err := r.eventService.GetEventsWithSecretsByUserId(selfId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

// создает новое событие
// доступно только админам
func (r *Router) CreateEvent(c *gin.Context) {
	// парсим данные из тела запроса
	var payload dto.CreateEventDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		// если данные невалидные, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем событие через сервис
	err := r.eventService.CreateEvent(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем успешный response
	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully"})
}

// обновляет информацию о событии
// доступно только админам
func (r *Router) UpdateEvent(c *gin.Context) {
	// достаем id события из параметров запроса
	eventId := c.Param("event_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	// парсим данные из тела запроса
	var payload dto.UpdateEventDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		// если данные невалидные, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// обновляем событие через сервис
	err = r.eventService.UpdateEvent(id, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем успешный response
	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// удаляет событие
// доступно только админам
func (r *Router) DeleteEvent(c *gin.Context) {
	// достаем id события из параметров запроса
	eventId := c.Param("event_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	// удаляем событие через сервис
	err = r.eventService.DeleteEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем успешный response
	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

// получает список всех событий курса
// доступно только авторизованным пользователям
func (r *Router) GetEventsByCourse(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// получаем список событий из сервиса
	events, err := r.eventService.GetEventsByCourseId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}
