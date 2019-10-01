package rx

import (
	"io"
)

type (
	Event struct {
		data    interface{}
		err     error
		control *Control
	}
	Observer interface {
		Push(*Event)
	}
	Observable   func(*Control)
	Operator     func(Observable) Observable
	ControlSet   map[*Control]interface{}
	ObserverFunc func(*Event)
	ObserverChan chan *Event
)

func (observer ObserverFunc) Push(event *Event) {
	observer(event)
}
func (observer ObserverChan) Push(event *Event) {
	observer <- event
}

var (
	Complete = io.EOF
)

func (set ControlSet) add(ctrl *Control) {
	set[ctrl] = nil
}
func (set ControlSet) remove(ctrl *Control) {
	delete(set, ctrl)
}
func (set ControlSet) isEmpty() bool {
	return len(set) == 0
}
func (ob Observable) Pipe(cbs ...Operator) Observable {
	for _, cb := range cbs {
		ob = cb(ob)
	}
	return ob
}
