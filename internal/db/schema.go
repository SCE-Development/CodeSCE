package db

import (
	"database/sql"
	"fmt"
)

var schemas = []string{
	// 1. Assessment templates
	`CREATE TABLE IF NOT EXISTS assessment_templates (
		id            BIGSERIAL PRIMARY KEY,
		title         TEXT NOT NULL,
		description   TEXT,
		duration_minutes INTEGER NOT NULL,
		created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
	)`,

	// 2. Questions
	`CREATE TABLE IF NOT EXISTS questions (
		id                      BIGSERIAL PRIMARY KEY,
		assessment_template_id  BIGINT NOT NULL REFERENCES assessment_templates(id) ON DELETE CASCADE,
		type                    TEXT NOT NULL CHECK (type IN ('mcq', 'multi_select', 'short_answer', 'code')),
		prompt                  TEXT NOT NULL,
		language                TEXT,
		options                 JSONB,
		correct_answer          JSONB,
		score                   INTEGER NOT NULL DEFAULT 0,
		position                INTEGER NOT NULL,
		created_at              TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at              TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE (assessment_template_id, position)
	)`,

	// 3. Test cases (for code questions)
	`CREATE TABLE IF NOT EXISTS test_cases (
		id              BIGSERIAL PRIMARY KEY,
		question_id     BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
		input           TEXT NOT NULL,
		expected_output TEXT NOT NULL,
		is_hidden       BOOLEAN NOT NULL DEFAULT FALSE,
		score           INTEGER NOT NULL DEFAULT 0,
		created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at      TIMESTAMP NOT NULL DEFAULT NOW()
	)`,

	// 4. Assessment invites
	`CREATE TABLE IF NOT EXISTS assessment_invites (
		id                      BIGSERIAL PRIMARY KEY,
		assessment_template_id  BIGINT NOT NULL REFERENCES assessment_templates(id) ON DELETE CASCADE,
		candidate_email         TEXT NOT NULL,
		invite_token            TEXT NOT NULL UNIQUE,
		expires_at              TIMESTAMP NOT NULL,
		status                  TEXT NOT NULL DEFAULT 'pending'
		                        CHECK (status IN ('pending', 'started', 'completed', 'expired', 'cancelled')),
		created_at              TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at              TIMESTAMP NOT NULL DEFAULT NOW()
	)`,

	// 5. Assessment attempts
	`CREATE TABLE IF NOT EXISTS assessment_attempts (
		id            BIGSERIAL PRIMARY KEY,
		invite_id     BIGINT NOT NULL REFERENCES assessment_invites(id) ON DELETE CASCADE,
		started_at    TIMESTAMP,
		completed_at  TIMESTAMP,
		status        TEXT NOT NULL DEFAULT 'in_progress'
		              CHECK (status IN ('in_progress', 'submitted', 'graded', 'abandoned')),
		total_score   INTEGER,
		created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at    TIMESTAMP NOT NULL DEFAULT NOW()
	)`,

	// 6. Answers (non-code responses)
	`CREATE TABLE IF NOT EXISTS answers (
		id             BIGSERIAL PRIMARY KEY,
		attempt_id     BIGINT NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
		question_id    BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
		response       JSONB,
		score_awarded  INTEGER,
		created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at     TIMESTAMP NOT NULL DEFAULT NOW(),
		UNIQUE (attempt_id, question_id)
	)`,

	// 7. Submissions (code responses)
	`CREATE TABLE IF NOT EXISTS submissions (
		id             BIGSERIAL PRIMARY KEY,
		attempt_id     BIGINT NOT NULL REFERENCES assessment_attempts(id) ON DELETE CASCADE,
		question_id    BIGINT NOT NULL REFERENCES questions(id) ON DELETE CASCADE,
		language       TEXT NOT NULL,
		source_code    TEXT NOT NULL,
		status         TEXT NOT NULL DEFAULT 'queued'
		               CHECK (status IN ('queued', 'running', 'passed', 'failed', 'runtime_error', 'compile_error')),
		stdout         TEXT,
		stderr         TEXT,
		execution_ms   INTEGER,
		memory_bytes   BIGINT,
		score_awarded  INTEGER,
		created_at     TIMESTAMP NOT NULL DEFAULT NOW(),
		updated_at     TIMESTAMP NOT NULL DEFAULT NOW()
	)`,

	// 8. Submission results (per test case)
	`CREATE TABLE IF NOT EXISTS submission_results (
		id             BIGSERIAL PRIMARY KEY,
		submission_id  BIGINT NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
		test_case_id   BIGINT NOT NULL REFERENCES test_cases(id) ON DELETE CASCADE,
		passed         BOOLEAN NOT NULL DEFAULT FALSE,
		actual_output  TEXT,
		execution_ms   INTEGER,
		created_at     TIMESTAMP NOT NULL DEFAULT NOW()
	)`,

	// Indexes on foreign keys
	`CREATE INDEX IF NOT EXISTS idx_questions_template ON questions(assessment_template_id)`,
	`CREATE INDEX IF NOT EXISTS idx_test_cases_question ON test_cases(question_id)`,
	`CREATE INDEX IF NOT EXISTS idx_invites_template ON assessment_invites(assessment_template_id)`,
	`CREATE INDEX IF NOT EXISTS idx_attempts_invite ON assessment_attempts(invite_id)`,
	`CREATE INDEX IF NOT EXISTS idx_answers_attempt ON answers(attempt_id)`,
	`CREATE INDEX IF NOT EXISTS idx_answers_question ON answers(question_id)`,
	`CREATE INDEX IF NOT EXISTS idx_submissions_attempt ON submissions(attempt_id)`,
	`CREATE INDEX IF NOT EXISTS idx_submissions_question ON submissions(question_id)`,
	`CREATE INDEX IF NOT EXISTS idx_submission_results_submission ON submission_results(submission_id)`,
	`CREATE INDEX IF NOT EXISTS idx_submission_results_test_case ON submission_results(test_case_id)`,
}

// InitSchema creates all tables and indexes for the assessment platform.
func InitSchema(db *sql.DB) error {
	for i, stmt := range schemas {
		if _, err := db.Exec(stmt); err != nil {
			return fmt.Errorf("failed to execute schema statement %d: %w", i, err)
		}
	}
	return nil
}
