package main

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/tracelog"
)

// DatabaseLogger is responsible for logging information about all Database operations
// by implementing tracelog.Logger interface.
type DatabaseLogger struct {
	// Logger from Golang standard library is used to print the info.
	logger *log.Logger
}

// NewDatabaseLogger creates an instance of the MovieService.
func NewDatabaseLogger(logger *log.Logger) tracelog.Logger {
	return DatabaseLogger{
		logger: logger,
	}
}

// Log is used to print information about completed operations.
func (dl DatabaseLogger) Log(
	ctx context.Context,
	level tracelog.LogLevel,
	msg string,
	data map[string]any,
) {
	dl.logger.Printf("Operation - %s", msg)
	for k, v := range data {
		dl.logger.Printf(" - %v: %v", k, v)
	}
	dl.logger.Println("-----------")
}

// LoggingService implements Service interface, which means, that this struct
// can be a wrapper around the another struct that implements Service interface
// For example: MovieService.
type LoggingService struct {
	logger *log.Logger
	next   Service
}

// NewLogginService creates an instance of the LoggingService.
func NewLoggingService(logger *log.Logger, next Service) Service {
	return LoggingService{
		logger: logger,
		next:   next,
	}
}

func (ls LoggingService) GetMovie(ctx context.Context, id string) (movie *Movie, err error) {
	defer func(start time.Time) {
		ls.logger.Println("GetMovie Results")
		if movie != nil {
			ls.logger.Printf(
				" - Id: %s, Name: %s, ReleaseYear: %d, Rating: %v, Genres: %v, Director: %s",
				movie.Id,
				movie.Name,
				movie.ReleaseYear,
				movie.Rating,
				movie.Genres,
				movie.Director,
			)
		}
		ls.logger.Printf(" - err: %v", err)
		ls.logger.Printf(" - time: %v", time.Since(start))
		ls.logger.Println("-----------")
	}(time.Now())

	return ls.next.GetMovie(ctx, id)
}

func (ls LoggingService) GetAllMovies(ctx context.Context) (movies []Movie, err error) {
	defer func(start time.Time) {
		ls.logger.Println("GetAllMovies Results")
		for _, m := range movies {
			ls.logger.Printf(
				" - Id: %s, Name: %s, ReleaseYear: %d, Rating: %v, Genres: %v, Director: %s",
				m.Id,
				m.Name,
				m.ReleaseYear,
				m.Rating,
				m.Genres,
				m.Director,
			)
		}
		ls.logger.Printf(" - err: %v", err)
		ls.logger.Printf(" - time: %v", time.Since(start))
		ls.logger.Println("-----------")
	}(time.Now())

	return ls.next.GetAllMovies(ctx)
}

func (ls LoggingService) CreateMovie(ctx context.Context, m *Movie) (id string, err error) {
	defer func(start time.Time) {
		ls.logger.Println("CreateMovie Results")
		ls.logger.Printf(" - id: %v", id)
		ls.logger.Printf(" - err: %v", err)
		ls.logger.Printf(" - time: %v", time.Since(start))
		ls.logger.Println("-----------")
	}(time.Now())

	return ls.next.CreateMovie(ctx, m)
}

func (ls LoggingService) UpdateMovie(ctx context.Context, id string, m *Movie) (err error) {
	defer func(start time.Time) {
		ls.logger.Println("UpdateMovie Results")
		ls.logger.Printf(" - err: %v", err)
		ls.logger.Printf(" - time: %v", time.Since(start))
		ls.logger.Println("-----------")
	}(time.Now())

	return ls.next.UpdateMovie(ctx, id, m)
}

func (ls LoggingService) DeleteMovie(ctx context.Context, id string) (err error) {
	defer func(start time.Time) {
		ls.logger.Println("DeleteMovie Results")
		ls.logger.Printf(" - err: %v", err)
		ls.logger.Printf(" - time: %v", time.Since(start))
		ls.logger.Println("-----------")
	}(time.Now())

	return ls.next.DeleteMovie(ctx, id)
}
