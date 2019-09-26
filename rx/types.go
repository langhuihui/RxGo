package rx

import (
	"io"
)

type (
	//Observer 观察者
	Observer func(interface{}, error)
	Sink     chan Observer
	//Observable 可观察对象
	Observable func(Sink)
	Operator   func(Observable) Observable
)

var (
	Complete = io.EOF
)

func (obs Observer) Next(data interface{}) {
	obs(data, nil)
}
func (obs Observer) Complete() {
	obs(nil, Complete)
}
func (obs Observer) Error(err error) {
	obs(nil, err)
}
