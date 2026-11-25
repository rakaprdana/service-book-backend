package seeders

import (
	models "service-book/domains/models/user"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func RunRoleSeeder(db *gorm.DB) {
	roles := []models.Role{
		{
			Code: "ADMIN",
			Name: "Administrator",
		},
	}

	for _, role := range roles {
		err := db.FirstOrCreate(&role, models.Role{Code: role.Code}).Error
		if err != nil {
			logrus.Errorf("Failed to seed role: %v", err)
			panic(err)
		}
		logrus.Infof("Seeded role: %s is success", role.Code)
	}

}
