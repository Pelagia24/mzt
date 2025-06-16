package router

import (
	"net/http"

	"mzt/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListCourses получает список всех курсов
// просто возвращает все курсы из базы
func (r *Router) ListCourses(c *gin.Context) {
	// получаем список курсов из сервиса
	courses, err := r.courseService.ListCourses()
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем список курсов клиенту
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

// GetCourse получает информацию о курсе
// берет курс из базы по его id
func (r *Router) GetCourse(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(courseId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// получаем информацию о курсе из сервиса
	course, err := r.courseService.GetCourse(id)
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем информацию о курсе клиенту
	c.JSON(http.StatusOK, gin.H{"course": course})
}

// CreateCourse создает новый курс
// доступно только админам
func (r *Router) CreateCourse(c *gin.Context) {
	// парсим данные из тела запроса
	var payload dto.CreateCourseDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		// если данные невалидные, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// создаем курс через сервис
	err := r.courseService.CreateCourse(&payload)
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем успешный response
	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully"})
}

// UpdateCourse обновляет информацию о курсе
// доступно только админам
func (r *Router) UpdateCourse(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(courseId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// парсим данные из тела запроса
	var payload dto.UpdateCourseDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		// если данные невалидные, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// обновляем курс через сервис
	err = r.courseService.UpdateCourse(id, &payload)
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем успешный response
	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

// DeleteCourse удаляет курс
// доступно только админам
func (r *Router) DeleteCourse(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(courseId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// удаляем курс через сервис
	err = r.courseService.DeleteCourse(id)
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем успешный response
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

// ListLessons получает список всех уроков курса
// берет все уроки из базы по id курса
func (r *Router) ListLessons(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(courseId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// получаем список уроков из сервиса
	lessons, err := r.courseService.ListLessons(id)
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем список уроков клиенту
	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
}

// GetLesson получает информацию об уроке
// берет урок из базы по его id
func (r *Router) GetLesson(c *gin.Context) {
	// достаем id урока из параметров запроса
	lessonId := c.Param("lesson_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(lessonId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}
	lesson, err := r.courseService.GetLesson(id)
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем информацию об уроке клиенту
	c.JSON(http.StatusOK, gin.H{"lesson": lesson})
}

// CreateLesson создает новый урок
// доступно только админам
func (r *Router) CreateLesson(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// парсим данные из тела запроса
	var payload dto.CreateLessonDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// создаем урок через сервис
	err = r.courseService.CreateLesson(id, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lesson created successfully"})
}

// UpdateLesson обновляет информацию об уроке
// доступно только админам
func (r *Router) UpdateLesson(c *gin.Context) {
	// достаем id урока из параметров запроса
	lessonId := c.Param("lesson_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}
	var payload dto.UpdateLessonDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// обновляем урок через сервис
	err = r.courseService.UpdateLesson(id, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lesson updated successfully"})
}

// DeleteLesson удаляет урок
// доступно только админам
func (r *Router) DeleteLesson(c *gin.Context) {
	// достаем id урока из параметров запроса
	lessonId := c.Param("lesson_id")
	id, err := uuid.Parse(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}
	// удаляем урок через сервис
	err = r.courseService.DeleteLesson(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lesson deleted successfully"})
}

// MyCourses получает список курсов пользователя
// берет все курсы на которые записан пользователь
func (r *Router) MyCourses(c *gin.Context) {
	// достаем id пользователя из контекста
	self, ok := c.Get("self")
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	// преобразуем id в uuid
	selfId, ok := self.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "Unknown sender"})
		return
	}

	// получаем список курсов пользователя из сервиса
	courses, err := r.courseService.ListUserCourses(selfId)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Can't get user courses"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Lesson deleted successfully",
		"courses": courses})
}

// ListUsersOnCourse получает список пользователей на курсе
// берет всех пользователей записанных на курс
func (r *Router) ListUsersOnCourse(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// получаем список пользователей из сервиса
	users, err := r.courseService.ListUsersOnCourse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

// RemoveUserFromCourse удаляет пользователя с курса
// удаляет запись о том что пользователь записан на курс
func (r *Router) RemoveUserFromCourse(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// достаем id пользователя из параметров запроса
	userId := c.Param("user_id")
	uid, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// удаляем пользователя с курса через сервис
	err = r.courseService.RemoveUserFromCourse(id, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User removed from course successfully"})
}

// GetProgress получает прогресс пользователя по курсу
// берет значение прогресса из базы
func (r *Router) GetProgress(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// достаем id пользователя из контекста
	userId := c.GetString("user_id")
	// преобразуем строку в uuid
	uid, err := uuid.Parse(userId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// получаем прогресс из сервиса
	progress, err := r.courseService.GetProgress(id, uid)
	if err != nil {
		// если что-то пошло не так, возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// отправляем прогресс клиенту
	c.JSON(http.StatusOK, gin.H{"progress": progress})
}

// UpdateProgress обновляет прогресс пользователя по курсу
// меняет значение прогресса в базе
func (r *Router) UpdateProgress(c *gin.Context) {
	// достаем id курса из параметров запроса
	courseId := c.Param("course_id")
	// преобразуем строку в uuid
	id, err := uuid.Parse(courseId)
	if err != nil {
		// если id невалидный, возвращаем ошибку
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	// достаем id пользователя из контекста
	userId := c.GetString("user_id")
	// преобразуем строку в uuid
	uid, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	// парсим данные из тела запроса
	var payload dto.UpdateProgressDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// обновляем прогресс через сервис
	err = r.courseService.UpdateProgress(id, uid, payload.Progress)
	if err != nil {
		// если что-то пошло не так возвращаем ошибку
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
}
