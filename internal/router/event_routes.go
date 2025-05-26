package router

import (
	"net/http"

	"mzt/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Router) ListEvents(c *gin.Context) {
	events, err := r.eventService.GetEvents()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (r *Router) GetEvent(c *gin.Context) {
	eventId := c.Param("event_id")
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	event, err := r.eventService.GetEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

func (r *Router) GetEventWithSecrets(c *gin.Context) {
	eventId := c.Param("event_id")
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	
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
	
	event, err := r.eventService.GetEventWithSecrets(id, selfId)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"event": event})
}

func (r *Router) GetMyEventsWithSecrets(c *gin.Context) {
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
	
	events, err := r.eventService.GetEventsWithSecretsByUserId(selfId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
}

func (r *Router) CreateEvent(c *gin.Context) {
	var payload dto.CreateEventDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.eventService.CreateEvent(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully"})
}

func (r *Router) UpdateEvent(c *gin.Context) {
	eventId := c.Param("event_id")
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	var payload dto.UpdateEventDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.eventService.UpdateEvent(id, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

func (r *Router) DeleteEvent(c *gin.Context) {
	eventId := c.Param("event_id")
	id, err := uuid.Parse(eventId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}
	err = r.eventService.DeleteEvent(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
}

func (r *Router) GetEventsByCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	events, err := r.eventService.GetEventsByCourseId(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"events": events})
} 