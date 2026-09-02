package events

// Publisher publishes application events and closes its underlying resources.
type Publisher interface {
	Publish(eventType string, payload interface{}, metadata map[string]string) error
	Close() error
}
