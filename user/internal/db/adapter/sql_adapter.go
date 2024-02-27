package adapter

import (
	"database/sql"
	"projects/LDmitryLD/task-service/user/internal/models"

	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SQLAdapterer interface {
	GetUserByID(id int) (models.User, error)
	CreateUser(user models.User) (int, error)
}

type SQLAdapter struct {
	db *sqlx.DB
}

func NewSQLAdapter(db *sqlx.DB) *SQLAdapter {
	return &SQLAdapter{
		db: db,
	}
}

func (s *SQLAdapter) GetUserByID(id int) (models.User, error) {
	q := `
	SELECT 
		first_name, last_name, email
	 FROM 
	 	users
	WHERE 
		id = $1
	`
	var user models.User
	err := s.db.Get(&user, q, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.User{}, status.Error(codes.NotFound, "not found")
		}
		return models.User{}, err
	}

	return user, err
}

func (s *SQLAdapter) CreateUser(user models.User) (int, error) {
	q := `
	INSERT INTO users
		(first_name, last_name, email)
	VALUES
		($1, $2, $3)	
	RETURNING id	
	`
	var id int
	if err := s.db.QueryRow(q, user.FirstName, user.LastName, user.Email).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}
