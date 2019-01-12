package movies

type CreateMovie struct {
	Title      string
	Year       int
	MinimalAge int
	Genre      string
}

type MovieDTO struct {
	ID         movieID
	Title      string
	Year       int
	MinimalAge int
	Genre      string
}

type CreatedMovieDTO struct {
	ID int
}

type PageInfo struct {
	HasNextPage bool `json:"hasNextPage"`
	LastCursor  int  `json:"lastCursor"`
}

type GenreListingRequest struct {
	Genre        string `json:"genre"`
	CursorOffset int    `json:"cursorOffset"`
	Limit        int    `json:"limit"`
}

type ListingDTO struct {
	TotalResults int `json:"totalResults"`
	Movies       []MovieDTO
	PageInfo     PageInfo `json:"pageInfo"`
}

type Facade interface {
	Add(movie CreateMovie) (CreatedMovieDTO, error)
	Get(movie movieID) (MovieDTO, error)
	ListGenre(request GenreListingRequest) (ListingDTO, error)
}
