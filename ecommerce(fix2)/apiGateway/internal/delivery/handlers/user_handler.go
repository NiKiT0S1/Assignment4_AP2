package handlers

import (
	grpcDelivery "apiGateway/internal/grpc"
	"apiGateway/internal/proto"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, userClient *grpcDelivery.UserClient, redisClient *grpcDelivery.RedisClient) {
	r.POST("/register", func(c *gin.Context) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}
		user, err := userClient.Register(body.Username, body.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "registration failed: " + err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"id": user.Id, "username": user.Username})
	})

	r.POST("/login", func(c *gin.Context) {
		var body struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// Создаем запрос на авторизацию
		req := &proto.AuthRequest{
			Username: body.Username,
			Password: body.Password,
		}

		user, err := userClient.Authenticate(c, req)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"user_id": user.Id, "token": "Bearer " + strconv.Itoa(int(user.Id))})
	})

	r.GET("/profile/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		// Try to get from cache first
		cacheKey := fmt.Sprintf("user:%d", id)
		var cachedUser struct {
			Id       int32  `json:"id"`
			Username string `json:"username"`
		}

		err = redisClient.Get(c.Request.Context(), cacheKey, &cachedUser)
		if err == nil {
			c.JSON(http.StatusOK, cachedUser)
			return
		}

		// If not in cache, get from service
		user, err := userClient.GetProfile(int32(id))
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
			return
		}

		// Cache the result
		userData := struct {
			Id       int32  `json:"id"`
			Username string `json:"username"`
		}{
			Id:       user.Id,
			Username: user.Username,
		}
		redisClient.Set(c.Request.Context(), cacheKey, userData, 30*time.Minute)

		c.JSON(http.StatusOK, userData)
	})

	// Add cache invalidation for user updates
	r.PUT("/profile/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
			return
		}

		var body struct {
			Username string `json:"username"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
			return
		}

		// Update user in service
		req := &proto.UpdateProfileRequest{
			Id:       int32(id),
			Username: body.Username,
		}
		user, err := userClient.UpdateProfile(c.Request.Context(), req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "update failed"})
			return
		}

		// Invalidate cache
		cacheKey := fmt.Sprintf("user:%d", id)
		redisClient.Delete(c.Request.Context(), cacheKey)

		c.JSON(http.StatusOK, gin.H{"id": user.Id, "username": user.Username})
	})
}
