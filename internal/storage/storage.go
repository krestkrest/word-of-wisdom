package storage

//go:generate mockgen -source=storage.go -destination=storage_mock.go -package=storage Storage

// Storage provides access to the stored list of quotes
type Storage interface {
	// Start initialize storage
	Start() error
	// GetQuote get random quote from the list
	GetQuote() string
}
