package rx

//SubscribeS 同步订阅Observable
func (ob Observable) SubscribeS(observer Observer, stop Stop) {
	ob(NewControl(observer, stop))
}

//SubscribeA 异步订阅Observable
func (ob Observable) SubscribeA(observer Observer, stop Stop) (control *Control) {
	control = NewControl(observer, stop)
	go ob(control)
	return
}

//SubscribeAsync 异步订阅模式，返回一个可用于终止事件流的函数
func (ob Observable) SubscribeAsync(onNext func(data interface{}), onComplete func(), onError func(error)) func() {
	source := ob.SubscribeA(func(event *Event) {
		if event.err == nil {
			onNext(event.data)
		} else if event.err == Complete {
			onComplete()
			event.control.Stop()
		} else {
			onError(event.err)
			event.control.Stop()
		}
	}, make(Stop))
	return source.Stop
}

//SubscribeOnCurrent 同步订阅模式，会阻塞当前goroutine，函数返回代表完成
func (ob Observable) SubscribeSync(onNext func(data interface{})) (err error) {
	next := make(chan interface{})
	complete := make(chan error)
	defer close(next)
	defer close(complete)
	source := ob.SubscribeA(func(event *Event) {
		if event.err == nil {
			next <- event.data
		} else if event.err == Complete {
			complete <- nil
		} else {
			complete <- err
		}
	}, make(Stop))
	for {
		select {
		case data := <-next:
			onNext(data)
		case err = <-complete:
			source.Stop()
			return
		}
	}
}
