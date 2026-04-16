package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func CreateAssessment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "create assessment route wired"})
}

func ListAssessments(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list assessments route wired"})
}

func GetAssessment(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "get assessment route wired",
		"id":      c.Param("id"),
	})
}

func AddQuestion(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "add question route wired",
		"id":      c.Param("id"),
	})
}

func AddTestCases(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "add test cases route wired",
		"id":      c.Param("id"),
		"qid":     c.Param("qid"),
	})
}
