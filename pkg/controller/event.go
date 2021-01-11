package controller


type EventType int

const (
	EventOnAdd = iota
	EventOnUpdate
	EventOnDelete
)

type Event struct {
	Type EventType
	Object interface{}
}


type IEventsHook interface {
	IHook
	GetEventChan() <- chan Event
}

type eventsHook struct {
	events chan Event
}

func (eh *eventsHook)OnAdd(obj interface{}){
	eh.events <- Event{
		Type: EventOnAdd,
		Object:  obj,
	}
}

func(eh *eventsHook)OnUpdate(obj interface{}){
	eh.events <- Event{
		Type:   EventOnUpdate,
		Object: obj,
	}
}
func(eh *eventsHook)OnDelete(obj interface{}){
	eh.events <- Event{
		Type:   EventOnDelete,
		Object: obj,
	}
}

func(eh *eventsHook)GetEventChan()<-chan Event{
	return eh.events
}

func NewEventsHook(channelSize int)IEventsHook{
	return &eventsHook{events: make(chan Event, channelSize)}
}