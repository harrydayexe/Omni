package utilities

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"

	"github.com/testcontainers/testcontainers-go/wait"
)

type ErrFailedToStartContainer struct {
	message    string
	underlying error
}

func (e ErrFailedToStartContainer) Error() string {
	return fmt.Sprintf("%s: %v", e.message, e.underlying)
}

// createDBTestContainer creates a MariaDB container for testing
func createDBTestContainer(ctx context.Context) (testcontainers.Container, error) {
	// Request a MariaDB test container
	req := testcontainers.ContainerRequest{
		Image:        "mariadb:latest",
		ExposedPorts: []string{"3306/tcp"},
		Env: map[string]string{
			"MYSQL_ROOT_PASSWORD": "Password1!",
			"MYSQL_DATABASE":      "omni",
			"MYSQL_USER":          "testuser",
			"MYSQL_PASSWORD":      "testpass",
		},
		WaitingFor: wait.ForLog("mariadbd: ready for connections.").WithStartupTimeout(30 * time.Second),
	}

	// Start the container
	mariaDBContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, ErrFailedToStartContainer{"failed to start container", err}
	}

	return mariaDBContainer, nil
}

// waitForDB ensures the database is ready before running migrations
func waitForDB(db *sql.DB) error {
	for i := 0; i < 10; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		time.Sleep(2 * time.Second)
	}
	return fmt.Errorf("database connection timeout")
}

// createDbInstance creates a database connection to the MariaDB container
func createDbInstance(ctx context.Context, container testcontainers.Container) (*sql.DB, error) {
	// Get the database connection details
	host, err := container.Host(ctx)
	if err != nil {
		return nil, ErrFailedToStartContainer{"failed to get container host", err}
	}

	port, err := container.MappedPort(ctx, "3306")
	if err != nil {
		return nil, ErrFailedToStartContainer{"failed to get container port", err}
	}

	dsn := fmt.Sprintf("testuser:testpass@tcp(%s:%s)/omni?multiStatements=true&parseTime=true", host, port.Port())
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, ErrFailedToStartContainer{"failed to open database", err}
	}

	err = waitForDB(db)
	if err != nil {
		return nil, ErrFailedToStartContainer{"database connection timeout", err}
	}

	return db, nil
}

// SetupTestContainer creates a MariaDB container, applies migrations, and seeds data.
// It returns a database connection and a cleanup function to terminate the container.
func SetupTestContainer(migrationsDir, seedSQLFile string) (*sql.DB, func(), error) {
	ctx := context.Background()

	// Create the MariaDB container
	mariaDBContainer, err := createDBTestContainer(ctx)
	if err != nil {
		return nil, nil, err
	}
	db, err := createDbInstance(ctx, mariaDBContainer)
	if err != nil {
		return nil, nil, err
	}

	// Run migrations
	if err := runMigrations(db, migrationsDir); err != nil {
		return nil, nil, err
	}

	// Seed the database
	if err := seedDatabase(db, seedSQLFile); err != nil {
		return nil, nil, err
	}

	// Cleanup function to terminate the container
	cleanup := func() {
		db.Close()
		mariaDBContainer.Terminate(ctx)
	}

	return db, cleanup, nil
}

// runMigrations applies the migrations from the given directory
func runMigrations(db *sql.DB, migrationsDir string) error {
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return ErrFailedToStartContainer{"failed to create migration driver: %w", err}
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+migrationsDir, "mysql", driver)
	if err != nil {
		return ErrFailedToStartContainer{"failed to create migration instance: %w", err}
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return ErrFailedToStartContainer{"migration failed: %w", err}
	}

	return nil
}

// seedDatabase executes an SQL file to populate test data
func seedDatabase(db *sql.DB, seedSQLFile string) error {
	sqlBytes, err := os.ReadFile(filepath.Clean(seedSQLFile))
	if err != nil {
		return ErrFailedToStartContainer{"failed to read seed file: %w", err}
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return ErrFailedToStartContainer{"failed to execute seed SQL: %w", err}
	}

	return nil
}
