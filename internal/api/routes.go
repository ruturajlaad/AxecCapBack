package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ruturajlaad/TradeExecEngine/internal/middleware"
	"github.com/ruturajlaad/TradeExecEngine/internal/trade"
	"github.com/ruturajlaad/TradeExecEngine/internal/user"
)

func SetupRoutes(r *gin.Engine, tradeHandler *trade.Handler, tradeService *trade.Service) {

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Welcome to AxeCap!"})
	})

	r.POST("/register", func(c *gin.Context) {
		var req struct{ Username, Password string }
		if c.BindJSON(&req) != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Output"})
			return
		}

		if err := user.CreateUser(req.Username, req.Password); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"User": "user created"})

	})

	r.POST("/login", func(c *gin.Context) {
		var req struct{ Username, Password string }

		c.BindJSON(&req)
		u, err := user.Authenticate(req.Username, req.Password)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials!"})
			return
		}
		token, _ := middleware.GenerateToken(u.Username)
		c.JSON(http.StatusOK, gin.H{"token": token})
	})

	auth := r.Group("/api", middleware.AuthMiddleware())
	auth.POST("/order", tradeHandler.PlaceOrder)
	auth.GET("/position", tradeHandler.GetOpenPositions)

	//user routes
	userRoutes := r.Group("user", middleware.AuthMiddleware())
	{
		userRoutes.GET("/balance", func(c *gin.Context) {
			userID := c.GetInt("UserID")

			balance, err := tradeService.GetUserBalance(userID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"user_id": userID,
				"balance": balance,
			})
		})
	}
}
