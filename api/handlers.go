package api

import (
	"net/http"
	"pelith-assignment/services"

	"github.com/gin-gonic/gin"
)

func getUserTasksStatus(c *gin.Context) {
	address := c.Param("address")

	err := services.CreateUser(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	status, err := services.GetUserTaskStatus(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, status)
}

func getUserPointsHistory(c *gin.Context) {
	address := c.Param("address")

	err := services.CreateUser(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	history, err := services.GetUserPointsHistory(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, history)
}

func getPointsLeaderBoard(c *gin.Context) {
	leaderBoard, err := services.GetPointsLeaderBoard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, leaderBoard)
}
