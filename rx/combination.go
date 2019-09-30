package rx

//Merge 合并多个事件流
func Merge(sources ...Observable) Observable {
	return func(sink *Control) {
		scs := make(ControlSet)
		onNext := func(event *Event) {
			if event.err == nil {
				sink.Push(event)
			} else {
				scs.remove(event.control)
				if scs.isEmpty() {
					sink.Complete()
				}
			}
		}
		for _, source := range sources {
			scs.add(source.SubscribeA(onNext, sink.stop))
		}
	}
}

//Concat 连接多个事件流
func Concat(sources ...Observable) Observable {
	return func(sink *Control) {
		remains := sources
		var source *Control
		ob := remains[0]
		remains = remains[1:]
		source = NewControl(func(event *Event) {
			if event.err == Complete && len(remains) > 0 {
				ob = remains[0]
				remains = remains[1:]
				ob(source)
			} else {
				sink.Push(event)
			}
		}, sink.stop)
		ob(source)
	}
}
func (ob Observable) share(childrenCtrl <-chan *Control) {
	children := make(ControlSet)
	fromSource := func(event *Event) {
		for sink := range children {
			sink.Push(event)
		}
	}
	var source *Control
	for child := range childrenCtrl {
		if child.IsStopped() {
			children.remove(child)
			if children.isEmpty() {
				source.Stop()
			}
		} else {
			children.add(child)
			if len(children) == 1 {
				source = ob.SubscribeA(fromSource, make(Stop))
			}
		}
	}
}

//Share 共享数据源
func (ob Observable) Share() Observable {
	childrenCtrl := make(chan *Control)
	go ob.share(childrenCtrl)
	return func(sink *Control) {
		childrenCtrl <- sink //加入观察者
		<-sink.stop
		childrenCtrl <- sink //移除观察者
	}
}

//StartWith 在订阅之前先发送一些数据
func (ob Observable) StartWith(xs ...interface{}) Observable {
	return func(sink *Control) {
		for _, data := range xs {
			sink.Next(data)
			if sink.IsStopped() {
				return
			}
		}
		ob(sink)
	}
}

//CombineLatest 合并多个流的最新数据
func CombineLatest(sources ...Observable) Observable {
	count := len(sources)
	return func(sink *Control) {
		remain := count //尚为有最新数据的源
		live := count   //尚没有完成的数据源
		buffer := make([]interface{}, count)
		e := &Event{
			data:    buffer,
			control: sink,
		}
		for i, ob := range sources {
			buffer[i] = Complete
			ob.SubscribeA(func(event *Event) {
				if event.err == nil {
					if buffer[i] == Complete {
						remain--
					}
					buffer[i] = event.data
					if remain == 0 {
						sink.Push(e)
					}
				} else {
					live--
					//所有数据源没有全部发送过数据或者已经全部都完成了
					if remain > 0 || live == 0 {
						sink.Complete()
					}
				}
			}, sink.stop)
		}
	}
}
