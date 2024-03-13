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
	"github.com/shopspring/decimal"
)

type Database interface {
	Get(context.Context, string) (*Movie, error)
	GetAll(context.Context) ([]Movie, error)
	Insert(context.Context, *Movie) (string, error)
	Update(context.Context, string, *Movie) error
	Delete(context.Context, string) error
}

type MovieDatabase struct {
	pool *pgxpool.Pool
}

func (mdb MovieDatabase) Get(ctx context.Context, id string) (*Movie, error) {
	rows, err := mdb.pool.Query(ctx, "select * from movie where id = $1", id)
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
	rows, err := mdb.pool.Query(ctx, "select * from movie")
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

func (mdb MovieDatabase) Insert(ctx context.Context, movie *Movie) (string, error) {
	tx, err := mdb.pool.Begin(context.Background())
	if err != nil {
		return "", err
	}

	defer tx.Rollback(ctx)

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
		return "", err
	}
	defer rows.Close()

	id, err := pgx.CollectExactlyOneRow(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return "", err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return "", err
	}
	return id.String(), nil
}
func (mdb MovieDatabase) Update(ctx context.Context, id string, movie *Movie) error {
	tx, err := mdb.pool.Begin(context.Background())
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	q, params := mdb.buildUpdateQuery(movie, id)
	ct, err := tx.Exec(context.Background(), q, params...)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("entity with such id does not exist")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (mdb MovieDatabase) buildUpdateQuery(movie *Movie, id string) (string, []any) {
	statements := []string{}
	params := []any{}
	counter := 1

	fmt.Println(movie)
	if movie.Name != "" {
		addToQuery(movie.Name, "name", &counter, &statements, &params)
	}

	if movie.ReleaseYear != 0 {
		addToQuery(movie.ReleaseYear, "release_year", &counter, &statements, &params)
	}

	if movie.Rating != *new(decimal.Decimal) {
		addToQuery(movie.Rating, "rating", &counter, &statements, &params)
	}

	if movie.Genres != nil {
		addToQuery(movie.Genres, "genres", &counter, &statements, &params)
	}

	if movie.Director != "" {
		addToQuery(movie.Director, "director", &counter, &statements, &params)
	}

	params = append(params, id)
	fmt.Println(params)
	q := fmt.Sprintf("update movie set %s where id = $%d", strings.Join(statements, " "), counter)
	fmt.Println(q)
	return q, params
}

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

func (mdb MovieDatabase) Delete(ctx context.Context, id string) error {
	tx, err := mdb.pool.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	ct, err := tx.Exec(ctx, "delete from movie where id = $1", id)
	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("entity with such id does not exist")
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}

func ConnectDB(dbUrl string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(os.Getenv(dbUrl))
	if err != nil {
		log.Fatalf("Unable to load database config: %v\n", err)
	}
	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		pgxuuid.Register(conn.TypeMap())
		pgxdecimal.Register(conn.TypeMap())
		return nil
	}
	return pgxpool.NewWithConfig(context.Background(), config)
}
