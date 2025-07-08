package testing

import (
	"context"
	"sync"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	postgresDriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global container management for shared test containers
var (
	globalTestContainers = make(map[string]*PostgresTestContainer)
	containerMutex       sync.Mutex
)

// PostgresTestContainer provides an optimized setup for PostgreSQL testcontainers with transaction support
// This can be used across different services for consistent database testing
type PostgresTestContainer struct {
	DB        *gorm.DB
	Container *postgres.PostgresContainer
	Context   context.Context
}

// TestContainerConfig holds configuration for the test container
type TestContainerConfig struct {
	Image    string          // PostgreSQL image to use (default: "postgres:15-alpine")
	Database string          // Database name (default: "testdb")
	Username string          // Username (default: "testuser")
	Password string          // Password (default: "testpass")
	LogLevel logger.LogLevel // GORM log level (default: logger.Silent)
}

// DefaultTestContainerConfig returns a default configuration for test containers
func DefaultTestContainerConfig() TestContainerConfig {
	return TestContainerConfig{
		Image:    "postgres:15-alpine",
		Database: "testdb",
		Username: "testuser",
		Password: "testpass",
		LogLevel: logger.Silent,
	}
}

// NewPostgresTestContainer creates a new PostgreSQL testcontainer with transaction support
// models: slice of GORM models to auto-migrate
func NewPostgresTestContainer(t *testing.T, models []interface{}, config ...TestContainerConfig) *PostgresTestContainer {
	t.Helper()

	ctx := context.Background()

	// Use default config if none provided
	cfg := DefaultTestContainerConfig()
	if len(config) > 0 {
		cfg = config[0]
	}

	// Start PostgreSQL container
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage(cfg.Image),
		postgres.WithDatabase(cfg.Database),
		postgres.WithUsername(cfg.Username),
		postgres.WithPassword(cfg.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2)),
	)
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	// Get connection details
	connectionString, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("Failed to get connection string: %v", err)
	}

	// Connect to database
	db, err := gorm.Open(postgresDriver.Open(connectionString), &gorm.Config{
		Logger: logger.Default.LogMode(cfg.LogLevel),
	})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate provided models
	if len(models) > 0 {
		err = db.AutoMigrate(models...)
		if err != nil {
			t.Fatalf("Failed to migrate database: %v", err)
		}
	}

	return &PostgresTestContainer{
		DB:        db,
		Container: postgresContainer,
		Context:   ctx,
	}
}

// Cleanup terminates the container and cleans up resources
func (tc *PostgresTestContainer) Cleanup(t *testing.T) {
	t.Helper()

	if tc.Container != nil {
		err := tc.Container.Terminate(tc.Context)
		if err != nil {
			t.Errorf("Failed to terminate container: %v", err)
		}
	}
}

// WithTransaction executes a function within a transaction and automatically rolls back
// This is the primary method for test isolation - all tests should use this
func (tc *PostgresTestContainer) WithTransaction(fn func(db *gorm.DB)) {
	tx := tc.DB.Begin()
	defer tx.Rollback() // Always rollback to ensure test isolation
	fn(tx)
}

// GetSharedTestContainer returns a shared test container for the given package/test suite
// This ensures only one container per test package, improving performance and resource usage
func GetSharedTestContainer(t *testing.T, packageName string, models []interface{}, config ...TestContainerConfig) *PostgresTestContainer {
	t.Helper()

	containerMutex.Lock()
	defer containerMutex.Unlock()

	// Check if container already exists for this package
	if container, exists := globalTestContainers[packageName]; exists {
		return container
	}

	// Create new container for this package
	container := NewPostgresTestContainer(t, models, config...)
	globalTestContainers[packageName] = container

	// Register cleanup to be called when the test finishes
	t.Cleanup(func() {
		containerMutex.Lock()
		defer containerMutex.Unlock()

		if container != nil {
			container.Cleanup(t)
			delete(globalTestContainers, packageName)
		}
	})

	return container
}
