CREATE TABLE rented_movies (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  movie_id int NOT NULL,
  rented_at TIME NOT NULL,
  should_return TIME NOT NULL
);