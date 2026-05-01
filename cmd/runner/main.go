package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"

	"CodeSCE/internal/execution"
	"CodeSCE/redisclient"
)

func main() {
	// connect to postgres
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}
	database, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to create postgres pool: %v", err)
	}
	defer database.Close()
	if err := database.Ping(context.Background()); err != nil {
		log.Fatalf("failed to ping database: %v", err)
	}

	// connect to redis
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is required")
	}
	if err := redisclient.Init(); err != nil {
		log.Fatalf("redis fail connect: %v", err)
	}
	defer redisclient.Close()
	rdb := redisclient.GetClient()

	workerCount := 4
	if raw := os.Getenv("WORKER_COUNT"); raw != "" {
		v, err := strconv.Atoi(raw)
		if err != nil || v <= 0 {
			log.Fatalf("invalid WORKER_COUNT: %q", raw)
		}
		workerCount = v
	}

	executionTimeout := 10 * time.Second
	if raw := os.Getenv("EXECUTION_TIMEOUT"); raw != "" {
		d, err := time.ParseDuration(raw)
		if err != nil {
			log.Fatalf("invalid EXECUTION_TIMEOUT: %q", raw)
		}
		executionTimeout = d
	}

	sandboxMemoryLimit := os.Getenv("SANDBOX_MEMORY_LIMIT")
	if sandboxMemoryLimit == "" {
		sandboxMemoryLimit = "256m"
	}

	// creates context, cancels on kill signals
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	// start worker pool
	for i := 0; i < workerCount; i++ {
		worker := &execution.Worker{
			ID:                 i + 1,
			Redis:              rdb,
			DB:                 database,
			ExecutionTimeout:   executionTimeout,
			SandboxMemoryLimit: sandboxMemoryLimit,
		}

		wg.Add(1)
		go func(w *execution.Worker) {
			defer wg.Done()
			w.Run(ctx)
		}(worker)
	}

	// block until shutdown signal
	<-ctx.Done()
	log.Println("runner shutting down, waiting for in-flight workers to finish...")

	// wait for workers to drain and exit
	wg.Wait()

	log.Println("runner shutdown complete")

}
