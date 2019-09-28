package rx

//Subscribe 订阅Observable，后续对sink发送Observer来拉取数据
func (ob Observable) Subscribe(observer Observer) (control *Control) {
	control = NewControl(observer)
	go ob(control)
	return
}

//SubscribeAsync 异步订阅模式，返回一个可用于终止事件流的函数
func (ob Observable) SubscribeAsync(onNext func(data interface{}), onComplete func(), onError func(error)) func() {
	source := ob.Subscribe(func(event *Event) {
		if event.err == nil {
			onNext(event.data)
		} else if event.err == Complete {
			onComplete()
		} else {
			onError(event.err)
		}
	})
	return source.Stop
}

//SubscribeOnCurrent 同步订阅模式，会阻塞当前goroutine，函数返回代表完成
func (ob Observable) SubscribeSync(onNext func(data interface{})) (err error) {
	next := make(chan interface{})
	complete := make(chan error)
	defer close(next)
	defer close(complete)
	ob.Subscribe(func(event *Event) {
		if event.err == nil {
			next <- event.data
		} else if event.err == Complete {
			complete <- nil
		} else {
			complete <- err
		}
	})
	for {
		select {
		case data := <-next:
			onNext(data)
		case err = <-complete:
			return
		}
	}
}
