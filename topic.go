package pubsub

import (
	"errors"
	"sync"
)

var topicIsClosed = errors.New("topic is closed")

type Topic[T any] struct {
	guard    sync.RWMutex
	closed   bool
	channels []chan T
}

// Subscribes to the topic using an unbuffered channel
func (t *Topic[T]) SubscribeUnbuffered() (<-chan T, error) {
	t.guard.Lock()
	defer t.guard.Unlock()

	if t.closed {
		return nil, topicIsClosed
	}

	channel := make(chan T)
	t.channels = append(t.channels, channel)

	return channel, nil
}

// Subscribes to the topic using a buff
func (t *Topic[T]) SubscribeBuffered(capacity int) (<-chan T, error) {
	t.guard.Lock()
	defer t.guard.Unlock()

	if t.closed {
		return nil, topicIsClosed
	}

	channel := make(chan T, capacity)
	t.channels = append(t.channels, channel)

	return channel, nil
}

// Unsubscribes a channel from a topic
func (t *Topic[T]) Unsubscribe(channel <-chan T) bool {
	t.guard.Lock()
	defer t.guard.Unlock()

	if t.closed {
		return true
	}

	for i, c := range t.channels {
		if c == channel {
			t.channels = append(t.channels[:i], t.channels[i+1:]...)
			close(c)

			return true
		}
	}

	return false
}

// Closes all channels. After this attempts to subscribe or publish will fail.
// If is same to call Close() multiple times.
func (t *Topic[T]) Close() {
	t.guard.Lock()
	defer t.guard.Unlock()

	if t.closed {
		return
	}

	t.closed = true

	for _, channel := range t.channels {
		close(channel)
	}

	t.channels = nil
}

// Sends data to all channels in the topic.
// Returns the number of channels that received data.
func (t *Topic[T]) Publish(data T) (int, error) {
	t.guard.RLock()
	defer t.guard.RUnlock()

	if t.closed {
		return 0, topicIsClosed
	}

	for _, channel := range t.channels {
		channel <- data
	}

	return len(t.channels), nil
}
