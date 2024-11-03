package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type HealthHandler interface {
	CheckHealth(ctx *gin.Context)
	CheckReidsConnect(ctx *gin.Context)
}

type healthHandler struct {
	redisClient *redis.ClusterClient
}

func NewHealthHandler(redisClient *redis.ClusterClient) HealthHandler {
	return &healthHandler{
		redisClient: redisClient,
	}
}

func (h *healthHandler) CheckHealth(ctx *gin.Context) {

	response := map[string]string{
		"message": "ok!",
	}
	ctx.JSON(http.StatusOK, response)
}

func (h *healthHandler) CheckReidsConnect(ctx *gin.Context) {

	var message string

	pong, err := h.redisClient.Ping(ctx).Result()
	if err != nil {
		message = "Ping Fail !"
	} else {
		message = pong
	}

	response := map[string]string{
		"message": message,
	}

	ctx.JSON(http.StatusOK, response)
}
