package fees

import "time"

type UserFeesDTO struct {
	Fees []OverrentFeeDTO
}

type OverrentFeeDTO struct {
	MovieID        int
	RentedAt       time.Time
	ShouldReturnAt time.Time
	Cash           float32
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
}
