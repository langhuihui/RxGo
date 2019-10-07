package rx

//Stop 定义一个用于取消订阅的channel，以close该channel为信号
type Stop chan bool

//Observer 用于沟通上下游之间的桥梁，可向下发送数据，向上取消订阅
type Observer struct {
	next     NextHandler //缓存当前的NextHandler，后续可以被替换
	dispose  Stop        //取消订阅的信号，只用来close
	complete Stop        //用于发出完成信号
	err      error       //缓存当前的错误
}

func NewObserver(next NextHandler, stop Stop) *Observer {
	return &Observer{next, stop, make(Stop), nil}
}

//Dispose 取消订阅
func (c *Observer) Dispose() {
	if !c.IsDisposed() {
		close(c.dispose)
		c.Complete()
	}
}

//Error 出错
func (c *Observer) Error(err error) {
	c.err = err
	c.Dispose()
}

//IsDisposed 判断是否已经取消订阅
func (c *Observer) IsDisposed() bool {
	select {
	case <-c.dispose:
		return true
	default:
		return false
	}
}

//Wait 等待stop信号
func (c *Observer) Wait(disposeDefer func()) error {
	select {
	case <-c.dispose:
		if disposeDefer != nil {
			disposeDefer()
		}
	case <-c.complete:
	}
	return c.err
}

//Next 推送数据
func (c *Observer) Next(data interface{}) {
	c.Push(&Event{Data: data})
}

//Push 推送数据
func (c *Observer) Push(event *Event) {
	event.Target = c //将事件中的control设置为当前的Control
	if !c.IsDisposed() {
		c.next.OnNext(event)
	}
}

//Complete 完成
func (c *Observer) Complete() {
	select {
	case <-c.complete:
	default:
		close(c.complete)
	}
}
