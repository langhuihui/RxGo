package rx

//Take 获取最多count数量的事件，然后完成
func (ob Observable) Take(count int) Observable {
	return func(sink *Control) {
		remain := count
		ob.SubscribeS(func(event *Event) {
			if remain--; remain < 0 {
				sink.Complete()
				sink.Stop()
			} else {
				sink.Push(event)
			}
		}, sink.stop) //复用下游的stop信号
	}
}

//Skip 跳过若干个数据
func (ob Observable) Skip(count int) Observable {
	return func(sink *Control) {
		remain := count
		ob.SubscribeS(func(event *Event) {
			if remain--; remain <= 0 {
				//使用下游的Observer代替本函数，使上游数据直接下发到下游
				event.control.observer = sink.observer
			}
		}, sink.stop) //复用下游的stop信号
	}
}

//TakeUntil 一直获取事件直到unitl传来事件为止
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Control) {
		until.SubscribeA(func(event *Event) {
			if event.err != nil {
				//获取到任何数据就让下游完成
				sink.Complete()
				sink.Stop()
			}
		}, sink.stop)
		ob(sink)
	}
}
