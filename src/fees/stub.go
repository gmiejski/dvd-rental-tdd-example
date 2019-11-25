package fees

import (
	"github.com/pkg/errors"
	"time"
)

type FacadeStub struct {
	fees map[int]UserFeesDTO
}

func (f *FacadeStub) GetFees(userID int) (UserFeesDTO, error) {
	fees, err := f.findFees(userID)
	if err != nil {
		return UserFeesDTO{Fees: []FeeDTO{}}, nil
	}
	return fees, nil
}

func (f *FacadeStub) AddFee(userID int, movieID int, rentedAt time.Time, shouldReturnAt time.Time, cash float32) {
	newFee := FeeDTO{MovieID: movieID, RentedAt: rentedAt, CurrentFee: cash, ShouldReturnAt: shouldReturnAt}
	_, err := f.findFees(userID)
	if err != nil {
		f.fees[userID] = UserFeesDTO{Fees: []FeeDTO{}}
	}
	f.fees[userID] = UserFeesDTO{Fees: append(f.fees[userID].Fees, newFee)}
}

func (f *FacadeStub) findFees(userID int) (UserFeesDTO, error) {
	for user, fee := range f.fees {
		if user == userID {
			return fee, nil
		}
	}
	return UserFeesDTO{}, errors.New("user not found")
}

func Build() Facade {
	return &FacadeStub{fees: make(map[int]UserFeesDTO)}
}
