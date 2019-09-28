package rx

type Control struct {
	observer Observer
	control  chan Observer
	closed   bool
}

func NewControl(observer Observer) *Control {
	return &Control{
		observer: observer,
		control:  make(chan Observer, 1),
		closed:   false,
	}
}

//Stop 取消订阅
func (c *Control) Stop() {
	if !c.closed {
		c.closed = true
		close(c.control)
	}
}

//Complete 事件流完成
func (c *Control) Complete() {
	c.Stop()
	c.observer(&Event{
		control: c,
		err:     Complete,
	})
}

//Next 推送数据
func (c *Control) Next(data interface{}) {
	c.observer(&Event{
		data:    data,
		control: c,
	})
	if len(c.control) == 0 {
		c.control <- c.observer
	}
}

//Push 推送数据
func (c *Control) Push(event *Event) {
	if event.err != nil {
		c.Stop()
		c.observer(event)
	} else {
		c.observer(event)
		if len(c.control) == 0 {
			c.control <- c.observer
		}
	}
}

//Check 检查是否已取消订阅，并更新观察者回调函数
func (c *Control) Check() (ok bool) {
	c.observer, ok = <-c.control
	return
}
