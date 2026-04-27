// internal/handlers/handlers.go
package handlers

import (
	"college-diary/internal/db"
	"college-diary/internal/models"
	"college-diary/internal/types"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte("VeRySeCrEtKey")

func Home(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		user = nil
	}
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Главная",
		"User":  user,
	})
}

func LoginPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"Title": "Вход"})
}

func RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"Title": "Регистрация"})
}

func Register(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")
	name := c.PostForm("name")

	if email == "" || password == "" || name == "" {
		c.String(http.StatusBadRequest, "Заполните все поля")
		return
	}

	var exists models.User
	if db.DB.Where("email = ?", email).First(&exists).Error == nil {
		c.String(http.StatusBadRequest, "Этот email уже занят")
		return
	}

	hashed, err := models.HashPassword(password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка сервера")
		return
	}

	user := models.User{
		Email:    email,
		Name:     name,
		Password: hashed,
		Role:     models.RoleStudent,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "Ошибка создания пользователя")
		return
	}

	c.Redirect(http.StatusFound, "/login")
}

func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	if err := db.DB.Where("email = ?", email).First(&user).Error; err != nil {
		c.String(http.StatusBadRequest, "Неверный email или пароль")
		return
	}

	if !user.CheckPassword(password) {
		c.String(http.StatusBadRequest, "Неверный email или пароль")
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &types.Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		c.String(http.StatusInternalServerError, "Ошибка авторизации")
		return
	}

	c.SetCookie("token", tokenString, 86400, "/", "", false, true)
	c.Redirect(http.StatusFound, "/dashboard")
}

func Dashboard(c *gin.Context) {
	user, _ := c.Get("user")
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Личный кабинет",
		"User":  user,
	})
}

func Schedule(c *gin.Context) {
	c.String(http.StatusOK, "Расписание (скоро будет)")
}

func Grades(c *gin.Context) {
	c.String(http.StatusOK, "Оценки (скоро будет)")
}

func GradesEdit(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title":    "Выставить оценки",
		"User":     user,
		"EditMode": true,
	})
}

func GradesSave(c *gin.Context) {
	c.String(http.StatusOK, "Оценки сохранены!")
}

func AdminPanel(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.HTML(http.StatusOK, "index.html", gin.H{
		"Title": "Админ-панель",
		"User":  user,
		"Admin": true,
	})
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("token")
		if err != nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		claims := &types.Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		var user models.User
		db.DB.First(&user, claims.UserID)
		c.Set("user", user)
		c.Set("claims", claims)
		c.Next()
	}
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "", false, true)
	c.Redirect(http.StatusFound, "/")
}
