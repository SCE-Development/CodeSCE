package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"

	"CodeSCE/internal/db"
	"CodeSCE/internal/handlers"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	database, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	if err := db.InitSchema(database); err != nil {
		log.Fatalf("failed to initialize schema: %v", err)
	}

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := r.Group("/api")
	api.POST("/assessments", handlers.CreateAssessment)
	api.GET("/assessments", handlers.ListAssessments)
	api.GET("/assessments/:id", handlers.GetAssessment)
	api.POST("/assessments/:id/questions", handlers.AddQuestion)
	api.POST("/assessments/:id/questions/:qid/test-cases", handlers.AddTestCases)

	api.POST("/assessments/:id/invites", handlers.CreateInvite)
	api.GET("/invites", handlers.ListInvites)
	api.GET("/invites/:token", handlers.ValidateInvite)
	api.POST("/invites/:token/start", handlers.StartAttempt)

	api.GET("/attempts", handlers.ListAttempts)
	api.GET("/attempts/:id", handlers.GetAttempt)
	api.GET("/attempts/:id/questions", handlers.GetAttemptQuestions)
	api.POST("/attempts/:id/answers", handlers.SaveAnswers)
	api.POST("/attempts/:id/submit", handlers.SubmitAttempt)

	api.POST("/attempts/:id/questions/:qid/submissions", handlers.CreateSubmission)
	api.GET("/submissions/:id", handlers.GetSubmission)
	api.POST("/attempts/:id/questions/:qid/run", handlers.RunCode)
	r.Run(":6767")
}
