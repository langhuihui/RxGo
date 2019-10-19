package rx

import (
	"sync/atomic"
	"time"
)

//Take 获取最多count数量的事件，然后完成
func (ob Observable) Take(count uint) Observable {
	return func(sink *Observer) error {
		remain := int32(count)
		if remain == 0 {
			return nil
		}
		return ob(FuncObserver(func(event *Event) {
			sink.Push(event)
			if atomic.AddInt32(&remain, -1) == 0 {
				event.Target.Dispose() //取消订阅上游事件流
			}
		}, sink))
	}
}

//TakeUntil 一直获取事件直到unitl传来事件为止
func (ob Observable) TakeUntil(until Observable) Observable {
	return func(sink *Observer) error {
		observer := &Observer{next: sink.next}
		utilObserver := &Observer{next: NextFunc(func(event *Event) {
			//获取到任何数据就让下游完成
			event.Target.Dispose()
			observer.Dispose()
		})}
		go until(utilObserver)
		sink.Defer(utilObserver, observer)
		return ob(observer)
	}
}

//TakeWhile 如果测试函数返回false则完成
func (ob Observable) TakeWhile(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob(FuncObserver(func(event *Event) {
			if f(event.Data) {
				sink.Push(event)
			} else {
				event.Target.Dispose() //取消订阅上游事件流
			}
		}, sink))
	}
}

//Skip 跳过若干个数据
func (ob Observable) Skip(count uint) Observable {
	return func(sink *Observer) error {
		remain := int32(count)
		if remain == 0 {
			return ob(sink)
		}
		return ob(FuncObserver(func(event *Event) {
			if atomic.AddInt32(&remain, -1) == 0 {
				//使用下游的Observer代替本函数，使上游数据直接下发到下游
				event.ChangeHandler(sink)
			}
		}, sink))
	}
}

//SkipWhile 如果测试函数返回false则开始传送
func (ob Observable) SkipWhile(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob(FuncObserver(func(event *Event) {
			if !f(event.Data) {
				event.ChangeHandler(sink)
			}
		}, sink))
	}
}

//SkipUntil 直到开关事件流发出事件前一直跳过事件
func (ob Observable) SkipUntil(until Observable) Observable {
	return func(sink *Observer) error {
		source := FuncObserver(EmptyNext, sink) //前期跳过所有数据
		untilc := FuncObserver(func(event *Event) {
			//获取到任何数据就对接上下游
			source.next = sink.next
			//本事件流历史使命已经完成，取消订阅
			event.Target.Dispose()
		}, sink)
		go until(untilc)
		defer untilc.Dispose() //上游完成后则终止这个订阅，如果已经终止重复Dispose没有影响
		return ob(source)
	}
}

//IgnoreElements 忽略所有元素
func (ob Observable) IgnoreElements() Observable {
	return func(sink *Observer) error {
		return ob(FuncObserver(EmptyNext, sink))
	}
}

//Filter 过滤一些元素
func (ob Observable) Filter(f func(data interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob(FuncObserver(func(event *Event) {
			if f(event.Data) {
				sink.Push(event)
			}
		}, sink))
	}
}

//Distinct 过滤掉重复出现的元素
func (ob Observable) Distinct() Observable {
	return func(sink *Observer) error {
		buffer := make(map[interface{}]bool)
		next := make(NextChan)
		go func() {
			for event := range next {
				if _, ok := buffer[event.Data]; !ok {
					buffer[event.Data] = true
					sink.Push(event)
				}
			}
		}()
		defer close(next)
		return ob(ChanObserver(next, sink))
	}
}

//DistinctUntilChanged 过滤掉和前一个元素相同的元素
func (ob Observable) DistinctUntilChanged() Observable {
	return func(sink *Observer) error {
		var lastData interface{}
		return ob(FuncObserver(func(event *Event) {
			if event.Data != lastData {
				lastData = event.Data
				sink.Push(event)
			}
		}, sink))
	}
}

//Debounce 防抖动
func (ob Observable) Debounce(f func(interface{}) Observable) Observable {
	return func(sink *Observer) error {
		throttles := make(chan *Event, 1) //一个缓冲，保证不会阻塞
		var throttle *Observer
		go func() {
			for event := range throttles {
				f(event.Data)(throttle)
				sink.Push(event)
				throttle.Dispose()
			}
		}()
		observer := FuncObserver(func(event *Event) {
			if throttle == nil || throttle.disposed {
				throttle = FuncObserver(func(event *Event) {
					event.Target.Dispose()
				}, nil)
				throttles <- event
			}
		}, nil)
		sink.AddDisposeFunc(func() {
			if throttle != nil {
				throttle.Dispose()
			}
			observer.Dispose()
		})
		return ob(observer)
	}
}

//DebounceTime 按时间防抖动
func (ob Observable) DebounceTime(duration time.Duration) Observable {
	return func(sink *Observer) error {
		debounce := false
		return ob(FuncObserver(func(event *Event) {
			if !debounce {
				debounce = true
				time.AfterFunc(duration, func() {
					sink.Push(event)
					debounce = false
				})
			}
		}, sink))
	}
}

//Throttle 节流阀
func (ob Observable) Throttle(f func(interface{}) Observable) Observable {
	return func(sink *Observer) error {
		throttles := make(chan *Event, 1) //一个缓冲，保证不会阻塞
		throttle := &Observer{
			next: NextFunc(func(event *Event) {
				event.Target.Dispose()
			}),
			disposed: true,
		}
		sink.Defer(throttle)
		go func() {
			for event := range throttles {
				f(event.Data)(throttle)
			}
		}()
		defer throttle.Dispose()
		return ob(FuncObserver(func(event *Event) {
			if throttle.disposed {
				throttle.disposed = false
				throttle.disposeList = nil
				sink.Push(event)
				throttles <- event
			}
		}, sink))
	}
}

//ThrottleTime 按照时间来节流
func (ob Observable) ThrottleTime(duration time.Duration) Observable {
	return func(sink *Observer) error {
		throttle := false
		restore := func() {
			throttle = false
		}
		return ob(FuncObserver(func(event *Event) {
			if !throttle {
				throttle = true
				sink.Push(event)
				time.AfterFunc(duration, restore)
			}
		}, sink))
	}
}

//ElementAt 取第几个元素
func (ob Observable) ElementAt(index uint) Observable {
	return func(sink *Observer) error {
		var count uint = 0
		return ob(FuncObserver(func(event *Event) {
			if count == index {
				sink.Push(event)
				event.Target.Dispose()
			} else {
				count++
			}
		}, sink))
	}
}

//Find 查询符合条件的元素
func (ob Observable) Find(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob(FuncObserver(func(event *Event) {
			if f(event.Data) {
				sink.Push(event)
				event.Target.Dispose()
			}
		}, sink))
	}
}

//FindIndex 查找符合条件的元素的序号
func (ob Observable) FindIndex(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		index := 0
		return ob(FuncObserver(func(event *Event) {
			if f(event.Data) {
				sink.Next(index)
				event.Target.Dispose()
			} else {
				index++
			}
		}, sink))
	}
}

//First 完成时返回第一个元素
func (ob Observable) First() Observable {
	return func(sink *Observer) error {
		return ob(FuncObserver(func(event *Event) {
			sink.Push(event)
			event.Target.Dispose()
		}, sink))
	}
}

//Last 完成时返回最后一个元素
func (ob Observable) Last() Observable {
	return func(sink *Observer) error {
		var last interface{}
		defer func() {
			sink.Next(last)
		}()
		return ob(FuncObserver(func(event *Event) {
			last = event.Data
		}, sink))
	}
}
