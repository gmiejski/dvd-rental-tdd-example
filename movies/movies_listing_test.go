package movies

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strconv"
	"testing"
)

const horror = "horror"
const family = "family"

func TestListAllGenreMovies(t *testing.T) {
	// given
	moviesFacade := buildMoviesFacade()
	totalMovies := 5
	movies := createMoviesWithGenre(totalMovies, horror)
	for _, movie := range movies {
		moviesFacade.Add(movie)
	}

	// when first part of listing is retrieved
	listingPartOne, err := moviesFacade.ListGenre(
		GenreListingRequest{Genre: horror, CursorOffset: 0, Limit: 3},
	)

	// then
	require.NoError(t, err)
	assert.EqualValues(t, PageInfo{HasNextPage: true, LastCursor: 3}, listingPartOne.PageInfo)
	assert.EqualValues(t, 3, len(listingPartOne.Movies))

	// when getting second part of listing results
	listingPartTwo, err := moviesFacade.ListGenre(
		GenreListingRequest{Genre: horror, Limit: 3, CursorOffset: listingPartOne.PageInfo.LastCursor},
	)

	// then all results are returned
	assert.EqualValues(t, 2, len(listingPartTwo.Movies))
	assert.EqualValues(t, PageInfo{HasNextPage: false, LastCursor: totalMovies}, listingPartTwo.PageInfo)
	assert.True(t, moviesMatches(movies[0], listingPartOne.Movies[0]))
	assert.True(t, moviesMatches(movies[1], listingPartOne.Movies[1]))
	assert.True(t, moviesMatches(movies[2], listingPartOne.Movies[2]))
	assert.True(t, moviesMatches(movies[3], listingPartTwo.Movies[0]))
	assert.True(t, moviesMatches(movies[4], listingPartTwo.Movies[1]))
}

func TestListOnlyMoviesFromSpecificGenre(t *testing.T) {
	// given
	moviesFacade := buildMoviesFacade()
	_, err := moviesFacade.Add(CreateMovie{Title: "Scary", MinimalAge: 18, Year: 2017, Genre: horror})
	require.NoError(t, err)

	// when
	listing, err := moviesFacade.ListGenre(
		GenreListingRequest{Genre: family, Limit: 3, CursorOffset: 0},
	)

	// then
	require.NoError(t, err)
	assert.Empty(t, listing.Movies)
	assert.EqualValues(t, PageInfo{HasNextPage: false, LastCursor: -1}, listing.PageInfo)
	assert.EqualValues(t, 0, listing.TotalResults)
}

func createMoviesWithGenre(totalMovies int, genre string) []CreateMovie {
	moviesToCreate := make([]CreateMovie, 0)
	for i := 1; i <= totalMovies; i++ {
		moviesToCreate = append(moviesToCreate,
			CreateMovie{Title: "movie " + strconv.Itoa(i), Year: 2000 + i, MinimalAge: 18, Genre: genre})
	}
	return moviesToCreate
}

func moviesMatches(expected CreateMovie, actual MovieDTO) bool {
	return expected.Title == actual.Title &&
		expected.Year == actual.Year &&
		expected.Genre == actual.Genre &&
		expected.MinimalAge == actual.MinimalAge
}
