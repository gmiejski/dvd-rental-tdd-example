package rental

import (
	"github.com/gmiejski/dvd-rental-tdd-example/fees"
	"github.com/gmiejski/dvd-rental-tdd-example/movies"
	"github.com/gmiejski/dvd-rental-tdd-example/users"
)

type config struct {
	maxRentedMoviesCount int
}

type testOptionFacade = func(*facade)

var withConfig = func(c config) testOptionFacade {
	return func(f *facade) {
		f.config = c
	}
}

var withFeesFacade = func(feesFacade fees.Facade) testOptionFacade {
	return func(f *facade) {
		f.fees = feesFacade
	}
}

func buildTestFacade(users users.UsersFacade, movies movies.Facade, options ...testOptionFacade) RentalFacade {
	feesStub := fees.NewFacadeStub()
	baseTestFacade := &facade{
		users:      users,
		movies:     movies,
		fees:       &feesStub,
		repository: newInMemoryRepository(),
		config:     config{maxRentedMoviesCount: 10},
	}

	for _, option := range options {
		option(baseTestFacade)
	}

	return baseTestFacade
}
