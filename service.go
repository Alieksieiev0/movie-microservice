package main

import (
	"context"
	"fmt"
)

// Service is responsible for to fetching/passing data.
type Service interface {
	// GetMovie fetches movie using provided id.
	GetMovie(ctx context.Context, id string) (*Movie, error)
	// GetAllMovies fetches all stored movies using provided id.
	GetAllMovies(ctx context.Context) ([]Movie, error)
	// CreateMovie creates movie using provided Movie struct.
	CreateMovie(ctx context.Context, movie *Movie) (string, error)
	// UpdateMovie updates movie with provided id using provided Movie struct.
	UpdateMovie(ctx context.Context, id string, movie *Movie) error
	// DeleteMovie deletes movie with provided id.
	DeleteMovie(ctx context.Context, id string) error
}

// MovieService is struct implementing Service interface.
type MovieService struct {
	// MovieService uses db interface, to fetch from or pass to the connected database.
	db Database
}

// NewServer creates an instance of the MovieService.
func NewMovieService(db Database) Service {
	return MovieService{
		db: db,
	}
}

func (ms MovieService) GetMovie(ctx context.Context, id string) (*Movie, error) {
	movie, err := ms.db.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error fetching by id: %v", err)
	}
	return movie, nil
}

func (ms MovieService) GetAllMovies(ctx context.Context) ([]Movie, error) {
	movies, err := ms.db.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("error fetching by id: %v", err)
	}
	return movies, nil
}

func (ms MovieService) CreateMovie(ctx context.Context, m *Movie) (string, error) {
	id, err := ms.db.Insert(ctx, m)
	if err != nil {
		return "", fmt.Errorf("error creating movie: %v", err)
	}
	return id, nil
}

func (ms MovieService) UpdateMovie(ctx context.Context, id string, m *Movie) error {
	err := ms.db.Update(ctx, id, m)
	if err != nil {
		return fmt.Errorf("error updating movie: %v", err)
	}
	return nil
}

func (ms MovieService) DeleteMovie(ctx context.Context, id string) error {
	err := ms.db.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting movie: %v", err)
	}
	return nil
}
