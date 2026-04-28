package execution

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type Result struct {
	Stdout      string `json:"stdout"`
	Stderr      string `json:"stderr"`
	Code        int    `json:"code"`
	ExecutionMS int64  `json:"execution_ms"`
}

func Run(ctx context.Context, language string, sourceCode string, stdin string) (*Result, error) {
	// create a temp dir
	tempDir, err := os.MkdirTemp("", "codehere")
	if err != nil {
		return nil, err
	}
	// makes sure temp directory and contents are deleted when run is done
	defer os.RemoveAll(tempDir)

	filePath := filepath.Join(tempDir, "main.py")

	// write source code to file named main.py
	err = os.WriteFile(filePath, []byte(sourceCode), 0644)
	if err != nil {
		return nil, err
	}

	// get executiontimeout variable to time
	timeoutStr := os.Getenv("EXECUTION_TIMEOUT")
	if timeoutStr == "" {
		timeoutStr = "10s"
	}
	timeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return nil, err
	}

	// enforce hard timeout
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	// get memory limit from env variable
	memoryLimit := os.Getenv("SANDBOX_MEMORY_LIMIT")
	if memoryLimit == "" {
		memoryLimit = "256m"
	}

	args := []string{
		"run", "--rm", "-i",
		"--network=none",
		"--read-only",
		"--memory=" + memoryLimit,
		"--memory-swap=" + memoryLimit,
		"--cpus=0.5",
		"--pids-limit=64",
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=64m",
		"-v", fmt.Sprintf("%s:/code:ro", tempDir),
		"codesce-sandbox-python:latest",
	}

	// run da command and pipe stdin to container i think
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdin = strings.NewReader(stdin)

	//capture studout/stderr from the container
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err = cmd.Run()
	duration := time.Since(start).Milliseconds()

	// do you know you have 10 seconds....
	// 10? 10, 10, yes....
	if ctx.Err() == context.DeadlineExceeded {
		return &Result{
			Stdout:      stdout.String(),
			Stderr:      "execution timed out (" + timeoutStr + " limit)",
			Code:        -1,
			ExecutionMS: duration,
		}, nil
	}

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, err
		}
	}
	// WHEN THE GO FILE FINALLY TURNS GREEN 🥹🥹🥹
	return &Result{
		Stdout:      stdout.String(),
		Stderr:      stderr.String(),
		Code:        exitCode,
		ExecutionMS: duration,
	}, nil
}
