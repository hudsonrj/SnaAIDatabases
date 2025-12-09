package database

import (
	"database/sql"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func GetDBPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dbDir := filepath.Join(homeDir, ".snip")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(dbDir, "notes.db"), nil
}

func Connect() (*sql.DB, error) {
	dbPath, err := GetDBPath()
	if err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := ensureDatabase(db); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureDatabase(db *sql.DB) error {
	query := `
    -- Main Table
    CREATE TABLE IF NOT EXISTS notes (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        content TEXT NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );


    CREATE TABLE IF NOT EXISTS tags (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS notes_tags (
        note_id INTEGER NOT NULL,
        tag_id INTEGER NOT NULL,
        PRIMARY KEY (note_id, tag_id),
        FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
        FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
    );

    -- Index
    CREATE INDEX IF NOT EXISTS idx_notes_title ON notes(title);
    CREATE INDEX IF NOT EXISTS idx_notes_created_at ON notes(created_at);
    
    -- FTS Table
    CREATE VIRTUAL TABLE IF NOT EXISTS notes_fts USING fts4(id, title, content);
    
    -- Populate FTS table with existing data (only if empty)
    INSERT OR IGNORE INTO notes_fts(id, title, content) 
    SELECT id, title, content FROM notes 
    WHERE id NOT IN (SELECT id FROM notes_fts);
    
    -- Triggers
    CREATE TRIGGER IF NOT EXISTS notes_fts_ai AFTER INSERT ON notes BEGIN
        INSERT INTO notes_fts(id, title, content) VALUES (new.id, new.title, new.content);
    END;
    
    CREATE TRIGGER IF NOT EXISTS notes_fts_au AFTER UPDATE ON notes BEGIN
        UPDATE notes_fts SET title = new.title, content = new.content WHERE id = old.id;
    END;
    
    CREATE TRIGGER IF NOT EXISTS notes_fts_ad AFTER DELETE ON notes BEGIN
        DELETE FROM notes_fts WHERE id = old.id;
    END;

    -- Projects Table
    CREATE TABLE IF NOT EXISTS projects (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT NOT NULL,
        description TEXT,
        status TEXT DEFAULT 'active',
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    -- Tasks Table
    CREATE TABLE IF NOT EXISTS tasks (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        project_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        status TEXT DEFAULT 'pending',
        priority TEXT DEFAULT 'medium',
        due_date DATETIME,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
    );

    -- Checklists Table
    CREATE TABLE IF NOT EXISTS checklists (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        task_id INTEGER,
        project_id INTEGER,
        title TEXT NOT NULL,
        description TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (task_id) REFERENCES tasks(id) ON DELETE CASCADE,
        FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
    );

    -- Checklist Items Table
    CREATE TABLE IF NOT EXISTS checklist_items (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        checklist_id INTEGER NOT NULL,
        title TEXT NOT NULL,
        description TEXT,
        completed INTEGER DEFAULT 0,
        item_order INTEGER DEFAULT 0,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        FOREIGN KEY (checklist_id) REFERENCES checklists(id) ON DELETE CASCADE
    );

    -- Indexes
    CREATE INDEX IF NOT EXISTS idx_projects_status ON projects(status);
    CREATE INDEX IF NOT EXISTS idx_tasks_project_id ON tasks(project_id);
    CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
    CREATE INDEX IF NOT EXISTS idx_checklists_task_id ON checklists(task_id);
    CREATE INDEX IF NOT EXISTS idx_checklists_project_id ON checklists(project_id);
    CREATE INDEX IF NOT EXISTS idx_checklist_items_checklist_id ON checklist_items(checklist_id);

    -- Database Analysis Tables
    CREATE TABLE IF NOT EXISTS db_analyses (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        title TEXT NOT NULL,
        database_type TEXT NOT NULL,
        analysis_type TEXT NOT NULL,
        connection_config TEXT NOT NULL,
        log_file_path TEXT,
        output_type TEXT NOT NULL,
        result TEXT,
        ai_insights TEXT,
        status TEXT DEFAULT 'pending',
        error_message TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    -- Error Knowledge Base Table
    CREATE TABLE IF NOT EXISTS error_knowledge_base (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        database_type TEXT NOT NULL,
        error_code TEXT,
        error_message TEXT NOT NULL,
        error_pattern TEXT,
        solution TEXT NOT NULL,
        category TEXT,
        severity TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
        updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );

    -- Indexes for DB Analysis
    CREATE INDEX IF NOT EXISTS idx_db_analyses_database_type ON db_analyses(database_type);
    CREATE INDEX IF NOT EXISTS idx_db_analyses_analysis_type ON db_analyses(analysis_type);
    CREATE INDEX IF NOT EXISTS idx_db_analyses_status ON db_analyses(status);
    CREATE INDEX IF NOT EXISTS idx_db_analyses_created_at ON db_analyses(created_at);
    CREATE INDEX IF NOT EXISTS idx_error_kb_database_type ON error_knowledge_base(database_type);
    CREATE INDEX IF NOT EXISTS idx_error_kb_error_code ON error_knowledge_base(error_code);
    `

	_, err := db.Exec(query)
	return err
}
