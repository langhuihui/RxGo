package rx

type (
	Event struct {
		Data    interface{}
		Context *Observer
	}
	NextHandler interface {
		OnNext(*Event)
	}

	Observable  func(*Observer) error
	Operator    func(Observable) Observable
	ObserverSet map[*Observer]interface{}
	NextFunc    func(*Event)
	NextCancel  func()
	NextChan    chan *Event
)

func (next NextFunc) OnNext(event *Event) {
	next(event)
}
func (next NextChan) OnNext(event *Event) {
	next <- event
}
func (next NextCancel) OnNext(event *Event) {
	next()
}

func (event *Event) ChangeHandler(observer *Observer) {
	event.Context.next = observer.next
}

var (
	EmptyNext = NextFunc(func(event *Event) {})
)

func (set ObserverSet) add(ctrl *Observer) {
	set[ctrl] = nil
}
func (set ObserverSet) remove(ctrl *Observer) {
	delete(set, ctrl)
}
func (set ObserverSet) isEmpty() bool {
	return len(set) == 0
}
func (ob Observable) Pipe(cbs ...Operator) Observable {
	for _, cb := range cbs {
		ob = cb(ob)
	}
	return ob
}
