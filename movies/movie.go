package movies

type Genre string

type Movie struct {
	ID         MovieID
	Title      string
	Year       int
	MinimalAge int
}

type MovieID int
