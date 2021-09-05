package migrations

import (
	"myapp/config"
	"myapp/entity"
)

func MigrateTable() {
	db := config.GetDB()

	var models = []interface{}{
		entity.User{},
	}

	db.AutoMigrate(models...)
}
