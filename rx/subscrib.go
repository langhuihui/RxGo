package rx

//subscribe 内部同步订阅Observable
func (ob Observable) subscribe(onNext NextHandler, stop Stop) error {
	return ob(NewObserver(onNext, stop))
}

//subscribeAsync 对外异步订阅模式，返回一个可用于终止事件流的控制器,可以在未收到数据时也能终止事件流
func (ob Observable) subscribeAsync(onNext NextHandler, onComplete func(error)) *Observer {
	source := NewObserver(onNext, make(Stop))
	go func() {
		onComplete(ob(source))
	}()
	return source
}

//Subscribe 对外同步订阅Observable，可以在收到的事件中访问Control的Stop函数来终止事件流
func (ob Observable) Subscribe(onNext NextHandler) error {
	return ob.subscribe(onNext, make(Stop))
}

//SubscribeAsync 对外异步订阅模式，返回一个可用于终止事件流的控制器,可以在未收到数据时也能终止事件流
func (ob Observable) SubscribeAsync(onNext NextHandler, onComplete func(), onError func(error)) *Observer {
	source := NewObserver(onNext, make(Stop))
	go func() {
		if err := ob(source); err != nil {
			onError(err)
		} else {
			onComplete()
		}
	}()
	return source
}
