package rx

//Take 获取最多count数量的事件，然后完成
func (ob Observable) Take(count int) Observable {
	return func(sink *Control) {
		remain := count
		source := ob.Subscribe(func(event *Event) {
			if remain--; remain < 0 {
				sink.Complete()
			} else {
				sink.Push(event)
			}
		})
		for sink.Check() {
		}
		source.Stop()
	}
}

//Skip 跳过若干个数据
func (ob Observable) Skip(count int) Observable {
	return func(sink *Control) {
		remain := count
		source := ob.Subscribe(func(event *Event) {
			if remain--; remain <= 0 {
				event.control.control <- sink.observer
			}
		})
		for sink.Check() {
		}
		source.Stop()
	}
}

//TakeUntil 一直获取事件直到unitl传来事件为止
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Control) {
		until.Subscribe(func(event *Event) {
			if event.err != nil {
				sink.Complete()
				event.control.Stop()
			}
		})
		source := ob.Subscribe(sink.Push)
		for sink.Check() {
		}
		source.Stop()
	}
}
