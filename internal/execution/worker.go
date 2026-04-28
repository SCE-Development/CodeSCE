package execution

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

const SubmissionsQueue = "codesce:submissions"

// create worker instance
type Worker struct {
	Redis *redis.Client
	DB    *pgxpool.Pool
	ID    int

	ExecutionTimeout   time.Duration
	SandboxMemoryLimit string
}

// payload from redis
type Job struct {
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

func (w *Worker) Run(ctx context.Context) {
	log.Printf("worker %d started", w.ID)

	for {
		// graceful shutdown shi
		select {
		case <-ctx.Done():
			log.Printf("worker %d stopping", w.ID)
			return
		default:
		}
		// get j*b from redis
		result, err := w.Redis.BRPop(ctx, 5*time.Second, SubmissionsQueue).Result()
		if err != nil {
			if err == redis.Nil || ctx.Err() != nil {
				continue
			}
			log.Printf("worker %d BRPOP error: %v", w.ID, err)
			time.Sleep(1 * time.Second)
			continue
		}

		if len(result) != 2 {
			log.Printf("worker %d unexpected BRPOP result: %#v", w.ID, result)
			continue
		}
		// get payload json string
		payload := result[1]
		// process job
		if err := w.handleJob(ctx, payload); err != nil {
			log.Printf("worker %d failed job: %v", w.ID, err)
		}
	}
}

func (w *Worker) handleJob(ctx context.Context, payload string) error {
	var job Job
	if err := json.Unmarshal([]byte(payload), &job); err != nil {
		return fmt.Errorf("unmarshal job: %w", err)
	}

	log.Printf("worker %d processing submission=%d language=%s", w.ID, job.SubmissionID, job.Language)

	if _, err := w.DB.Exec(ctx, `
		UPDATE submissions
		SET status = 'running', updated_at = NOW()
		WHERE id = $1
	`, job.SubmissionID); err != nil {
		return fmt.Errorf("update submission running: %w", err)
	}

	if _, err := w.DB.Exec(ctx, `
		DELETE FROM submission_results
		WHERE submission_id = $1
	`, job.SubmissionID); err != nil {
		return fmt.Errorf("delete old submission results: %w", err)
	}

	totalScore := 0
	scoreAwarded := 0
	totalExecutionMS := int64(0)

	finalStatus := "passed"
	finalStdout := ""
	finalStderr := ""

	// loop through test cases
	for i, tc := range job.TestCases {
		totalScore += tc.Score

		// run code in sandbox
		runCtx, cancel := context.WithTimeout(ctx, w.ExecutionTimeout)
		execResult, err := Run(runCtx, job.Language, job.SourceCode, tc.Input)
		cancel()

		passed := false
		caseStatus := "failed"
		stdout := ""
		stderr := ""
		exitCode := -1
		executionMS := int64(0)

		// get result
		if execResult != nil {
			stdout = execResult.Stdout
			stderr = execResult.Stderr
			exitCode = execResult.Code
		}

		// runtime error
		if err != nil {
			stderr = err.Error()
			caseStatus = "runtime_error"
			if finalStatus == "passed" {
				finalStatus = "runtime_error"
			}
		} else {
				passed = strings.TrimSpace(stdout) == strings.TrimSpace(tc.ExpectedOutput)			if exitCode != 0 {
				caseStatus = classifyStatus(stderr)
				if finalStatus == "passed" {
					finalStatus = caseStatus
				}
			} else if passed {
				caseStatus = "passed"
				scoreAwarded += tc.Score
			} else {
				caseStatus = "failed"
				if finalStatus == "passed" {
					finalStatus = "failed"
				}
			}
		}

		if finalStdout == "" {
			finalStdout = stdout
		}
		if finalStderr == "" && stderr != "" {
			finalStderr = stderr
		}
		totalExecutionMS += executionMS

		if _, err := w.DB.Exec(ctx, `
			INSERT INTO submission_results (
				submission_id,
				test_case_index,
				stdout,
				stderr,
				expected_output,
				actual_output,
				passed,
				score_awarded,
				execution_ms,
				created_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW())
		`,
			job.SubmissionID,
			i,
			stdout,
			stderr,
			tc.ExpectedOutput,
			stdout,
			passed,
			func() int {
				if passed {
					return tc.Score
				}
				return 0
			}(),
			executionMS,
		); err != nil {
			return fmt.Errorf("insert submission_result: %w", err)
		}
	}

	if len(job.TestCases) == 0 {
		finalStatus = "failed"
	}

	if _, err := w.DB.Exec(ctx, `
		UPDATE submissions
		SET status = $2,
		    stdout = $3,
		    stderr = $4,
		    execution_ms = $5,
		    score_awarded = $6,
		    updated_at = NOW()
		WHERE id = $1
	`,
		job.SubmissionID,
		finalStatus,
		finalStdout,
		finalStderr,
		totalExecutionMS,
		scoreAwarded,
	); err != nil {
		return fmt.Errorf("update final submission: %w", err)
	}

	log.Printf(
		"worker %d completed submission=%d status=%s score=%d/%d",
		w.ID,
		job.SubmissionID,
		finalStatus,
		scoreAwarded,
		totalScore,
	)

	return nil
}

func classifyStatus(stderr string) string {
	lower := strings.ToLower(stderr)
	if strings.Contains(lower, "syntaxerror") || strings.Contains(lower, "compile") {
		return "compile_error"
	}
	return "runtime_error"
}
