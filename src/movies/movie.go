package movies

type movie struct {
	ID         MovieID
	Title      string
	Year       int
	MinimalAge int
	Genre      string
}

type MovieID int
