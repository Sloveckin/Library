package postgres

import (
	"Library/internal/model"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	NoSuchAuthor = errors.New("no such author")
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

func (a *AuthorRepositoryPostgres) Create(name string) (*model.Author, error) {
	var author model.Author
	err := a.pool.QueryRow(context.Background(), "INSERT INTO authors (name) VALUES ($1) RETURNING id, name", name).Scan(&author.Id, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (a *AuthorRepositoryPostgres) Get(id string) (*model.Author, error) {
	var author model.Author
	err := a.pool.QueryRow(context.Background(), "SELECT id, name FROM authors WHERE id = $1", id).Scan(&author.Id, &author.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, NoSuchAuthor
		}
		return nil, err
	}

	return &author, nil
}

func (a *AuthorRepositoryPostgres) Delete(id string) error {
	ctx := context.Background()

	tx, err := a.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		_ = tx.Commit(ctx)
	}()

	_, err = tx.Exec(ctx, "delete from AuthorToBook where BookId = ($1)", id)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "delete from Authors where id = ($1)", id)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (a *AuthorRepositoryPostgres) ExistsById(id string) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM authors WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a *AuthorRepositoryPostgres) ExistsByName(name string) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM authors WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
