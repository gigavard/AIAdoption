package http

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/gigavard/AIAdoption/ExTodoGolang/internal/domain"
	"github.com/gigavard/AIAdoption/ExTodoGolang/pkg/errors"
)

type Handler struct {
	log  *slog.Logger
	repo domain.Repository
}

func NewHandler(log *slog.Logger, repo domain.Repository) *Handler {
	return &Handler{log: log, repo: repo}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /todos", h.ListTodos)
	mux.HandleFunc("POST /todos", h.CreateTodo)
	mux.HandleFunc("GET /todos/{id}", h.GetTodo)
	mux.HandleFunc("PUT /todos/{id}", h.UpdateTodo)
	mux.HandleFunc("PATCH /todos/{id}/complete", h.CompleteTodo)
	mux.HandleFunc("DELETE /todos/{id}", h.DeleteTodo)
	mux.HandleFunc("GET /openapi.json", h.OpenAPISpec)
}

func (h *Handler) ListTodos(w http.ResponseWriter, r *http.Request) {
	todos, err := h.repo.List()
	if err != nil {
		h.log.Error("failed to list todos", "err", err)
		errors.Respond(w, errors.NewError("list_failed", "Failed to list todos", http.StatusInternalServerError))
		return
	}

	if todos == nil {
		todos = []domain.Todo{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

func (h *Handler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var input domain.CreateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request", "err", err)
		errors.Respond(w, errors.NewError("invalid_input", "Invalid request body", http.StatusBadRequest))
		return
	}

	if strings.TrimSpace(input.Title) == "" {
		errors.Respond(w, errors.NewError("validation_error", "Title is required", http.StatusBadRequest))
		return
	}

	todo := &domain.Todo{
		Title:   input.Title,
		Content: input.Content,
		Status:  domain.StatusPending,
	}

	id, err := h.repo.Create(todo)
	if err != nil {
		h.log.Error("failed to create todo", "err", err)
		errors.Respond(w, errors.NewError("create_failed", "Failed to create todo", http.StatusInternalServerError))
		return
	}

	todo.ID = id
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) GetTodo(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		errors.Respond(w, errors.NewError("invalid_id", "Invalid todo ID", http.StatusBadRequest))
		return
	}

	todo, err := h.repo.GetByID(id)
	if err != nil {
		if err.Error() == "todo not found" {
			errors.Respond(w, errors.NewError("not_found", "Todo not found", http.StatusNotFound))
			return
		}
		h.log.Error("failed to get todo", "err", err)
		errors.Respond(w, errors.NewError("get_failed", "Failed to get todo", http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		errors.Respond(w, errors.NewError("invalid_id", "Invalid todo ID", http.StatusBadRequest))
		return
	}

	var input domain.UpdateTodoInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.log.Error("failed to decode request", "err", err)
		errors.Respond(w, errors.NewError("invalid_input", "Invalid request body", http.StatusBadRequest))
		return
	}

	todo, err := h.repo.GetByID(id)
	if err != nil {
		errors.Respond(w, errors.NewError("not_found", "Todo not found", http.StatusNotFound))
		return
	}

	if input.Title != nil {
		if strings.TrimSpace(*input.Title) == "" {
			errors.Respond(w, errors.NewError("validation_error", "Title cannot be empty", http.StatusBadRequest))
			return
		}
		todo.Title = *input.Title
	}

	if input.Content != nil {
		todo.Content = *input.Content
	}

	if input.Status != nil {
		todo.Status = *input.Status
	}

	if err := h.repo.Update(id, todo); err != nil {
		h.log.Error("failed to update todo", "err", err)
		errors.Respond(w, errors.NewError("update_failed", "Failed to update todo", http.StatusInternalServerError))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) CompleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		errors.Respond(w, errors.NewError("invalid_id", "Invalid todo ID", http.StatusBadRequest))
		return
	}

	if err := h.repo.Complete(id); err != nil {
		errors.Respond(w, errors.NewError("not_found", "Todo not found", http.StatusNotFound))
		return
	}

	todo, _ := h.repo.GetByID(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

func (h *Handler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := h.parseID(r)
	if err != nil {
		errors.Respond(w, errors.NewError("invalid_id", "Invalid todo ID", http.StatusBadRequest))
		return
	}

	if err := h.repo.Delete(id); err != nil {
		errors.Respond(w, errors.NewError("not_found", "Todo not found", http.StatusNotFound))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) OpenAPISpec(w http.ResponseWriter, r *http.Request) {
	spec := map[string]interface{}{
		"openapi": "3.0.0",
		"info": map[string]string{
			"title":   "Todo API",
			"version": "1.0.0",
		},
		"servers": []map[string]string{
			{"url": "http://localhost:8080"},
		},
		"paths": map[string]interface{}{
			"/todos": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "List all todos",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "List of todos",
						},
					},
				},
				"post": map[string]interface{}{
					"summary": "Create a new todo",
					"responses": map[string]interface{}{
						"201": map[string]interface{}{
							"description": "Created todo",
						},
					},
				},
			},
			"/todos/{id}": map[string]interface{}{
				"get": map[string]interface{}{
					"summary": "Get todo by ID",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Todo details",
						},
					},
				},
				"put": map[string]interface{}{
					"summary": "Update todo",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Updated todo",
						},
					},
				},
				"delete": map[string]interface{}{
					"summary": "Delete todo",
					"responses": map[string]interface{}{
						"204": map[string]interface{}{
							"description": "Todo deleted",
						},
					},
				},
			},
			"/todos/{id}/complete": map[string]interface{}{
				"patch": map[string]interface{}{
					"summary": "Mark todo as completed",
					"responses": map[string]interface{}{
						"200": map[string]interface{}{
							"description": "Completed todo",
						},
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(spec)
}

func (h *Handler) parseID(r *http.Request) (int64, error) {
	idStr := r.PathValue("id")
	return strconv.ParseInt(idStr, 10, 64)
}
