package postgres

import (
	"errors"
	"regexp"
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

const (
	authorOneName                 = "Author One"
	updatedAuthorName             = "New Name"
	dbErrorMessage                = "db error"
	errUnexpectedFmt              = "unexpected error: %v"
	errUnexpectedAuthorFmt        = "unexpected author: %+v"
	errExpectedErrorNil           = "expected error, got nil"
	errExpectedNoSuchAuthorFmt    = "expected ErrNoSuchAuthor, got %v"
	errExpectedGenericErrorFmt    = "expected generic error, got %v"
	queryAuthorByID               = "SELECT Id, Name FROM Authors WHERE Id = $1"
	queryAuthorByName             = "SELECT Id, Name FROM Authors WHERE Name = $1"
	queryUpdateAuthor             = "UPDATE Authors SET Name = $1 WHERE Id = $2 RETURNING Id, Name"
	queryDeleteAuthorToBookByID   = "DELETE FROM AuthorToBook WHERE AuthorId = $1"
	queryExistsAuthorByID         = "SELECT EXISTS(SELECT 1 FROM Authors WHERE Id = $1)"
)

// ── Create ────────────────────────────────────────────────────────────────────

func TestAuthorCreateSuccess(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO Authors (Name) VALUES ($1) RETURNING Id, Name")).
		WithArgs(authorOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", authorOneName))

	author, err := repo.Create(authorOneName)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if author.Id != "a1" || author.Name != authorOneName {
		t.Errorf(errUnexpectedAuthorFmt, author)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestAuthorCreateError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO Authors (Name) VALUES ($1) RETURNING Id, Name")).
		WithArgs(authorOneName).
		WillReturnError(errors.New("insert error"))

	_, err := repo.Create(authorOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestAuthorGetSuccess(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryAuthorByID)).
		WithArgs("a1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", authorOneName))

	author, err := repo.Get("a1")
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if author.Id != "a1" {
		t.Errorf(errUnexpectedAuthorFmt, author)
	}
}

func TestAuthorGetNotFound(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryAuthorByID)).
		WithArgs("a1").
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.Get("a1")
	if !errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf(errExpectedNoSuchAuthorFmt, err)
	}
}

func TestAuthorGetError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryAuthorByID)).
		WithArgs("a1").
		WillReturnError(errors.New(dbErrorMessage))

	_, err := repo.Get("a1")
	if err == nil || errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf(errExpectedGenericErrorFmt, err)
	}
}

// ── GetByName ─────────────────────────────────────────────────────────────────

func TestAuthorGetByNameSuccess(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryAuthorByName)).
		WithArgs(authorOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", authorOneName))

	author, err := repo.GetByName(authorOneName)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if author.Name != authorOneName {
		t.Errorf(errUnexpectedAuthorFmt, author)
	}
}

func TestAuthorGetByNameNotFound(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryAuthorByName)).
		WithArgs(authorOneName).
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.GetByName(authorOneName)
	if !errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf(errExpectedNoSuchAuthorFmt, err)
	}
}

func TestAuthorGetByNameError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryAuthorByName)).
		WithArgs(authorOneName).
		WillReturnError(errors.New(dbErrorMessage))

	_, err := repo.GetByName(authorOneName)
	if err == nil || errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf(errExpectedGenericErrorFmt, err)
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestAuthorUpdateSuccess(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAuthor)).
		WithArgs(updatedAuthorName, "a1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", updatedAuthorName))

	author, err := repo.Update("a1", updatedAuthorName)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if author.Name != updatedAuthorName {
		t.Errorf("unexpected author name: %v", author.Name)
	}
}

func TestAuthorUpdateNotFound(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAuthor)).
		WithArgs(updatedAuthorName, "a1").
		WillReturnError(pgx.ErrNoRows)

	_, err := repo.Update("a1", updatedAuthorName)
	if !errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf(errExpectedNoSuchAuthorFmt, err)
	}
}

func TestAuthorUpdateError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateAuthor)).
		WithArgs(updatedAuthorName, "a1").
		WillReturnError(errors.New(dbErrorMessage))

	_, err := repo.Update("a1", updatedAuthorName)
	if err == nil || errors.Is(err, ErrNoSuchAuthor) {
		t.Errorf(errExpectedGenericErrorFmt, err)
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestAuthorDeleteSuccess(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("a1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Authors WHERE Id = $1")).
		WithArgs("a1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	if err := repo.Delete("a1"); err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestAuthorDeleteBeginError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New("begin error"))

	if err := repo.Delete("a1"); err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestAuthorDeleteDeleteAuthorToBookError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("a1").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	if err := repo.Delete("a1"); err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestAuthorDeleteDeleteAuthorError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("a1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM Authors WHERE Id = $1")).
		WithArgs("a1").
		WillReturnError(errors.New("delete error"))
	mock.ExpectRollback()

	if err := repo.Delete("a1"); err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── ExistsById ────────────────────────────────────────────────────────────────

func TestAuthorExistsByIdTrue(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsAuthorByID)).
		WithArgs("a1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsById("a1")
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestAuthorExistsByIdFalse(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsAuthorByID)).
		WithArgs("a1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := repo.ExistsById("a1")
	if err != nil || exists {
		t.Errorf("expected false, got %v, %v", exists, err)
	}
}

func TestAuthorExistsByIdError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsAuthorByID)).
		WithArgs("a1").
		WillReturnError(errors.New(dbErrorMessage))

	_, err := repo.ExistsById("a1")
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── ExistsByName ──────────────────────────────────────────────────────────────

func TestAuthorExistsByNameTrue(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Authors WHERE Name = $1)")).
		WithArgs(authorOneName).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByName(authorOneName)
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestAuthorExistsByNameError(t *testing.T) {
	repo, mock := newMockAuthorRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta("SELECT EXISTS(SELECT 1 FROM Authors WHERE Name = $1)")).
		WithArgs(authorOneName).
		WillReturnError(errors.New(dbErrorMessage))

	_, err := repo.ExistsByName(authorOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}
