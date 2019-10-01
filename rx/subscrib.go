package rx

//subscribe 内部同步订阅Observable
func (ob Observable) subscribe(observer Observer, stop Stop) error {
	return ob(NewControl(observer, stop))
}

//Subscribe 对外同步订阅Observable，可以在收到的事件中访问Control的Stop函数来终止事件流
func (ob Observable) Subscribe(observer Observer) error {
	return ob.subscribe(observer, make(Stop))
}

//SubscribeAsync 对外异步订阅模式，返回一个可用于终止事件流的控制器,可以在未收到数据时也能终止事件流
func (ob Observable) SubscribeAsync(onNext Observer, onComplete func(), onError func(error)) *Control {
	source := NewControl(onNext, make(Stop))
	go func() {
		if err := ob(source); err != nil {
			onError(err)
		} else {
			onComplete()
		}
	}()
	return source
}
