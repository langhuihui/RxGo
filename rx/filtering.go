package rx

//Take 获取最多count数量的事件，然后完成
func (ob Observable) Take(count int) Observable {
	return func(sink *Control) error {
		remain := count
		return ob.subscribe(ObserverFunc(func(event *Event) {
			if remain--; remain < 0 {
				sink.Stop()
			} else {
				sink.Push(event)
			}
		}), sink.stop) //复用下游的stop信号
	}
}

//TakeUntil 一直获取事件直到unitl传来事件为止
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Control) error {
		go until(NewControl(ObserverFunc(func(event *Event) {
			//获取到任何数据就让下游完成
			sink.Stop()
		}), sink.stop))
		return ob(sink)
	}
}

//TakeWhile 如果测试函数返回false则完成
func (ob Observable) TakeWhile(f func(interface{}) bool) Observable {
	return func(sink *Control) error {
		return ob.subscribe(ObserverFunc(func(event *Event) {
			if f(event.Data) {
				sink.Push(event)
			} else {
				sink.Stop()
				//sink.Complete()
			}
		}), sink.stop)
	}
}

//Skip 跳过若干个数据
func (ob Observable) Skip(count int) Observable {
	return func(sink *Control) error {
		remain := count
		return ob.subscribe(ObserverFunc(func(event *Event) {
			if remain--; remain <= 0 {
				//使用下游的Observer代替本函数，使上游数据直接下发到下游
				event.Control.observer = sink.observer
			}
		}), sink.stop) //复用下游的stop信号
	}
}

//SkipWhile 如果测试函数返回false则开始传送
func (ob Observable) SkipWhile(f func(interface{}) bool) Observable {
	return func(sink *Control) error {
		return ob.subscribe(ObserverFunc(func(event *Event) {
			if !f(event.Data) {
				event.Control.observer = sink.observer
			}
		}), sink.stop)
	}
}

//SkipUntil 直到开关事件流发出事件前一直跳过事件
func (ob Observable) SkipUntil(until Observable) Observable {
	return func(sink *Control) error {
		source := NewControl(EmptyObserver, sink.stop)
		untilc := NewControl(ObserverFunc(func(event *Event) {
			//获取到任何数据就对接上下游
			source.observer = sink.observer
			//本事件流历史使命已经完成，取消订阅
			event.Control.Stop()
		}), make(Stop))
		go until(untilc)
		defer untilc.Stop() //上游完成后则终止这个订阅，如果已经终止重复Stop没有影响
		return ob(source)
	}
}

//IgnoreElements 忽略所有元素
func (ob Observable) IgnoreElements() Observable {
	return func(sink *Control) error {
		return ob.subscribe(EmptyObserver, sink.stop)
	}
}
