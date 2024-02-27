package adapter

import (
	"context"
	"log"
	"projects/LDmitryLD/task-service/task/internal/models"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SQLAdapterer interface {
	List(userID int) ([]models.Task, error)
	Create(task models.Task) (int, error)
	Delete(ctx context.Context, userId int, taskId int) error
}

type SQLAdapter struct {
	db *sqlx.DB
}

func NewSQLAdapter(db *sqlx.DB) *SQLAdapter {
	return &SQLAdapter{
		db: db,
	}
}

func (s *SQLAdapter) List(userID int) ([]models.Task, error) {
	var tasks []models.Task
	q := `
	SELECT
		id, task_name, description
	FROM
		tasks
	WHERE user_id = $1		
	`

	err := s.db.Select(&tasks, q, userID)
	if err != nil {
		log.Println("SQLAdapter.List err:", err)
		return nil, err
	}

	return tasks, nil
}

func (s *SQLAdapter) Create(task models.Task) (int, error) {
	q := `
	INSERT INTO tasks
		(user_id, task_name, description)
	VALUES
		($1, $2, $3)	
	RETURNING id	
	`
	var taskID int
	err := s.db.QueryRow(q, task.UserId, task.TaskName, task.Description).Scan(&taskID)
	if err != nil {
		return 0, err
	}

	return taskID, nil
}

func (s *SQLAdapter) Delete(ctx context.Context, userId int, taskId int) error {
	q := `
	DELETE FROM
		tasks
	WHERE 
		id = $1 AND user_id = $2	
	`

	result, err := s.db.ExecContext(ctx, q, taskId, userId)
	if err != nil {
		return err
	}

	rowsAffecred, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffecred == 0 {
		return status.Error(codes.NotFound, "task not found")
	}

	return nil
}
