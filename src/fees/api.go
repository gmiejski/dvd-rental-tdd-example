package fees

import "time"

type UserFeesDTO struct {
	Fees []FeeDTO
}

type FeeDTO struct {
	MovieID        int
	RentedAt       time.Time
	ShouldReturnAt time.Time
	CurrentFee     float32
}

func (fees UserFeesDTO) OverrentMovieIDs() []int {
	var movieIDs []int
	for _, fee := range fees.Fees {
		movieIDs = append(movieIDs, fee.MovieID)
	}
	return movieIDs
}

type Facade interface {
	GetFees(userID int) (UserFeesDTO, error)
	AddFee(userID int, movieID int, rentedAt time.Time, shouldReturnAt time.Time, cash float32)
}
