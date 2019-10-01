package rx

import "errors"

//Merge 合并多个事件流
func Merge(sources ...Observable) Observable {
	count := len(sources)
	return func(sink *Control) (err error) {
		remain := count
		source := NewControl(ObserverFunc(sink.Push), make(Stop))
		for _, ob := range sources {
			go func() {
				ob(source)
				remain--
				if remain == 0 {
					sink.Stop()
				}
			}()
		}
		defer source.Stop()
		return sink.Wait()
	}
}

//Concat 连接多个事件流
func Concat(sources ...Observable) Observable {
	return func(sink *Control) (err error) {
		for _, ob := range sources {
			err = ob(sink)
			//如果出现取消订阅或者出现错误，都进行完成动作，否则是正常完成，接着订阅下一个数据源
			if sink.IsStopped() || err != nil {
				return err
			}
		}
		return //全部完成都没有错误，则正常完成
	}
}
func (ob Observable) share(childrenCtrl <-chan *Control) {
	children := make(ControlSet)
	eventChan := make(ObserverChan)
	var source *Control
	var sourceError error
	for {
		select {
		case child := <-childrenCtrl:
			if child.IsStopped() {
				children.remove(child)
				if children.isEmpty() {
					source.Stop()
				}
			} else {
				children.add(child)
				if len(children) == 1 {
					source = NewControl(eventChan, make(Stop))
					go func() {
						sourceError = ob(source)
						close(eventChan)
					}()
				}
			}
		case event, ok := <-eventChan:
			if ok {
				for sink := range children {
					sink.Push(event)
				}
			} else {
				for sink := range children {
					sink.Error(sourceError)
				}
			}
		}
	}
}

//Share 共享数据源
func (ob Observable) Share() Observable {
	childrenCtrl := make(chan *Control)
	go ob.share(childrenCtrl)
	return func(sink *Control) error {
		childrenCtrl <- sink //加入观察者
		defer func() {
			childrenCtrl <- sink //移除观察者
		}()
		return sink.Wait()
	}
}

//StartWith 在订阅之前先发送一些数据
func (ob Observable) StartWith(xs ...interface{}) Observable {
	return func(sink *Control) error {
		for _, data := range xs {
			sink.Next(data)
			if sink.IsStopped() {
				return sink.err
			}
		}
		return ob(sink)
	}
}

//CombineLatest 合并多个流的最新数据
func CombineLatest(sources ...Observable) Observable {
	count := len(sources)
	NoData := errors.New("NoData")
	return func(sink *Control) error {
		remain := count                      //尚未有最新数据的源
		live := count                        //尚没有完成的数据源
		buffer := make([]interface{}, count) //待发送的数据缓存
		e := &Event{buffer, sink}
		controls := make([]*Control, count)
		for i, ob := range sources {
			buffer[i] = NoData //将该坑位设置为尚未填充数据状态
			controls[i] = NewControl(ObserverFunc(func(event *Event) {
				if buffer[i] == NoData {
					remain--
				}
				buffer[i] = event.Data
				if remain == 0 {
					sink.Push(e)
				}
			}), make(Stop))
			go func(source *Control) {
				ob(source)
				if live--; remain > 0 || live == 0 {
					sink.Stop()
				}
			}(controls[i])
		}
		defer func() {
			for _, source := range controls {
				source.Stop()
			}
		}()
		return sink.Wait()
	}
}
