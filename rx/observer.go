package rx

//NewSink 订阅Observable，后续对sink发送Observer来拉取数据
func (ob Observable) NewSink() (sink Sink) {
	sink = make(Sink)
	go ob(sink)
	return sink
}

//Subscribe 异步订阅模式，返回一个可用于终止事件流的函数
func (ob Observable) Subscribe(onNext func(data interface{}), onComplete func(), onError func(error)) func() {
	source := ob.NewSink()
	var loop Observer
	loop = func(data interface{}, err error) {
		if err == nil {
			onNext(data)
			if loop != nil {
				source <- loop
			}
		} else if err == Complete {
			onComplete()
		} else {
			onError(err)
		}
	}
	source <- loop
	return func() {
		loop = nil
		close(source)
	}
}

//SubscribeOnCurrent 同步订阅模式，会阻塞当前goroutine，函数返回代表完成
func (ob Observable) SubscribeOnCurrent(onNext func(data interface{})) (err error) {
	source := ob.NewSink()
	var loop Observer
	next := make(chan interface{})
	complete := make(chan error)
	defer close(next)
	defer close(complete)
	defer close(source)
	loop = func(data interface{}, err error) {
		if err == nil {
			next <- data
			if loop != nil {
				source <- loop
			}
		} else if err == Complete {
			complete <- nil
		} else {
			complete <- err
		}
	}
	source <- loop
	for {
		select {
		case data := <-next:
			onNext(data)
		case err = <-complete:
			return
		}
	}
}
