package services

import (
	"CodeSCE/internal/core/models"
	"context"
	"errors"
	"fmt"
)

var (
	ErrAssessmentNotFound = errors.New("assessment not found")
	ErrQuestionNotFound   = errors.New("question not found")
	ErrWrongAssessment    = errors.New("question does not belong to this assessment")
)

type FieldError struct {
	Field   string
	Message string
}

func (e FieldError) Error() string {
	return e.Field + ": " + e.Message
}

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

type QuestionWithTestCases struct {
	models.Question
	TestCases []models.TestCase `json:"testCases,omitempty"`
}

type AssessmentWithQuestions struct {
	models.AssessmentTemplate
	Questions []QuestionWithTestCases `json:"questions"`
}

type AssessmentService struct {
	assessmentRepo AssessmentRepository
	questionRepo   QuestionRepository
	testCaseRepo   TestCaseRepository
}

func NewAssessmentService(
	ar AssessmentRepository,
	qr QuestionRepository,
	tr TestCaseRepository,
) *AssessmentService {
	return &AssessmentService{
		assessmentRepo: ar,
		questionRepo:   qr,
		testCaseRepo:   tr,
	}
}

func (s *AssessmentService) CreateAssessment(ctx context.Context, assessment *models.AssessmentTemplate) error {
	if assessment == nil {
		return FieldError{Field: "assessment", Message: "required"}
	}
	if err := validateAssessment(assessment); err != nil {
		return err
	}
	return s.assessmentRepo.Create(ctx, assessment)
}

func (s *AssessmentService) ListAssessments(ctx context.Context) ([]models.AssessmentTemplate, error) {
	return s.assessmentRepo.List(ctx)
}

func (s *AssessmentService) GetAssessment(ctx context.Context, id int64) (*AssessmentWithQuestions, error) {
	a, err := s.assessmentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, ErrAssessmentNotFound
	}
	questions, err := s.questionRepo.ListByAssessmentID(ctx, id)
	if err != nil {
		return nil, err
	}
	result := make([]QuestionWithTestCases, len(questions))
	for i := range questions {
		result[i].Question = questions[i]
		if questions[i].Type == "code" {
			tcs, err := s.testCaseRepo.ListByQuestionID(ctx, questions[i].ID)
			if err != nil {
				return nil, err
			}
			result[i].TestCases = tcs
		}
	}
	return &AssessmentWithQuestions{AssessmentTemplate: *a, Questions: result}, nil
}

func (s *AssessmentService) AddQuestion(ctx context.Context, assessmentID int64, question *models.Question) error {
	if question == nil {
		return FieldError{Field: "question", Message: "required"}
	}
	question.AssessmentTemplateID = assessmentID
	if err := validateQuestion(question); err != nil {
		return err
	}
	a, err := s.assessmentRepo.GetByID(ctx, assessmentID)
	if err != nil {
		return err
	}
	if a == nil {
		return ErrAssessmentNotFound
	}
	return s.questionRepo.Create(ctx, question)
}

// AddTestCases persists test cases for a question. assessmentID must match the question's assessment (route :id) so ownership can be enforced.
func (s *AssessmentService) AddTestCases(ctx context.Context, assessmentID, questionID int64, testCases []models.TestCase) error {
	if err := validateTestCases(testCases); err != nil {
		return err
	}
	q, err := s.questionRepo.GetByID(ctx, questionID)
	if err != nil {
		return err
	}
	if q == nil {
		return ErrQuestionNotFound
	}
	if q.AssessmentTemplateID != assessmentID {
		return ErrWrongAssessment
	}
	return s.testCaseRepo.CreateMany(ctx, questionID, testCases)
}

func validateAssessment(a *models.AssessmentTemplate) error {
	if a.Title == "" {
		return FieldError{Field: "title", Message: "required"}
	}
	if a.DurationMinutes <= 0 {
		return FieldError{Field: "duration_minutes", Message: "must be greater than 0"}
	}
	return nil
}

func validateQuestion(q *models.Question) error {
	if q.Prompt == "" {
		return FieldError{Field: "prompt", Message: "required"}
	}
	if err := validateQuestionType(q.Type); err != nil {
		return err
	}
	if err := validatePositiveScore(q.Score, "score"); err != nil {
		return err
	}
	if q.Type == "code" {
		if q.Language == nil || *q.Language == "" {
			return FieldError{Field: "language", Message: "required for code questions"}
		}
	}
	if questionTypeRequiresOptions(q.Type) && q.Options == nil {
		return FieldError{Field: "options", Message: "required for MCQ questions"}
	}
	return nil
}

func validateQuestionType(typ string) error {
	switch typ {
	case "mcq", "multi_select", "short_answer", "code":
		return nil
	default:
		return FieldError{Field: "type", Message: "must be mcq, multi_select, short_answer, or code"}
	}
}

func questionTypeRequiresOptions(typ string) bool {
	switch typ {
	case "mcq", "multi_select":
		return true
	default:
		return false
	}
}

func validateTestCases(cases []models.TestCase) error {
	for i := range cases {
		if err := validatePositiveScore(cases[i].Score, fmt.Sprintf("test_cases[%d].score", i)); err != nil {
			return err
		}
	}
	return nil
}

func validatePositiveScore(score int, field string) error {
	if score <= 0 {
		return FieldError{Field: field, Message: "must be greater than 0"}
	}
	return nil
}
