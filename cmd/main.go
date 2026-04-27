package main

import (
	"college-diary/internal/handlers"
	"college-diary/internal/middleware"
	"college-diary/internal/models"
	"github.com/gin-gonic/gin"
	"college-diary/internal/db"

)

func main() {

	//gin.SetMode(gin.ReleaseMode)
	db.Init()

	r := gin.Default()

	r.LoadHTMLGlob("templates/index.html")
	r.Static("/static", "./static")

	r.GET("/", handlers.Home)
	r.GET("/login", handlers.LoginPage)
	r.POST("/login", handlers.Login)
	r.GET("/register", handlers.RegisterPage)
	r.POST("/register", handlers.Register)
	r.GET("/logout", handlers.Logout)

	authorized := r.Group("/")
	authorized.Use(handlers.AuthRequired())

	{
		authorized.GET("/dashboard", handlers.Dashboard)
		authorized.GET("/schedule", handlers.Schedule)
		authorized.GET("/grades", handlers.Grades)

		authorized.GET("/grades/edit", middleware.RequireRole(models.RoleTeacher, models.RoleAdmin), handlers.GradesEdit)
		authorized.GET("/grades/save", middleware.RequireRole(models.RoleTeacher, models.RoleAdmin), handlers.GradesSave)

		authorized.GET("/admin", middleware.RequireRole(models.RoleAdmin), handlers.AdminPanel)
	}

	r.Run(":8080")
}
