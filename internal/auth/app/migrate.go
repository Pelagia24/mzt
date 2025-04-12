package app

import (
	"mzt/internal/auth"
	"mzt/internal/auth/entity"
)

func Migrate(r *auth.UserRepo) {
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

}
