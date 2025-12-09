package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/snip/internal/project"
)

var ErrProjectNotFound = errors.New("project not found")

type ProjectRepository interface {
	Create(p *project.Project) error
	GetByID(id int) (*project.Project, error)
	GetAll(status string) ([]*project.Project, error)
	Update(id int, name, description, status string) error
	Delete(id int) error
	Close() error
}

type projectRepository struct {
	db *sql.DB
}

func NewProjectRepository(db *sql.DB) (ProjectRepository, error) {
	return &projectRepository{db: db}, nil
}

func (r *projectRepository) Close() error {
	return r.db.Close()
}

func (r *projectRepository) Create(p *project.Project) error {
	query := `
		INSERT INTO projects (name, description, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`
	result, err := r.db.Exec(query, p.Name, p.Description, p.Status, p.CreatedAt, p.UpdatedAt)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	p.ID = int(id)
	return nil
}

func (r *projectRepository) GetByID(id int) (*project.Project, error) {
	query := `SELECT id, name, description, status, created_at, updated_at FROM projects WHERE id = ?`
	
	p := &project.Project{}
	err := r.db.QueryRow(query, id).Scan(
		&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrProjectNotFound
		}
		return nil, err
	}

	return p, nil
}

func (r *projectRepository) GetAll(status string) ([]*project.Project, error) {
	var query string
	var args []interface{}

	if status != "" {
		query = `SELECT id, name, description, status, created_at, updated_at FROM projects WHERE status = ? ORDER BY created_at DESC`
		args = []interface{}{status}
	} else {
		query = `SELECT id, name, description, status, created_at, updated_at FROM projects ORDER BY created_at DESC`
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []*project.Project
	for rows.Next() {
		p := &project.Project{}
		err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.Status, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}

	return projects, nil
}

func (r *projectRepository) Update(id int, name, description, status string) error {
	query := `
		UPDATE projects 
		SET name = ?, description = ?, status = ?, updated_at = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, name, description, status, time.Now(), id)
	return err
}

func (r *projectRepository) Delete(id int) error {
	query := `DELETE FROM projects WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

