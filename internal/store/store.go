package store

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

// Session represents a debugging session for a problem.
type Session struct {
	ID               int        `db:"id"`
	ProblemName      string     `db:"problem_name"`
	StartedAt        time.Time  `db:"started_at"`
	CompletedAt      *time.Time `db:"completed_at"`
	TimeTakenSeconds *int       `db:"time_taken_seconds"`
	Score            *int       `db:"score"`
	TestsPassed      int        `db:"tests_passed"`
	TestsFailed      int        `db:"tests_failed"`
	FilesChanged     int        `db:"files_changed"`
	HintsUsed        int        `db:"hints_used"`
	Status           string     `db:"status"` // "started", "completed", "failed"
}

// Store wraps the database connection.
type Store struct {
	db *sqlx.DB
}

const schema = `
CREATE TABLE IF NOT EXISTS sessions (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	problem_name TEXT NOT NULL,
	started_at DATETIME NOT NULL,
	completed_at DATETIME,
	time_taken_seconds INTEGER,
	score INTEGER,
	tests_passed INTEGER DEFAULT 0,
	tests_failed INTEGER DEFAULT 0,
	files_changed INTEGER DEFAULT 0,
	hints_used INTEGER DEFAULT 0,
	status TEXT NOT NULL DEFAULT 'started'
);
`

// DefaultDBPath returns the default database path: ./dblab-workspace/dblab.db
func DefaultDBPath() (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot get current directory: %w", err)
	}
	dir := filepath.Join(cwd, "dblab-workspace")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("cannot create dblab-workspace directory: %w", err)
	}
	return filepath.Join(dir, "dblab.db"), nil
}

// NewStore opens the SQLite database and runs migrations.
func NewStore(dbPath string) (*Store, error) {
	db, err := sqlx.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if _, err := db.Exec(schema); err != nil {
		return nil, fmt.Errorf("failed to run schema migration: %w", err)
	}

	return &Store{db: db}, nil
}

// Close closes the database connection.
func (s *Store) Close() error {
	return s.db.Close()
}

// CreateSession creates a new debugging session.
func (s *Store) CreateSession(problemName string) (*Session, error) {
	now := time.Now()
	result, err := s.db.Exec(
		`INSERT INTO sessions (problem_name, started_at, status) VALUES (?, ?, 'started')`,
		problemName, now,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}

	id, _ := result.LastInsertId()
	return &Session{
		ID:          int(id),
		ProblemName: problemName,
		StartedAt:   now,
		Status:      "started",
	}, nil
}

// GetSession retrieves a session by ID.
func (s *Store) GetSession(id int) (*Session, error) {
	var session Session
	err := s.db.Get(&session, `SELECT * FROM sessions WHERE id = ?`, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return &session, nil
}

// GetActiveSessionByProblem retrieves the most recent active (started) session for a problem.
func (s *Store) GetActiveSessionByProblem(problemName string) (*Session, error) {
	var session Session
	err := s.db.Get(&session,
		`SELECT * FROM sessions WHERE problem_name = ? AND status = 'started' ORDER BY started_at DESC LIMIT 1`,
		problemName,
	)
	if err != nil {
		return nil, fmt.Errorf("no active session found for %q: %w", problemName, err)
	}
	return &session, nil
}

// CompleteSession marks a session as completed with results.
func (s *Store) CompleteSession(id int, timeTaken int, score int, testsPassed int, testsFailed int, filesChanged int) error {
	now := time.Now()
	_, err := s.db.Exec(
		`UPDATE sessions SET completed_at = ?, time_taken_seconds = ?, score = ?, tests_passed = ?, tests_failed = ?, files_changed = ?, status = 'completed' WHERE id = ?`,
		now, timeTaken, score, testsPassed, testsFailed, filesChanged, id,
	)
	if err != nil {
		return fmt.Errorf("failed to complete session: %w", err)
	}
	return nil
}

// FailSession marks a session as failed.
func (s *Store) FailSession(id int, testsPassed int, testsFailed int) error {
	_, err := s.db.Exec(
		`UPDATE sessions SET tests_passed = ?, tests_failed = ?, status = 'failed' WHERE id = ?`,
		testsPassed, testsFailed, id,
	)
	if err != nil {
		return fmt.Errorf("failed to fail session: %w", err)
	}
	return nil
}

// ListSessions retrieves all sessions.
func (s *Store) ListSessions() ([]Session, error) {
	var sessions []Session
	err := s.db.Select(&sessions, `SELECT * FROM sessions ORDER BY started_at DESC`)
	if err != nil {
		return nil, fmt.Errorf("failed to list sessions: %w", err)
	}
	return sessions, nil
}

// GetSolvedProblems returns a set of problem names that have been solved.
func (s *Store) GetSolvedProblems() (map[string]bool, error) {
	var names []string
	err := s.db.Select(&names,
		`SELECT DISTINCT problem_name FROM sessions WHERE status = 'completed'`,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get solved problems: %w", err)
	}
	solved := make(map[string]bool)
	for _, name := range names {
		solved[name] = true
	}
	return solved, nil
}
