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

type bookScenario struct {
	name       string
	prepare    func(mock pgxmock.PgxPoolIface)
	run        func(repo *BookPostgresRepository) (*model.Book, error)
	assertBook func(t *testing.T, book *model.Book)
	wantErr    bool
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestCreateContract(t *testing.T) {
	runBookScenarioContract(t, []bookScenario{
		{
			name: "success",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertBook)).
					WithArgs(bookOneName).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
				mock.ExpectExec(regexp.QuoteMeta(queryInsertAuthorToBook)).
					WithArgs("a1", "b1").
					WillReturnResult(pgxmock.NewResult("INSERT", 1))
				mock.ExpectCommit()
			},
			assertBook: func(t *testing.T, book *model.Book) {
				t.Helper()
				if book.Id != "b1" || book.Name != bookOneName {
					t.Errorf(errUnexpectedBookFmt, book)
				}
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Create(bookOneName, testAuthor)
			},
		},
		{
			name: beginErrorMessage,
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Create(bookOneName, testAuthor)
			},
			wantErr: true,
		},
		{
			name: "insert book error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertBook)).
					WithArgs(bookOneName).
					WillReturnError(errors.New("insert error"))
				mock.ExpectRollback()
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Create(bookOneName, testAuthor)
			},
			wantErr: true,
		},
		{
			name: "insert author error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryInsertBook)).
					WithArgs(bookOneName).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
				mock.ExpectExec(regexp.QuoteMeta(queryInsertAuthorToBook)).
					WithArgs("a1", "b1").
					WillReturnError(errors.New("author insert error"))
				mock.ExpectRollback()
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Create(bookOneName, testAuthor)
			},
			wantErr: true,
		},
	})
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

	type fetchScenario struct {
		name       string
		prepare    func(mock pgxmock.PgxPoolIface)
		assertBook func(t *testing.T, book *model.Book)
		wantErr    bool
	}

	scenarios := []fetchScenario{
		{
			name: "success",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
					WithArgs(lookupValue).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
				mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
					WithArgs("b1").
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("a1", authorOneName))
				mock.ExpectCommit()
			},
			assertBook: assertBook,
		},
		{
			name: beginErrorMessage,
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))
			},
			wantErr: true,
		},
		{
			name: "book query error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
					WithArgs(lookupValue).
					WillReturnError(errors.New(queryErrorMessage))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "authors query error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
					WithArgs(lookupValue).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
				mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
					WithArgs("b1").
					WillReturnError(errors.New("authors query error"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "author scan error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
					WithArgs(lookupValue).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
				mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
					WithArgs("b1").
					WillReturnRows(pgxmock.NewRows([]string{"Id"}).AddRow("a1"))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: commitErrorMessage,
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(bookQuery)).
					WithArgs(lookupValue).
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", bookOneName))
				mock.ExpectQuery(regexp.QuoteMeta(queryBookAuthorsByBookID)).
					WithArgs("b1").
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}))
				mock.ExpectCommit().WillReturnError(errors.New(commitErrorMessage))
			},
			wantErr: true,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			repo, mock := newMockRepo(t)
			scenario.prepare(mock)

			book, err := fetch(repo)

			if scenario.wantErr {
				if err == nil {
					t.Fatal(errExpectedErrorNil)
				}
				return
			}

			if err != nil {
				t.Fatalf(errUnexpectedFmt, err)
			}
			if scenario.assertBook != nil {
				scenario.assertBook(t, book)
			}
		})
	}
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestDeleteContract(t *testing.T) {
	type deleteScenario struct {
		name    string
		prepare func(mock pgxmock.PgxPoolIface)
		wantErr bool
	}

	scenarios := []deleteScenario{
		{
			name: "success",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
					WithArgs("b1").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteBookByID)).
					WithArgs("b1").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectCommit()
			},
		},
		{
			name: "delete author to book error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
					WithArgs("b1").
					WillReturnError(errors.New(deleteErrorMessage))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "delete book error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
					WithArgs("b1").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteBookByID)).
					WithArgs("b1").
					WillReturnError(errors.New(deleteErrorMessage))
				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: beginErrorMessage,
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))
			},
			wantErr: true,
		},
		{
			name: commitErrorMessage,
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
					WithArgs("b1").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteBookByID)).
					WithArgs("b1").
					WillReturnResult(pgxmock.NewResult("DELETE", 1))
				mock.ExpectCommit().WillReturnError(errors.New(commitErrorMessage))
			},
			wantErr: true,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			repo, mock := newMockRepo(t)
			scenario.prepare(mock)

			err := repo.Delete("b1")
			if scenario.wantErr {
				if err == nil {
					t.Fatal(errExpectedErrorNil)
				}
				return
			}
			if err != nil {
				t.Fatalf(errUnexpectedFmt, err)
			}
		})
	}
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestUpdateContract(t *testing.T) {
	runBookScenarioContract(t, []bookScenario{
		{
			name: "success",
			prepare: func(mock pgxmock.PgxPoolIface) {
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
			},
			assertBook: func(t *testing.T, book *model.Book) {
				t.Helper()
				if book.Name != updatedBookName {
					t.Errorf("unexpected book name: %v", book.Name)
				}
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Update("b1", updatedBookName, testAuthor)
			},
		},
		{
			name: "update error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBookByID)).
					WithArgs(updatedBookName, "b1").
					WillReturnError(errors.New("update error"))
				mock.ExpectRollback()
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Update("b1", updatedBookName, testAuthor)
			},
			wantErr: true,
		},
		{
			name: "delete authors error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin()
				mock.ExpectQuery(regexp.QuoteMeta(queryUpdateBookByID)).
					WithArgs(updatedBookName, "b1").
					WillReturnRows(pgxmock.NewRows([]string{"Id", "Name"}).AddRow("b1", updatedBookName))
				mock.ExpectExec(regexp.QuoteMeta(queryDeleteAuthorToBookByID)).
					WithArgs("b1").
					WillReturnError(errors.New(deleteErrorMessage))
				mock.ExpectRollback()
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Update("b1", updatedBookName, testAuthor)
			},
			wantErr: true,
		},
		{
			name: beginErrorMessage,
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectBegin().WillReturnError(errors.New(beginErrorMessage))
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Update("b1", updatedBookName, testAuthor)
			},
			wantErr: true,
		},
		{
			name: "insert author error",
			prepare: func(mock pgxmock.PgxPoolIface) {
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
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Update("b1", updatedBookName, testAuthor)
			},
			wantErr: true,
		},
		{
			name: commitErrorMessage,
			prepare: func(mock pgxmock.PgxPoolIface) {
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
			},
			run: func(repo *BookPostgresRepository) (*model.Book, error) {
				return repo.Update("b1", updatedBookName, testAuthor)
			},
			wantErr: true,
		},
	})
}

