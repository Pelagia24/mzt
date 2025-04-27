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
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "HTML Structure and Semantics",
					Desc:       "Learn about HTML5 elements and semantic markup",
					VideoURL:   "http://example.com/html-semantics",
					SummaryURL: "http://example.com/html-semantics-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "CSS Layouts and Flexbox",
					Desc:       "Master CSS layouts using Flexbox and Grid",
					VideoURL:   "http://example.com/css-layouts",
					SummaryURL: "http://example.com/css-layouts-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "Responsive Design",
					Desc:       "Create responsive websites that work on all devices",
					VideoURL:   "http://example.com/responsive-design",
					SummaryURL: "http://example.com/responsive-design-summary",
				},
			}
		case "JavaScript Essentials":
			lessons = []*entity.Lesson{
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "JavaScript Basics",
					Desc:       "Learn JavaScript fundamentals and syntax",
					VideoURL:   "http://example.com/js-basics",
					SummaryURL: "http://example.com/js-basics-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "DOM Manipulation",
					Desc:       "Work with the Document Object Model",
					VideoURL:   "http://example.com/dom-manipulation",
					SummaryURL: "http://example.com/dom-manipulation-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "Async JavaScript",
					Desc:       "Master Promises, Async/Await, and Event Loop",
					VideoURL:   "http://example.com/async-js",
					SummaryURL: "http://example.com/async-js-summary",
				},
			}
		case "React.js Development":
			lessons = []*entity.Lesson{
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "React Components",
					Desc:       "Learn about React components and props",
					VideoURL:   "http://example.com/react-components",
					SummaryURL: "http://example.com/react-components-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "State Management",
					Desc:       "Manage application state with React hooks",
					VideoURL:   "http://example.com/react-state",
					SummaryURL: "http://example.com/react-state-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "React Router",
					Desc:       "Implement routing in React applications",
					VideoURL:   "http://example.com/react-router",
					SummaryURL: "http://example.com/react-router-summary",
				},
			}
		case "Vue.js Mastery":
			lessons = []*entity.Lesson{
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "Vue Components",
					Desc:       "Create and use Vue components",
					VideoURL:   "http://example.com/vue-components",
					SummaryURL: "http://example.com/vue-components-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "Vuex State Management",
					Desc:       "Manage state with Vuex",
					VideoURL:   "http://example.com/vuex",
					SummaryURL: "http://example.com/vuex-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "Vue Router",
					Desc:       "Implement routing in Vue applications",
					VideoURL:   "http://example.com/vue-router",
					SummaryURL: "http://example.com/vue-router-summary",
				},
			}
		case "TypeScript for Frontend":
			lessons = []*entity.Lesson{
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "TypeScript Basics",
					Desc:       "Learn TypeScript fundamentals",
					VideoURL:   "http://example.com/ts-basics",
					SummaryURL: "http://example.com/ts-basics-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "TypeScript with React",
					Desc:       "Use TypeScript in React applications",
					VideoURL:   "http://example.com/ts-react",
					SummaryURL: "http://example.com/ts-react-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "TypeScript with Vue",
					Desc:       "Use TypeScript in Vue applications",
					VideoURL:   "http://example.com/ts-vue",
					SummaryURL: "http://example.com/ts-vue-summary",
				},
			}
		case "Frontend Testing":
			lessons = []*entity.Lesson{
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "Unit Testing",
					Desc:       "Write unit tests for frontend code",
					VideoURL:   "http://example.com/unit-testing",
					SummaryURL: "http://example.com/unit-testing-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "Integration Testing",
					Desc:       "Test component interactions",
					VideoURL:   "http://example.com/integration-testing",
					SummaryURL: "http://example.com/integration-testing-summary",
				},
				{
					LessonID:   uuid.New(),
					CourseID:   course.CourseID,
					Title:      "E2E Testing",
					Desc:       "End-to-end testing with Cypress",
					VideoURL:   "http://example.com/e2e-testing",
					SummaryURL: "http://example.com/e2e-testing-summary",
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
