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

const (
	authorOneName               = "Author One"
	bookOneName                 = "Book One"
	updatedBookName             = "New Name"
	beginErrorMessage           = "begin error"
	commitErrorMessage          = "commit error"
	queryErrorMessage           = "query error"
	deleteErrorMessage          = "delete error"
	errUnexpectedFmt            = "unexpected error: %v"
	errUnexpectedBookFmt        = "unexpected book: %+v"
	errExpectedErrorNil         = "expected error, got nil"
	queryInsertBook             = "INSERT INTO Books (Name) VALUES ($1) RETURNING Id, Name"
	queryInsertAuthorToBook     = "INSERT INTO AuthorToBook (AuthorId, BookId) VALUES ($1, $2)"
	queryBookByID               = "SELECT Id, Name FROM Books WHERE Id = $1"
	queryBookByName             = "SELECT Id, Name FROM Books WHERE Name = $1"
	queryBookAuthorsByBookID    = "SELECT s.Id, s.Name FROM AuthorToBook AS f LEFT JOIN Authors AS s ON f.AuthorId = s.Id WHERE f.BookId = $1"
	queryDeleteAuthorToBookByID = "DELETE FROM AuthorToBook WHERE BookId = $1"
	queryDeleteBookByID         = "DELETE FROM Books WHERE Id = $1"
	queryUpdateBookByID         = "UPDATE Books SET Name = $1 WHERE Id = $2 RETURNING Id, Name"
	queryExistsBookByID         = "SELECT EXISTS(SELECT FROM Books WHERE id = $1)"
	queryExistsBookByName       = "SELECT EXISTS(SELECT FROM Books WHERE name = $1)"
)

var testAuthor = model.Author{Id: "a1", Name: authorOneName}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCreateSuccess(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryInsertBook)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectExec(regexp.QuoteMeta(queryInsertAuthorToBook)).
		WithArgs("a1", "b1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	book, err := repo.Create(bookOneName, testAuthor)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if book.Id != "b1" || book.Name != bookOneName {
		t.Errorf(errUnexpectedBookFmt, book)
	}
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
	}
}

func TestCreateBeginError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))

	_, err := repo.Create(bookOneName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestCreateInsertBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryInsertBook)).
		WithArgs(bookOneName).
		WillReturnError(errors.New("insert error"))
	mock.ExpectRollback()

	_, err := repo.Create(bookOneName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestCreateInsertAuthorError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryInsertBook)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectExec(regexp.QuoteMeta(queryInsertAuthorToBook)).
		WithArgs("a1", "b1").
		WillReturnError(errors.New("author insert error"))
	mock.ExpectRollback()

	_, err := repo.Create(bookOneName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── Get ───────────────────────────────────────────────────────────────────────

func TestGetSuccess(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", authorOneName))
	mock.ExpectCommit()

	book, err := repo.Get("b1")
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if book.Id != "b1" || len(book.Authors) != 1 {
		t.Errorf(errUnexpectedBookFmt, book)
	}
}

func TestGetBeginError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetQueryBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByID)).
		WithArgs("b1").
		WillReturnError(errors.New(queryErrorMessage))
	mock.ExpectRollback()

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetQueryAuthorsError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnError(errors.New("authors query error"))
	mock.ExpectRollback()

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetCommitError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}))
	mock.ExpectCommit().WillReturnError(errors.New(commitErrorMessage))

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetAuthorScanError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	// One column returned, but repository scans two fields.
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id"}).AddRow("a1"))
	mock.ExpectRollback()

	_, err := repo.Get("b1")
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── GetByName ─────────────────────────────────────────────────────────────────

func TestGetByNameSuccess(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByName)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", authorOneName))
	mock.ExpectCommit()

	book, err := repo.GetByName(bookOneName)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if book.Name != bookOneName {
		t.Errorf(errUnexpectedBookFmt, book)
	}
}

