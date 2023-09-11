package initializers

import "github.com/jenniekibiri/jwt-go/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
