package main

import (
	"context"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

func TestGetMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	ms := NewMovieService(NewMovieDatabase(mock))

	rows := pgxmock.NewRows(testMovieColumn()).AddRow(testMovieRow(id)...)
	mock.ExpectQuery("select *").WithArgs(id.String()).WillReturnRows(rows)
	movie, err := ms.GetMovie(context.Background(), id.String())
	if err != nil {
		t.Errorf("error fetching: %v", err)
	}

	if movie.Id != id {
		t.Errorf("wrong fields returned; expected: %v, got: %v", movie.Id, id)
	}
}

func TestGetAllMovies(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	ms := NewMovieService(NewMovieDatabase(mock))

	rows := pgxmock.NewRows(testMovieColumn())
	values := [][]any{}
	for i := 0; i < 10; i++ {
		values = append(values, testMovieRow(testUUID(t)))
	}
	rows.AddRows(values...)

	mock.ExpectQuery("select *").WillReturnRows(rows)
	movies, err := ms.GetAllMovies(context.Background())
	if err != nil {
		t.Errorf("error fetching: %v", err)
	}

	if len(movies) != 10 {
		t.Errorf(
			"wrong number of movies returned; expected: 10, got: %d",
			len(movies),
		)
	}
}

func TestCreateMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	movie := testMovie()
	ms := NewMovieService(NewMovieDatabase(mock))

	rows := mock.NewRows([]string{"id"}).AddRow(id)
	value := testMovieRow(id)[1:]
	mock.ExpectBegin()
	mock.ExpectQuery("insert into").WithArgs(value...).WillReturnRows(rows)
	mock.ExpectCommit()
	resId, err := ms.CreateMovie(context.Background(), movie)
	if err != nil {
		t.Errorf("error creating: %v", err)
	}

	if resId != id.String() {
		t.Errorf("wrong id returned; expected: %s, got: %s", id.String(), resId)
	}
}

func TestUpdateMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	movie := &Movie{Name: "updateTest"}
	ms := NewMovieService(NewMovieDatabase(mock))

	mock.ExpectBegin()
	mock.ExpectExec("update movie").
		WithArgs(movie.Name, id.String()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	err := ms.UpdateMovie(context.Background(), id.String(), movie)
	if err != nil {
		t.Errorf("error updating: %v", err)
	}
}

func TestDeleteMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)
	s := NewMovieService(NewMovieDatabase(mock))

	mock.ExpectBegin()
	mock.ExpectExec("delete from movie").
		WithArgs(id.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	err := s.DeleteMovie(context.Background(), id.String())
	if err != nil {
		t.Errorf("error deleting: %v", err)
	}
}
