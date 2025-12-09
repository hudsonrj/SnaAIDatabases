package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/snip/internal/checklist"
)

var ErrChecklistNotFound = errors.New("checklist not found")

type ChecklistRepository interface {
	Create(c *checklist.Checklist) error
	GetByID(id int) (*checklist.Checklist, error)
	GetByTaskID(taskID int) ([]*checklist.Checklist, error)
	GetByProjectID(projectID int) ([]*checklist.Checklist, error)
	GetAll() ([]*checklist.Checklist, error)
	Update(id int, title, description string) error
	Delete(id int) error
	Close() error
}

type ChecklistItemRepository interface {
	Create(item *checklist.ChecklistItem) error
	GetByChecklistID(checklistID int) ([]*checklist.ChecklistItem, error)
	Update(id int, title, description string, completed bool) error
	ToggleComplete(id int) error
	Delete(id int) error
	Close() error
}

type checklistRepository struct {
	db *sql.DB
}

type checklistItemRepository struct {
	db *sql.DB
}

func NewChecklistRepository(db *sql.DB) (ChecklistRepository, error) {
	return &checklistRepository{db: db}, nil
}

func NewChecklistItemRepository(db *sql.DB) (ChecklistItemRepository, error) {
	return &checklistItemRepository{db: db}, nil
}

func (r *checklistRepository) Close() error {
	return r.db.Close()
}

func (r *checklistItemRepository) Close() error {
	return r.db.Close()
}

func (r *checklistRepository) Create(c *checklist.Checklist) error {
	query := `
		INSERT INTO checklists (task_id, project_id, title, description, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?)
	`
	var taskID, projectID interface{}
	if c.TaskID != nil {
		taskID = *c.TaskID
	}
	if c.ProjectID != nil {
		projectID = *c.ProjectID
	}

	result, err := r.db.Exec(query, taskID, projectID, c.Title, c.Description, c.CreatedAt, c.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	c.ID = int(id)
	return nil
}

func (r *checklistRepository) GetByID(id int) (*checklist.Checklist, error) {
	query := `SELECT id, task_id, project_id, title, description, created_at, updated_at FROM checklists WHERE id = ?`
	
	c := &checklist.Checklist{}
	var taskID, projectID sql.NullInt64
	err := r.db.QueryRow(query, id).Scan(
		&c.ID, &taskID, &projectID, &c.Title, &c.Description, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrChecklistNotFound
		}
		return nil, err
	}

	if taskID.Valid {
		tid := int(taskID.Int64)
		c.TaskID = &tid
	}
	if projectID.Valid {
		pid := int(projectID.Int64)
		c.ProjectID = &pid
	}

	return c, nil
}

func (r *checklistRepository) GetByTaskID(taskID int) ([]*checklist.Checklist, error) {
	query := `SELECT id, task_id, project_id, title, description, created_at, updated_at 
		FROM checklists WHERE task_id = ? ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanChecklists(rows)
}

func (r *checklistRepository) GetByProjectID(projectID int) ([]*checklist.Checklist, error) {
	query := `SELECT id, task_id, project_id, title, description, created_at, updated_at 
		FROM checklists WHERE project_id = ? ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query, projectID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanChecklists(rows)
}

func (r *checklistRepository) GetAll() ([]*checklist.Checklist, error) {
	query := `SELECT id, task_id, project_id, title, description, created_at, updated_at 
		FROM checklists ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanChecklists(rows)
}

func (r *checklistRepository) scanChecklists(rows *sql.Rows) ([]*checklist.Checklist, error) {
	var checklists []*checklist.Checklist
	for rows.Next() {
		c := &checklist.Checklist{}
		var taskID, projectID sql.NullInt64
		err := rows.Scan(&c.ID, &taskID, &projectID, &c.Title, &c.Description, &c.CreatedAt, &c.UpdatedAt)
		if err != nil {
			return nil, err
		}
		if taskID.Valid {
			tid := int(taskID.Int64)
			c.TaskID = &tid
		}
		if projectID.Valid {
			pid := int(projectID.Int64)
			c.ProjectID = &pid
		}
		checklists = append(checklists, c)
	}
	return checklists, nil
}

func (r *checklistRepository) Update(id int, title, description string) error {
	query := `UPDATE checklists SET title = ?, description = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, title, description, time.Now(), id)
	return err
}

func (r *checklistRepository) Delete(id int) error {
	query := `DELETE FROM checklists WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

// ChecklistItemRepository methods

func (r *checklistItemRepository) Create(item *checklist.ChecklistItem) error {
	query := `
		INSERT INTO checklist_items (checklist_id, title, description, completed, item_order, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	completed := 0
	if item.Completed {
		completed = 1
	}

	result, err := r.db.Exec(query, item.ChecklistID, item.Title, item.Description, completed, item.Order, item.CreatedAt, item.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	item.ID = int(id)
	return nil
}

func (r *checklistItemRepository) GetByChecklistID(checklistID int) ([]*checklist.ChecklistItem, error) {
	query := `SELECT id, checklist_id, title, description, completed, item_order, created_at, updated_at 
		FROM checklist_items WHERE checklist_id = ? ORDER BY item_order ASC, created_at ASC`
	
	rows, err := r.db.Query(query, checklistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*checklist.ChecklistItem
	for rows.Next() {
		item := &checklist.ChecklistItem{}
		var completed int
		err := rows.Scan(&item.ID, &item.ChecklistID, &item.Title, &item.Description, &completed, &item.Order, &item.CreatedAt, &item.UpdatedAt)
		if err != nil {
			return nil, err
		}
		item.Completed = completed == 1
		items = append(items, item)
	}

	return items, nil
}

func (r *checklistItemRepository) Update(id int, title, description string, completed bool) error {
	completedVal := 0
	if completed {
		completedVal = 1
	}
	query := `UPDATE checklist_items SET title = ?, description = ?, completed = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, title, description, completedVal, time.Now(), id)
	return err
}

func (r *checklistItemRepository) ToggleComplete(id int) error {
	query := `
		UPDATE checklist_items 
		SET completed = CASE WHEN completed = 1 THEN 0 ELSE 1 END,
		    updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, time.Now(), id)
	return err
}

func (r *checklistItemRepository) Delete(id int) error {
	query := `DELETE FROM checklist_items WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

