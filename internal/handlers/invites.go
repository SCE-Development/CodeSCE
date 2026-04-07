package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
)

func CreateInvite(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "create invite route wired",
		"id":      c.Param("id"),
	})
}

func ListInvites(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "list invites route wired"})
}

func ValidateInvite(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "validate invite route wired",
		"token":   c.Param("token"),
	})
}

func StartAttempt(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "start attempt route wired",
		"token":   c.Param("token"),
	})
}