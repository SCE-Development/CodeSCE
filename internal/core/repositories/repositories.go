package repositories

import (
	"CodeSCE/internal/core/models"
	"context"
)

type AssessmentRepository interface {
	Create(ctx context.Context, assessment *models.AssessmentTemplate) error
	List(ctx context.Context) ([]models.AssessmentTemplate, error)
	GetByID(ctx context.Context, id int64) (*models.AssessmentTemplate, error)
}

type QuestionRepository interface {
	Create(ctx context.Context, question *models.Question) error
	GetByID(ctx context.Context, id int64) (*models.Question, error)
	ListByAssessmentID(ctx context.Context, assessmentID int64) ([]models.Question, error)
}

type TestCaseRepository interface {
	CreateMany(ctx context.Context, questionID int64, testCases []models.TestCase) error
	ListByQuestionID(ctx context.Context, questionID int64) ([]models.TestCase, error)
}
