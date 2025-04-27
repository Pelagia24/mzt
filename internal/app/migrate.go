package app

import (
	"errors"
	"fmt"
	"mzt/internal/entity"
	"mzt/internal/repository"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func Migrate(r *repository.UserRepo) {
	r.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	err := r.DB.AutoMigrate(
		&entity.User{},
		&entity.UserData{},
		&entity.Auth{},
		&entity.Course{},
		&entity.CourseAssignment{},
		&entity.Lesson{},
	)
	if err != nil {
		panic(err)
	}

	adminPasswdHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	users := []*entity.User{
		{
			ID:         uuid.New(),
			PasswdHash: string(adminPasswdHash),
			Role:       1,
		},
		{
			ID:         uuid.New(),
			PasswdHash: "test_hash2",
			Role:       0,
		},
		{
			ID:         uuid.New(),
			PasswdHash: "test_hash3",
			Role:       0,
		},
	}

	createdUsers := make(map[string]*entity.User)
	for i, user := range users {
		email := "test" + fmt.Sprint(i) + "@example.com"

		var existingUserData entity.UserData
		err = r.DB.Where("email = ?", email).First(&existingUserData).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				userData := &entity.UserData{
					UserID:      user.ID,
					Email:       email,
					Name:        "Test User " + fmt.Sprint(i),
					Birthdate:   time.Now(),
					PhoneNumber: "+1234567890",
					Telegram:    "@test" + fmt.Sprint(i),
					City:        "Test City",
					Age:         25,
					Employment:  "Test Company",
				}
				auth := &entity.Auth{
					UserID: user.ID,
					Key:    "test_key_" + fmt.Sprint(i),
				}

				err = r.CreateUser(user, userData, auth)
				if err != nil {
					panic(err)
				}
				createdUsers[email] = user
			} else {
				panic(err)
			}
		} else {
			var existingUser entity.User
			err = r.DB.Where("id = ?", existingUserData.UserID).First(&existingUser).Error
			if err != nil {
				panic(err)
			}
			createdUsers[email] = &existingUser
		}
	}

	courses := []*entity.Course{
		{
			CourseID: uuid.New(),
			Title:    "HTML & CSS Fundamentals",
			Desc:     "Learn the building blocks of web development with HTML5 and CSS3",
		},
		{
			CourseID: uuid.New(),
			Title:    "JavaScript Essentials",
			Desc:     "Master JavaScript programming from basics to advanced concepts",
		},
		{
			CourseID: uuid.New(),
			Title:    "React.js Development",
			Desc:     "Build modern web applications with React.js",
		},
		{
			CourseID: uuid.New(),
			Title:    "Vue.js Mastery",
			Desc:     "Create reactive user interfaces with Vue.js",
		},
		{
			CourseID: uuid.New(),
			Title:    "TypeScript for Frontend",
			Desc:     "Add type safety to your JavaScript code with TypeScript",
		},
		{
			CourseID: uuid.New(),
			Title:    "Frontend Testing",
			Desc:     "Learn testing strategies for frontend applications",
		},
	}

	createdCourses := make(map[string]*entity.Course)
	for _, course := range courses {
		var existingCourse entity.Course
		err = r.DB.Where("title = ?", course.Title).First(&existingCourse).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = r.DB.Create(course).Error
				if err != nil {
					panic(err)
				}
				createdCourses[course.Title] = course
			} else {
				panic(err)
			}
		} else {
			createdCourses[course.Title] = &existingCourse
		}
	}

	for _, course := range createdCourses {
		var lessons []*entity.Lesson

		switch course.Title {
		case "HTML & CSS Fundamentals":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "HTML Structure and Semantics",
					Summery:  "Learn about HTML5 elements and semantic markup",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "CSS Layouts and Flexbox",
					Summery:  "Master CSS layouts using Flexbox and Grid",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Responsive Design",
					Summery:  "Create responsive websites that work on all devices",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
			}
		case "JavaScript Essentials":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "JavaScript Basics",
					Summery:  "Learn JavaScript fundamentals and syntax",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "DOM Manipulation",
					Summery:  "Work with the Document Object Model",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Async JavaScript",
					Summery:  "Master Promises, Async/Await, and Event Loop",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
			}
		case "React.js Development":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "React Components",
					Summery:  "Learn about React components and props",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "State Management",
					Summery:  "Manage application state with React hooks",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "React Router",
					Summery:  "Implement routing in React applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
			}
		case "Vue.js Mastery":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Vue Components",
					Summery:  "Create and use Vue components",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Vuex State Management",
					Summery:  "Manage state with Vuex",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Vue Router",
					Summery:  "Implement routing in Vue applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
			}
		case "TypeScript for Frontend":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "TypeScript Basics",
					Summery:  "Learn TypeScript fundamentals",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "TypeScript with React",
					Summery:  "Use TypeScript in React applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "TypeScript with Vue",
					Summery:  "Use TypeScript in Vue applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
			}
		case "Frontend Testing":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Unit Testing",
					Summery:  "Write unit tests for frontend code",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Integration Testing",
					Summery:  "Test component interactions",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "E2E Testing",
					Summery:  "End-to-end testing with Cypress",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Lorem Ipsum is simply dummy text of the printing and typesetting industry. Lorem Ipsum has been the industry's standard dummy text ever since the 1500s, when an unknown printer took a galley of type and scrambled it to make a type specimen book. It has survived not only five centuries, but also the leap into electronic typesetting, remaining essentially unchanged. It was popularised in the 1960s with the release of Letraset sheets containing Lorem Ipsum passages, and more recently with desktop publishing software like Aldus PageMaker including versions of Lorem Ipsum.",
				},
			}
		}

		for _, lesson := range lessons {
			var existingLesson entity.Lesson
			err = r.DB.Where("title = ? AND course_id = ?", lesson.Title, course.CourseID).First(&existingLesson).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = r.DB.Create(lesson).Error
					if err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
		}
	}

	for _, user := range createdUsers {
		for _, course := range createdCourses {
			progress := uint(0)
			switch course.Title {
			case "HTML & CSS Fundamentals":
				progress = 75
			case "JavaScript Essentials":
				progress = 50
			case "React.js Development":
				progress = 25
			case "Vue.js Mastery":
				progress = 0
			case "TypeScript for Frontend":
				progress = 100
			case "Frontend Testing":
				progress = 30
			}

			var existingAssignment entity.CourseAssignment
			err = r.DB.Where("user_id = ? AND course_id = ?", user.ID, course.CourseID).First(&existingAssignment).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					assignment := &entity.CourseAssignment{
						CaID:     uuid.New(),
						UserID:   user.ID,
						CourseID: course.CourseID,
						Progress: progress,
					}
					err = r.DB.Create(assignment).Error
					if err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
		}
	}

	fmt.Println("Template data migration completed successfully")
}
