package movies

type movie struct {
	ID         movieID
	Title      string
	Year       int
	MinimalAge int
	Genre      string
}

type movieID int
