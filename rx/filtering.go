package rx

//Take 获取最多count数量的事件，然后完成
func (ob Observable) Take(count uint) Observable {
	return func(sink *Observer) error {
		remain := count
		if remain == 0 {
			return nil
		}
		return ob.subscribe(NextFunc(func(event *Event) {
			sink.Push(event)
			if remain--; remain == 0 {
				sink.Stop()
			}
		}), sink.stop) //复用下游的stop信号
	}
}

//TakeUntil 一直获取事件直到unitl传来事件为止
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Observer) error {
		go until(NewObserver(NextFunc(func(event *Event) {
			//获取到任何数据就让下游完成
			sink.Stop()
		}), sink.stop))
		return ob(sink)
	}
}

//TakeWhile 如果测试函数返回false则完成
func (ob Observable) TakeWhile(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob.subscribe(NextFunc(func(event *Event) {
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
func (ob Observable) Skip(count uint) Observable {
	return func(sink *Observer) error {
		remain := count
		if remain == 0 {
			return ob(sink)
		}
		return ob.subscribe(NextFunc(func(event *Event) {
			if remain--; remain == 0 {
				//使用下游的Observer代替本函数，使上游数据直接下发到下游
				event.ChangeHandler(sink)
			}
		}), sink.stop) //复用下游的stop信号
	}
}

//SkipWhile 如果测试函数返回false则开始传送
func (ob Observable) SkipWhile(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob.subscribe(NextFunc(func(event *Event) {
			if !f(event.Data) {
				event.ChangeHandler(sink)
			}
		}), sink.stop)
	}
}

//SkipUntil 直到开关事件流发出事件前一直跳过事件
func (ob Observable) SkipUntil(until Observable) Observable {
	return func(sink *Observer) error {
		source := NewObserver(EmptyNext, sink.stop) //前期跳过所有数据
		untilc := NewObserver(NextFunc(func(event *Event) {
			//获取到任何数据就对接上下游
			source.next = sink.next
			//本事件流历史使命已经完成，取消订阅
			event.Target.Stop()
		}), make(Stop))
		go until(untilc)
		defer untilc.Stop() //上游完成后则终止这个订阅，如果已经终止重复Stop没有影响
		return ob(source)
	}
}

//IgnoreElements 忽略所有元素
func (ob Observable) IgnoreElements() Observable {
	return func(sink *Observer) error {
		return ob.subscribe(EmptyNext, sink.stop)
	}
}
