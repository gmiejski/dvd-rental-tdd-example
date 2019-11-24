package domain_es

type Repository interface {
	Get(user int) (*UserRents, error)
	Save(userID int, eventsToSave []Event) error
}
