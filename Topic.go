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
func (t *Topic[T]) Subscribe() (chan<- T, error) {
	t.guard.Lock()
	defer t.guard.Unlock()

	if t.closed {
		return nil, topicIsClosed
	}

	channel := make(chan T)
	t.channels = append(t.channels, channel)

	return channel, nil
}

// Subscribes to the topic using an existing channel.
// The topic will take ownership of the channel
func (t *Topic[T]) SubscribeWith(channel chan T) error {
	t.guard.Lock()
	defer t.guard.Unlock()

	if t.closed {
		return topicIsClosed
	}

	t.channels = append(t.channels, channel)

	return nil
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

func (t *Topic[T]) Publish(data T) error {
	t.guard.RLock()
	defer t.guard.RUnlock()

	if t.closed {
		return topicIsClosed
	}

	for _, channel := range t.channels {
		channel <- data
	}

	return nil
}
