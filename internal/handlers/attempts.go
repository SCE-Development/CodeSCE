package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func ListAttempts(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":      "list attempts route wired",
		"assessmentId": c.Query("assessmentId"),
	})
}

func GetAttempt(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get attempt route wired",
		"id":      c.Param("id"),
	})
}

func GetAttemptQuestions(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get attempt questions route wired",
		"id":      c.Param("id"),
	})
}

func SaveAnswers(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "save answers route wired",
		"id":      c.Param("id"),
	})
}

func SubmitAttempt(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "submit attempt route wired",
		"id":      c.Param("id"),
	})
}

func RunCode(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"attempt_id":  c.Param("id"),
		"question_id": c.Param("qid"),
		"passed_test_cases": []int{},
	})
}
