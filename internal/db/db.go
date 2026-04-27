package db

import (
	"college-diary/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() {

	dsn := "host=localhost user=postgres password= dbname= port= sslmode=disable TimeZone=Europe/Moscow"

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Ошибка подключения к базе данных")
	}

	db.AutoMigrate(&models.User{})

	var count int64 
	db.Model(&models.User{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := models.HashPassword("160423")
		admin := models.User{
			Email : "admin@college.ru",
			Name : "Админ",
			Password : hashedPassword,
			Role : models.RoleAdmin,
		}

		db.Create(&admin)
	}

	DB = db
}
