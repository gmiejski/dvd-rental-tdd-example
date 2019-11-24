package infrastructure

import (
	"database/sql"
	"github.com/gmiejski/dvd-rental-tdd-example/src/users"
)

func NewPostgresRepository(dbDSN string) users.Repository {
	db, err := sql.Open("postgres", dbDSN)
	if err != nil {
		panic(err.Error())
	}
	return &postgresRepository{db: db}
}

type postgresRepository struct {
	db *sql.DB
}

func (r *postgresRepository) Save(user users.User) (users.User, error) {
	lastInsertId := 0
	err := r.db.QueryRow(
		"INSERT INTO users (name, age) VALUES($1, $2) RETURNING id",
		user.Name,
		user.Age).Scan(&lastInsertId)
	if err != nil {
		return users.User{}, err
	}
	createdUser := user
	createdUser.ID = lastInsertId
	return createdUser, nil
}

func (r *postgresRepository) Find(userID int) (*users.User, error) {
	var user users.User
	row := r.db.QueryRow(`SELECT id, name, age FROM users WHERE id = $1`, userID)
	err := row.Scan(&user.ID, &user.Name, &user.Age)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
