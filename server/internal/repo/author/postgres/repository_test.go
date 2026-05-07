package postgres

import (
	"errors"
	"regexp"
	"server/internal/model"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
)

func newMockAuthorRepo(t *testing.T) (*AuthorRepositoryPostgres, pgxmock.PgxPoolIface) {
	t.Helper()
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatalf("failed to create pgxmock pool: %v", err)
	}
	return &AuthorRepositoryPostgres{pool: mock}, mock
}

var testAuthorModel = model.Author{Id: "a1", Name: "Author One"}

// ── Create ────────────────────────────────────────────────────────────────────

func TestAuthorCreate_Success(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO Authors (Name) VALUES ($1) RETURNING Id, Name")).
		WithArgs("Author One").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", "Author One"))

	author, err := repo.Create("Author One")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if author.Id != "a1" || author.Name != "Author One" {
		t.Errorf("unexpected author: %+v", author)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestAuthorCreate_Error(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO Authors (Name) VALUES ($1) RETURNING Id, Name")).
		WithArgs("Author One").
		WillReturnError(errors.New("insert error"))

	_, err := repo.Create("Author One")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestAuthorGet_Success(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Authors WHERE Id = $1")).
		WithArgs("a1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", "Author One"))

	author, err := repo.Get("a1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if author.Id != "a1" {
		t.Errorf("unexpected author: %+v", author)
	}
}

func TestAuthorGet_NotFound(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Authors WHERE Id = $1")).
		WithArgs("a1").
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.Get("a1")
	if !errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf("expected ErrNoSuchAuthor, got %v", err)
	}
}

func TestAuthorGet_Error(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Authors WHERE Id = $1")).
		WithArgs("a1").
		WillReturnError(errors.New("db error"))

	_, err := repo.Get("a1")
	if err == nil || errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf("expected generic error, got %v", err)
	}
}

// ── GetByName ─────────────────────────────────────────────────────────────────

func TestAuthorGetByName_Success(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Authors WHERE Name = $1")).
		WithArgs("Author One").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", "Author One"))

	author, err := repo.GetByName("Author One")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if author.Name != "Author One" {
		t.Errorf("unexpected author: %+v", author)
	}
}

func TestAuthorGetByName_NotFound(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Authors WHERE Name = $1")).
		WithArgs("Author One").
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByName("Author One")
	if !errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf("expected ErrNoSuchAuthor, got %v", err)
	}
}

func TestAuthorGetByName_Error(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT Id, Name FROM Authors WHERE Name = $1")).
		WithArgs("Author One").
		WillReturnError(errors.New("db error"))

	_, err := repo.GetByName("Author One")
	if err == nil || errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf("expected generic error, got %v", err)
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestAuthorUpdate_Success(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("UPDATE Authors SET Name = $1 WHERE Id = $2 RETURNING Id, Name")).
		WithArgs("New Name", "a1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", "New Name"))

	author, err := repo.Update("a1", "New Name")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if author.Name != "New Name" {
		t.Errorf("unexpected author name: %v", author.Name)
	}
}

func TestAuthorUpdate_NotFound(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("UPDATE Authors SET Name = $1 WHERE Id = $2 RETURNING Id, Name")).
		WithArgs("New Name", "a1").
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.Update("a1", "New Name")
	if !errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf("expected ErrNoSuchAuthor, got %v", err)
	}
}

func TestAuthorUpdate_Error(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("UPDATE Authors SET Name = $1 WHERE Id = $2 RETURNING Id, Name")).
		WithArgs("New Name", "a1").
		WillReturnError(errors.New("db error"))

	_, err := repo.Update("a1", "New Name")
	if err == nil || errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf("expected generic error, got %v", err)
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestAuthorDelete_Success(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE AuthorId = $1")).
		WithArgs("a1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Authors WHERE Id = $1")).
		WithArgs("a1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	if err := repo.Delete("a1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestAuthorDelete_BeginError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	if err := repo.Delete("a1"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAuthorDelete_DeleteAuthorToBookError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE AuthorId = $1")).
		WithArgs("a1").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	if err := repo.Delete("a1"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestAuthorDelete_DeleteAuthorError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM AuthorToBook WHERE AuthorId = $1")).
		WithArgs("a1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Authors WHERE Id = $1")).
		WithArgs("a1").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	if err := repo.Delete("a1"); err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── ExistsById ────────────────────────────────────────────────────────────────

func TestAuthorExistsById_True(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Authors WHERE Id = $1)")).
		WithArgs("a1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsById("a1")
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestAuthorExistsById_False(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Authors WHERE Id = $1)")).
		WithArgs("a1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := repo.ExistsById("a1")
	if err != nil || exists {
		t.Errorf("expected false, got %v, %v", exists, err)
	}
}

func TestAuthorExistsById_Error(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Authors WHERE Id = $1)")).
		WithArgs("a1").
		WillReturnError(errors.New("db error"))

	_, err := repo.ExistsById("a1")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// ── ExistsByName ──────────────────────────────────────────────────────────────

func TestAuthorExistsByName_True(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Authors WHERE Name = $1)")).
		WithArgs("Author One").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByName("Author One")
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestAuthorExistsByName_Error(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Authors WHERE Name = $1)")).
		WithArgs("Author One").
		WillReturnError(errors.New("db error"))

	_, err := repo.ExistsByName("Author One")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}