package main

import (
	"fmt"
	"mzt/config"
	"mzt/internal/entity"
	"mzt/internal/repository"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	// Get database credentials from PostgreSQL standard environment variables
	dbHost := getEnvOrDefault("PGHOST", "localhost")
	dbPort := getEnvOrDefault("PGPORT", "5432")
	dbUser := getEnvOrDefault("PGUSER", "postgres")
	dbPass := getEnvOrDefault("PGPASSWORD", "postgres")
	dbName := getEnvOrDefault("PGDATABASE", "mzt")

	fmt.Printf("Connecting to database %s on %s:%s as %s\n", dbName, dbHost, dbPort, dbUser)

	cfg := &config.Config{
		DB: config.DB{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPass,
			Name:     dbName,
		},
	}

	userRepo := repository.NewUserRepo(cfg)
	courseRepo := repository.NewCourseRepo(cfg)
	eventRepo := repository.NewEventRepo(cfg)

	// Enable UUID extension
	userRepo.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")

	// Auto migrate all tables
	err := userRepo.DB.AutoMigrate(
		&entity.User{},
		&entity.UserData{},
		&entity.Auth{},
		&entity.Payment{},
		&entity.Course{},
		&entity.CourseAssignment{},
		&entity.Lesson{},
		&entity.Event{},
		&entity.CoursePrice{},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate database: %v", err))
	}

	// Create test users
	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	if err != nil {
		panic(fmt.Sprintf("Failed to hash password: %v", err))
	}

	testUsers := []struct {
		user     *entity.User
		userData *entity.UserData
		auth     *entity.Auth
	}{
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: string(passwordHash),
				Role:       1, // Admin
			},
			userData: &entity.UserData{
				Email:           "admin@example.com",
				Name:            "Admin User",
				Birthdate:       time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
				PhoneNumber:     "+79001234567",
				Telegram:        "@admin_user",
				City:            "Moscow",
				Age:             33,
				Employment:      "Tech Lead",
				IsBusinessOwner: "No",
				PositionAtWork:  "Senior Developer",
				MonthIncome:     150000,
			},
			auth: &entity.Auth{
				Key: "admin_key",
			},
		},
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: string(passwordHash),
				Role:       0, // Student
			},
			userData: &entity.UserData{
				Email:           "ivan@example.com",
				Name:            "Иван Петров",
				Birthdate:       time.Date(1995, 5, 15, 0, 0, 0, 0, time.UTC),
				PhoneNumber:     "+79157894561",
				Telegram:        "@ivan_petrov",
				City:            "Saint Petersburg",
				Age:             28,
				Employment:      "Junior Developer",
				IsBusinessOwner: "No",
				PositionAtWork:  "Frontend Developer",
				MonthIncome:     80000,
			},
			auth: &entity.Auth{
				Key: "ivan_key",
			},
		},
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: string(passwordHash),
				Role:       0,
			},
			userData: &entity.UserData{
				Email:           "anna@example.com",
				Name:            "Anna Smith",
				Birthdate:       time.Date(1992, 8, 23, 0, 0, 0, 0, time.UTC),
				PhoneNumber:     "+79269874563",
				Telegram:        "@anna_smith",
				City:            "Kazan",
				Age:             31,
				Employment:      "Business Owner",
				IsBusinessOwner: "Yes",
				PositionAtWork:  "CEO",
				MonthIncome:     250000,
			},
			auth: &entity.Auth{
				Key: "anna_key",
			},
		},
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: string(passwordHash),
				Role:       0,
			},
			userData: &entity.UserData{
				Email:           "maria@example.com",
				Name:            "Мария Иванова",
				Birthdate:       time.Date(1998, 3, 10, 0, 0, 0, 0, time.UTC),
				PhoneNumber:     "+79631234567",
				Telegram:        "@maria_iv",
				City:            "Novosibirsk",
				Age:             25,
				Employment:      "Student",
				IsBusinessOwner: "No",
				PositionAtWork:  "Intern",
				MonthIncome:     45000,
			},
			auth: &entity.Auth{
				Key: "maria_key",
			},
		},
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: string(passwordHash),
				Role:       0,
			},
			userData: &entity.UserData{
				Email:           "alex@example.com",
				Name:            "Алексей Смирнов",
				Birthdate:       time.Date(1988, 11, 30, 0, 0, 0, 0, time.UTC),
				PhoneNumber:     "+79567891234",
				Telegram:        "@alex_sm",
				City:            "Yekaterinburg",
				Age:             35,
				Employment:      "Senior Manager",
				IsBusinessOwner: "No",
				PositionAtWork:  "Project Manager",
				MonthIncome:     180000,
			},
			auth: &entity.Auth{
				Key: "alex_key",
			},
		},
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: string(passwordHash),
				Role:       0,
			},
			userData: &entity.UserData{
				Email:           "elena@example.com",
				Name:            "Елена Кузнецова",
				Birthdate:       time.Date(1993, 7, 20, 0, 0, 0, 0, time.UTC),
				PhoneNumber:     "+79876543210",
				Telegram:        "@elena_k",
				City:            "Nizhny Novgorod",
				Age:             30,
				Employment:      "Business Owner",
				IsBusinessOwner: "Yes",
				PositionAtWork:  "Founder",
				MonthIncome:     300000,
			},
			auth: &entity.Auth{
				Key: "elena_key",
			},
		},
	}

	createdUsers := make(map[string]*entity.User)
	for _, tu := range testUsers {
		var existingUserData entity.UserData
		err = userRepo.DB.Where("email = ?", tu.userData.Email).First(&existingUserData).Error
		if err != nil {
			if err.Error() == "record not found" {
				tu.userData.UserID = tu.user.ID
				tu.auth.UserID = tu.user.ID

				err = userRepo.CreateUser(tu.user, tu.userData, tu.auth)
				if err != nil {
					panic(fmt.Sprintf("Failed to create user %s: %v", tu.userData.Email, err))
				}
				createdUsers[tu.userData.Email] = tu.user
				fmt.Printf("Created user: %s\n", tu.userData.Email)
			} else {
				panic(fmt.Sprintf("Failed to check existing user %s: %v", tu.userData.Email, err))
			}
		} else {
			var existingUser entity.User
			err = userRepo.DB.Where("id = ?", existingUserData.UserID).First(&existingUser).Error
			if err != nil {
				panic(fmt.Sprintf("Failed to get existing user %s: %v", tu.userData.Email, err))
			}
			createdUsers[tu.userData.Email] = &existingUser
			fmt.Printf("User already exists: %s\n", tu.userData.Email)
		}
	}

	// Create test courses
	courses := []*entity.Course{
		{
			CourseID: uuid.New(),
			Title:    "Эмоциональный интеллект и коммуникация",
			Desc:     "Развитие эмоционального интеллекта, управление негативными эмоциями и развитие коммуникативных навыков лидера",
		},
		{
			CourseID: uuid.New(),
			Title:    "Стратегическое лидерство",
			Desc:     "Формирование траектории личного развития руководителя и основы установления долгосрочных взаимоотношений с ЛПР",
		},
		{
			CourseID: uuid.New(),
			Title:    "Стресс и Energy Management",
			Desc:     "Управление стрессом и психо-эмоциональным состоянием, техники повышения личной эффективности",
		},
		{
			CourseID: uuid.New(),
			Title:    "Самодисциплина и мотивация",
			Desc:     "Система личной мотивации, самодисциплина и эффективное планирование времени",
		},
		{
			CourseID: uuid.New(),
			Title:    "Высокоэффективные команды",
			Desc:     "8 ключей личной эффективности руководителя и создание сильных команд",
		},
		{
			CourseID: uuid.New(),
			Title:    "Коучинг в управлении",
			Desc:     "Коучинговый подход в управлении командой и развитие психоэмоциональной устойчивости лидера",
		},
		{
			CourseID: uuid.New(),
			Title:    "Управление изменениями",
			Desc:     "Методология управления изменениями в организации и лидерство в условиях неопределенности",
		},
		{
			CourseID: uuid.New(),
			Title:    "Деловые переговоры",
			Desc:     "Техники ведения эффективных переговоров и достижения взаимовыгодных соглашений",
		},
		{
			CourseID: uuid.New(),
			Title:    "Управление конфликтами",
			Desc:     "Стратегии разрешения конфликтов и построение конструктивного диалога в команде",
		},
		{
			CourseID: uuid.New(),
			Title:    "Публичные выступления",
			Desc:     "Искусство публичных выступлений и эффективная презентация идей",
		},
	}

	createdCourses := make(map[string]*entity.Course)
	for _, course := range courses {
		var existingCourse entity.Course
		err = courseRepo.DB.Where("title = ?", course.Title).First(&existingCourse).Error
		if err != nil {
			if err.Error() == "record not found" {
				err = courseRepo.AddCourse(course)
				if err != nil {
					panic(fmt.Sprintf("Failed to create course %s: %v", course.Title, err))
				}
				createdCourses[course.Title] = course
				fmt.Printf("Created course: %s\n", course.Title)
			} else {
				panic(fmt.Sprintf("Failed to check existing course %s: %v", course.Title, err))
			}
		} else {
			createdCourses[course.Title] = &existingCourse
			fmt.Printf("Course already exists: %s\n", course.Title)
		}
	}

	// Create test events
	events := []*entity.Event{
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Эмоциональный интеллект и коммуникация"].CourseID,
			Title:       "Введение в эмоциональный интеллект",
			Description: "Основы эмоционального интеллекта и его роль в лидерстве",
			EventDate:   time.Now().Add(48 * time.Hour),
			SecretInfo:  "Zoom Meeting ID: 123-456-789, Password: eqleader",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Эмоциональный интеллект и коммуникация"].CourseID,
			Title:       "Управление эмоциями в стрессовых ситуациях",
			Description: "Практические техники управления эмоциями в сложных ситуациях",
			EventDate:   time.Now().Add(72 * time.Hour),
			SecretInfo:  "Zoom Meeting ID: 987-654-321, Password: stress",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Стратегическое лидерство"].CourseID,
			Title:       "Основы стратегического мышления",
			Description: "Развитие стратегического мышления и видения",
			EventDate:   time.Now().Add(96 * time.Hour),
			SecretInfo:  "Google Meet Link: meet.google.com/strategy-123",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Стресс и Energy Management"].CourseID,
			Title:       "Техники управления энергией",
			Description: "Практические методы управления личной энергией и ресурсами",
			EventDate:   time.Now().Add(120 * time.Hour),
			SecretInfo:  "Discord Server: discord.gg/energy-management",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Самодисциплина и мотивация"].CourseID,
			Title:       "Системы личной мотивации",
			Description: "Построение эффективной системы личной мотивации",
			EventDate:   time.Now().Add(144 * time.Hour),
			SecretInfo:  "Slack Channel: #motivation",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Высокоэффективные команды"].CourseID,
			Title:       "Формирование сильных команд",
			Description: "Принципы создания и развития высокоэффективных команд",
			EventDate:   time.Now().Add(168 * time.Hour),
			SecretInfo:  "Microsoft Teams Link: teams.microsoft.com/team-building",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Коучинг в управлении"].CourseID,
			Title:       "Основы коучингового подхода",
			Description: "Введение в коучинговый подход в управлении",
			EventDate:   time.Now().Add(192 * time.Hour),
			SecretInfo:  "Workshop Materials: leadership-coaching.com/materials",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Управление изменениями"].CourseID,
			Title:       "Методология управления изменениями",
			Description: "Практические инструменты управления изменениями в организации",
			EventDate:   time.Now().Add(216 * time.Hour),
			SecretInfo:  "Zoom Meeting ID: 456-789-012, Password: change",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Деловые переговоры"].CourseID,
			Title:       "Стратегии ведения переговоров",
			Description: "Техники эффективных переговоров и достижения соглашений",
			EventDate:   time.Now().Add(240 * time.Hour),
			SecretInfo:  "Google Meet Link: meet.google.com/negotiations",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Управление конфликтами"].CourseID,
			Title:       "Разрешение конфликтов",
			Description: "Стратегии разрешения конфликтов в команде",
			EventDate:   time.Now().Add(264 * time.Hour),
			SecretInfo:  "Discord Server: discord.gg/conflict-resolution",
		},
		{
			EventID:     uuid.New(),
			CourseID:    createdCourses["Публичные выступления"].CourseID,
			Title:       "Искусство публичных выступлений",
			Description: "Основы эффективных публичных выступлений",
			EventDate:   time.Now().Add(288 * time.Hour),
			SecretInfo:  "Zoom Meeting ID: 789-012-345, Password: public",
		},
	}

	for _, event := range events {
		var existingEvent entity.Event
		err = eventRepo.DB.Where("title = ?", event.Title).First(&existingEvent).Error
		if err != nil {
			if err.Error() == "record not found" {
				err = eventRepo.AddEvent(event)
				if err != nil {
					panic(fmt.Sprintf("Failed to create event %s: %v", event.Title, err))
				}
				fmt.Printf("Created event: %s\n", event.Title)
			} else {
				panic(fmt.Sprintf("Failed to check existing event %s: %v", event.Title, err))
			}
		} else {
			fmt.Printf("Event already exists: %s\n", event.Title)
		}
	}

	// Create course assignments for users
	for _, user := range createdUsers {
		for _, course := range createdCourses {
			progress := uint(0)
			switch course.Title {
			case "Эмоциональный интеллект и коммуникация":
				progress = 75
			case "Стратегическое лидерство":
				progress = 50
			case "Стресс и Energy Management":
				progress = 25
			case "Самодисциплина и мотивация":
				progress = 0
			case "Высокоэффективные команды":
				progress = 100
			case "Коучинг в управлении":
				progress = 30
			case "Управление изменениями":
				progress = 45
			case "Деловые переговоры":
				progress = 60
			case "Управление конфликтами":
				progress = 15
			case "Публичные выступления":
				progress = 90
			}

			var existingAssignment entity.CourseAssignment
			err = courseRepo.DB.Where("user_id = ? AND course_id = ?", user.ID, course.CourseID).First(&existingAssignment).Error
			if err != nil {
				if err.Error() == "record not found" {
					assignment := &entity.CourseAssignment{
						CaID:     uuid.New(),
						UserID:   user.ID,
						CourseID: course.CourseID,
						Progress: progress,
					}
					err = courseRepo.CreateCourseAssignment(assignment)
					if err != nil {
						panic(fmt.Sprintf("Failed to create course assignment for user %s and course %s: %v", user.ID, course.Title, err))
					}
					fmt.Printf("Created course assignment for user %s and course %s\n", user.ID, course.Title)
				} else {
					panic(fmt.Sprintf("Failed to check existing course assignment for user %s and course %s: %v", user.ID, course.Title, err))
				}
			} else {
				fmt.Printf("Course assignment already exists for user %s and course %s\n", user.ID, course.Title)
			}
		}
	}

	// Create course prices
	for _, course := range createdCourses {
		var existingPrice entity.CoursePrice
		err = courseRepo.DB.Where("course_id = ?", course.CourseID).First(&existingPrice).Error
		if err != nil {
			if err.Error() == "record not found" {
				var price float64
				switch course.Title {
				case "Эмоциональный интеллект и коммуникация":
					price = 5000.00
				case "Стратегическое лидерство":
					price = 7500.00
				case "Стресс и Energy Management":
					price = 10000.00
				case "Самодисциплина и мотивация":
					price = 9500.00
				case "Высокоэффективные команды":
					price = 8500.00
				case "Коучинг в управлении":
					price = 6000.00
				case "Управление изменениями":
					price = 8000.00
				case "Деловые переговоры":
					price = 7000.00
				case "Управление конфликтами":
					price = 6500.00
				case "Публичные выступления":
					price = 5500.00
				default:
					price = 5000.00
				}

				coursePrice := &entity.CoursePrice{
					CourseID:     course.CourseID,
					Amount:       price,
					CurrencyCode: "RUB",
				}
				err = courseRepo.DB.Create(coursePrice).Error
				if err != nil {
					panic(fmt.Sprintf("Failed to create price for course %s: %v", course.Title, err))
				}
				fmt.Printf("Created price %f for course %s\n", price, course.Title)
			} else {
				panic(fmt.Sprintf("Failed to check existing price for course %s: %v", course.Title, err))
			}
		} else {
			fmt.Printf("Price already exists for course %s\n", course.Title)
		}
	}

	fmt.Println("Database migration and seeding completed successfully")
}

// Helper function to get environment variable with default value
func getEnvOrDefault(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
