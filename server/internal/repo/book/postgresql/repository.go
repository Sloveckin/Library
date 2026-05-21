package postgres

import (
	"context"
	"log"
	"server/internal/model"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	logBeginTxFailed   = "Failed to begin transaction: %v"
	logRollbackFailed  = "Failed to rollback transaction: %v"
	logQueryBookFailed = "Failed to query book: %v"
	logCommitTxFailed  = "Failed to commit transaction: %v"
)

type pgxPool interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type BookPostgresRepository struct {
	pool pgxPool
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
		log.Printf(logBeginTxFailed, err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf(logRollbackFailed, err)
		}
	}()

	var book model.Book
	book.Authors = make([]model.Author, 0, len(authors))

	err = tx.QueryRow(ctx, "INSERT INTO Books (Name) VALUES ($1) RETURNING Id, Name", name).Scan(&book.Id, &book.Name)
	if err != nil {
		log.Printf(logQueryBookFailed, err)
		_ = tx.Rollback(ctx)
		return nil, err
	}

	for _, author := range authors {
		_, err = tx.Exec(ctx, "INSERT INTO AuthorToBook (AuthorId, BookId) VALUES ($1, $2)", author.Id, book.Id)
		if err != nil {
			log.Printf("Failed to insert into AuthorToBook: %v", err)
			_ = tx.Rollback(ctx)
			return nil, err
		}
	}

	book.Authors = append(book.Authors, authors...)

	if err = tx.Commit(ctx); err != nil {
		log.Printf(logCommitTxFailed, err)
		return nil, err
	}

	return &book, nil
}

func (b BookPostgresRepository) Get(id string) (*model.Book, error) {
	return b.getByField("SELECT Id, Name FROM Books WHERE Id = $1", id)
}

func (b BookPostgresRepository) GetByName(name string) (*model.Book, error) {
	return b.getByField("SELECT Id, Name FROM Books WHERE Name = $1", name)
}

func (b BookPostgresRepository) getByField(query, value string) (*model.Book, error) {
	ctx := context.Background()

	tx, err := b.pool.Begin(ctx)
	if err != nil {
		log.Printf(logBeginTxFailed, err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf(logRollbackFailed, err)
		}
	}()

	var book model.Book
	err = tx.QueryRow(ctx, query, value).Scan(&book.Id, &book.Name)
	if err != nil {
		log.Printf(logQueryBookFailed, err)
		_ = tx.Rollback(ctx)
		return nil, err
	}

	rows, err := tx.Query(
		ctx,
		"SELECT s.Id, s.Name FROM AuthorToBook AS f LEFT JOIN Authors AS s ON f.AuthorId = s.Id WHERE f.BookId = $1",
		book.Id,
	)
	if err != nil {
		log.Printf("Failed to query authors: %v", err)
		_ = tx.Rollback(ctx)
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var author model.Author
		err = rows.Scan(&author.Id, &author.Name)
		if err != nil {
			log.Printf("Failed to scan author: %v", err)
			_ = tx.Rollback(ctx)
			return nil, err
		}

		book.Authors = append(book.Authors, author)
	}

	if err = rows.Err(); err != nil {
		log.Printf("Rows iteration error: %v", err)
		_ = tx.Rollback(ctx)
		return nil, err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf(logCommitTxFailed, err)
		return nil, err
	}

	return &book, nil
}

func (b BookPostgresRepository) Delete(id string) error {
	ctx := context.Background()

	tx, err := b.pool.Begin(ctx)
	if err != nil {
		log.Printf(logBeginTxFailed, err)
		return err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf(logRollbackFailed, err)
		}
	}()

	_, err = tx.Exec(ctx, "DELETE FROM AuthorToBook WHERE BookId = $1", id)
	if err != nil {
		log.Printf("Failed to delete from AuthorToBook: %v", err)
		_ = tx.Rollback(ctx)
		return err
	}

	_, err = tx.Exec(ctx, "DELETE FROM Books WHERE Id = $1", id)
	if err != nil {
		log.Printf("Failed to delete from Books: %v", err)
		_ = tx.Rollback(ctx)
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf(logCommitTxFailed, err)
		return err
	}

	return nil
}

func (b BookPostgresRepository) Update(id, name string, authors ...model.Author) (*model.Book, error) {
	ctx := context.Background()

	tx, err := b.pool.Begin(ctx)
	if err != nil {
		log.Printf(logBeginTxFailed, err)
		return nil, err
	}

	defer func() {
		if err := tx.Rollback(ctx); err != nil && err != pgx.ErrTxClosed {
			log.Printf(logRollbackFailed, err)
		}
	}()

	var book model.Book
	book.Authors = make([]model.Author, 0, len(authors))

	err = tx.QueryRow(ctx, "UPDATE Books SET Name = $1 WHERE Id = $2 RETURNING Id, Name", name, id).Scan(&book.Id, &book.Name)
	if err != nil {
		log.Printf("Failed to update book: %v", err)
		_ = tx.Rollback(ctx)
		return nil, err
	}

	_, err = tx.Exec(ctx, "DELETE FROM AuthorToBook WHERE BookId = $1", id)
	if err != nil {
		log.Printf("Failed to delete old authors: %v", err)
		_ = tx.Rollback(ctx)
		return nil, err
	}

	for _, author := range authors {
		_, err = tx.Exec(ctx, "INSERT INTO AuthorToBook (AuthorId, BookId) VALUES ($1, $2)", author.Id, book.Id)
		if err != nil {
			log.Printf("Failed to insert into AuthorToBook: %v", err)
			_ = tx.Rollback(ctx)
			return nil, err
		}

		book.Authors = append(book.Authors, author)
	}

	if err = tx.Commit(ctx); err != nil {
		log.Printf(logCommitTxFailed, err)
		return nil, err
	}

	return &book, nil
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
