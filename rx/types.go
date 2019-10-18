package rx

type (
	Event struct {
		Data   interface{}
		Target *Observer
	}
	NextHandler interface {
		OnNext(*Event)
	}
	DisposeHandler interface {
		Dispose()
	}
	Observable  func(*Observer) error
	Operator    func(Observable) Observable
	ObserverSet map[*Observer]interface{}
	NextFunc    func(*Event)
	NextChan    chan *Event
	DisposeFunc func()
	DisposeChan chan Stop
	DisposeList struct {
		list []DisposeHandler
	}
)

func (d *DisposeList) Dispose() {
	for _, handler := range d.list {
		handler.Dispose()
	}
}
func (d *DisposeList) Add(handler ...DisposeHandler) {
	d.list = append(d.list, handler...)
}
func (next NextFunc) OnNext(event *Event) {
	next(event)
}
func (next NextChan) OnNext(event *Event) {
	next <- event
}
func (dispose DisposeFunc) Dispose() {
	dispose()
}
func (dispose DisposeChan) Dispose() {
	select {
	case <-dispose:
	default:
		close(dispose)
	}
}
func (event *Event) ChangeHandler(observer *Observer) {
	event.Target.next = observer.next
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
