package rental

type Config struct {
	MaxRentedMoviesCount int
}

func StandardConfig() Config {
	return Config{MaxRentedMoviesCount: 2}
}

type RentalOption = func(*Config)

var MaxRentedMoviesCount = func(maxCount int) RentalOption {
	return func(f *Config) {
		f.MaxRentedMoviesCount = maxCount
	}
}
