package main

import (
	"context"
	"fmt"
)

type Service interface {
	GetMovie(context.Context, string) (*Movie, error)
	GetAllMovies(context.Context) ([]Movie, error)
	CreateMovie(context.Context, *Movie) (string, error)
	UpdateMovie(context.Context, string, *Movie) error
	DeleteMovie(context.Context, string) error
}

type MovieService struct {
	db Database
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
