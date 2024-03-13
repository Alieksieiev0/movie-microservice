package main

import (
	"context"
	"strings"
	"testing"

	"github.com/gofrs/uuid/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/shopspring/decimal"
)

func TestGet(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	mdb := NewMovieDatabase(mock)

	rows := pgxmock.NewRows(testMovieColumn()).AddRow(testMovieRow(id)...)
	mock.ExpectQuery("select *").WithArgs(id.String()).WillReturnRows(rows)
	if _, err := mdb.Get(context.Background(), id.String()); err != nil {
		t.Errorf("error was not expected while querying: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestGetAll(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	mdb := NewMovieDatabase(mock)

	rows := pgxmock.NewRows(testMovieColumn())
	values := [][]any{}
	for i := 0; i < 10; i++ {
		values = append(values, testMovieRow(testUUID(t)))
	}
	rows.AddRows(values...)

	mock.ExpectQuery("select *").WillReturnRows(rows)
	movies, err := mdb.GetAll(context.Background())
	if err != nil {
		t.Errorf("error was not expected while querying: %s", err)
	}

	if len(movies) != 10 {
		t.Fatal("incorrect number of entities was returned")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestInsert(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	mdb := NewMovieDatabase(mock)

	rows := mock.NewRows([]string{"id"}).AddRow(id)
	value := testMovieRow(id)[1:]
	mock.ExpectBegin()
	mock.ExpectQuery("insert into").WithArgs(value...).WillReturnRows(rows)
	mock.ExpectCommit()

	if _, err := mdb.Insert(context.Background(), testMovie()); err != nil {
		t.Errorf("error was not expected while inserting: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestUpdate(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	mdb := NewMovieDatabase(mock)

	movie := &Movie{Name: "updateTest"}
	mock.ExpectBegin()
	mock.ExpectExec("update movie").
		WithArgs(movie.Name, id.String()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()

	if err := mdb.Update(context.Background(), id.String(), movie); err != nil {
		t.Errorf("error was not expected while updating: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestBuildUpdateQuery(t *testing.T) {
	id := testUUID(t)
	movie := testMovie()
	q, params, err := MovieDatabase{}.buildUpdateQuery(movie, id.String())
	if err != nil {
		t.Errorf("error during building update query: %v", err)
	}
	if !strings.Contains(q, "name = $1") && params[0] != "test" {
		t.Errorf("name was not appended to update query")
	}

	if !strings.Contains(q, "release_year = $2") && params[1] != 2024 {
		t.Errorf("release year was not appended to update query")
	}

	if !strings.Contains(q, "rating = $3") && params[2] != decimal.NewFromInt(10) {
		t.Errorf("rating year was not appended to update query")
	}

	if !strings.Contains(q, "genres = $4") && params[3] == nil {
		t.Errorf("rating year was not appended to update query")
	}

	if !strings.Contains(q, "director = $5") && params[4] != "someone" {
		t.Errorf("someone year was not appended to update query")
	}
}

func TestDelete(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	mdb := NewMovieDatabase(mock)

	mock.ExpectBegin()
	mock.ExpectExec("delete from movie").
		WithArgs(id.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()

	if err := mdb.Delete(context.Background(), id.String()); err != nil {
		t.Errorf("error was not expected while deleting: %s", err)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func testPoolMock(t *testing.T) pgxmock.PgxPoolIface {
	mock, err := pgxmock.NewPool()
	if err != nil {
		t.Fatal(err)
	}
	return mock
}

func testUUID(t *testing.T) uuid.UUID {
	id, err := uuid.NewV4()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func testMovieColumn() []string {
	return []string{"id", "name", "release_year", "rating", "genres", "director"}
}

func testMovieRow(id uuid.UUID) []any {
	return []any{id, "test", 2024, decimal.NewFromInt(10), []string{"test"}, "someone"}
}

func testMovie() *Movie {
	return &Movie{
		Name:        "test",
		ReleaseYear: 2024,
		Rating:      decimal.NewFromInt(10),
		Genres:      []string{"test"},
		Director:    "someone",
	}
}
