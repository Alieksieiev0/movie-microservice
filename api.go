package main

import (
	"context"
	"encoding/json"
	"net/http"
)

// Server contains handlers for all supportend endpoints, registers handlers and starts server.
type Server struct {
	// Supplied service is used to perform the appropriate operation for each endpoint.
	svc Service
}

// NewServer creates an instance of the Server.
func NewServer(svc Service) Server {
	return Server{
		svc: svc,
	}
}

// Start registers all handlers and starts the server using the provided address.
func (s Server) Start(addr string) error {
	http.HandleFunc("GET /movies/{id}", s.handleGetMovie)
	http.HandleFunc("GET /movies", s.handleGetAllMovies)
	http.HandleFunc("POST /movies", s.handleCreateMovie)
	http.HandleFunc("PUT /movies/{id}", s.handleUpdateMovie)
	http.HandleFunc("DELETE /movies/{id}", s.handleDeleteMovie)
	return http.ListenAndServe(addr, nil)
}

// handleGetMovie call Service to get movie with provided by path value id.
// If successful, the fetched movie is written to the response body.
func (s Server) handleGetMovie(w http.ResponseWriter, r *http.Request) {
	movie, err := s.svc.GetMovie(context.Background(), r.PathValue("id"))
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, movie)
}

// handleGetAllMovies calls Service to get all the movies that are currently stored.
// If successful, the fetched movies are written to the response body.
func (s Server) handleGetAllMovies(w http.ResponseWriter, _ *http.Request) {
	movies, err := s.svc.GetAllMovies(context.Background())
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusOK, movies)
}

// handleCreateMovie calls Service to create a new movie, using the data provided.
// in the request body. If successful, id of the created movie is written to the response body.
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

// handleUpdateMovie calls Service to update a movie, using the data provided in the body.
// and id provided in the request path value. If successful, nothing will be returned.
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

// handleDeleteMovie calls Service to delete a movie, using id provided in the request path value.
// If successful, nothing will be returned.
func (s Server) handleDeleteMovie(w http.ResponseWriter, r *http.Request) {
	err := s.svc.DeleteMovie(context.Background(), r.PathValue("id"))
	if err != nil {
		writeJson(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error()})
		return
	}
	writeJson(w, http.StatusNoContent, map[string]any{})
}

// writeJson is responsible for writing status code and response body.
func writeJson(w http.ResponseWriter, s int, v any) {
	w.WriteHeader(s)
	w.Header().Add("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(v)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
