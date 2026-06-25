package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/config"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/domain"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/http"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/storage"
	"github.com/gigavard/AIAdoption/ExTodoGolang/pkg/logger"
)

func setupServer(t *testing.T) *http.Server {
	log := logger.New()
	cfg := &config.Config{HTTPAddr: ":8080"}

	repo, err := storage.NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	return http.NewServer(log, cfg, repo)
}

func TestCreateTodoHandler(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	input := domain.CreateTodoInput{
		Title:   "Test Todo",
		Content: "Test Content",
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("expected status 201, got %d", w.Code)
	}

	var todo domain.Todo
	json.NewDecoder(w.Body).Decode(&todo)

	if todo.Title != "Test Todo" {
		t.Errorf("expected title 'Test Todo', got %s", todo.Title)
	}
	if todo.ID == 0 {
		t.Fatal("expected non-zero ID")
	}
}

func TestCreateTodoMissingTitle(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	input := domain.CreateTodoInput{
		Content: "Test Content",
	}

	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("expected Content-Type application/problem+json, got %s", ct)
	}
}

func TestListTodos(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	// Create a todo first
	input := domain.CreateTodoInput{Title: "Test"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	// List todos
	req = httptest.NewRequest("GET", "/todos", nil)
	w = httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var todos []domain.Todo
	json.NewDecoder(w.Body).Decode(&todos)

	if len(todos) != 1 {
		t.Errorf("expected 1 todo, got %d", len(todos))
	}
}

func TestGetTodo(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	// Create a todo
	input := domain.CreateTodoInput{Title: "Test"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var created domain.Todo
	json.NewDecoder(w.Body).Decode(&created)

	// Get the todo
	req = httptest.NewRequest("GET", "/todos/"+string(rune(created.ID)), nil)
	w = httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var retrieved domain.Todo
	json.NewDecoder(w.Body).Decode(&retrieved)

	if retrieved.Title != "Test" {
		t.Errorf("expected title 'Test', got %s", retrieved.Title)
	}
}

func TestGetTodoNotFound(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	req := httptest.NewRequest("GET", "/todos/999", nil)
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", w.Code)
	}

	if ct := w.Header().Get("Content-Type"); ct != "application/problem+json" {
		t.Errorf("expected Content-Type application/problem+json, got %s", ct)
	}
}

func TestUpdateTodo(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	// Create a todo
	input := domain.CreateTodoInput{Title: "Original"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var created domain.Todo
	json.NewDecoder(w.Body).Decode(&created)

	// Update the todo
	newTitle := "Updated"
	updateInput := domain.UpdateTodoInput{Title: &newTitle}
	body, _ = json.Marshal(updateInput)

	req = httptest.NewRequest("PUT", "/todos/"+string(rune(created.ID)), bytes.NewReader(body))
	w = httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var updated domain.Todo
	json.NewDecoder(w.Body).Decode(&updated)

	if updated.Title != "Updated" {
		t.Errorf("expected title 'Updated', got %s", updated.Title)
	}
}

func TestCompleteTodo(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	// Create a todo
	input := domain.CreateTodoInput{Title: "Test"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var created domain.Todo
	json.NewDecoder(w.Body).Decode(&created)

	// Complete the todo
	req = httptest.NewRequest("PATCH", "/todos/"+string(rune(created.ID))+"/complete", nil)
	w = httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var completed domain.Todo
	json.NewDecoder(w.Body).Decode(&completed)

	if completed.Status != domain.StatusCompleted {
		t.Errorf("expected status 'completed', got %s", completed.Status)
	}
}

func TestDeleteTodo(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	// Create a todo
	input := domain.CreateTodoInput{Title: "Test"}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var created domain.Todo
	json.NewDecoder(w.Body).Decode(&created)

	// Delete the todo
	req = httptest.NewRequest("DELETE", "/todos/"+string(rune(created.ID)), nil)
	w = httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Errorf("expected status 204, got %d", w.Code)
	}
}

func TestOpenAPISpec(t *testing.T) {
	server := setupServer(t)
	defer server.Close()

	req := httptest.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()

	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", w.Code)
	}

	var spec map[string]interface{}
	json.NewDecoder(w.Body).Decode(&spec)

	if spec["openapi"] != "3.0.0" {
		t.Errorf("expected openapi version 3.0.0")
	}
}
