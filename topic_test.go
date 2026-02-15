package pubsub

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestCheckInterfaces(t *testing.T) {
	var _ Publisher[int32] = (*Topic[int32])(nil)
	var _ Subscriber[int32] = (*Topic[int32])(nil)
}

func TestDefault(t *testing.T) {
	var topic Topic[string]

	count, err := topic.Publish("hello")
	if count != 0 {
		t.Error("count should be zero")
	}

	if err != nil {
		t.Error("err should be nil")
	}
}

func TestUnsubscribe(t *testing.T) {
	var topic Topic[string]

	channel, _ := topic.SubscribeUnbuffered()
	success := topic.Unsubscribe(channel)

	if !success {
		t.Error("failed to unsubscribe channel")
	}
}

func TestUnsubscribeTwice(t *testing.T) {
	var topic Topic[string]

	channel, _ := topic.SubscribeUnbuffered()
	success := topic.Unsubscribe(channel)

	if !success {
		t.Error("failed to unsubscribe channel")
	}

	success = topic.Unsubscribe(channel)
	if success {
		t.Error("should have failed to unsubscribe the channel a second time")
	}
}

func TestSubscribeAfterClose(t *testing.T) {
	var topic Topic[string]

	topic.Close()

	channel, err := topic.SubscribeUnbuffered()
	if channel != nil {
		t.Error("no channel should be returned as we are closed")
	}

	if err == nil {
		t.Error("we should get back an error")
	}
}

func TestSubscribeUnbuffered(t *testing.T) {
	var topic Topic[string]

	channel, err := topic.SubscribeUnbuffered()
	if channel == nil {
		t.Error("no channel returned")
	}

	if err != nil {
		t.Error("err should be nil")
	}
}

func TestSubscribeUnbuffered_OneReceiver(t *testing.T) {
	var topic Topic[int32]

	channel, _ := topic.SubscribeUnbuffered()

	var wg sync.WaitGroup
	var counter atomic.Int32

	wg.Go(func() {
		for value := range channel {
			counter.Add(value)
		}
	})

	topic.Publish(10)
	topic.Publish(20)

	topic.Close()
	wg.Wait()

	if counter.Load() != 30 {
		t.Error("invalid count")
	}
}

func TestSubscribeUnbuffered_ManyReceivers(t *testing.T) {
	var topic Topic[int32]

	var wg sync.WaitGroup
	var counter atomic.Int32

	for i := 0; i < 5; i += 1 {
		channel, _ := topic.SubscribeUnbuffered()
		wg.Go(func() {
			for value := range channel {
				counter.Add(value)
			}
		})
	}

	topic.Publish(10)
	topic.Publish(20)

	topic.Close()
	wg.Wait()

	if counter.Load() != 150 {
		t.Error("invalid count")
	}
}

func TestSubscribeUnbuffered_ManyReceivers_KeepOpen(t *testing.T) {
	var topic Topic[int32]

	var counter atomic.Int32

	for i := 0; i < 5; i += 1 {
		channel, _ := topic.SubscribeUnbuffered()
		go func() {
			for value := range channel {
				if value == 0 {
					break
				}
				counter.Add(value)
			}
		}()
	}

	topic.Publish(10)
	topic.Publish(20)
	topic.Publish(0)

	if counter.Load() != 150 {
		t.Error("invalid count")
	}

	topic.Close()
}

func TestSubscribeBuffered_OneReceiver(t *testing.T) {
	var topic Topic[int32]

	channel, _ := topic.SubscribeBuffered(10)

	var wg sync.WaitGroup
	var counter atomic.Int32

	topic.Publish(10)
	topic.Publish(20)
	topic.Close()

	wg.Go(func() {
		for value := range channel {
			counter.Add(value)
		}
	})

	wg.Wait()

	if counter.Load() != 30 {
		t.Error("invalid count")
	}
}

func TestSubscribeBuffered_ManyReceivers(t *testing.T) {
	var topic Topic[int32]

	var wg sync.WaitGroup
	var counter atomic.Int32

	for i := 0; i < 5; i += 1 {
		channel, _ := topic.SubscribeBuffered(10)
		wg.Go(func() {
			for value := range channel {
				if value == 0 {
					break
				}
				counter.Add(value)
			}
		})
	}

	topic.Publish(10)
	topic.Publish(20)
	topic.Publish(0)

	wg.Wait()

	if counter.Load() != 150 {
		t.Error("invalid count")
	}

	topic.Close()
}
