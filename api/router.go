package api

import (
	"github.com/gin-gonic/gin"
)

func InitAPIService() {
	router := gin.Default()

	user := router.Group("/api/user")
	{
		user.GET("/:address/tasks", getUserTasksStatus)
		user.GET("/:address/points-history", getUserPointsHistory)
		user.GET("/leaderboard", getPointsLeaderBoard)
	}

	router.Run(":8080")
}
