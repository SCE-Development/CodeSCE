package models

import (
	"encoding/json"
	"time"
)

type AssessmentTemplate struct {
	ID              int64     `json:"id" db:"id"`
	Title           string    `json:"title" db:"title"`
	Description     *string   `json:"description,omitempty" db:"description"`
	DurationMinutes int       `json:"durationMinutes" db:"duration_minutes"`
	CreatedAt       time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt       time.Time `json:"updatedAt" db:"updated_at"`
}

type Question struct {
	ID                   int64            `json:"id" db:"id"`
	AssessmentTemplateID int64            `json:"assessmentTemplateId" db:"assessment_template_id"`
	Type                 string           `json:"type" db:"type"`
	Prompt               string           `json:"prompt" db:"prompt"`
	Options              *json.RawMessage `json:"options,omitempty" db:"options"`
	CorrectAnswer        *json.RawMessage `json:"correctAnswer,omitempty" db:"correct_answer"`
	Score                int              `json:"score" db:"score"`
	Language             *string          `json:"language,omitempty" db:"language"`
}

type TestCase struct {
	ID             int64  `json:"id" db:"id"`
	QuestionID     int64  `json:"questionId" db:"question_id"`
	Input          string `json:"input" db:"input"`
	ExpectedOutput string `json:"expectedOutput" db:"expected_output"`
	IsHidden       bool   `json:"isHidden" db:"is_hidden"`
	Score          int    `json:"score" db:"score"`
}

type AssessmentInvite struct {
	ID                   int64     `json:"id" db:"id"`
	AssessmentTemplateID int64     `json:"assessmentTemplateId" db:"assessment_template_id"`
	CandidateEmail       string    `json:"candidateEmail" db:"candidate_email"`
	InviteToken          string    `json:"inviteToken" db:"invite_token"`
	ExpiresAt            time.Time `json:"expiresAt" db:"expires_at"`
	Status               string    `json:"status" db:"status"`
	CreatedAt            time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt            time.Time `json:"updatedAt" db:"updated_at"`
}

type AssessmentAttempt struct {
	ID          int64      `json:"id" db:"id"`
	InviteID    int64      `json:"inviteId" db:"invite_id"`
	StartedAt   *time.Time `json:"startedAt,omitempty" db:"started_at"`
	CompletedAt *time.Time `json:"completedAt,omitempty" db:"completed_at"`
	Status      string     `json:"status" db:"status"`
	TotalScore  *int       `json:"totalScore,omitempty" db:"total_score"`
	CreatedAt   time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" db:"updated_at"`
}

type Answer struct {
	ID           int64           `json:"id" db:"id"`
	AttemptID    int64           `json:"attemptId" db:"attempt_id"`
	QuestionID   int64           `json:"questionId" db:"question_id"`
	Response     json.RawMessage `json:"response" db:"response"`
	ScoreAwarded *int            `json:"scoreAwarded,omitempty" db:"score_awarded"`
	CreatedAt    time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time       `json:"updatedAt" db:"updated_at"`
}

type Submission struct {
	ID           int64     `json:"id" db:"id"`
	AttemptID    int64     `json:"attemptId" db:"attempt_id"`
	QuestionID   int64     `json:"questionId" db:"question_id"`
	Language     string    `json:"language" db:"language"`
	SourceCode   string    `json:"sourceCode" db:"source_code"`
	Status       string    `json:"status" db:"status"`
	Stdout       *string   `json:"stdout,omitempty" db:"stdout"`
	Stderr       *string   `json:"stderr,omitempty" db:"stderr"`
	ExecutionMS  *int      `json:"executionMs,omitempty" db:"execution_ms"`
	MemoryBytes  *int64    `json:"memoryBytes,omitempty" db:"memory_bytes"`
	ScoreAwarded *int      `json:"scoreAwarded,omitempty" db:"score_awarded"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
}

type SubmissionResult struct {
	ID           int64     `json:"id" db:"id"`
	SubmissionID int64     `json:"submissionId" db:"submission_id"`
	TestCaseID   int64     `json:"testCaseId" db:"test_case_id"`
	Passed       bool      `json:"passed" db:"passed"`
	ActualOutput string    `json:"actualOutput" db:"actual_output"`
	ExecutionMS  *int      `json:"executionMs,omitempty" db:"execution_ms"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
}

// code execution job from the Redis stream
// shoutout goderpad/models/job.go for da code yo!
type ExecuteJob struct {
	JobID    string `json:"jobId"`
	RoomID   string `json:"roomId"`
	UserID   string `json:"userId"`
	Code     string `json:"code"`
	Language string `json:"language"`
}

// result of a code execution
type ExecuteResult struct {
	JobID  string `json:"jobId"`
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
	Stdout string `json:"stdout"`
	Stderr string `json:"stderr"`
	Code   int    `json:"code"`
}
