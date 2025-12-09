package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/snip/internal/task"
)

var ErrTaskNotFound = errors.New("task not found")

type TaskRepository interface {
	Create(t *task.Task) error
	GetByID(id int) (*task.Task, error)
	GetByProjectID(projectID int, status string) ([]*task.Task, error)
	GetAll(status string) ([]*task.Task, error)
	Update(id int, title, description, status, priority string, dueDate *time.Time) error
	Delete(id int) error
	ToggleComplete(id int) error
	Close() error
}

type taskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) (TaskRepository, error) {
	return &taskRepository{db: db}, nil
}

func (r *taskRepository) Close() error {
	return r.db.Close()
}

func (r *taskRepository) Create(t *task.Task) error {
	query := `
		INSERT INTO tasks (project_id, title, description, status, priority, due_date, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`
	var dueDate interface{}
	if t.DueDate != nil {
		dueDate = t.DueDate
	}

	result, err := r.db.Exec(query, t.ProjectID, t.Title, t.Description, t.Status, t.Priority, dueDate, t.CreatedAt, t.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	t.ID = int(id)
	return nil
}

func (r *taskRepository) GetByID(id int) (*task.Task, error) {
	query := `SELECT id, project_id, title, description, status, priority, due_date, created_at, updated_at FROM tasks WHERE id = ?`
	
	t := &task.Task{}
	var dueDate sql.NullTime
	err := r.db.QueryRow(query, id).Scan(
		&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrTaskNotFound
		}
		return nil, err
	}

	if dueDate.Valid {
		t.DueDate = &dueDate.Time
	}

	return t, nil
}

func (r *taskRepository) GetByProjectID(projectID int, status string) ([]*task.Task, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT id, project_id, title, description, status, priority, due_date, created_at, updated_at 
			FROM tasks WHERE project_id = ? AND status = ? ORDER BY created_at DESC`
		args = []interface{}{projectID, status}
	} else {
		query = `SELECT id, project_id, title, description, status, priority, due_date, created_at, updated_at 
			FROM tasks WHERE project_id = ? ORDER BY created_at DESC`
		args = []interface{}{projectID}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		t := &task.Task{}
		var dueDate sql.NullTime
		err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepository) GetAll(status string) ([]*task.Task, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT id, project_id, title, description, status, priority, due_date, created_at, updated_at 
			FROM tasks WHERE status = ? ORDER BY created_at DESC`
		args = []interface{}{status}
	} else {
		query = `SELECT id, project_id, title, description, status, priority, due_date, created_at, updated_at 
			FROM tasks ORDER BY created_at DESC`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*task.Task
	for rows.Next() {
		t := &task.Task{}
		var dueDate sql.NullTime
		err := rows.Scan(&t.ID, &t.ProjectID, &t.Title, &t.Description, &t.Status, &t.Priority, &dueDate, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if dueDate.Valid {
			t.DueDate = &dueDate.Time
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *taskRepository) Update(id int, title, description, status, priority string, dueDate *time.Time) error {
	query := `
		UPDATE tasks 
		SET title = ?, description = ?, status = ?, priority = ?, due_date = ?, updated_at = ?
		WHERE id = ?
	`
	var dueDateVal interface{}
	if dueDate != nil {
		dueDateVal = dueDate
	}
	_, err := r.db.Exec(query, title, description, status, priority, dueDateVal, time.Now(), id)
	return err
}

func (r *taskRepository) Delete(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *taskRepository) ToggleComplete(id int) error {
	query := `
		UPDATE tasks 
		SET status = CASE WHEN status = 'completed' THEN 'pending' ELSE 'completed' END,
		    updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

