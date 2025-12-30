package build

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/kmdkuk/nfrecap/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCache implements store.Cache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Get(workTitle string, typ string) (model.Metadata, bool, error) {
	args := m.Called(workTitle, typ)
	return args.Get(0).(model.Metadata), args.Bool(1), args.Error(2)
}

func (m *MockCache) Put(workTitle string, typ string, md model.Metadata) error {
	args := m.Called(workTitle, typ, md)
	return args.Error(0)
}

// MockProvider implements provider.Provider
type MockProvider struct {
	mock.Mock
}

func (m *MockProvider) Lookup(workTitle string, typ string) (model.Metadata, bool, error) {
	args := m.Called(workTitle, typ)
	return args.Get(0).(model.Metadata), args.Bool(1), args.Error(2)
}

func TestRun(t *testing.T) {
	// Fixed date for testing
	recordDate := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)
	recordStr := "2023-01-01"

	tests := []struct {
		name          string
		records       []model.ViewingRecord
		opts          Options
		setupMocks    func(*MockCache, *MockProvider)
		expectedSum   Summary
		expectedError string
	}{
		{
			name: "Cache Hit",
			records: []model.ViewingRecord{
				{Title: "Inception: Season 1", Date: recordDate}, // Normalized -> WorkTitle: Inception, Type: tv
			},
			opts: Options{Fetch: true},
			setupMocks: func(c *MockCache, p *MockProvider) {
				// Verify normalized title logic (e.g. "Inception: Season 1" -> workTitle "Inception", type "tv")
				// based on title/normalize.go in previous steps
				c.On("Get", "Inception", "tv").Return(model.Metadata{Title: "Inception TV"}, true, nil)
			},
			expectedSum: Summary{
				CacheHits:   1,
				CacheMisses: 0,
				Fetched:     0,
				Unresolved:  0,
			},
		},
		{
			name: "Cache Miss, Fetch=True, Provider Found",
			records: []model.ViewingRecord{
				{Title: "New Movie", Date: recordDate},
			},
			opts: Options{Fetch: true},
			setupMocks: func(c *MockCache, p *MockProvider) {
				c.On("Get", "New Movie", "movie").Return(model.Metadata{}, false, nil)
				p.On("Lookup", "New Movie", "movie").Return(model.Metadata{Title: "New Movie Found"}, true, nil)
				c.On("Put", "New Movie", "movie", model.Metadata{Title: "New Movie Found"}).Return(nil)
			},
			expectedSum: Summary{
				CacheHits:   0,
				CacheMisses: 1,
				Fetched:     1,
				Unresolved:  0,
			},
		},
		{
			name: "Cache Miss, Fetch=False",
			records: []model.ViewingRecord{
				{Title: "Skipped Movie", Date: recordDate},
			},
			opts: Options{Fetch: false},
			setupMocks: func(c *MockCache, p *MockProvider) {
				c.On("Get", "Skipped Movie", "movie").Return(model.Metadata{}, false, nil)
				// Provider should NOT be called
			},
			expectedSum: Summary{
				CacheHits:   0,
				CacheMisses: 1,
				Fetched:     0,
				Unresolved:  1,
			},
		},
		{
			name: "Cache Miss, Provider Not Found",
			records: []model.ViewingRecord{
				{Title: "Unknown", Date: recordDate},
			},
			opts: Options{Fetch: true},
			setupMocks: func(c *MockCache, p *MockProvider) {
				c.On("Get", "Unknown", "movie").Return(model.Metadata{}, false, nil)
				p.On("Lookup", "Unknown", "movie").Return(model.Metadata{}, false, nil)
			},
			expectedSum: Summary{
				CacheHits:   0,
				CacheMisses: 1,
				Fetched:     0,
				Unresolved:  1,
			},
		},
		{
			name: "Provider Error",
			records: []model.ViewingRecord{
				{Title: "Error Case", Date: recordDate},
			},
			opts: Options{Fetch: true},
			setupMocks: func(c *MockCache, p *MockProvider) {
				c.On("Get", "Error Case", "movie").Return(model.Metadata{}, false, nil)
				p.On("Lookup", "Error Case", "movie").Return(model.Metadata{}, false, errors.New("network error"))
			},
			expectedSum:   Summary{},
			expectedError: "network error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCache := new(MockCache)
			mockProvider := new(MockProvider)

			if tt.setupMocks != nil {
				tt.setupMocks(mockCache, mockProvider)
			}

			output, sum, err := Run(tt.records, mockCache, mockProvider, tt.opts)

			if tt.expectedError != "" {
				assert.ErrorContains(t, err, tt.expectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedSum, sum)
				assert.NotNil(t, output)

				// Basic check on output JSON structure
				var built Built
				err := json.Unmarshal(output, &built)
				assert.NoError(t, err)
				assert.Len(t, built.Items, len(tt.records))
				assert.Equal(t, recordStr, built.Items[0].Date)
			}

			mockCache.AssertExpectations(t)
			mockProvider.AssertExpectations(t)
		})
	}
}
