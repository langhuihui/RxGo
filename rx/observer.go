package rx

import "sync"

//Stop 定义一个用于取消订阅的channel，以close该channel为信号
type Stop chan bool

//Observer 用于沟通上下游之间的桥梁，可向下发送数据，向上取消订阅
type Observer struct {
	next NextHandler //缓存当前的NextHandler，后续可以被替换
	//dispose  DisposeHandler //取消订阅的信号，只用来close
	disposed    bool             //是否已经取消
	disposeList []DisposeHandler //缓存的DisposeHandler
	err         error            //缓存当前的错误
	lock        sync.Mutex       //用于Dispose的锁
}

func sinkObserver(next NextHandler, sink *Observer) (result *Observer) {
	result = &Observer{next: next}
	if sink != nil {
		sink.Defer(result)
	}
	return
}
func FuncObserver(next NextFunc, sink *Observer) (result *Observer) {
	return sinkObserver(next, sink)
}
func ChanObserver(next NextChan, sink *Observer) (result *Observer) {
	return sinkObserver(next, sink)
}

//Dispose 取消订阅
func (c *Observer) Dispose() {
	c.lock.Lock()
	if c.disposed {
		return
	}
	c.disposed = true
	c.lock.Unlock()
	for _, handler := range c.disposeList {
		handler.Dispose()
	}
}

//Next 推送数据
func (c *Observer) Next(data interface{}) {
	c.Push(&Event{Data: data, Target: c})
}

//Push 推送数据
func (c *Observer) Push(event *Event) {
	if !c.disposed {
		if event.Target != c {
			event = &Event{Data: event.Data, Target: c}
		}
		c.next.OnNext(event)
	}
}
func (c *Observer) Defer(handler ...DisposeHandler) {
	c.disposeList = append(c.disposeList, handler...)
}
func (c *Observer) AddDisposeFunc(disposeFunc DisposeFunc) DisposeFunc {
	c.Defer(disposeFunc)
	return disposeFunc
}
func (c *Observer) AddDisposeChan() (disposeChan DisposeChan) {
	disposeChan = make(DisposeChan)
	c.Defer(disposeChan)
	return
}
func (c *Observer) Wait() error {
	<-c.AddDisposeChan()
	return c.err
}
