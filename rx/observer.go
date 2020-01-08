package rx

import (
	"context"
)

//Stop 定义一个用于取消订阅的channel，以close该channel为信号
type Stop chan bool

//Observer 用于沟通上下游之间的桥梁，可向下发送数据，向上取消订阅
type Observer struct {
	context.Context //组合继承方式
	cancel          context.CancelFunc
	next            NextHandler //缓存当前的NextHandler，后续可以被替换
}

func (c *Observer) CreateFuncObserver(next NextFunc) *Observer {
	return &Observer{c, c.cancel, next}
}
func (c *Observer) NewFuncObserver(next NextFunc) (result *Observer) {
	result = &Observer{nil, nil, next}
	result.Context, result.cancel = context.WithCancel(c)
	return
}
func (c *Observer) CreateChanObserver(next NextChan) *Observer {
	return &Observer{c, c.cancel, next}
}
func (c *Observer) NewChanObserver(next NextChan) (result *Observer) {
	result = &Observer{nil, nil, next}
	result.Context, result.cancel = context.WithCancel(c)
	return
}
func (c *Observer) IsDisposed() bool {
	return c.Err() != nil
}

//Next 推送数据
func (c *Observer) Next(data interface{}) {
	c.Push(&Event{Data: data, Context: c})
}

//Push 推送数据
func (c *Observer) Push(event *Event) {
	if !c.IsDisposed() {
		if event.Context != c {
			event = &Event{Data: event.Data, Context: c}
		}
		c.next.OnNext(event)
	}
}
