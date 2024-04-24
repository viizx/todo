package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Health() (map[string]string, error)
	Init() error
	CreateToDo(t ToDoRequest) error
	UpdateToDo(id string, t ToDoRequest) error
	GetAllToDos(sort string) ([]ToDo, error)
	GetToDo(id string) (*ToDo, error)
}

type service struct {
	db *sql.DB
}

var (
	database   = os.Getenv("DB_DATABASE")
	password   = os.Getenv("DB_PASSWORD")
	username   = os.Getenv("DB_USERNAME")
	port       = os.Getenv("DB_PORT")
	host       = os.Getenv("DB_HOST")
	dbInstance *service
)

func New() Service {
	if dbInstance != nil {
		return dbInstance
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Failed to open database connection: %v", err)
	}

	dbInstance = &service{db: db}
	return dbInstance
}

func (s *service) Health() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := s.db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("database health check failed: %w", err)
	}

	return map[string]string{"message": "It's healthy"}, nil
}

func (s *service) Init() error {
	if _, err := s.db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";"); err != nil {
		return fmt.Errorf("error creating uuid extension: %w", err)
	}
	return s.createToDoTable()
}

func (s *service) createToDoTable() error {
	query := `
    CREATE TABLE IF NOT EXISTS todo (
        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
        text TEXT,
        done BOOLEAN,
        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
    );
    `
	if _, err := s.db.Exec(query); err != nil {
		return fmt.Errorf("error creating todos table: %w", err)
	}
	return nil
}

func (s *service) CreateToDo(t ToDoRequest) error {
	query := `INSERT INTO todo (text, done) VALUES ($1, $2)`
	if _, err := s.db.Exec(query, t.Text, t.Done); err != nil {
		return fmt.Errorf("failed to create todo: %w", err)
	}
	return nil
}

func (s *service) UpdateToDo(id string, t ToDoRequest) error {
	query := `UPDATE todo SET 
        text = COALESCE($1, text),
        done = COALESCE($2, done),
	      updated_at = CURRENT_TIMESTAMP
    WHERE id = $3`

	var newText interface{}
	if t.Text != nil {
		newText = *t.Text
	}

	var newDone interface{}
	if t.Done != nil {
		newDone = *t.Done
	}

	if _, err := s.db.Exec(query, newText, newDone, id); err != nil {
		return fmt.Errorf("failed to update todo: %w", err)
	}

	return nil
}

func (s *service) GetAllToDos(sort string) ([]ToDo, error) {
	var todos []ToDo = []ToDo{}
	query := fmt.Sprintf(`SELECT id, text, done, created_at, updated_at FROM todo ORDER BY created_at %s;`, sort)

	rows, err := s.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve todos: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r ToDo
		if err := rows.Scan(&r.ID, &r.Text, &r.Done, &r.CreatedAt, &r.UpdatedAt); err != nil {
			return nil, fmt.Errorf("failed to scan todo: %w", err)
		}
		todos = append(todos, r)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return todos, nil
}

func (s *service) GetToDo(id string) (*ToDo, error) {
	var t ToDo
	query := `SELECT id, text, done, created_at, updated_at FROM todo WHERE id = $1`
	if err := s.db.QueryRow(query, id).Scan(&t.ID, &t.Text, &t.Done, &t.CreatedAt, &t.UpdatedAt); err != nil {
		return nil, fmt.Errorf("failed to get todo: %w", err)
	}
	return &t, nil
}

type ToDo struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ToDoRequest struct {
	Text *string `json:"text"`
	Done *bool   `json:"done"`
}
