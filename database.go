package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gofrs/uuid/v5"
	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	pgxdecimal "github.com/jackc/pgx-shopspring-decimal"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/tracelog"
	"github.com/shopspring/decimal"
)

// Database is responsible for providing database CRUD interface.
type Database interface {
	// Get fetches movie from the DB using provided id.
	Get(ctx context.Context, id string) (*Movie, error)
	// GetAllMovies fetches all movies stored in DB using provided id.
	GetAll(ctx context.Context) ([]Movie, error)
	// CreateMovie creates movie row in DB using provided Movie struct.
	Insert(ctx context.Context, movie *Movie) (string, error)
	// UpdateMovie updates movie row in DB with provided id using provided Movie struct.
	Update(ctx context.Context, id string, movie *Movie) error
	// DeleteMovie deletes movie row from DB with provided id.
	Delete(ctx context.Context, id string) error
}

// databaseConn is responsible for providing methods for communicating with DB.
type databaseConn interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, optionsAndArgs ...interface{}) pgx.Row
}

// MovieDatabase is a struct implementing Database interface.
type MovieDatabase struct {
	conn databaseConn
}

// NewMovieDatabase creates an instance of the MovieDatabase.
func NewMovieDatabase(conn databaseConn) Database {
	return MovieDatabase{
		conn: conn,
	}
}

func (mdb MovieDatabase) Get(ctx context.Context, id string) (*Movie, error) {
	rows, err := mdb.conn.Query(ctx, "select * from movie where id = $1", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movie, err := pgx.CollectExactlyOneRow(rows, pgx.RowToAddrOfStructByName[Movie])
	if err != nil {
		return nil, err
	}
	return movie, nil
}

func (mdb MovieDatabase) GetAll(ctx context.Context) ([]Movie, error) {
	rows, err := mdb.conn.Query(ctx, "select * from movie")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	movies, err := pgx.CollectRows(rows, pgx.RowToStructByName[Movie])
	if err != nil {
		return nil, err
	}
	return movies, nil
}

func (mdb MovieDatabase) Insert(ctx context.Context, movie *Movie) (id string, err error) {
	tx, err := mdb.conn.Begin(context.Background())
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			err = tx.Rollback(ctx)
		}
	}()

	q := `
	insert into movie(name, release_year, rating, genres, director)
	values ($1, $2, $3, $4, $5)
	returning id
	`
	rows, err := tx.Query(
		context.Background(),
		q,
		movie.Name,
		movie.ReleaseYear,
		movie.Rating,
		movie.Genres,
		movie.Director,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	movieId, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return
	}
	return movieId.String(), nil
}
func (mdb MovieDatabase) Update(ctx context.Context, id string, movie *Movie) (err error) {
	tx, err := mdb.conn.Begin(context.Background())
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			err = tx.Rollback(ctx)
		}
	}()

	q, params, err := mdb.buildUpdateQuery(movie, id)
	if err != nil {
		return
	}

	ct, err := tx.Exec(context.Background(), q, params...)
	if err != nil {
		return
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("entity with such id does not exist")
	}

	return nil
}

// buildUpdateQuery dynamically adds statements into the query string for fields
// with values different from types zero values.
func (mdb MovieDatabase) buildUpdateQuery(movie *Movie, id string) (string, []any, error) {
	statements := []string{}
	params := []any{}
	counter := 1

	if movie.Name != "" {
		addToQuery(movie.Name, "name", &counter, &statements, &params)
	}

	if movie.ReleaseYear != 0 {
		addToQuery(movie.ReleaseYear, "release_year", &counter, &statements, &params)
	}

	ratingVal, err := movie.Rating.Value()
	if err != nil {
		return "", nil, err
	}
	zeroVal, err := new(decimal.Decimal).Value()
	if err != nil {
		return "", nil, err
	}

	if ratingVal != zeroVal {
		addToQuery(movie.Rating, "rating", &counter, &statements, &params)
	}

	if movie.Genres != nil {
		addToQuery(movie.Genres, "genres", &counter, &statements, &params)
	}

	if movie.Director != "" {
		addToQuery(movie.Director, "director", &counter, &statements, &params)
	}

	params = append(params, id)
	q := fmt.Sprintf("update movie set %s where id = $%d", strings.Join(statements, ", "), counter)
	return q, params, nil
}

// addToQuery dynamically appends statements and params with provided name and field,
// increasing the counter of added statements.
func addToQuery[T any](
	field T,
	name string,
	counter *int,
	statements *[]string,
	params *[]any,
) {
	*statements = append(*statements, fmt.Sprintf("%s = $%d", name, *counter))
	*params = append(*params, field)
	*counter++
}

func (mdb MovieDatabase) Delete(ctx context.Context, id string) (err error) {
	tx, err := mdb.conn.Begin(ctx)
	if err != nil {
		return
	}

	defer func() {
		switch err {
		case nil:
			err = tx.Commit(ctx)
		default:
			err = tx.Rollback(ctx)
		}
	}()

	ct, err := tx.Exec(ctx, "delete from movie where id = $1", id)
	if err != nil {
		return
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("entity with such id does not exist")
	}

	return nil
}

// ConnectDB connects to the PostgreSql database, using provided DB url.
func ConnectDB(dbUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(os.Getenv(dbUrl))
	if err != nil {
		log.Fatalf("Unable to load database config: %v\n", err)
	}

	logger := log.New(os.Stdout, "SQL INFO: ", log.LstdFlags)
	config.ConnConfig.Tracer = &tracelog.TraceLog{
		Logger:   NewDatabaseLogger(logger),
		LogLevel: tracelog.LogLevelTrace,
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	return pgxpool.NewWithConfig(context.Background(), config)
}
