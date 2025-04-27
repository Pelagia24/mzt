package router

import (
	"net/http"

	"mzt/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Router) ListCourses(c *gin.Context) {
	courses, err := r.courseService.ListCourses()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"courses": courses})
}

func (r *Router) GetCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	course, err := r.courseService.GetCourse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"course": course})
}

func (r *Router) CreateCourse(c *gin.Context) {
	var payload dto.CreateCourseDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := r.courseService.CreateCourse(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "Course created successfully"})
}

func (r *Router) UpdateCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	var payload dto.UpdateCourseDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.courseService.UpdateCourse(id, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Course updated successfully"})
}

func (r *Router) DeleteCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	err = r.courseService.DeleteCourse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Course deleted successfully"})
}

func (r *Router) ListLessons(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	lessons, err := r.courseService.ListLessons(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lessons": lessons})
}

func (r *Router) GetLesson(c *gin.Context) {
	lessonId := c.Param("lesson_id")
	id, err := uuid.Parse(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}
	lesson, err := r.courseService.GetLesson(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"lesson": lesson})
}

func (r *Router) CreateLesson(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	var payload dto.CreateLessonDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.courseService.CreateLesson(id, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lesson created successfully"})
}

func (r *Router) UpdateLesson(c *gin.Context) {
	lessonId := c.Param("lesson_id")
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
	err = r.courseService.UpdateLesson(id, &payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lesson updated successfully"})
}

func (r *Router) DeleteLesson(c *gin.Context) {
	lessonId := c.Param("lesson_id")
	id, err := uuid.Parse(lessonId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid lesson ID"})
		return
	}
	err = r.courseService.DeleteLesson(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Lesson deleted successfully"})
}

func (r *Router) AssignUserToCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	var payload dto.AssignUserToCourseDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId, err := uuid.Parse(payload.UserId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	err = r.courseService.AssignUserToCourse(id, userId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User enrolled to course successfully"})
}

func (r *Router) ListUsersOnCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	users, err := r.courseService.ListUsersOnCourse(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"users": users})
}

func (r *Router) RemoveUserFromCourse(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	userId := c.Param("user_id")
	uid, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	err = r.courseService.RemoveUserFromCourse(id, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User removed from course successfully"})
}

func (r *Router) GetProgress(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	userId := c.GetString("user_id")
	uid, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	progress, err := r.courseService.GetProgress(id, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"progress": progress})
}

func (r *Router) UpdateProgress(c *gin.Context) {
	courseId := c.Param("course_id")
	id, err := uuid.Parse(courseId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid course ID"})
		return
	}
	userId := c.GetString("user_id")
	uid, err := uuid.Parse(userId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	var payload dto.UpdateProgressDto
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err = r.courseService.UpdateProgress(id, uid, payload.Progress)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Progress updated successfully"})
}