func runBookScenarioContract(t *testing.T, scenarios []bookScenario) {
	t.Helper()

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			repo, mock := newMockRepo(t)
			scenario.prepare(mock)

			book, err := scenario.run(repo)
			if !assertBookScenarioResult(t, scenario, book, err) {
				return
			}

			assertMockExpectations(t, mock)
		})
	}
}

func assertBookScenarioResult(
	t *testing.T,
	scenario bookScenario,
	book *model.Book,
	err error,
) bool {
	t.Helper()

	if scenario.wantErr {
		if err == nil {
			t.Fatal(errExpectedErrorNil)
		}
		return false
	}

	if err != nil {
		t.Fatalf(errUnexpectedFmt, err)
	}

	if scenario.assertBook != nil {
		scenario.assertBook(t, book)
	}

	return true
}

func assertMockExpectations(t *testing.T, mock pgxmock.PgxPoolIface) {
	t.Helper()

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unmet expectations: %v", err)
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

	type existsScenario struct {
		name     string
		prepare  func(mock pgxmock.PgxPoolIface)
		wantErr  bool
		wantBool bool
	}

	scenarios := []existsScenario{
		{
			name: "returns true",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(arg).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(true))
			},
			wantBool: true,
		},
		{
			name: "returns false",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(arg).
					WillReturnRows(pgxmock.NewRows([]string{"exists"}).AddRow(false))
			},
		},
		{
			name: "returns error",
			prepare: func(mock pgxmock.PgxPoolIface) {
				mock.ExpectQuery(regexp.QuoteMeta(query)).
					WithArgs(arg).
					WillReturnError(errors.New(queryErrorMessage))
			},
			wantErr: true,
		},
	}

	for _, scenario := range scenarios {
		scenario := scenario
		t.Run(scenario.name, func(t *testing.T) {
			repo, mock := newMockRepo(t)
			scenario.prepare(mock)

			exists, err := existsCheck(repo)
			if scenario.wantErr {
				if err == nil {
					t.Fatal(errExpectedErrorNil)
				}
				return
			}
			if err != nil || exists != scenario.wantBool {
				t.Errorf("expected %v, got %v, %v", scenario.wantBool, exists, err)
			}
		})
	}
}
