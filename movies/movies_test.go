package movies

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRetrievingAddedMovie(t *testing.T) {
	// given
	moviesFacade := buildMoviesFacade()
	createdMovie, err := moviesFacade.Add(CreateMovie{Title: "Saw", Year: 2017, MinimalAge: 18})
	require.NoError(t, err)

	// when
	movie, err := moviesFacade.Get(MovieID(createdMovie.ID))

	// then
	require.NoError(t, err)
	assert.EqualValues(t, MovieDTO{ID: 1, Title: "Saw", MinimalAge: 18, Year: 2017}, movie)
}

func TestErrorWhenMovieNotFound(t *testing.T) {
	// given
	moviesFacade := buildMoviesFacade()

	// when
	_, err := moviesFacade.Get(10)

	// then
	require.Error(t, err)
	require.IsType(t, errors.Cause(err), MovieNotFound{})
}
