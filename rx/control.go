package rx

//Stop 定义一个用于取消订阅的channel，以close该channel为信号
type Stop chan bool

//Control 用于沟通上下游之间的桥梁，可向下发送数据，向上取消订阅
type Control struct {
	observer Observer //缓存当前的Observer，后续可以被替换
	stop     Stop     //取消订阅的信号，只用来close
}

func NewControl(observer Observer, stop Stop) *Control {
	return &Control{
		observer: observer,
		stop:     stop,
	}
}

//Stop 取消订阅
func (c *Control) Stop() {
	if !c.IsStopped() {
		close(c.stop)
	}
}

//IsStopped 判断是否已经取消订阅
func (c *Control) IsStopped() bool {
	select {
	case <-c.stop:
		return true
	default:
		return false
	}
}

//Complete 事件流完成
func (c *Control) Complete() {
	c.Push(&Event{err: Complete})
}

//Next 推送数据
func (c *Control) Next(data interface{}) {
	c.Push(&Event{data: data})
}

//Push 推送数据
func (c *Control) Push(event *Event) {
	event.control = c //将事件中的control设置为当前的Control
	if !c.IsStopped() {
		c.observer.Push(event)
	}
}
