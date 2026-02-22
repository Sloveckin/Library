package postgres

import (
	"context"
	"errors"
	"log"

	"Library/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNoSuchAuthor = errors.New("no such author")

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

	err := a.pool.QueryRow(context.Background(), "INSERT INTO Authors (Name) VALUES ($1) RETURNING Id, Name", name).
		Scan(&author.Id, &author.Name)
	if err != nil {
		return nil, err
	}

	return &author, nil
}

func (a *AuthorRepositoryPostgres) Get(id string) (*model.Author, error) {
	var author model.Author

	err := a.pool.QueryRow(context.Background(), "SELECT Id, Name FROM Authors WHERE Id = $1", id).
		Scan(&author.Id, &author.Name)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNoSuchAuthor
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
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf("Failed to rollback transaction: %v", err)
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM AuthorToBook WHERE AuthorId = $1", id)
	if err != nil {
		log.Printf("Failed to delete from AuthorToBook: %v", err)
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM Authors WHERE Id = $1", id)
	if err != nil {
		log.Printf("Failed to delete from Authors: %v", err)
		_ = tx.Rollback(ctx)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf("Failed to commit transaction: %v", err)
		return err
	}

	return nil
}

func (a *AuthorRepositoryPostgres) ExistsById(id string) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM Authors WHERE Id = $1)", id).
		Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (a *AuthorRepositoryPostgres) ExistsByName(name string) (bool, error) {
	var exists bool
	err := a.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT 1 FROM Authors WHERE Name = $1)", name).
		Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
