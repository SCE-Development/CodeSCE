package handlers

import (
	"CodeSCE/internal/core/models"
	"CodeSCE/internal/core/services"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

var assessmentService *services.AssessmentService

func SetAssessmentService(s *services.AssessmentService) {
	assessmentService = s
}

func CreateAssessment(c *gin.Context) {
	if assessmentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "assessment service not configured"})
		return
	}

	var req struct {
		Title           string  `json:"title"`
		Description     *string `json:"description"`
		DurationMinutes int     `json:"duration_minutes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(req.Title) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title: required"})
		return
	}
	if req.DurationMinutes <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "duration_minutes: must be greater than 0"})
		return
	}

	assessment := models.AssessmentTemplate{
		Title:           req.Title,
		Description:     req.Description,
		DurationMinutes: req.DurationMinutes,
	}

	if err := assessmentService.CreateAssessment(c.Request.Context(), &assessment); err != nil {
		handleAssessmentError(c, err)
		return
	}

	c.JSON(http.StatusCreated, assessment)
}

func ListAssessments(c *gin.Context) {
	if assessmentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "assessment service not configured"})
		return
	}

	assessments, err := assessmentService.ListAssessments(c.Request.Context())
	if err != nil {
		handleAssessmentError(c, err)
		return
	}

	c.JSON(http.StatusOK, assessments)
}

func GetAssessment(c *gin.Context) {
	if assessmentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "assessment service not configured"})
		return
	}

	assessmentID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assessment, svcErr := assessmentService.GetAssessment(c.Request.Context(), assessmentID)
	if svcErr != nil {
		handleAssessmentError(c, svcErr)
		return
	}

	c.JSON(http.StatusOK, assessment)
}

func AddQuestion(c *gin.Context) {
	if assessmentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "assessment service not configured"})
		return
	}

	assessmentID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var req struct {
		Type          string           `json:"type"`
		Prompt        string           `json:"prompt"`
		Options       *json.RawMessage `json:"options"`
		CorrectAnswer *json.RawMessage `json:"correct_answer"`
		Score         int              `json:"score"`
		Language      *string          `json:"language"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}
	if strings.TrimSpace(req.Type) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type: required"})
		return
	}
	if !isValidQuestionType(req.Type) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "type: must be mcq_single, mcq_multi, short_answer, or code"})
		return
	}
	if strings.TrimSpace(req.Prompt) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "prompt: required"})
		return
	}
	if req.Score <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "score: must be greater than 0"})
		return
	}

	question := models.Question{
		Type:          req.Type,
		Prompt:        req.Prompt,
		Options:       req.Options,
		CorrectAnswer: req.CorrectAnswer,
		Score:         req.Score,
		Language:      req.Language,
	}

	if err := assessmentService.AddQuestion(c.Request.Context(), assessmentID, &question); err != nil {
		handleAssessmentError(c, err)
		return
	}

	c.JSON(http.StatusCreated, question)
}

func AddTestCases(c *gin.Context) {
	if assessmentService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "assessment service not configured"})
		return
	}

	assessmentID, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	questionID, err := parseIDParam(c, "qid")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var payload struct {
		TestCases []models.TestCase `json:"testCases"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	if err := assessmentService.AddTestCases(c.Request.Context(), assessmentID, questionID, payload.TestCases); err != nil {
		handleAssessmentError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"testCases": payload.TestCases})
}

func parseIDParam(c *gin.Context, name string) (int64, error) {
	value := c.Param(name)
	id, err := strconv.ParseInt(value, 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("invalid " + name)
	}
	return id, nil
}

func handleAssessmentError(c *gin.Context, err error) {
	var fieldErr services.FieldError
	if errors.As(err, &fieldErr) {
		c.JSON(http.StatusBadRequest, gin.H{"error": fieldErr.Error()})
		return
	}

	switch {
	case errors.Is(err, services.ErrAssessmentNotFound), errors.Is(err, services.ErrQuestionNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	case errors.Is(err, services.ErrWrongAssessment):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func isValidQuestionType(questionType string) bool {
	switch questionType {
	case "mcq_single", "mcq_multi", "short_answer", "code":
		return true
	default:
		return false
	}
}
