package repository

import (
	"mzt/config"
	"mzt/internal/entity"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	cfg := &config.Config{
		DB: config.DB{
			Host:     "localhost",
			Port:     "5433",
			User:     "postgres",
			Password: "postgres",
			Name:     "mzt_test",
		},
	}

	db := connectDB(cfg)

	time.Sleep(time.Second)

	err := db.Migrator().DropTable(
		&entity.Course{},
		&entity.Lesson{},
		&entity.CourseAssignment{},
		&entity.User{},
		&entity.UserData{},
		&entity.Auth{},
	)
	require.NoError(t, err)

	err = db.AutoMigrate(
		&entity.Course{},
		&entity.Lesson{},
		&entity.CourseAssignment{},
		&entity.User{},
		&entity.UserData{},
		&entity.Auth{},
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := db.Migrator().DropTable(
			&entity.Course{},
			&entity.Lesson{},
			&entity.CourseAssignment{},
			&entity.User{},
			&entity.UserData{},
			&entity.Auth{},
		)
		require.NoError(t, err)
	})

	return db
}
