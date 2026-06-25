package integration_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/config"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/domain"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/http"
	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/storage"
	"github.com/gigavard/AIAdoption/ExTodoGolang/pkg/logger"
)

func setupTestServer(t *testing.T) *http.Server {
	log := logger.New()
	cfg := &config.Config{HTTPAddr: ":8080"}

	repo, err := storage.NewSQLiteRepository(":memory:")
	if err != nil {
		t.Fatalf("failed to create repository: %v", err)
	}

	return http.NewServer(log, cfg, repo)
}

func TestCRUDWorkflow(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Create a todo
	input := domain.CreateTodoInput{
		Title:   "Buy groceries",
		Content: "Milk, eggs, bread",
	}
	body, _ := json.Marshal(input)
	req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var created domain.Todo
	json.NewDecoder(w.Body).Decode(&created)
	id := created.ID

	// Read the todo
	req = httptest.NewRequest("GET", "/todos/"+string(rune(id)), nil)
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var retrieved domain.Todo
	json.NewDecoder(w.Body).Decode(&retrieved)
	if retrieved.Title != "Buy groceries" {
		t.Fatalf("expected 'Buy groceries', got %s", retrieved.Title)
	}

	// Update the todo
	newTitle := "Buy groceries and cook"
	updateInput := domain.UpdateTodoInput{Title: &newTitle}
	body, _ = json.Marshal(updateInput)
	req = httptest.NewRequest("PUT", "/todos/"+string(rune(id)), bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var updated domain.Todo
	json.NewDecoder(w.Body).Decode(&updated)
	if updated.Title != "Buy groceries and cook" {
		t.Fatalf("expected updated title")
	}

	// Complete the todo
	req = httptest.NewRequest("PATCH", "/todos/"+string(rune(id))+"/complete", nil)
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var completed domain.Todo
	json.NewDecoder(w.Body).Decode(&completed)
	if completed.Status != domain.StatusCompleted {
		t.Fatalf("expected completed status")
	}

	// Delete the todo
	req = httptest.NewRequest("DELETE", "/todos/"+string(rune(id)), nil)
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected 204 after delete")
	}

	// Verify deletion
	req = httptest.NewRequest("GET", "/todos/"+string(rune(id)), nil)
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for deleted todo")
	}
}

func TestMultipleTodos(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Create multiple todos
	titles := []string{"Task 1", "Task 2", "Task 3"}
	ids := []int64{}

	for _, title := range titles {
		input := domain.CreateTodoInput{Title: title}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.Handler.ServeHTTP(w, req)

		var todo domain.Todo
		json.NewDecoder(w.Body).Decode(&todo)
		ids = append(ids, todo.ID)
	}

	// List all
	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var todos []domain.Todo
	json.NewDecoder(w.Body).Decode(&todos)

	if len(todos) != 3 {
		t.Fatalf("expected 3 todos, got %d", len(todos))
	}

	// Complete one
	newStatus := domain.StatusCompleted
	updateInput := domain.UpdateTodoInput{Status: &newStatus}
	body, _ := json.Marshal(updateInput)
	req = httptest.NewRequest("PUT", "/todos/"+string(rune(ids[0])), bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	// Delete one
	req = httptest.NewRequest("DELETE", "/todos/"+string(rune(ids[1])), nil)
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	// List remaining
	req = httptest.NewRequest("GET", "/todos", nil)
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	json.NewDecoder(w.Body).Decode(&todos)
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos after delete, got %d", len(todos))
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		body           interface{}
		expectedStatus int
		expectedType   string
	}{
		{
			name:           "Missing title on create",
			method:         "POST",
			path:           "/todos",
			body:           domain.CreateTodoInput{Content: "content"},
			expectedStatus: http.StatusBadRequest,
			expectedType:   "application/problem+json",
		},
		{
			name:           "Get non-existent todo",
			method:         "GET",
			path:           "/todos/999",
			body:           nil,
			expectedStatus: http.StatusNotFound,
			expectedType:   "application/problem+json",
		},
		{
			name:           "Invalid ID format",
			method:         "GET",
			path:           "/todos/invalid",
			body:           nil,
			expectedStatus: http.StatusBadRequest,
			expectedType:   "application/problem+json",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := setupTestServer(t)
			defer server.Close()

			var body []byte
			if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, tt.path, bytes.NewReader(body))
			w := httptest.NewRecorder()
			server.Handler.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			if ct := w.Header().Get("Content-Type"); ct != tt.expectedType {
				t.Errorf("expected Content-Type %s, got %s", tt.expectedType, ct)
			}
		})
	}
}

func TestTimestampsAndOrdering(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	// Create todos with small delay
	ids := []int64{}
	for i := 0; i < 3; i++ {
		input := domain.CreateTodoInput{Title: "Todo " + string(rune('1'+i))}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/todos", bytes.NewReader(body))
		w := httptest.NewRecorder()
		server.Handler.ServeHTTP(w, req)

		var todo domain.Todo
		json.NewDecoder(w.Body).Decode(&todo)
		ids = append(ids, todo.ID)

		if !todo.CreatedAt.Equal(todo.UpdatedAt) {
			t.Errorf("expected CreatedAt == UpdatedAt on creation")
		}

		time.Sleep(10 * time.Millisecond)
	}

	// List and verify ordering (newest first)
	req := httptest.NewRequest("GET", "/todos", nil)
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var todos []domain.Todo
	json.NewDecoder(w.Body).Decode(&todos)

	if len(todos) != 3 {
		t.Fatalf("expected 3 todos")
	}

	// Should be in reverse order (newest first)
	if todos[0].ID != ids[2] || todos[1].ID != ids[1] || todos[2].ID != ids[0] {
		t.Errorf("todos not in expected order (newest first)")
	}

	// Update a todo and verify UpdatedAt changes
	time.Sleep(20 * time.Millisecond)
	newTitle := "Updated"
	updateInput := domain.UpdateTodoInput{Title: &newTitle}
	body, _ := json.Marshal(updateInput)
	req = httptest.NewRequest("PUT", "/todos/"+string(rune(ids[0])), bytes.NewReader(body))
	w = httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	var updated domain.Todo
	json.NewDecoder(w.Body).Decode(&updated)

	if updated.UpdatedAt.Before(updated.CreatedAt) {
		t.Errorf("expected UpdatedAt >= CreatedAt")
	}
}

func TestOpenAPIEndpoint(t *testing.T) {
	server := setupTestServer(t)
	defer server.Close()

	req := httptest.NewRequest("GET", "/openapi.json", nil)
	w := httptest.NewRecorder()
	server.Handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var spec map[string]interface{}
	err := json.NewDecoder(w.Body).Decode(&spec)
	if err != nil {
		t.Fatalf("failed to decode spec: %v", err)
	}

	if spec["openapi"] != "3.0.0" {
		t.Errorf("expected openapi 3.0.0")
	}

	if info, ok := spec["info"].(map[string]interface{}); !ok || info["title"] != "Todo API" {
		t.Errorf("expected Title in info")
	}

	if paths, ok := spec["paths"].(map[string]interface{}); !ok || len(paths) == 0 {
		t.Errorf("expected paths in spec")
	}
}
