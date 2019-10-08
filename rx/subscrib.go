package rx

//Subscribe 对外同步订阅Observable，可以在收到的事件中访问Control的Stop函数来终止事件流
func (ob Observable) Subscribe(onNext NextHandler) error {
	return ob(NewObserver(onNext, make(Stop), make(Stop)))
}

//SubscribeAsync 对外异步订阅模式，返回一个可用于终止事件流的控制器,可以在未收到数据时也能终止事件流
func (ob Observable) SubscribeAsync(onNext NextHandler, onComplete func(), onError func(error)) *Observer {
	source := NewObserver(onNext, make(Stop), make(Stop))
	go func() {
		err := ob(source)
		select {
		case <-source.dispose:
			return //由于取消订阅而退出的，不调用onComplete和onError
		default:
		}
		if err != nil {
			onError(err)
		} else {
			onComplete()
		}
	}()
	return source
}
