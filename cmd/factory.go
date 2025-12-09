package cmd

import (
	"fmt"
	"sync"

	"github.com/snip/internal/database"
	"github.com/snip/internal/handler"
	"github.com/snip/internal/repository"
)

var (
	globalNoteRepo          repository.NoteRepository
	globalTagRepo           repository.TagRepository
	globalProjectRepo       repository.ProjectRepository
	globalTaskRepo          repository.TaskRepository
	globalChecklistRepo     repository.ChecklistRepository
	globalChecklistItemRepo repository.ChecklistItemRepository
	globalDBAnalysisRepo    repository.DBAnalysisRepository
	repoOnce                sync.Once
)

func getRepository() (repository.NoteRepository, repository.TagRepository, error) {
	var err error
	repoOnce.Do(func() {
		db, connectErr := database.Connect()
		if connectErr != nil {
			err = connectErr
			return
		}
		globalNoteRepo, err = repository.NewNoteRepository(db)
		if err != nil {
			return
		}
		globalTagRepo, err = repository.NewTagRepository(db)
		if err != nil {
			return
		}
		globalProjectRepo, err = repository.NewProjectRepository(db)
		if err != nil {
			return
		}
		globalTaskRepo, err = repository.NewTaskRepository(db)
		if err != nil {
			return
		}
		globalChecklistRepo, err = repository.NewChecklistRepository(db)
		if err != nil {
			return
		}
		globalChecklistItemRepo, err = repository.NewChecklistItemRepository(db)
		if err != nil {
			return
		}
		globalDBAnalysisRepo, err = repository.NewDBAnalysisRepository(db)
	})
	return globalNoteRepo, globalTagRepo, err
}

func setupHandler() (handler.Handler, error) {
	noteRepo, tagRepo, err := getRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h := handler.NewHandler(noteRepo, tagRepo)

	return h, nil
}

func setupProjectHandler() (handler.ProjectHandler, error) {
	_, _, err := getRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h := handler.NewProjectHandler(globalProjectRepo, globalTaskRepo)
	return h, nil
}

func setupTaskHandler() (handler.TaskHandler, error) {
	_, _, err := getRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h := handler.NewTaskHandler(globalTaskRepo, globalProjectRepo)
	return h, nil
}

func setupChecklistHandler() (handler.ChecklistHandler, error) {
	_, _, err := getRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h := handler.NewChecklistHandler(globalChecklistRepo, globalChecklistItemRepo)
	return h, nil
}

func executeWithHandler(fn func(handler.Handler) error) error {
	h, err := setupHandler()
	if err != nil {
		return fmt.Errorf("failed to setup handler: %w", err)
	}

	return fn(h)
}

func executeWithProjectHandler(fn func(handler.ProjectHandler) error) error {
	h, err := setupProjectHandler()
	if err != nil {
		return fmt.Errorf("failed to setup project handler: %w", err)
	}

	return fn(h)
}

func executeWithTaskHandler(fn func(handler.TaskHandler) error) error {
	h, err := setupTaskHandler()
	if err != nil {
		return fmt.Errorf("failed to setup task handler: %w", err)
	}

	return fn(h)
}

func executeWithChecklistHandler(fn func(handler.ChecklistHandler) error) error {
	h, err := setupChecklistHandler()
	if err != nil {
		return fmt.Errorf("failed to setup checklist handler: %w", err)
	}

	return fn(h)
}

func setupDBAnalysisHandler() (handler.DBAnalysisHandler, error) {
	_, _, err := getRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	h, err := handler.NewDBAnalysisHandler(globalDBAnalysisRepo)
	if err != nil {
		return nil, fmt.Errorf("failed to setup DB analysis handler: %w", err)
	}

	return h, nil
}

func executeWithDBAnalysisHandler(fn func(handler.DBAnalysisHandler) error) error {
	h, err := setupDBAnalysisHandler()
	if err != nil {
		return fmt.Errorf("failed to setup DB analysis handler: %w", err)
	}

	return fn(h)
}
