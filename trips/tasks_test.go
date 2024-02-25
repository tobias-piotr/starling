package trips

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"starling/internal/domain"

	"github.com/stretchr/testify/assert"
)

// MockedTripRepository implements TripRepository interface with the most basic interactions
// needed for the tests.
type MockedTripRepository struct {
	trip *Trip
}

func (m *MockedTripRepository) Create(data *TripData) (*Trip, error) {
	return &Trip{}, nil
}

func (m *MockedTripRepository) GetAll(page int, perPage int) ([]*TripOverview, error) {
	return []*TripOverview{}, nil
}

func (m *MockedTripRepository) Get(id string) (*Trip, error) {
	return m.trip, nil
}

func (m *MockedTripRepository) Update(id string, data map[string]any) error {
	// Only handle status updates
	status := data["status"].(string)
	if status == CompletedStatus.String() {
		m.trip.Status = CompletedStatus
	}
	if status == FailedStatus.String() {
		m.trip.Status = FailedStatus
	}
	return nil
}

func (m *MockedTripRepository) AddResult(tripID string) (*TripResult, error) {
	tr := &TripResult{}
	m.trip.Result = tr
	return tr, nil
}

func (m *MockedTripRepository) UpdateResult(resultID string, data map[string]any) error {
	// Update trip result object based on data map
	// Make keys uppercased to match the struct fields
	structData := make(map[string]any, len(data))
	for key, val := range data {
		chars := strings.Split(key, "")
		chars[0] = strings.ToUpper(chars[0])
		upKey := strings.Join(chars, "")
		structData[upKey] = val
	}

	rv := reflect.ValueOf(m.trip.Result).Elem()
	for key, val := range structData {
		fv := rv.FieldByName(key)
		fv.SetString(val.(string))
	}

	return nil
}

type MockedAiClient struct {
	response string
	err      error
}

func (m *MockedAiClient) Send(prompt string) (string, error) {
	return m.response, m.err
}

func TestRequestTrip(t *testing.T) {
	repo := &MockedTripRepository{&Trip{}}
	resp := "test response"
	aiClient := &MockedAiClient{resp, nil}
	err := RequestTrip(repo, aiClient, "1")
	assert.Nil(t, err)
	assert.Equal(t, CompletedStatus, repo.trip.Status)
	assert.Equal(t, resp, repo.trip.Result.Summary)
	assert.Equal(t, resp, repo.trip.Result.Attractions)
	assert.Equal(t, resp, repo.trip.Result.Weather)
	assert.Equal(t, resp, repo.trip.Result.Prices)
	assert.Equal(t, resp, repo.trip.Result.Luggage)
	assert.Equal(t, resp, repo.trip.Result.Documents)
	assert.Equal(t, resp, repo.trip.Result.Commuting)
}

func TestRequestTripFailed(t *testing.T) {
	repo := &MockedTripRepository{&Trip{}}
	aiClient := &MockedAiClient{"", errors.New("test error")}
	err := RequestTrip(repo, aiClient, "1")
	assert.NotNil(t, err)
	assert.Equal(t, FailedStatus, repo.trip.Status)
	cerr := err.(*domain.CompositeErr)
	// All the fields should have failed
	assert.Equal(t, len(requestResultConf), len(cerr.Errs))
}
