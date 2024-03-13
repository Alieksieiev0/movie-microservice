package main

import (
	"context"
	"encoding/json"
	"net/http"
)

type Server struct {
	svc Service
}

func (s Server) Start(addr string) error {
	http.HandleFunc("GET /movies/{id}", s.handleGetMovie)
	http.HandleFunc("GET /movies", s.handleGetAllMovies)
	http.HandleFunc("POST /movies", s.handleCreateMovie)
	http.HandleFunc("PUT /movies/{id}", s.handleUpdateMovie)
	http.HandleFunc("DELETE /movies/{id}", s.handleDeleteMovie)
	return http.ListenAndServe(addr, nil)
}

func (s Server) handleGetMovie(w http.ResponseWriter, r *http.Request) {
	movie, err := s.svc.GetMovie(context.Background(), r.PathValue("id"))
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, movie)
}

func (s Server) handleGetAllMovies(w http.ResponseWriter, _ *http.Request) {
	movies, err := s.svc.GetAllMovies(context.Background())
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, movies)
}

func (s Server) handleCreateMovie(w http.ResponseWriter, r *http.Request) {
	movie := &Movie{}
	err := json.NewDecoder(r.Body).Decode(movie)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	id, err := s.svc.CreateMovie(context.Background(), movie)
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusCreated, map[string]any{"id": id})
}

func (s Server) handleUpdateMovie(w http.ResponseWriter, r *http.Request) {
	movie := &Movie{}
	err := json.NewDecoder(r.Body).Decode(movie)
	if err != nil {
		writeJson(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
		return
	}

	err = s.svc.UpdateMovie(context.Background(), r.PathValue("id"), movie)
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusNoContent, map[string]any{})
}

func (s Server) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	err := s.svc.DeleteMovie(context.Background(), r.PathValue("id"))
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusNoContent, map[string]any{})
}

func writeJson(w http.ResponseWriter, s int, v any) {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
