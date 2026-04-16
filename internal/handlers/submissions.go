package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func CreateSubmission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message":       "create submission route wired",
		"attempt_id":    c.Param("id"),
		"question_id":   c.Param("qid"),
		"submission_id": "mock-submission-id",
	})
}

func GetSubmission(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get submission route wired",
		"id":      c.Param("id"),
		"status":  "queued",
	})
}