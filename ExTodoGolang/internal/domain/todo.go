package domain

import "time"

// Status represents the lifecycle state of a Todo
type Status string

const (
	StatusPending   Status = "pending"
	StatusCompleted Status = "completed"
)

// Todo represents a single todo item
type Todo struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateTodoInput represents input for creating a new todo
type CreateTodoInput struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// UpdateTodoInput represents input for updating an existing todo
type UpdateTodoInput struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
	Status  *Status `json:"status,omitempty"`
}

// Repository defines the todo persistence interface
type Repository interface {
	Create(todo *Todo) (int64, error)
	GetByID(id int64) (*Todo, error)
	List() ([]Todo, error)
	Update(id int64, todo *Todo) error
	Delete(id int64) error
	Complete(id int64) error
}
