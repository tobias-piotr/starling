package trips

type TripRepository interface {
	Create(data *TripData) (*Trip, error)
	GetAll(page int, perPage int) ([]*TripOverview, error)
	Get(id string) (*Trip, error)
	Update(id string, data map[string]any) error
}
