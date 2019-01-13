package rental

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
)

type postgresRepository struct {
	db *sql.DB
}

func (p *postgresRepository) Save(rents UserRents) error {
	_, err := p.db.Exec("DELETE FROM rented_movies WHERE user_id = $1", rents.UserID)

	if err != nil {
		return err
	}

	for _, movie := range rents.RentedMovies {
		query := `INSERT INTO rented_movies(user_id, movie_id, rented_at, should_return)
	VALUES($1, $2, $3, $4)`
		_, err := p.db.Exec(query, rents.UserID, movie.MovieID, movie.RentedAt, movie.ReturnAt)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *postgresRepository) Get(userID int) (*UserRents, error) {
	rows, err := p.db.Query(`SELECT movie_id, rented_at, should_return FROM rented_movies WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	rentedMovies, err := p.toRentedMovies(rows)
	if err != nil {
		return nil, err
	}
	return &UserRents{UserID: userID, RentedMovies: rentedMovies}, nil
}

func (p *postgresRepository) toRentedMovies(rows *sql.Rows) ([]RentedMovie, error) {
	var movies []RentedMovie
	defer rows.Close()
	for rows.Next() {
		var movie = RentedMovie{}
		// Here `Scan` performs the data type conversions for us
		// based on the type of the destination variable.
		// If an error occur in the conversion, `Scan` will return
		// that error for you.
		err := rows.Scan(
			&movie.MovieID, &movie.RentedAt, &movie.ReturnAt)
		if err != nil {
			return nil, errors.Wrapf(err, "row scanning failed")
		}

		movies = append(movies, movie)
	}
	if err := rows.Err(); err != nil {
		return nil, errors.Wrapf(err, "Error while looping through rented movies")
	}
	return movies, nil
}

func NewPostgresRepository(dbDSN string) Repository {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err.Error())
	}
	return &postgresRepository{db: db}
}
