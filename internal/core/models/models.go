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
