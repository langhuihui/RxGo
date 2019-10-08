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

func NewObserver(next NextHandler, dispose Stop, complete Stop) *Observer {
	return &Observer{next, dispose, complete, nil}
}

//New 派生一个Observer，flag 为复用策略，0：不复用，1：只复用dispose，2：只复用complete，3：都复用
func (c *Observer) New(next NextHandler, flag int) *Observer {
	dispose := c.dispose
	complete := c.complete
	switch {
	case flag&1 == 0: //0、2
		dispose = make(Stop)
	case flag&2 == 0: //0、1
		complete = make(Stop)
	}
	return &Observer{next, dispose, complete, nil}
}
func (c *Observer) New0(next NextHandler) *Observer {
	return c.New(next, 0)
}
func (c *Observer) New1(next NextHandler) *Observer {
	return c.New(next, 1)
}
func (c *Observer) New2(next NextHandler) *Observer {
	return c.New(next, 2)
}

//New3，dispose和complete信号都复用，即下游可以迫使上行链路完成和取消
func (c *Observer) New3(next NextHandler) *Observer {
	return c.New(next, 3)
}

//Dispose 取消订阅
func (c *Observer) Dispose() {
	select {
	case <-c.dispose:
	default:
		close(c.dispose)
	}
}

//Error 出错
func (c *Observer) Error(err error) {
	c.err = err
	c.Complete()
}

//Aborted 判断是否已经取消订阅或者已完成
func (c *Observer) Aborted() bool {
	select {
	case <-c.dispose:
		return true
	case <-c.complete:
		return true
	default:
		return false
	}
}

//Wait 等待stop信号
func (c *Observer) Wait() error {
	select {
	case <-c.dispose:
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
	if !c.Aborted() {
		c.next.OnNext(event)
	}
}

//Complete 迫使Observable完成
func (c *Observer) Complete() {
	select {
	case <-c.complete:
	default:
		close(c.complete)
	}
}
func (c *Observer) Completed() bool {
	select {
	case <-c.complete:
		return true
	default:
		return false
	}
}
