package trips

type TripRepository interface {
	Create(data *TripData) (*Trip, error)
	GetAll(page int, perPage int) ([]*TripOverview, error)
	Get(id string) (*Trip, error)
	Update(id string, data map[string]any) error
	AddResult(tripID string) (*TripResult, error)
	UpdateResult(resultID string, data map[string]any) error
}
