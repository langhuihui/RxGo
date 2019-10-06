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
		return ob.subscribe(NextFunc(func(event *Event) {
			sink.Push(event)
			if atomic.AddInt32(&remain, -1) == 0 {
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
		remain := int32(count)
		if remain == 0 {
			return ob(sink)
		}
		return ob.subscribe(NextFunc(func(event *Event) {
			if atomic.AddInt32(&remain, -1) == 0 {
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

//Filter 过滤一些元素
func (ob Observable) Filter(f func(data interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob.subscribe(NextFunc(func(event *Event) {
			if f(event.Data) {
				sink.Push(event)
			}
		}), sink.stop)
	}
}

//Distinct 过滤掉重复出现的元素
func (ob Observable) Distinct() Observable {
	return func(sink *Observer) error {
		buffer := make(map[interface{}]bool)
		return ob.subscribe(NextFunc(func(event *Event) {
			if _, ok := buffer[event.Data]; !ok {
				buffer[event.Data] = true
				sink.Push(event)
			}
		}), sink.stop)
	}
}

//DistinctUntilChanged 过滤掉和前一个元素相同的元素
func (ob Observable) DistinctUntilChanged() Observable {
	return func(sink *Observer) error {
		var lastData interface{}
		return ob.subscribe(NextFunc(func(event *Event) {
			if event.Data != lastData {
				lastData = event.Data
				sink.Push(event)
			}
		}), sink.stop)
	}
}

//Debounce 防抖动
func (ob Observable) Debounce(f func(interface{}) Observable) Observable {
	return func(sink *Observer) error {
		dstop := make(Stop)
		close(dstop)
		//最后要关闭中间可能订阅的防抖动Observable
		defer func() {
			select {
			case <-dstop:
			default:
				close(dstop)
			}
		}()
		return ob.subscribe(NextFunc(func(event *Event) {
			select {
			case <-dstop:
				//开始订阅防抖Observable，等到防抖Observable发出事件或者完成后，就发出现在的事件，期间则忽略任何元素
				dstop = make(Stop)
				go func(event *Event) {
					f(event.Data).subscribe(NextFunc(justStop), dstop)
					sink.Push(event)
				}(event)
			default:
			}
		}), sink.stop)
	}
}

//DebounceTime 按时间防抖动
func (ob Observable) DebounceTime(duration time.Duration) Observable {
	return func(sink *Observer) error {
		debounce := false
		return ob.subscribe(NextFunc(func(event *Event) {
			if !debounce {
				debounce = true
				time.AfterFunc(duration, func() {
					sink.Push(event)
					debounce = false
				})
			}
		}), sink.stop)
	}
}

//Throttle 节流阀
func (ob Observable) Throttle(f func(interface{}) Observable) Observable {
	return func(sink *Observer) error {
		dstop := make(Stop)
		close(dstop)
		//最后要关闭中间可能订阅的防抖动Observable
		defer func() {
			select {
			case <-dstop:
			default:
				close(dstop)
			}
		}()
		return ob.subscribe(NextFunc(func(event *Event) {
			select {
			case <-dstop:
				sink.Push(event)
				dstop = make(Stop)
				go f(event.Data).subscribe(NextFunc(justStop), dstop)
			default:
			}
		}), sink.stop)
	}
}

//ThrottleTime 按照时间来节流
func (ob Observable) ThrottleTime(duration time.Duration) Observable {
	return func(sink *Observer) error {
		throttle := false
		restore := func() {
			throttle = false
		}
		return ob.subscribe(NextFunc(func(event *Event) {
			if !throttle {
				throttle = true
				sink.Push(event)
				time.AfterFunc(duration, restore)
			}
		}), sink.stop)
	}
}

//ElementAt 取第几个元素
func (ob Observable) ElementAt(index uint) Observable {
	return func(sink *Observer) error {
		var count uint = 0
		return ob.subscribe(NextFunc(func(event *Event) {
			if count == index {
				sink.Push(event)
				sink.Stop()
			} else {
				count++
			}
		}), sink.stop)
	}
}

//Find 查询符合条件的元素
func (ob Observable) Find(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		return ob.subscribe(NextFunc(func(event *Event) {
			if f(event.Data) {
				sink.Push(event)
				sink.Stop()
			}
		}), sink.stop)
	}
}

//FindIndex 查找符合条件的元素的序号
func (ob Observable) FindIndex(f func(interface{}) bool) Observable {
	return func(sink *Observer) error {
		index := 0
		return ob.subscribe(NextFunc(func(event *Event) {
			if f(event.Data) {
				sink.Next(index)
				sink.Stop()
			} else {
				index++
			}
		}), sink.stop)
	}
}

//First 完成时返回第一个元素
func (ob Observable) First() Observable {
	return func(sink *Observer) error {
		return ob.subscribe(NextFunc(func(event *Event) {
			sink.Push(event)
			sink.Stop()
		}), sink.stop)
	}
}

//Last 完成时返回最后一个元素
func (ob Observable) Last() Observable {
	return func(sink *Observer) error {
		var last interface{}
		defer func() {
			sink.Next(last)
			sink.Stop()
		}()
		return ob.subscribe(NextFunc(func(event *Event) {
			last = event.Data
		}), sink.stop)
	}
}
