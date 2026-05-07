package postgres

import (
	"errors"
	"regexp"
	"server/internal/model"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

func newMockRepo(t *testing.T) (*BookPostgresRepository, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	return &BookPostgresRepository{pool: mock}, mock
}

var testAuthor = model.Author{Id: "a1", Name: "Author One"}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCreate_Success(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO Books (Name) VALUES ($1) RETURNING Id, Name")).
		WithArgs("Book One").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", "Book One"))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO AuthorToBook (AuthorId, BookId) VALUES ($1, $2)")).
		WithArgs("a1", "b1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	book, err := repo.Create("Book One", testAuthor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book.Id != "b1" || book.Name != "Book One" {
		t.Errorf("unexpected book: %+v", book)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreate_BeginError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	_, err := repo.Create("Book One", testAuthor)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreate_InsertBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO Books (Name) VALUES ($1) RETURNING Id, Name")).
		WithArgs("Book One").
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	_, err := repo.Create("Book One", testAuthor)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestCreate_InsertAuthorError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO Books (Name) VALUES ($1) RETURNING Id, Name")).
		WithArgs("Book One").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", "Book One"))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO AuthorToBook (AuthorId, BookId) VALUES ($1, $2)")).
		WithArgs("a1", "b1").
		WillReturnError(errors.New("author insert error"))
	mock.ExpectRollback()

	_, err := repo.Create("Book One", testAuthor)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestGet_Success(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Books WHERE Id = $1")).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", "Book One"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT s.Id, s.Name FROM AuthorToBook AS f LEFT JOIN Authors AS s ON f.AuthorId = s.Id WHERE f.BookId = $1")).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", "Author One"))
	mock.ExpectCommit()

	book, err := repo.Get("b1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book.Id != "b1" || len(book.Authors) != 1 {
		t.Errorf("unexpected book: %+v", book)
	}
}

func TestGet_BeginError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGet_QueryBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Books WHERE Id = $1")).
		WithArgs("b1").
		WillReturnError(errors.New("query error"))
	mock.ExpectRollback()

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestGet_QueryAuthorsError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Books WHERE Id = $1")).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", "Book One"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT s.Id, s.Name FROM AuthorToBook AS f LEFT JOIN Authors AS s ON f.AuthorId = s.Id WHERE f.BookId = $1")).
		WithArgs("b1").
		WillReturnError(errors.New("authors query error"))
	mock.ExpectRollback()

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── GetByName ─────────────────────────────────────────────────────────────────

func TestGetByName_Success(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Books WHERE Name = $1")).
		WithArgs("Book One").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", "Book One"))
	mock.ExpectQuery(regexp.QuoteMeta("SELECT s.Id, s.Name FROM AuthorToBook AS f LEFT JOIN Authors AS s ON f.AuthorId = s.Id WHERE f.BookId = $1")).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", "Author One"))
	mock.ExpectCommit()

	book, err := repo.GetByName("Book One")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book.Name != "Book One" {
		t.Errorf("unexpected book: %+v", book)
	}
}

func TestGetByName_QueryError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Books WHERE Name = $1")).
		WithArgs("Book One").
		WillReturnError(errors.New("query error"))
	mock.ExpectRollback()

	_, err := repo.GetByName("Book One")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDelete_Success(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE BookId = $1")).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Books WHERE Id = $1")).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	if err := repo.Delete("b1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDelete_DeleteAuthorToBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE BookId = $1")).
		WithArgs("b1").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	if err := repo.Delete("b1"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestDelete_DeleteBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE BookId = $1")).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Books WHERE Id = $1")).
		WithArgs("b1").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	if err := repo.Delete("b1"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdate_Success(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("UPDATE Books SET Name = $1 WHERE Id = $2 RETURNING Id, Name")).
		WithArgs("New Name", "b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", "New Name"))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE BookId = $1")).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta("INSERT INTO AuthorToBook (AuthorId, BookId) VALUES ($1, $2)")).
		WithArgs("a1", "b1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	book, err := repo.Update("b1", "New Name", testAuthor)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if book.Name != "New Name" {
		t.Errorf("unexpected book name: %v", book.Name)
	}
}

func TestUpdate_UpdateError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("UPDATE Books SET Name = $1 WHERE Id = $2 RETURNING Id, Name")).
		WithArgs("New Name", "b1").
		WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	_, err := repo.Update("b1", "New Name", testAuthor)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestUpdate_DeleteAuthorsError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("UPDATE Books SET Name = $1 WHERE Id = $2 RETURNING Id, Name")).
		WithArgs("New Name", "b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", "New Name"))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE BookId = $1")).
		WithArgs("b1").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	_, err := repo.Update("b1", "New Name", testAuthor)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── ExistsById ────────────────────────────────────────────────────────────────

func TestExistsById_True(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT FROM Books WHERE id = $1)")).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsById("b1")
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestExistsById_False(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT FROM Books WHERE id = $1)")).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := repo.ExistsById("b1")
	if err != nil || exists {
		t.Errorf("expected false, got %v, %v", exists, err)
	}
}

func TestExistsById_Error(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT FROM Books WHERE id = $1)")).
		WithArgs("b1").
		WillReturnError(errors.New("query error"))

	_, err := repo.ExistsById("b1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── ExistsByName ──────────────────────────────────────────────────────────────

func TestExistsByName_True(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT FROM Books WHERE name = $1)")).
		WithArgs("Book One").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByName("Book One")
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestExistsByName_Error(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT FROM Books WHERE name = $1)")).
		WithArgs("Book One").
		WillReturnError(errors.New("query error"))

	_, err := repo.ExistsByName("Book One")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}