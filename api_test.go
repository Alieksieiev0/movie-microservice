package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pashagolub/pgxmock/v3"
)

func TestHandleGetMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/movies/%v", id), nil)
	r.SetPathValue("id", id.String())
	s := NewServer(NewMovieService(NewMovieDatabase(mock)))

	rows := pgxmock.NewRows(testMovieColumn()).AddRow(testMovieRow(id)...)
	mock.ExpectQuery("select *").WithArgs(id.String()).WillReturnRows(rows)
	s.handleGetMovie(w, r)

	movie := &Movie{}
	err := json.NewDecoder(w.Body).Decode(movie)
	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}

	if movie.Id != id {
		t.Errorf("wrong fields returned in response; expected: %v, got: %v", movie.Id, id)
	}
}

func TestHandleGetAllMovies(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/movies", nil)
	s := NewServer(NewMovieService(NewMovieDatabase(mock)))

	rows := pgxmock.NewRows(testMovieColumn())
	values := [][]any{}
	for i := 0; i < 10; i++ {
		values = append(values, testMovieRow(testUUID(t)))
	}
	rows.AddRows(values...)

	mock.ExpectQuery("select *").WillReturnRows(rows)
	s.handleGetAllMovies(w, r)
	movies := []Movie{}
	err := json.NewDecoder(w.Body).Decode(&movies)
	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}

	if len(movies) != 10 {
		t.Errorf(
			"wrong number of movies returned in response; expected: 10, got: %d",
			len(movies),
		)
	}
}

func TestHandleCreateMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)

	movie := testMovie()
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(movie)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPost, "/movies", &buf)
	s := NewServer(NewMovieService(NewMovieDatabase(mock)))

	rows := mock.NewRows([]string{"id"}).AddRow(id)
	value := testMovieRow(id)[1:]
	mock.ExpectBegin()
	mock.ExpectQuery("insert into").WithArgs(value...).WillReturnRows(rows)
	mock.ExpectCommit()
	s.handleCreateMovie(w, r)

	res := map[string]string{}
	err = json.NewDecoder(w.Body).Decode(&res)
	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}

	if resId, ok := res["id"]; !ok || resId != id.String() {
		t.Errorf("wrong id returned in response; expected: %s, got: %s", id.String(), resId)
	}
}

func TestHandleUpdateMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)

	movie := &Movie{Name: "updateTest"}
	buf := bytes.Buffer{}
	err := json.NewEncoder(&buf).Encode(movie)
	if err != nil {
		t.Fatal(err)
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/movies/%v", id), &buf)
	r.SetPathValue("id", id.String())
	s := NewServer(NewMovieService(NewMovieDatabase(mock)))

	mock.ExpectBegin()
	mock.ExpectExec("update movie").
		WithArgs(movie.Name, id.String()).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))
	mock.ExpectCommit()
	s.handleUpdateMovie(w, r)

	if w.Result().StatusCode != http.StatusNoContent {
		res := map[string]string{}
		err := json.NewDecoder(w.Body).Decode(&res)
		if err != nil {
			t.Errorf("error reading response body: %v", err)
		}
		t.Errorf("error returned in response: %v", res["error"])
	}
}

func TestHandleDeleteMovie(t *testing.T) {
	mock := testPoolMock(t)
	defer mock.Close()
	id := testUUID(t)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/movies/%v", id), nil)
	r.SetPathValue("id", id.String())
	s := NewServer(NewMovieService(NewMovieDatabase(mock)))

	mock.ExpectBegin()
	mock.ExpectExec("delete from movie").
		WithArgs(id.String()).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))
	mock.ExpectCommit()
	s.handleDeleteMovie(w, r)

	if w.Result().StatusCode != http.StatusNoContent {
		res := map[string]string{}
		err := json.NewDecoder(w.Body).Decode(&res)
		if err != nil {
			t.Errorf("error reading response body: %v", err)
		}
		t.Errorf("error returned in response: %v", res["error"])
	}
}

func TestWriteJson(t *testing.T) {
	w := httptest.NewRecorder()
	writeJson(w, http.StatusOK, map[string]string{"status": "ok"})

	if w.Result().StatusCode != http.StatusOK {
		t.Errorf(
			"error status code returned in response; expected: %d, got: %d",
			http.StatusOK,
			w.Result().StatusCode,
		)
	}
	res := map[string]string{}
	err := json.NewDecoder(w.Body).Decode(&res)
	if err != nil {
		t.Errorf("error reading response body: %v", err)
	}
	if resStatus, ok := res["status"]; !ok || resStatus != "ok" {
		t.Errorf("wrong status value returned in response; expected: ok, got: %s", resStatus)
	}
}
