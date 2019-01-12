package movies

type CreateMovie struct {
	Title      string
	Year       int
	MinimalAge int
}

type MovieDTO struct {
	ID         MovieID
	Title      string
	Year       int
	MinimalAge int
}

type CreatedMovieDTO struct {
	ID int
}

type MoviesFacade interface {
	Add(Movie CreateMovie) (CreatedMovieDTO, error)
	Get(Movie MovieID) (MovieDTO, error)
	ListGenre()
}
