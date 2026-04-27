package middleware

import (
	"college-diary/internal/models"
	"college-diary/internal/types"
	"net/http"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowed ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, exists := c.Get("claims")
		if !exists {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}

		userRole := claims.(*types.Claims).Role

		for _, role := range allowed {
			if userRole == role {
				c.Next()
				return
			}
		}

		c.HTML(http.StatusForbidden, "index.html", gin.H{
			"Title": "Доступ запрещен!",
			"Error": "Нету прав",
		})

		c.Abort()
	}
}
