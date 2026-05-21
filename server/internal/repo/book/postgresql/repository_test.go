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

func TestGetContract(t *testing.T) {
	runFetchBookContract(
		t,
		queryBookByID,
		"b1",
		func(repo *BookPostgresRepository) (*model.Book, error) { return repo.Get("b1") },
		func(t *testing.T, book *model.Book) {
			t.Helper()
			if book.Id != "b1" || len(book.Authors) != 1 {
				t.Errorf(errUnexpectedBookFmt, book)
			}
		},
	)
}

// ── GetByName ─────────────────────────────────────────────────────────────────

func TestGetByNameContract(t *testing.T) {
	runFetchBookContract(
		t,
		queryBookByName,
		bookOneName,
		func(repo *BookPostgresRepository) (*model.Book, error) { return repo.GetByName(bookOneName) },
		func(t *testing.T, book *model.Book) {
			t.Helper()
			if book.Name != bookOneName {
				t.Errorf(errUnexpectedBookFmt, book)
			}
		},
	)
}

func runFetchBookContract(
	t *testing.T,
	bookQuery string,
	lookupValue string,
	fetch func(repo *BookPostgresRepository) (*model.Book, error),
	assertBook func(t *testing.T, book *model.Book),
) {
	t.Helper()

	t.Run("success", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
			WithArgs(lookupValue).
			WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
		mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
			WithArgs("b1").
			WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", authorOneName))
		mock.ExpectCommit()

		book, err := fetch(repo)
		if err != nil {
			t.Fatalf(errUnexpectedFmt, err)
		}
		assertBook(t, book)
	})

	t.Run("begin error", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))

		_, err := fetch(repo)
		if err == nil {
			t.Fatal(errExpectedErrorNil)
		}
	})

	t.Run("book query error", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
			WithArgs(lookupValue).
			WillReturnError(errors.New(queryErrorMessage))
		mock.ExpectRollback()

		_, err := fetch(repo)
		if err == nil {
			t.Fatal(errExpectedErrorNil)
		}
	})

	t.Run("authors query error", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
			WithArgs(lookupValue).
			WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
		mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
			WithArgs("b1").
			WillReturnError(errors.New("authors query error"))
		mock.ExpectRollback()

		_, err := fetch(repo)
		if err == nil {
			t.Fatal(errExpectedErrorNil)
		}
	})

	t.Run("author scan error", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
			WithArgs(lookupValue).
			WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
		mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
			WithArgs("b1").
			WillReturnRows(pgxmock.NewRows([]string{"Id"}).AddRow("a1"))
		mock.ExpectRollback()

		_, err := fetch(repo)
		if err == nil {
			t.Fatal(errExpectedErrorNil)
		}
	})

	t.Run("commit error", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
			WithArgs(lookupValue).
			WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
		mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
			WithArgs("b1").
			WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}))
		mock.ExpectCommit().WillReturnError(errors.New(commitErrorMessage))

		_, err := fetch(repo)
		if err == nil {
			t.Fatal(errExpectedErrorNil)
		}
	})
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

func TestExistsByIDContract(t *testing.T) {
	runExistsContract(
		t,
		queryExistsBookByID,
		"b1",
		func(repo *BookPostgresRepository) (bool, error) { return repo.ExistsById("b1") },
	)
}

// ── ExistsByName ──────────────────────────────────────────────────────────────

func TestExistsByNameContract(t *testing.T) {
	runExistsContract(
		t,
		queryExistsBookByName,
		bookOneName,
		func(repo *BookPostgresRepository) (bool, error) { return repo.ExistsByName(bookOneName) },
	)
}

func runExistsContract(
	t *testing.T,
	query string,
	arg string,
	existsCheck func(repo *BookPostgresRepository) (bool, error),
) {
	t.Helper()

	t.Run("returns true", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(arg).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))

		exists, err := existsCheck(repo)
		if err != nil || !exists {
			t.Errorf("expected true, got %v, %v", exists, err)
		}
	})

	t.Run("returns false", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(arg).
			WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))

		exists, err := existsCheck(repo)
		if err != nil || exists {
			t.Errorf("expected false, got %v, %v", exists, err)
		}
	})

	t.Run("returns error", func(t *testing.T) {
		repo, mock := newMockRepo(t)

		mock.ExpectQuery(regexp.QuoteMeta(query)).
			WithArgs(arg).
			WillReturnError(errors.New(queryErrorMessage))

		_, err := existsCheck(repo)
		if err == nil {
			t.Fatal(errExpectedErrorNil)
		}
	})
}
