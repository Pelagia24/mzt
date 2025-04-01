package app

import (
	"mzt/internal/auth"
	"mzt/internal/auth/entity"
)

func Migrate(r *auth.RefreshTokensRepo) {
	r.DB.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"")
	err := r.DB.AutoMigrate(
		&entity.User{},
		&entity.UserJWT{},
	)
	if err != nil {
		panic(err)
	}

}
