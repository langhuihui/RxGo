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
func (ob Observable) share(childrenCtrl chan *Control) {
	children := make(ControlSet)
	fromSource := func(event *Event) {
		for sink := range children {
			sink.Push(event)
		}
	}
	var source *Control
	for child := range childrenCtrl {
		if child.IsClosed() {
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
		childrenCtrl <- sink
		<-sink.stop
		childrenCtrl <- sink
	}
}

//StartWith 在订阅之前先发送一些数据
func (ob Observable) StartWith(xs ...interface{}) Observable {
	return func(sink *Control) {
		for _, data := range xs {
			sink.Next(data)
			if sink.IsClosed() {
				return
			}
		}
		ob(sink)
	}
}
