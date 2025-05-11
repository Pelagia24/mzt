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
					Text:     "HTML5 introduced a set of semantic elements that provide meaning to the structure of web pages. These elements help both browsers and developers understand the purpose of different sections of content. Key semantic elements include:\n\n- <header>: Represents introductory content\n- <nav>: Defines navigation links\n- <main>: Specifies the main content\n- <article>: Represents self-contained content\n- <section>: Defines a section in a document\n- <aside>: Represents content that is tangentially related\n- <footer>: Represents a footer for a section\n\nUsing semantic HTML improves accessibility, SEO, and makes your code more maintainable.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "CSS Layouts and Flexbox",
					Summery:  "Master CSS layouts using Flexbox and Grid",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Flexbox is a one-dimensional layout method for arranging items in rows or columns. Key concepts include:\n\n- flex-direction: Determines the main axis\n- justify-content: Aligns items along the main axis\n- align-items: Aligns items along the cross axis\n- flex-wrap: Controls whether items wrap to new lines\n- flex-grow: Determines how much an item can grow\n- flex-shrink: Determines how much an item can shrink\n- flex-basis: Sets the initial main size of an item\n\nFlexbox is perfect for creating responsive layouts and centering content.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Responsive Design",
					Summery:  "Create responsive websites that work on all devices",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Responsive design ensures your website looks great on all devices. Key techniques include:\n\n- Media queries: Apply different styles based on screen size\n- Fluid layouts: Use relative units (%, em, rem)\n- Flexible images: Set max-width: 100%\n- Mobile-first approach: Design for mobile first, then enhance for larger screens\n- Viewport meta tag: Control the viewport's size and scale\n\nRemember to test your website on various devices and screen sizes.",
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
					Text:     "JavaScript is a versatile programming language. Key concepts include:\n\n- Variables and data types (let, const, var)\n- Operators and expressions\n- Control flow (if/else, switch, loops)\n- Functions and scope\n- Arrays and objects\n- Error handling (try/catch)\n\nUnderstanding these fundamentals is crucial for building interactive web applications.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "DOM Manipulation",
					Summery:  "Work with the Document Object Model",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "The DOM (Document Object Model) represents the HTML document as a tree of objects. Key operations include:\n\n- Selecting elements (querySelector, getElementById)\n- Modifying content (innerHTML, textContent)\n- Changing attributes (setAttribute, classList)\n- Creating and removing elements (createElement, appendChild)\n- Event handling (addEventListener)\n\nDOM manipulation is essential for creating dynamic web applications.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Async JavaScript",
					Summery:  "Master Promises, Async/Await, and Event Loop",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Asynchronous JavaScript is crucial for handling operations that take time. Key concepts include:\n\n- Callbacks and callback hell\n- Promises and Promise chaining\n- Async/await syntax\n- Event loop and microtasks\n- Error handling in async code\n\nUnderstanding async JavaScript is vital for building responsive applications.",
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
					Text:     "React components are the building blocks of React applications. Key concepts include:\n\n- Functional and class components\n- Props and prop types\n- Component composition\n- State and lifecycle methods\n- Conditional rendering\n\nComponents help create reusable, maintainable UI elements.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "State Management",
					Summery:  "Manage application state with React hooks",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "State management is crucial in React applications. Key concepts include:\n\n- useState hook for local state\n- useEffect for side effects\n- useContext for global state\n- useReducer for complex state logic\n- Custom hooks for reusable logic\n\nProper state management leads to predictable and maintainable applications.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "React Router",
					Summery:  "Implement routing in React applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "React Router enables navigation in single-page applications. Key concepts include:\n\n- Route configuration\n- Navigation with Link and useNavigate\n- Route parameters\n- Nested routes\n- Protected routes\n\nRouting is essential for creating multi-page experiences in single-page applications.",
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
					Text:     "Vue components are the foundation of Vue applications. Key concepts include:\n\n- Single-file components\n- Props and events\n- Component lifecycle\n- Computed properties\n- Watchers\n\nComponents help create modular and maintainable applications.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Vuex State Management",
					Summery:  "Manage state with Vuex",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Vuex is Vue's official state management solution. Key concepts include:\n\n- State, getters, mutations, and actions\n- Modules for large applications\n- Vuex store configuration\n- State persistence\n- DevTools integration\n\nVuex helps manage complex application state effectively.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Vue Router",
					Summery:  "Implement routing in Vue applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Vue Router provides routing capabilities to Vue applications. Key concepts include:\n\n- Route configuration\n- Navigation guards\n- Route meta fields\n- Nested routes\n- Dynamic route matching\n\nRouting is essential for creating multi-page experiences in Vue applications.",
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
					Text:     "TypeScript adds static typing to JavaScript. Key concepts include:\n\n- Basic types (string, number, boolean)\n- Interfaces and type aliases\n- Generics\n- Type assertions\n- Type inference\n\nTypeScript helps catch errors early and improves code maintainability.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "TypeScript with React",
					Summery:  "Use TypeScript in React applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "TypeScript enhances React development. Key concepts include:\n\n- Typing props and state\n- Generic components\n- Type definitions for hooks\n- Event handling types\n- Third-party library types\n\nTypeScript with React provides better type safety and developer experience.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "TypeScript with Vue",
					Summery:  "Use TypeScript in Vue applications",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "TypeScript integration with Vue 3. Key concepts include:\n\n- Component typing\n- Props and emits typing\n- Composition API with TypeScript\n- Type definitions for Vuex\n- Type definitions for Vue Router\n\nTypeScript with Vue provides better type safety and developer experience.",
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
					Text:     "Unit testing is essential for frontend development. Key concepts include:\n\n- Jest testing framework\n- Test structure (describe, it, expect)\n- Mocking and stubbing\n- Snapshot testing\n- Test coverage\n\nUnit tests help ensure code quality and prevent regressions.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Integration Testing",
					Summery:  "Test component interactions",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "Integration testing verifies component interactions. Key concepts include:\n\n- Testing component integration\n- Mocking API calls\n- Testing user interactions\n- Testing state changes\n- Testing routing\n\nIntegration tests ensure components work together correctly.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "E2E Testing",
					Summery:  "End-to-end testing with Cypress",
					VideoURL: "https://www.youtube.com/embed/MLpmiywRNzY",
					Text:     "End-to-end testing verifies the entire application. Key concepts include:\n\n- Cypress testing framework\n- Test structure and commands\n- Custom commands\n- Fixtures and stubs\n- Visual testing\n\nE2E tests ensure the application works as expected from a user's perspective.",
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
