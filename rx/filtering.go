package rx

//Take 获取最多count数量的事件，然后完成
func (ob Observable) Take(count int) Observable {
	return func(sink *Control) {
		remain := count
		ob.SubscribeS(func(event *Event) {
			if remain--; remain < 0 {
				sink.Complete()
			} else {
				sink.Push(event)
			}
		}, sink.stop) //复用下游的stop信号
	}
}

//TakeUntil 一直获取事件直到unitl传来事件为止
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Control) {
		until.SubscribeA(func(event *Event) {
			if event.err == nil {
				//获取到任何数据就让下游完成
				sink.Complete()
			}
		}, sink.stop)
		ob(sink)
	}
}

//TakeWhile 如果测试函数返回false则完成
func (ob Observable) TakeWhile(f func(interface{}) bool) Observable {
	return func(sink *Control) {
		ob.SubscribeS(func(event *Event) {
			if event.err != nil || f(event.data) {
				sink.Push(event)
			} else {
				sink.Complete()
			}
		}, sink.stop)
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

//SkipWhile 如果测试函数返回false则开始传送
func (ob Observable) SkipWhile(f func(interface{}) bool) Observable {
	return func(sink *Control) {
		ob.SubscribeS(func(event *Event) {
			if event.err != nil {
				sink.Push(event)
			} else if !f(event.data) {
				event.control.observer = sink.observer
			}
		}, sink.stop)
	}
}

//SkipUntil 直到开关事件流发出事件前一直跳过事件
func (ob Observable) SkipUntil(until Observable) Observable {
	return func(sink *Control) {
		source := ob.SubscribeA(func(event *Event) {
			//最初忽略除了完成和错误的事件以外的所有事件
			if event.err != nil {
				sink.Push(event)
			}
		}, sink.stop)
		untilc := until.SubscribeA(func(event *Event) {
			if event.err == nil {
				//获取到任何数据就对接上下游
				source.observer = sink.observer
			}
			//本事件流历史使命已经完成，取消订阅
			event.control.Stop()
		}, make(Stop))
		<-sink.stop
		untilc.Stop()
	}
}

//IgnoreElements 忽略所有元素
func (ob Observable) IgnoreElements() Observable {
	return func(sink *Control) {
		ob.SubscribeS(func(event *Event) {
			//完成和错误事件不忽略，其他都忽略
			if event.err != nil {
				sink.Push(event)
			}
		}, sink.stop)
	}
}
