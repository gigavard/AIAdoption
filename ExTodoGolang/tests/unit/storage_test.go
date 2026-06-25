package storage_test

import (
	"os"
	"testing"
	"time"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/domain"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/storage"
)

func setupTestDB(t *testing.T) *storage.SQLiteRepository {
	dbFile := ":memory:"
	repo, err := storage.NewSQLiteRepository(dbFile)
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}
	return repo
}

func TestCreateTodo(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	todo := &domain.Todo{
		Title:   "Test Todo",
		Content: "Test Content",
		Status:  domain.StatusPending,
	}

	id, err := repo.Create(todo)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	if id == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestGetByID(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	todo := &domain.Todo{
		Title:   "Test Todo",
		Content: "Test Content",
		Status:  domain.StatusPending,
	}

	id, err := repo.Create(todo)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	retrieved, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get todo: %v", err)
	}

	if retrieved.Title != "Test Todo" {
		t.Errorf("expected title 'Test Todo', got %s", retrieved.Title)
	}
	if retrieved.Status != domain.StatusPending {
		t.Errorf("expected status 'pending', got %s", retrieved.Status)
	}
}

func TestGetByIDNotFound(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	_, err := repo.GetByID(999)
	if err != storage.ErrNotFound {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestList(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	// Create multiple todos
	for i := 0; i < 3; i++ {
		todo := &domain.Todo{
			Title:   "Todo " + string(rune('1'+i)),
			Content: "Content",
			Status:  domain.StatusPending,
		}
		repo.Create(todo)
	}

	todos, err := repo.List()
	if err != nil {
		t.Fatalf("failed to list todos: %v", err)
	}

	if len(todos) != 3 {
		t.Errorf("expected 3 todos, got %d", len(todos))
	}
}

func TestUpdate(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	todo := &domain.Todo{
		Title:   "Original Title",
		Content: "Original Content",
		Status:  domain.StatusPending,
	}

	id, err := repo.Create(todo)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	// Update the todo
	updated := &domain.Todo{
		Title:   "Updated Title",
		Content: "Updated Content",
		Status:  domain.StatusPending,
	}

	err = repo.Update(id, updated)
	if err != nil {
		t.Fatalf("failed to update todo: %v", err)
	}

	retrieved, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get updated todo: %v", err)
	}

	if retrieved.Title != "Updated Title" {
		t.Errorf("expected title 'Updated Title', got %s", retrieved.Title)
	}
}

func TestComplete(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	todo := &domain.Todo{
		Title:   "Test Todo",
		Content: "Content",
		Status:  domain.StatusPending,
	}

	id, err := repo.Create(todo)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	err = repo.Complete(id)
	if err != nil {
		t.Fatalf("failed to complete todo: %v", err)
	}

	retrieved, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get todo: %v", err)
	}

	if retrieved.Status != domain.StatusCompleted {
		t.Errorf("expected status 'completed', got %s", retrieved.Status)
	}
}

func TestDelete(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	todo := &domain.Todo{
		Title:   "Test Todo",
		Content: "Content",
		Status:  domain.StatusPending,
	}

	id, err := repo.Create(todo)
	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	err = repo.Delete(id)
	if err != nil {
		t.Fatalf("failed to delete todo: %v", err)
	}

	_, err = repo.GetByID(id)
	if err != storage.ErrNotFound {
		t.Errorf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestTimestamps(t *testing.T) {
	repo := setupTestDB(t)
	defer repo.Close()

	todo := &domain.Todo{
		Title:   "Test Todo",
		Content: "Content",
		Status:  domain.StatusPending,
	}

	beforeCreate := time.Now()
	id, err := repo.Create(todo)
	afterCreate := time.Now()

	if err != nil {
		t.Fatalf("failed to create todo: %v", err)
	}

	retrieved, err := repo.GetByID(id)
	if err != nil {
		t.Fatalf("failed to get todo: %v", err)
	}

	if retrieved.CreatedAt.Before(beforeCreate) || retrieved.CreatedAt.After(afterCreate) {
		t.Errorf("CreatedAt timestamp outside expected range")
	}

	if retrieved.UpdatedAt.Before(beforeCreate) || retrieved.UpdatedAt.After(afterCreate) {
		t.Errorf("UpdatedAt timestamp outside expected range")
	}
}
