package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/domain"
	_ "github.com/mattn/go-sqlite3"
)

var ErrNotFound = errors.New("todo not found")

type SQLiteRepository struct {
	db *sql.DB
}

func NewSQLiteRepository(dbPath string) (*SQLiteRepository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Enable foreign keys
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	repo := &SQLiteRepository{db: db}

	// Initialize schema
	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

func (r *SQLiteRepository) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS todos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		content TEXT,
		status TEXT NOT NULL DEFAULT 'pending',
		created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_todos_status ON todos(status);
	CREATE INDEX IF NOT EXISTS idx_todos_created_at ON todos(created_at);
	`

	_, err := r.db.Exec(schema)
	return err
}

func (r *SQLiteRepository) Create(todo *domain.Todo) (int64, error) {
	now := time.Now()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	result, err := r.db.Exec(
		"INSERT INTO todos (title, content, status, created_at, updated_at) VALUES (?, ?, ?, ?, ?)",
		todo.Title, todo.Content, todo.Status, todo.CreatedAt, todo.UpdatedAt,
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert todo: %w", err)
	}

	return result.LastInsertId()
}

func (r *SQLiteRepository) GetByID(id int64) (*domain.Todo, error) {
	todo := &domain.Todo{}

	err := r.db.QueryRow(
		"SELECT id, title, content, status, created_at, updated_at FROM todos WHERE id = ?",
		id,
	).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt)

	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query todo: %w", err)
	}

	return todo, nil
}

func (r *SQLiteRepository) List() ([]domain.Todo, error) {
	rows, err := r.db.Query("SELECT id, title, content, status, created_at, updated_at FROM todos ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to query todos: %w", err)
	}
	defer rows.Close()

	todos := []domain.Todo{}
	for rows.Next() {
		todo := domain.Todo{}
		if err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Status, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating todos: %w", err)
	}

	return todos, nil
}

func (r *SQLiteRepository) Update(id int64, todo *domain.Todo) error {
	todo.UpdatedAt = time.Now()

	result, err := r.db.Exec(
		"UPDATE todos SET title = ?, content = ?, status = ?, updated_at = ? WHERE id = ?",
		todo.Title, todo.Content, todo.Status, todo.UpdatedAt, id,
	)
	if err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SQLiteRepository) Delete(id int64) error {
	result, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete todo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SQLiteRepository) Complete(id int64) error {
	result, err := r.db.Exec(
		"UPDATE todos SET status = ?, updated_at = ? WHERE id = ?",
		domain.StatusCompleted, time.Now(), id,
	)
	if err != nil {
		return fmt.Errorf("failed to complete todo: %w", err)
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}

	if rows == 0 {
		return ErrNotFound
	}

	return nil
}

func (r *SQLiteRepository) Close() error {
	return r.db.Close()
}
