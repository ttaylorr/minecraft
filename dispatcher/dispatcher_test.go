package dispatcher_test

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ttaylorr/minecraft/dispatcher"
)

func TestConstruction(t *testing.T) {
	d := dispatcher.New()
	assert.IsType(t, d, &dispatcher.Dispatcher{})
}

func TestSubscribedReceiversGetCorrectEvents(t *testing.T) {
	d := dispatcher.New()
	receiver := make(chan interface{})
	d.Subscribe("some-type", receiver)

	d.Events <- &dispatcher.Event{
		Type: "some-type",
		Data: "data 1",
	}

	d.Events <- &dispatcher.Event{
		Type: "other-type",
		Data: "data 2",
	}

	go func(receiver chan interface{}) {
		e := <-receiver

		assert.Equal(t, e, "data 1")
		assert.Equal(t, len(receiver), 0)
	}(receiver)
}

func TestAllSubscribersGetEvents(t *testing.T) {
	wg := new(sync.WaitGroup)

	wg.Add(2)

	receiver := func(wg *sync.WaitGroup) chan interface{} {
		c := make(chan interface{})

		go func() {
			e := <-c
			assert.Equal(t, e, "data")
			wg.Done()
		}()

		return c
	}

	d := dispatcher.New()
	d.Subscribe("some-type", receiver(wg))
	d.Subscribe("some-type", receiver(wg))

	d.Events <- &dispatcher.Event{
		Type: "some-type",
		Data: "data",
	}

	wg.Wait()
}
