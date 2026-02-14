package postgres

import (
	"Library/internal/model"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type BookPostgresRepository struct {
	pool *pgxpool.Pool
}

func NewBookPostgresRepository(connectionString string) (*BookPostgresRepository, error) {
	pool, err := pgxpool.New(context.Background(), connectionString)
	if err != nil {
		return nil, err
	}

	return &BookPostgresRepository{pool: pool}, nil
}

func (b BookPostgresRepository) Create(name string, authors ...model.Author) (*model.Book, error) {
	ctx := context.Background()

	tx, err := b.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Commit(ctx)
	}()

	var book model.Book
	book.Authors = make([]model.Author, 0, len(authors))

	err = tx.QueryRow(ctx, "INSERT INTO books (name) VALUES ($1) RETURNING id, name", name).Scan(&book.Id, &book.Name)
	if err != nil {
		return nil, err
	}

	for _, author := range authors {
		_, err = tx.Exec(ctx, "INSERT INTO AuthorToBook (AuthorId, BookId) VALUES ($1, $2)", author.Id, book.Id)
		if err != nil {
			return nil, err
		}
	}

	book.Authors = append(book.Authors, authors...)

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookPostgresRepository) Get(id string) (*model.Book, error) {
	ctx := context.Background()

	tx, err := b.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = tx.Commit(ctx)
	}()

	var book model.Book
	err = tx.QueryRow(ctx, "SELECT id, name FROM books WHERE id = $1", id).Scan(&book.Id, &book.Name)
	if err != nil {
		return nil, err
	}

	rows, err := tx.Query(ctx, "SELECT s.Id, s.Name FROM AuthorToBook as f left join Authors as s on f.AuthorId = s.Id  WHERE BookId = $1", id)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var author model.Author
		err = rows.Scan(&author.Id, &author.Name)
		if err != nil {
			return nil, err
		}

		book.Authors = append(book.Authors, author)
	}

	if err = tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &book, nil
}

func (b BookPostgresRepository) Delete(id string) error {
	ctx := context.Background()

	tx, err := b.pool.Begin(ctx)
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

	_, err = tx.Exec(ctx, "delete from Books where id = ($1)", id)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (b BookPostgresRepository) ExistsById(id string) (bool, error) {
	var exists bool
	err := b.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM Books WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (b BookPostgresRepository) ExistsByName(name string) (bool, error) {
	var exists bool
	err := b.pool.QueryRow(context.Background(), "SELECT EXISTS(SELECT FROM Books WHERE name = $1)", name).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}
