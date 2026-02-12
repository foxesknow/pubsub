package pubsub

type Subscriber[T any] interface {
	// Subscribes to the topic using an unbuffered channel
	SubscribeUnbuffered() (<-chan T, error)

	// Subscribes to the topic using a buffered channel
	SubscribeBuffered(capacity int) (<-chan T, error)

	// Unsubscribes a channel from a topic
	Unsubscribe(channel <-chan T) bool
}
