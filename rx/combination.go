package rx

import (
	"errors"
	"sync"
	"sync/atomic"
)

//Merge 合并多个事件流
func Merge(sources ...Observable) Observable {
	count := len(sources)
	return func(sink *Observer) (err error) {
		source := FuncObserver(sink.Push, sink)
		wg := &sync.WaitGroup{}
		wg.Add(count)
		for _, ob := range sources {
			go func(ob Observable) {
				ob(source)
				wg.Done()
			}(ob)
		}
		wg.Wait()
		source.Dispose()
		return sink.err
	}
}

//Concat 连接多个事件流
func Concat(sources ...Observable) Observable {
	return func(sink *Observer) (err error) {
		for _, ob := range sources {
			err = ob(sink)
			//如果出现取消订阅或者出现错误，都进行完成动作，否则是正常完成，接着订阅下一个数据源
			if sink.disposed || err != nil {
				return err
			}
		}
		return //全部完成都没有错误，则正常完成
	}
}
func (ob Observable) share(childrenCtrl <-chan *Observer) {
	children := make(ObserverSet)
	eventChan := make(NextChan)
	completeChan := make(chan error)
	var source *Observer //上游观察者
	for {
		select {
		case child := <-childrenCtrl:
			if child.disposed {
				children.remove(child)
				if children.isEmpty() { //最后一个子观察者退出则取消订阅上游事件流
					source.Dispose()
				}
			} else {
				children.add(child)
				if len(children) == 1 {
					source = &Observer{next: eventChan}
					go func() {
						completeChan <- ob(source)
					}()
				}
			}
		case event := <-eventChan:
			for sink := range children {
				sink.Push(event)
			}
		case sourceError := <-completeChan:
			//上游事件流完成，则使得所有子观察者全部完成
			for sink := range children {
				sink.err = sourceError
				sink.Dispose()
				//children.remove(sink)
			}
		}
	}
}

//Share 共享数据源
func (ob Observable) Share() Observable {
	childrenCtrl := make(chan *Observer)
	go ob.share(childrenCtrl)
	return func(sink *Observer) error {
		childrenCtrl <- sink //加入观察者
		<-sink.AddDisposeChan()
		childrenCtrl <- sink //移除观察者
		return sink.err
	}
}

//StartWith 在订阅之前先发送一些数据
func (ob Observable) StartWith(xs ...interface{}) Observable {
	return func(sink *Observer) error {
		for _, data := range xs {
			sink.Next(data)
			//如果在响应事件中终止了事件流，则函数返回退出
			if sink.disposed {
				return sink.err
			}
		}
		//如果用户没有终止当前事件流，则订阅上游事件流，并移交控制权
		return ob(sink)
	}
}

//CombineLatest 合并多个流的最新数据
func CombineLatest(sources ...Observable) Observable {
	count := len(sources)
	NoData := errors.New("NoData")
	return func(sink *Observer) error {
		remain := int32(count) //尚未有最新数据的源
		wg := sync.WaitGroup{}
		wg.Add(count)
		buffer := make([]interface{}, count) //待发送的数据缓存
		e := &Event{buffer, sink}
		for i := range sources {
			buffer[i] = NoData //将该坑位设置为尚未填充数据状态
			go func(j int) {
				sources[j](FuncObserver(func(event *Event) {
					if buffer[j] == NoData {
						atomic.AddInt32(&remain, -1)
					}
					buffer[j] = event.Data //缓存最新数据
					if remain == 0 {       //当所有事件流都产生了事件后，就发送事件
						sink.Push(e)
					}
				}, sink))
				wg.Done()
				//当有事件流完成后，如果此时还有事件流从未发送过事件或者所有事件流都已经完成，则完成当前事件流
				if remain > 0 {
					sink.Dispose()
				}
			}(i)
		}
		wg.Wait()
		return sink.err
	}
}

//Zip 将多个事件源的事件按顺序组合
func Zip(sources ...Observable) Observable {
	count := len(sources)
	return func(sink *Observer) error {
		remain := int32(count)
		input := make([]chan interface{}, count) //缓存每一个事件源的一个事件
		buffer := make([]interface{}, count)     //待发送的数据缓存
		e := &Event{buffer, sink}                //每次发送的事件对象
		for i := range sources {
			input[i] = make(chan interface{}, 1)
			go func(i int) {
				sink.err = sources[i](FuncObserver(func(event *Event) {
					input[i] <- event.Data //每次只缓存一个数据，然后阻塞等待一下一组
					if atomic.AddInt32(&remain, -1) == 0 {
						remain = int32(count)
						for j, dataChan := range input {
							buffer[j] = <-dataChan //把每一个事件源缓存的数据组合
						}
						sink.Push(e)
					}
				}, sink))
				sink.Dispose() //任意一个事件源完成就导致整个事件流停止
			}(i)
		}
		defer func() {
			//清除所有的临时用的channel
			for _, dataChan := range input {
				close(dataChan)
			}
		}()
		return sink.Wait() //用户取消订阅，会导致取消所有事件流的订阅
	}
}

//Race 使用最先收到事件的那个事件流
func Race(sources ...Observable) Observable {
	count := len(sources)
	return func(sink *Observer) error {
		observers := make([]*Observer, count)
		var winner *Observer
		next := make(NextChan)
		for i := range sources {
			observers[i] = &Observer{next: next}
			sink.Defer(observers[i])
		}
		//分成两次循环的目的，是防止并发读写observers
		for i := range sources {
			go func(i int) {
				sink.err = sources[i](observers[i])
				//如果有一个事件流完成，但从未发出过事件比如Empty，或者这个事件流就是最快的事件流，那么就导致事件流完成
				if winner == nil || winner == observers[i] {
					sink.Dispose()
				}
			}(i)
		}
		defer close(next)
		dispose := sink.AddDisposeChan()
		for {
			select {
			case <-dispose:
				return nil
			case event := <-next:
				if winner == nil { //有人率先到达了
					winner = event.Target
					//其他事件流此时可以统统取消
					for _, ob := range observers {
						if ob != winner {
							ob.Dispose()
						}
					}
					//转发事件
					sink.Push(event)
					//移交接收器
					winner.next = NextFunc(sink.Push)
				} else if winner == event.Target {
					sink.Push(event)
				}
			}
		}
		return sink.err
	}
}
