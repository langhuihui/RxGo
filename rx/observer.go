package rx

//Stop 定义一个用于取消订阅的channel，以close该channel为信号
type Stop chan bool

//Observer 用于沟通上下游之间的桥梁，可向下发送数据，向上取消订阅
type Observer struct {
	next NextHandler //缓存当前的NextHandler，后续可以被替换
	stop Stop        //取消订阅的信号，只用来close
	err  error       //缓存当前的错误
}

func NewObserver(next NextHandler, stop Stop) *Observer {
	return &Observer{next, stop, nil}
}

//Stop 取消订阅
func (c *Observer) Stop() {
	if !c.IsStopped() {
		close(c.stop)
	}
}

//Error 出错
func (c *Observer) Error(err error) {
	c.err = err
	c.Stop()
}

//IsStopped 判断是否已经取消订阅
func (c *Observer) IsStopped() bool {
	select {
	case <-c.stop:
		return true
	default:
		return false
	}
}

//Wait 等待stop信号
func (c *Observer) Wait() error {
	<-c.stop
	return c.err
}

//Next 推送数据
func (c *Observer) Next(data interface{}) {
	c.Push(&Event{Data: data})
}

//Push 推送数据
func (c *Observer) Push(event *Event) {
	event.Target = c //将事件中的control设置为当前的Control
	if !c.IsStopped() {
		c.next.OnNext(event)
	}
}
