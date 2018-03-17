package dispatcher

type Dispatcher struct {
	Events    chan *Event
	Listeners map[string][]chan interface{}
}

func New() *Dispatcher {
	d := &Dispatcher{
		Events:    make(chan *Event),
		Listeners: make(map[string][]chan interface{}),
	}

	go d.start()
	return d
}

func (d *Dispatcher) start() {
	for {
		go d.distributeEvent(<-d.Events)
	}
}

func (d *Dispatcher) distributeEvent(e *Event) {
	listeners := d.Listeners[e.Type]
	for _, l := range listeners {
		l <- e.Data
	}
}

func (d *Dispatcher) Subscribe(key string, listener chan interface{}) {
	listeners := d.Listeners[key]
	d.Listeners[key] = append(listeners, listener)
}
