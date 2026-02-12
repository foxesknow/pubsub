package pubsub

// Defines the behaviour of a publisher
type Publisher[T any] interface {
	// Sends data to all subscribers in the topic.
	// Returns the number of subscribers that received data.
	Publish(data T) (int, error)
}
