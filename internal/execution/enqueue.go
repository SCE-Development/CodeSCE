package execution

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

const SubmissionsQueue = "codesce:submissions"

// same as worker.go Job
type EnqueueJob struct {
	SubmissionID int64      `json:"submission_id"`
	QuestionID   int64      `json:"question_id"`
	Language     string     `json:"language"`
	SourceCode   string     `json:"source_code"`
	TestCases    []TestCase `json:"test_cases"`
}

type TestCase struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	Score          int    `json:"score"`
}

// push a job onto redis queue
func Enqueue(ctx context.Context, rdb *redis.Client, job EnqueueJob) error {
	data, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("marshal job: %w", err)
	}

	if err := rdb.LPush(ctx, SubmissionsQueue, data).Err(); err != nil {
		return fmt.Errorf("enqueue job: %w", err)
	}

	return nil
}