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
		&entity.Payment{},
		&entity.Course{},
		&entity.CourseAssignment{},
		&entity.Lesson{},
		&entity.Event{},
		&entity.CoursePrice{},
	)
	if err != nil {
		panic(err)
	}

	adminPasswdHash, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	testUsers := []struct {
		user     *entity.User
		userData *entity.UserData
		auth     *entity.Auth
	}{
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: string(adminPasswdHash),
				Role:       1,
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
				PasswdHash: "test_hash2",
				Role:       0,
			},
			userData: &entity.UserData{
				Email:           "student1@example.com",
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
				Key: "student1_key",
			},
		},
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: "test_hash3",
				Role:       0,
			},
			userData: &entity.UserData{
				Email:           "student2@example.com",
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
				Key: "student2_key",
			},
		},
		{
			user: &entity.User{
				ID:         uuid.New(),
				PasswdHash: "test_hash4",
				Role:       0,
			},
			userData: &entity.UserData{
				Email:           "student3@example.com",
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
				Key: "student3_key",
			},
		},
	}

	createdUsers := make(map[string]*entity.User)
	for _, tu := range testUsers {
		var existingUserData entity.UserData
		err = r.DB.Where("email = ?", tu.userData.Email).First(&existingUserData).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tu.userData.UserID = tu.user.ID
				tu.auth.UserID = tu.user.ID

				err = r.CreateUser(tu.user, tu.userData, tu.auth)
				if err != nil {
					panic(err)
				}
				createdUsers[tu.userData.Email] = tu.user
			} else {
				panic(err)
			}
		} else {
			var existingUser entity.User
			err = r.DB.Where("id = ?", existingUserData.UserID).First(&existingUser).Error
			if err != nil {
				panic(err)
			}
			createdUsers[tu.userData.Email] = &existingUser
		}
	}

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
	}

	events := []*entity.Event{
		{
			EventID:     uuid.New(),
			CourseID:    courses[0].CourseID,
			Title:       "Эмоциональный интеллект и коммуникация",
			Description: "Управление негативными эмоциями, развитие коммуникативных навыков лидера. Практические инструменты для развития EQ.",
			EventDate:   time.Now().Add(48 * time.Hour),
			SecretInfo:  "Zoom Meeting ID: 123-456-789, Password: eqleader",
		},
		{
			EventID:     uuid.New(),
			CourseID:    courses[0].CourseID,
			Title:       "Стратегия личного развития лидера",
			Description: "Формирование траектории личного развития руководителя без морального выгорания и потери фокуса внимания",
			EventDate:   time.Now().Add(72 * time.Hour),
			SecretInfo:  "Zoom Meeting ID: 987-654-321, Password: strategy",
		},
		{
			EventID:     uuid.New(),
			CourseID:    courses[1].CourseID,
			Title:       "Управление влиянием и стресс-менеджмент",
			Description: "Основы установления долгосрочных взаимоотношений с ЛПР и техники управления стрессом",
			EventDate:   time.Now().Add(96 * time.Hour),
			SecretInfo:  "Google Meet Link: meet.google.com/influence-123",
		},
		{
			EventID:     uuid.New(),
			CourseID:    courses[2].CourseID,
			Title:       "Energy Management и самодисциплина",
			Description: "Осознанный подход к управлению психо-эмоциональным состоянием и система личной мотивации",
			EventDate:   time.Now().Add(120 * time.Hour),
			SecretInfo:  "Discord Server: discord.gg/energy-management",
		},
		{
			EventID:     uuid.New(),
			CourseID:    courses[3].CourseID,
			Title:       "Эффективное планирование и тайм-менеджмент",
			Description: "Система эффективного планирования и ключевые инструменты тайм-менеджмента",
			EventDate:   time.Now().Add(144 * time.Hour),
			SecretInfo:  "Slack Channel: #time-management",
		},
		{
			EventID:     uuid.New(),
			CourseID:    courses[4].CourseID,
			Title:       "Создание высокоэффективных команд",
			Description: "8 ключей личной эффективности руководителя и принципы создания сильных команд",
			EventDate:   time.Now().Add(168 * time.Hour),
			SecretInfo:  "Microsoft Teams Link: teams.microsoft.com/team-building",
		},
		{
			EventID:     uuid.New(),
			CourseID:    courses[5].CourseID,
			Title:       "Коучинговый подход в лидерстве",
			Description: "Коучинговый подход в управлении командой и развитие психоэмоциональной устойчивости лидера",
			EventDate:   time.Now().Add(192 * time.Hour),
			SecretInfo:  "Workshop Materials: leadership-coaching.com/materials",
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

	for _, event := range events {
		var existingEvent entity.Event
		err = r.DB.Where("title = ?", event.Title).First(&existingEvent).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				err = r.DB.Create(event).Error
				if err != nil {
					panic(err)
				}
			} else {
				panic(err)
			}
		}
	}

	for _, course := range createdCourses {
		var lessons []*entity.Lesson

		switch course.Title {
		case "Эмоциональный интеллект и коммуникация":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Основы эмоционального интеллекта",
					Summery:  "Понимание и развитие эмоционального интеллекта в лидерстве",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Эмоциональный интеллект (EQ) - ключевой навык современного лидера. Основные компоненты:\n\n- Самосознание: понимание своих эмоций\n- Самоконтроль: управление своими эмоциями\n- Социальная осведомленность: понимание эмоций других\n- Управление отношениями: выстраивание эффективных коммуникаций\n\nРазвитый EQ помогает лидеру лучше понимать свою команду и принимать более взвешенные решения.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Управление негативными эмоциями",
					Summery:  "Техники работы с негативными эмоциями в лидерстве",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Управление негативными эмоциями - важный навык лидера. Ключевые техники:\n\n- Распознавание триггеров\n- Техники быстрой саморегуляции\n- Конструктивное выражение эмоций\n- Превращение негатива в мотивацию\n- Работа с эмоциональным выгоранием\n\nЭффективное управление негативными эмоциями позволяет сохранять ясность мышления и принимать взвешенные решения.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Развитие коммуникативных навыков",
					Summery:  "Эффективные коммуникации в лидерстве",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Коммуникативные навыки - основа успешного лидерства. Ключевые аспекты:\n\n- Активное слушание\n- Эмпатическое общение\n- Обратная связь\n- Невербальная коммуникация\n- Публичные выступления\n\nРазвитые коммуникативные навыки помогают выстраивать доверительные отношения в команде.",
				},
			}
		case "Стратегическое лидерство":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Личная стратегия развития",
					Summery:  "Формирование траектории личного развития руководителя",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Стратегия личного развития - фундамент успешного лидерства. Основные элементы:\n\n- Анализ текущих компетенций\n- Постановка целей развития\n- Планирование этапов роста\n- Выбор инструментов развития\n- Оценка прогресса\n\nЧеткая стратегия развития помогает лидеру расти без потери фокуса и выгорания.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Управление влиянием",
					Summery:  "Построение долгосрочных отношений с ЛПР",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Управление влиянием - ключевой навык в работе с ЛПР. Основные аспекты:\n\n- Построение авторитета\n- Networking стратегии\n- Техники убеждения\n- Управление репутацией\n- Развитие личного бренда\n\nЭффективное управление влиянием позволяет достигать целей через выстраивание прочных деловых отношений.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Стратегическое мышление",
					Summery:  "Развитие стратегического мышления лидера",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Стратегическое мышление - необходимый навык современного лидера. Ключевые компоненты:\n\n- Системный анализ\n- Прогнозирование\n- Принятие решений\n- Управление рисками\n- Стратегическое планирование\n\nРазвитое стратегическое мышление позволяет видеть большую картину и принимать эффективные решения.",
				},
			}
		case "Стресс и Energy Management":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Управление стрессом",
					Summery:  "Техники управления стрессом для лидеров",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Управление стрессом - важный навык современного лидера. Основные техники:\n\n- Распознавание стресса\n- Техники быстрой релаксации\n- Управление энергией\n- Профилактика стресса\n- Восстановление после стресса\n\nЭффективное управление стрессом позволяет сохранять высокую производительность.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Energy Management",
					Summery:  "Управление энергией и психо-эмоциональным состоянием",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Energy Management - система управления личной энергией. Ключевые аспекты:\n\n- Физическая энергия\n- Эмоциональная энергия\n- Ментальная энергия\n- Духовная энергия\n- Баланс работы и отдыха\n\nГрамотное управление энергией - основа высокой эффективности лидера.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Психо-эмоциональная устойчивость",
					Summery:  "Развитие психо-эмоциональной устойчивости лидера",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Психо-эмоциональная устойчивость - ключ к стабильной работе лидера. Основные компоненты:\n\n- Эмоциональная гибкость\n- Стрессоустойчивость\n- Адаптивность\n- Ментальная прочность\n- Работа с неопределенностью\n\nРазвитая психо-эмоциональная устойчивость помогает сохранять эффективность в сложных ситуациях.",
				},
			}
		case "Самодисциплина и мотивация":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Система личной мотивации",
					Summery:  "Построение эффективной системы личной мотивации",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Личная мотивация - двигатель развития лидера. Ключевые элементы:\n\n- Определение ценностей\n- Постановка целей\n- Создание системы вознаграждений\n- Работа с демотивацией\n- Поддержание долгосрочной мотивации\n\nЭффективная система мотивации помогает достигать поставленных целей.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Развитие самодисциплины",
					Summery:  "Техники развития самодисциплины лидера",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Самодисциплина - основа личной эффективности. Основные аспекты:\n\n- Формирование привычек\n- Управление временем\n- Постановка приоритетов\n- Преодоление прокрастинации\n- Работа с отвлечениями\n\nРазвитая самодисциплина позволяет достигать целей и поддерживать высокую продуктивность.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Эффективное планирование",
					Summery:  "Система эффективного планирования времени",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Эффективное планирование - ключ к достижению целей. Основные техники:\n\n- Стратегическое планирование\n- Тактическое планирование\n- Приоритизация задач\n- Делегирование\n- Контроль выполнения\n\nГрамотное планирование помогает достигать максимальных результатов при оптимальных затратах ресурсов.",
				},
			}
		case "Высокоэффективные команды":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "8 ключей эффективности",
					Summery:  "Ключевые принципы личной эффективности руководителя",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "8 ключей эффективности руководителя:\n\n1. Проактивность\n2. Целеполагание\n3. Приоритизация\n4. Мышление Win-Win\n5. Эмпатическое слушание\n6. Синергия\n7. Непрерывное развитие\n8. Баланс жизни\n\nПрименение этих принципов помогает достигать максимальной эффективности в управлении.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Создание сильных команд",
					Summery:  "Принципы формирования высокоэффективных команд",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Создание эффективной команды - ключевая задача лидера. Основные этапы:\n\n- Подбор участников\n- Формирование общего видения\n- Распределение ролей\n- Развитие коммуникации\n- Управление конфликтами\n\nПравильно построенная команда способна достигать выдающихся результатов.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Развитие команды",
					Summery:  "Методы развития и повышения эффективности команды",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Развитие команды - непрерывный процесс. Ключевые аспекты:\n\n- Обучение и развитие\n- Мотивация команды\n- Построение доверия\n- Управление результативностью\n- Создание культуры успеха\n\nПостоянное развитие команды - залог долгосрочного успеха организации.",
				},
			}
		case "Коучинг в управлении":
			lessons = []*entity.Lesson{
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Основы коучингового подхода",
					Summery:  "Применение коучинга в управлении командой",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Коучинговый подход в управлении - современный инструмент лидера. Основные принципы:\n\n- Партнерские отношения\n- Раскрытие потенциала\n- Постановка сильных вопросов\n- Фокус на решениях\n- Ответственность и осознанность\n\nКоучинговый подход помогает развивать самостоятельность и инициативность в команде.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Инструменты коучинга",
					Summery:  "Практические инструменты коучинга в управлении",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Инструменты коучинга для эффективного управления. Ключевые техники:\n\n- Модель GROW\n- Шкалирование\n- Колесо баланса\n- Линия времени\n- Техника SMART\n\nПрименение коучинговых инструментов повышает эффективность управления командой.",
				},
				{
					LessonID: uuid.New(),
					CourseID: course.CourseID,
					Title:    "Развитие через коучинг",
					Summery:  "Развитие команды через коучинговый подход",
					VideoURL: "https://www.youtube.com/embed/re5QbW8-Zz4",
					Text:     "Развитие команды через коучинг - эффективный подход к управлению. Основные аспекты:\n\n- Создание среды развития\n- Поддержка инициативы\n- Работа с целями\n- Обратная связь\n- Празднование успехов\n\nКоучинговый подход к развитию создает сильную и самостоятельную команду.",
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

	for _, course := range createdCourses {
		var existingPrice entity.CoursePrice
		err = r.DB.Where("course_id = ?", course.CourseID).First(&existingPrice).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
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
				default:
					price = 5000.00
				}

				coursePrice := &entity.CoursePrice{
					CourseID:     course.CourseID,
					Amount:       price,
					CurrencyCode: "RUB",
				}
				err = r.DB.Create(coursePrice).Error
				if err != nil {
					panic(err)
				}
				fmt.Printf("Created price %f for course %s\n", price, course.Title)
			}
		}
	}

	fmt.Println("Course price initialization completed")

	payments := []struct {
		userEmail   string
		courseTitle string
		amount      float64
		status      string
		paymentRef  string
		createdDays int
	}{
		{
			userEmail:   "student1@example.com",
			courseTitle: "Эмоциональный интеллект и коммуникация",
			amount:      5000.00,
			status:      "completed",
			paymentRef:  "pay_ref_1",
			createdDays: -30,
		},
		{
			userEmail:   "student1@example.com",
			courseTitle: "Стратегическое лидерство",
			amount:      7500.00,
			status:      "completed",
			paymentRef:  "pay_ref_2",
			createdDays: -20,
		},
		{
			userEmail:   "student2@example.com",
			courseTitle: "Высокоэффективные команды",
			amount:      10000.00,
			status:      "completed",
			paymentRef:  "pay_ref_3",
			createdDays: -15,
		},
		{
			userEmail:   "student2@example.com",
			courseTitle: "Коучинг в управлении",
			amount:      8500.00,
			status:      "completed",
			paymentRef:  "pay_ref_4",
			createdDays: -10,
		},
		{
			userEmail:   "student3@example.com",
			courseTitle: "Эмоциональный интеллект и коммуникация",
			amount:      5000.00,
			status:      "completed",
			paymentRef:  "pay_ref_5",
			createdDays: -5,
		},
		{
			userEmail:   "student3@example.com",
			courseTitle: "Самодисциплина и мотивация",
			amount:      9500.00,
			status:      "pending",
			paymentRef:  "pay_ref_6",
			createdDays: -1,
		},
	}

	for _, p := range payments {
		user := createdUsers[p.userEmail]
		course := createdCourses[p.courseTitle]
		if user != nil && course != nil {
			payment := &entity.Payment{
				PaymentID:    uuid.New(),
				UserID:       user.ID,
				CourseID:     course.CourseID,
				Amount:       p.amount,
				CurrencyCode: "RUB",
				CreatedAt:    time.Now().AddDate(0, 0, p.createdDays),
				Status:       p.status,
				PaymentRef:   p.paymentRef,
			}

			var existingPayment entity.Payment
			err = r.DB.Where("payment_ref = ?", p.paymentRef).First(&existingPayment).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					err = r.DB.Create(payment).Error
					if err != nil {
						panic(err)
					}
				} else {
					panic(err)
				}
			}
		}
	}
}