func TestGetByNameQueryError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByName)).
		WithArgs(bookOneName).
		WillReturnError(errors.New(queryErrorMessage))
	mock.ExpectRollback()

	_, err := repo.GetByName(bookOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetByNameBeginError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))

	_, err := repo.GetByName(bookOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetByNameQueryAuthorsError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByName)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnError(errors.New("authors query error"))
	mock.ExpectRollback()

	_, err := repo.GetByName(bookOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetByNameAuthorScanError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByName)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id"}).AddRow("a1"))
	mock.ExpectRollback()

	_, err := repo.GetByName(bookOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestGetByNameCommitError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryBookByName)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
	mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}))
	mock.ExpectCommit().WillReturnError(errors.New(commitErrorMessage))

	_, err := repo.GetByName(bookOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDeleteSuccess(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	if err := repo.Delete("b1"); err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
}

func TestDeleteDeleteAuthorToBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnError(errors.New(deleteErrorMessage))
	mock.ExpectRollback()

	if err := repo.Delete("b1"); err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestDeleteDeleteBookError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteBookByID)).
		WithArgs("b1").
		WillReturnError(errors.New(deleteErrorMessage))
	mock.ExpectRollback()

	if err := repo.Delete("b1"); err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestDeleteBeginError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))

	if err := repo.Delete("b1"); err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestDeleteCommitError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit().WillReturnError(errors.New(commitErrorMessage))

	if err := repo.Delete("b1"); err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdateSuccess(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBookByID)).
		WithArgs(updatedBookName, "b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", updatedBookName))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryInsertAuthorToBook)).
		WithArgs("a1", "b1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit()

	book, err := repo.Update("b1", updatedBookName, testAuthor)
	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}
	if book.Name != updatedBookName {
		t.Errorf("unexpected book name: %v", book.Name)
	}
}

func TestUpdateUpdateError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBookByID)).
		WithArgs(updatedBookName, "b1").
		WillReturnError(errors.New("update error"))
	mock.ExpectRollback()

	_, err := repo.Update("b1", updatedBookName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestUpdateDeleteAuthorsError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBookByID)).
		WithArgs(updatedBookName, "b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", updatedBookName))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnError(errors.New(deleteErrorMessage))
	mock.ExpectRollback()

	_, err := repo.Update("b1", updatedBookName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestUpdateBeginError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))

	_, err := repo.Update("b1", updatedBookName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestUpdateInsertAuthorError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBookByID)).
		WithArgs(updatedBookName, "b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", updatedBookName))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryInsertAuthorToBook)).
		WithArgs("a1", "b1").
		WillReturnError(errors.New("author insert error"))
	mock.ExpectRollback()

	_, err := repo.Update("b1", updatedBookName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

func TestUpdateCommitError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBookByID)).
		WithArgs(updatedBookName, "b1").
		WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", updatedBookName))
	mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
		WithArgs("b1").
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectExec(regexp.QuoteMeta(queryInsertAuthorToBook)).
		WithArgs("a1", "b1").
		WillReturnResult(pgxmock.NewResult("INSERT", 1))
	mock.ExpectCommit().WillReturnError(errors.New(commitErrorMessage))

	_, err := repo.Update("b1", updatedBookName, testAuthor)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── ExistsById ────────────────────────────────────────────────────────────────

func TestExistsByIdTrue(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsBookByID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsById("b1")
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestExistsByIdFalse(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsBookByID)).
		WithArgs("b1").
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := repo.ExistsById("b1")
	if err != nil || exists {
		t.Errorf("expected false, got %v, %v", exists, err)
	}
}

func TestExistsByIdError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsBookByID)).
		WithArgs("b1").
		WillReturnError(errors.New(queryErrorMessage))

	_, err := repo.ExistsById("b1")
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}

// ── ExistsByName ──────────────────────────────────────────────────────────────

func TestExistsByNameTrue(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsBookByName)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

	exists, err := repo.ExistsByName(bookOneName)
	if err != nil || !exists {
		t.Errorf("expected true, got %v, %v", exists, err)
	}
}

func TestExistsByNameFalse(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsBookByName)).
		WithArgs(bookOneName).
		WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

	exists, err := repo.ExistsByName(bookOneName)
	if err != nil || exists {
		t.Errorf("expected false, got %v, %v", exists, err)
	}
}

func TestExistsByNameError(t *testing.T) {
	repo, mock := newMockRepo(t)

	mock.ExpectQuery(regexp.QuoteMeta(queryExistsBookByName)).
		WithArgs(bookOneName).
		WillReturnError(errors.New(queryErrorMessage))

	_, err := repo.ExistsByName(bookOneName)
	if err == nil {
		t.Fatal(errExpectedErrorNil)
	}
}
