package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/snip/internal/dbanalysis"
)

type DBAnalysisRepository interface {
	Create(analysis *dbanalysis.DBAnalysis) error
	GetByID(id int) (*dbanalysis.DBAnalysis, error)
	GetAll(limit int, dbType dbanalysis.DatabaseType, analysisType dbanalysis.AnalysisType) ([]*dbanalysis.DBAnalysis, error)
	Update(analysis *dbanalysis.DBAnalysis) error
	Delete(id int) error
	GetRecent(limit int) ([]*dbanalysis.DBAnalysis, error)
	Close() error
}

type dbAnalysisRepository struct {
	db *sql.DB
}

func NewDBAnalysisRepository(db *sql.DB) (DBAnalysisRepository, error) {
	return &dbAnalysisRepository{db: db}, nil
}

func (r *dbAnalysisRepository) Close() error {
	return r.db.Close()
}

func (r *dbAnalysisRepository) Create(analysis *dbanalysis.DBAnalysis) error {
	query := `
		INSERT INTO db_analyses (
			title, database_type, analysis_type, connection_config, 
			log_file_path, output_type, result, ai_insights, 
			status, error_message, created_at, updated_at
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	result, err := r.db.Exec(
		query,
		analysis.Title,
		string(analysis.DatabaseType),
		string(analysis.AnalysisType),
		analysis.ConnectionConfig,
		analysis.LogFilePath,
		string(analysis.OutputType),
		analysis.Result,
		analysis.AIInsights,
		analysis.Status,
		analysis.ErrorMessage,
		analysis.CreatedAt,
		analysis.UpdatedAt,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	analysis.ID = int(id)
	return nil
}

func (r *dbAnalysisRepository) GetByID(id int) (*dbanalysis.DBAnalysis, error) {
	query := `
		SELECT id, title, database_type, analysis_type, connection_config,
		       log_file_path, output_type, result, ai_insights, status,
		       error_message, created_at, updated_at
		FROM db_analyses
		WHERE id = ?
	`

	analysis := &dbanalysis.DBAnalysis{}
	var dbType, analysisType, outputType, status string
	var logPath, result, aiInsights, errorMsg sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&analysis.ID,
		&analysis.Title,
		&dbType,
		&analysisType,
		&analysis.ConnectionConfig,
		&logPath,
		&outputType,
		&result,
		&aiInsights,
		&status,
		&errorMsg,
		&analysis.CreatedAt,
		&analysis.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("not found")
		}
		return nil, err
	}

	analysis.DatabaseType = dbanalysis.DatabaseType(dbType)
	analysis.AnalysisType = dbanalysis.AnalysisType(analysisType)
	analysis.OutputType = dbanalysis.OutputType(outputType)
	analysis.Status = status

	if logPath.Valid {
		analysis.LogFilePath = logPath.String
	}
	if result.Valid {
		analysis.Result = result.String
	}
	if aiInsights.Valid {
		analysis.AIInsights = aiInsights.String
	}
	if errorMsg.Valid {
		analysis.ErrorMessage = errorMsg.String
	}

	return analysis, nil
}

func (r *dbAnalysisRepository) GetAll(limit int, dbType dbanalysis.DatabaseType, analysisType dbanalysis.AnalysisType) ([]*dbanalysis.DBAnalysis, error) {
	query := `
		SELECT id, title, database_type, analysis_type, connection_config,
		       log_file_path, output_type, result, ai_insights, status,
		       error_message, created_at, updated_at
		FROM db_analyses
		WHERE 1=1
	`
	args := []interface{}{}

	if dbType != "" {
		query += " AND database_type = ?"
		args = append(args, string(dbType))
	}

	if analysisType != "" {
		query += " AND analysis_type = ?"
		args = append(args, string(analysisType))
	}

	query += " ORDER BY created_at DESC"

	if limit > 0 {
		query += " LIMIT ?"
		args = append(args, limit)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var analyses []*dbanalysis.DBAnalysis
	for rows.Next() {
		analysis := &dbanalysis.DBAnalysis{}
		var dbType, analysisType, outputType, status string
		var logPath, result, aiInsights, errorMsg sql.NullString

		err := rows.Scan(
			&analysis.ID,
			&analysis.Title,
			&dbType,
			&analysisType,
			&analysis.ConnectionConfig,
			&logPath,
			&outputType,
			&result,
			&aiInsights,
			&status,
			&errorMsg,
			&analysis.CreatedAt,
			&analysis.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		analysis.DatabaseType = dbanalysis.DatabaseType(dbType)
		analysis.AnalysisType = dbanalysis.AnalysisType(analysisType)
		analysis.OutputType = dbanalysis.OutputType(outputType)
		analysis.Status = status

		if logPath.Valid {
			analysis.LogFilePath = logPath.String
		}
		if result.Valid {
			analysis.Result = result.String
		}
		if aiInsights.Valid {
			analysis.AIInsights = aiInsights.String
		}
		if errorMsg.Valid {
			analysis.ErrorMessage = errorMsg.String
		}

		analyses = append(analyses, analysis)
	}

	return analyses, nil
}

func (r *dbAnalysisRepository) Update(analysis *dbanalysis.DBAnalysis) error {
	analysis.UpdatedAt = time.Now()
	query := `
		UPDATE db_analyses
		SET title = ?, result = ?, ai_insights = ?, status = ?, 
		    error_message = ?, updated_at = ?
		WHERE id = ?
	`

	_, err := r.db.Exec(
		query,
		analysis.Title,
		analysis.Result,
		analysis.AIInsights,
		analysis.Status,
		analysis.ErrorMessage,
		analysis.UpdatedAt,
		analysis.ID,
	)

	return err
}

func (r *dbAnalysisRepository) Delete(id int) error {
	query := `DELETE FROM db_analyses WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *dbAnalysisRepository) GetRecent(limit int) ([]*dbanalysis.DBAnalysis, error) {
	return r.GetAll(limit, "", "")
}

