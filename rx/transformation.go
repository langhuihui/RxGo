package rx

import (
	"context"
	"sync"
)

//Map 对元素做一层映射（转换）
func (ob Observable) Map(f func(interface{}) interface{}) Observable {
	return func(sink *Observer) error {
		return ob(sink.CreateFuncObserver(func(event *Event) {
			sink.Next(f(event.Data))
		}))
	}
}

//MapTo 映射成一个固定值
func (ob Observable) MapTo(data interface{}) Observable {
	return func(sink *Observer) error {
		return ob(sink.CreateFuncObserver(func(event *Event) {
			sink.Next(data)
		}))
	}
}

//MergeMap 将元素映射成事件流并合并其中的元素
func (ob Observable) MergeMap(f func(interface{}) Observable, resultSelector func(interface{}, interface{}) interface{}) Observable {
	return func(sink *Observer) (err error) {
		wg := sync.WaitGroup{}
		wg.Add(1) //母事件流等待信号
		subMap := func(data interface{}) {
			f(data)(sink.CreateFuncObserver(sink.Push))
			wg.Done()
		}
		if resultSelector != nil { //使用结果选择器
			subMap = func(data interface{}) {
				f(data)(sink.CreateFuncObserver(func(event *Event) {
					sink.Next(resultSelector(data, event.Data))
				}))
				wg.Done()
			}
		}
		go func() {
			err = ob(sink.CreateFuncObserver(func(event *Event) {
				wg.Add(1)
				go subMap(event.Data)
			}))
			wg.Done()
		}()
		wg.Wait() //等待所有子事件流完成后再完成
		return
	}
}

//MergeMapTo  将元素映射成指定的事件流并合并其中的元素
func (ob Observable) MergeMapTo(then Observable, resultSelector func(interface{}, interface{}) interface{}) Observable {
	return ob.MergeMap(func(i interface{}) Observable {
		return then
	}, resultSelector)
}

//SwitchMap 将元素映射成事件流然后发送其中的元素
func (ob Observable) SwitchMap(f func(interface{}) Observable, resultSelector func(interface{}, interface{}) interface{}) Observable {
	return func(sink *Observer) (err error) {
		currentSub := new(Observer)
		ctxSub := context.Background()
		wg := sync.WaitGroup{}
		subMap := func(data interface{}, currentSub *Observer) {
			f(data)(currentSub)
			wg.Done() //关闭Observer用来触发主事件流完成
		}
		wg.Add(1)
		go func() {
			err = ob(sink.CreateFuncObserver(func(event *Event) {
				if currentSub.cancel != nil { //关闭上一个子事件流
					currentSub.cancel()
				}
				if resultSelector != nil {
					data := event.Data
					currentSub.next = NextFunc(func(event *Event) {
						sink.Next(resultSelector(data, event.Data))
					})
				} else {
					currentSub.next = NextFunc(sink.Push)
				}
				currentSub.Context, currentSub.cancel = context.WithCancel(ctxSub)
				wg.Add(1)
				go subMap(event.Data, currentSub)
			}))
			wg.Done()
		}()
		wg.Wait()
		if currentSub.cancel != nil {
			currentSub.cancel()
		}
		return
	}
}

//SwitchMapTo  将元素映射成指定的事件流然后发送其中的元素
func (ob Observable) SwitchMapTo(then Observable, resultSelector func(interface{}, interface{}) interface{}) Observable {
	return ob.SwitchMap(func(i interface{}) Observable {
		return then
	}, resultSelector)
}

//Scan 类似于Reduce，只是在每次元素达到后进行计算和发射
func (ob Observable) Scan(f func(interface{}, interface{}) interface{}) Observable {
	return func(sink *Observer) error {
		var aac interface{}
		aacNext := func(event *Event) {
			aac = f(aac, event.Data)
			sink.Next(aac)
		}
		return ob(sink.CreateFuncObserver(func(event *Event) {
			aac = event.Data
			sink.Push(event)
			event.next = NextFunc(aacNext)
		}))
	}
}

//Repeat 重复若干次上游的事件流元素
func (ob Observable) Repeat(count int) Observable {
	return func(sink *Observer) (err error) {
		var cache []*Event
		if err = ob(sink.CreateFuncObserver(func(event *Event) {
			cache = append(cache, event)
			sink.Push(event)
		})); err == nil {
			for remain := count; remain >= 0; remain-- {
				for _, event := range cache {
					sink.Push(event)
				}
			}
		}
		return
	}
}

//PairWise 两两组成一组
func (ob Observable) PairWise() Observable {
	return func(sink *Observer) error {
		var cache []interface{}
		return ob(sink.CreateFuncObserver(func(event *Event) {
			if cache == nil {
				cache = append(cache, event.Data)
			} else {
				sink.Next(append(cache, event.Data))
				cache = []interface{}{event.Data}
			}
		}))
	}
}

//Buffer 缓冲元素直到被另一个事件触发再发送缓冲中的元素
func (ob Observable) Buffer(closingNotifier Observable) Observable {
	return func(sink *Observer) error {
		var buffer []interface{}
		closingNext := make(NextChan)
		obNext := make(NextChan)
		closingObs := sink.CreateChanObserver(closingNext)
		obs := sink.CreateChanObserver(obNext)
		go func() { //多路复用防止同时读写buffer对象
			for {
				select {
				case event, ok := <-obNext:
					if ok {
						buffer = append(buffer, event.Data)
					} else {
						return
					}
				case _, ok := <-closingNext:
					if ok {
						sink.Next(buffer)
						buffer = make([]interface{}, 0)
					} else {
						return
					}
				}
			}
		}()
		go func() {
			closingNotifier(closingObs)
			close(closingNext)
			obs.cancel() //如果控制事件流先完成，则终止当前事件流
		}()
		defer closingObs.cancel() //如果当前事件流先完成，则终止控制事件流
		defer close(obNext)
		return ob(obs)
	}
}
