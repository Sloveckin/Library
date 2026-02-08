package postgres

import (
	"Library/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthorRepositoryPostgres struct {
	pool *pgxpool.Pool
}

func NewAuthorRepositoryPostgres(connectionString string) (*AuthorRepositoryPostgres, error) {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	return &AuthorRepositoryPostgres{pool: pool}, nil
}

func (s *AuthorRepositoryPostgres) Create(name string) (*model.Author, error) {
	var author model.Author
	err := s.pool.QueryRow(context.Background(), "INSERT INTO authors (name) VALUES ($1) RETURNING id, name", name).Scan(&author.Id, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (s *AuthorRepositoryPostgres) Get(id string) (*model.Author, error) {
	var author model.Author
	err := s.pool.QueryRow(context.Background(), "SELECT * FROM authors WHERE id = $1", id).Scan(&author)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (s *AuthorRepositoryPostgres) Delete(id string) error {
	_, err := s.pool.Exec(context.Background(), "DELETE FROM authors WHERE id = $1", id)

	return err
}

func (s *AuthorRepositoryPostgres) ExistsById(id string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM authors WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (s *AuthorRepositoryPostgres) ExistsByName(name string) (bool, error) {
	var exists bool
	err := s.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM authors WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
